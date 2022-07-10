"""
初始化工具。单独一个文件，用于交互生成配置。
需要在$BASE/conf目录下，生成
config.yaml
initialize.log
共2个文件。
"""
import asyncio
import io
from pathlib import Path
from typing import Dict, Optional
from urllib.parse import quote

import aiomysql
from sanic import Request, response
from sqlalchemy.engine import Engine
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

from api.utility.constant import Constant
from api.utility.data.data_source import MySQL
from api.utility.external.functions import object_id
from api.utility.external.google_token import generate_token_and_qrcode
from api.utility.external.pynacl_util import password_hash

STATIC = Path(Constant.BASE).joinpath('resource', 'static_initialize')
CONF = Path(Constant.BASE).joinpath('conf')


class Initialization:
    def __init__(self, conf_dict: Dict):
        self.cfg: Dict = conf_dict
        self.logs: io.StringIO = io.StringIO()
        self.loop = asyncio.get_event_loop()
        self.pool: Optional[aiomysql.Pool] = None
        self.engine: Optional[Engine] = None

    async def aio_mysql_connect(self) -> bool:
        """
        创建数据库连接池（异步）
        由于需要异步创建，无法写在__init__函数中。需要单独进行。
        原生连接用于建立数据库和数据库操作用户。
        """
        m = self.cfg['mysql']
        try:
            self.logs.write('To connect database with aio_mysql\n')
            self.pool = await aiomysql.create_pool(
                host=m['host'],
                port=m['port'], user=m['username'],
                password=m['password'],
                loop=self.loop,
                # cursorclass=aiomysql.DictCursor,
                cursorclass=aiomysql.SSCursor,
                echo=True,
            )
            self.logs.write('success\n')
            return True
        except Exception as e:
            self.logs.write(f'aio_mysql_connect failed: {e}\n')
            return False

    async def alchemy_connect(self):
        """
        创建异步SQLAlchemy连接。
        SQLAlchemy用于创建表和管理员信息。
        """
        m = self.cfg['mysql']
        if m['createDatabase']:
            username = m['opsUsername']
            password = m['opsPassword']
        else:
            username = m['username']
            password = m['password']
        dsn: str = f'''
                        mysql+aiomysql://
        {quote(username)}:{quote(password)}@
        {m['host']}:{m['port']}/{m['database']}?charset=utf8mb4
        '''.strip().replace('\n', '').replace(' ', '')
        try:
            self.logs.write('To connect database with SQLAlchemy\n')
            self.engine = create_async_engine(dsn, echo=True)
            self.logs.write('success\n')
            return True
        except Exception as e:
            self.logs.write(f'alchemy_connect failed: {e}\n')
            return False

    def __del__(self):
        """
        关闭数据库连接池
        """
        if self.pool is not None:
            self.pool.close()

    async def create_database(self) -> bool:
        m = self.cfg['mysql']
        if m['createDatabase']:
            async def run_sql(_sql) -> bool:
                try:
                    async with self.pool.acquire() as conn:
                        async with conn.cursor(aiomysql.SSCursor) as cur:
                            await cur.execute(_sql)
                            self.logs.write(f'success\n')
                            return True
                except Exception as e:
                    self.logs.write(f'error: {e}\n')
                    return False

            ret = True
            self.logs.write('To create database\n')
            sql = f'''CREATE DATABASE IF NOT EXISTS `{m['database']}` /*!40100 COLLATE 'utf8mb4_bin' */;'''
            self.logs.write(f'create database with sql: {sql}\n')
            ret &= await run_sql(sql)

            sql = f'''
CREATE USER IF NOT EXISTS '{m['opsUsername']}'@'{m['opsUsernameIP']}' IDENTIFIED BY '{m['opsPassword']}';
'''.strip()
            sql_mask = f'''
CREATE USER IF NOT EXISTS '{m['opsUsername']}'@'{m['opsUsernameIP']}' IDENTIFIED BY '**************';
'''.strip()
            self.logs.write(f'create user with sql: {sql_mask}\n')
            ret &= await run_sql(sql)

            sql = f'''GRANT ALL PRIVILEGES ON `{m['database']}`.* TO '{m['opsUsername']}'@'{m['opsUsernameIP']}';'''
            self.logs.write(f'grant privileges with sql: {sql}\n')
            ret &= await run_sql(sql)
            return ret

    async def create_tables(self) -> bool:
        self.logs.write('To create tables\n')

        import api.v1.model.platform as _
        import api.v1.platform.model.__base__
        base = api.v1.platform.model.__base__.Base
        try:
            async with self.engine.begin() as conn:
                await conn.run_sync(base.metadata.create_all)
                self.logs.write(f'success\n')
                return True
        except Exception as e:
            self.logs.write(f'error: {e}\n')
            return False

    async def insert_admin(self) -> bool:
        self.logs.write('To insert admin\n')
        u = self.cfg['admin']
        import api.v1.model.platform

        role = api.v1.model.platform.platform.Role()
        role.id = object_id()
        role.name = '系统管理员'
        role.level = 9999

        admin = api.v1.model.platform.platform.Account()
        admin.id = object_id()
        admin.code = u['loginName']
        admin.login_name = u['loginName']
        admin.password = password_hash(u['password'])
        admin.token = u['token'].replace(' ', '')
        admin.comment = '系统管理员'
        admin.name = '系统管理员'

        account_role = api.v1.model.platform.platform.AccountRole()
        account_role.id = object_id()
        account_role.account_id = admin.id
        account_role.role_id = role.id

        try:
            async_session = sessionmaker(self.engine, expire_on_commit=False, class_=AsyncSession)
            async with async_session() as session:
                async with session.begin():
                    session.add_all([role, admin, account_role])
            self.logs.write(f'success\n')
            return True
        except Exception as e:
            self.logs.write(f'error: {e}\n')
            return False

    def write_logs(self) -> bool:
        try:
            self.logs.write(f'To write {Constant.INIT}\n')
            if not Constant.INIT.parent.exists():
                Constant.INIT.parent.mkdir(0o755, True, True)
            with open(Constant.INIT, 'w') as f:
                f.write(self.logs.getvalue())
            self.logs.write('success\n')
            return True
        except Exception as e:
            self.logs.write(f'error: {e}\n')
            return False

    def write_config(self) -> bool:
        try:
            self.logs.write(f'To write {Constant.CONF}\n')
            if not Constant.CONF.parent.exists():
                Constant.CONF.parent.mkdir(0o755, True, True)
            cfg: Config = Config()
            cfg.use_default_values()
            cfg.HTTP.__dict__['HOST'] = self.cfg['http']['host']
            cfg.HTTP.__dict__['PORT'] = self.cfg['http']['port']
            if not self.cfg['http']['randomSecret']:
                cfg.SESSION.__dict__['SECRET_HEX'] = self.cfg['http']['secretHex']

            m = self.cfg['mysql']
            if m['createDatabase']:
                username = m['opsUsername']
                password = m['opsPassword']
            else:
                username = m['username']
                password = m['password']
            cfg.DATA_SOURCES.clear()
            cfg.DATA_SOURCES.append(MySQL(dict(
                host=m['host'],
                port=m['port'],
                username=username,
                password=password,
                database=m['database'],
                option='?charset=utf8mb4',
                maxIdle=10,
                maxOpen=100,
                showSql=False,
                showExecTime=False,
            )))
            cfg.write(Constant.CONF)
            f = Constant.CONF.parent.joinpath('block.list')
            if not f.exists():
                with open(f, 'w') as fp:
                    fp.write('\n')
            f = Constant.CONF.parent.joinpath('allow.list')
            if not f.exists():
                with open(f, 'w') as fp:
                    fp.write('127.0.0.1\n')
            self.logs.write('success\n')
            return True
        except Exception as e:
            self.logs.write(f'error: {e}\n')
            return False


def create_and_run():
    from sanic import Sanic

    app = Sanic("initialize")
    app.static('/', str(STATIC), stream_large_files=True)

    @app.get('/')
    async def main_page(_: Request):
        return response.html(open(str(STATIC.joinpath('index.html')), 'rb').read())

    @app.post('/api/token')
    async def generate_token(req: Request):
        login_name = req.body.decode()
        token, img = generate_token_and_qrcode('星轨', login_name, '温州信通')
        return response.json(dict(token=token, qrcode=img))

    @app.post('/api/initialization')
    async def initialize(req: Request):
        cfg = req.json

        import json
        print(json.dumps(cfg, ensure_ascii=False, indent=2))

        init = Initialization(cfg)
        ret = await init.aio_mysql_connect()
        if not ret:
            return response.json(dict(code=1, message='aio_mysql_connect', logs=init.logs.getvalue()))
        ret = await init.create_database()
        if not ret:
            return response.json(dict(code=2, message='create_database', logs=init.logs.getvalue()))
        ret = await init.alchemy_connect()
        if not ret:
            return response.json(dict(code=3, message='alchemy_connect', logs=init.logs.getvalue()))
        ret = await init.create_tables()
        if not ret:
            return response.json(dict(code=4, message='create_tables', logs=init.logs.getvalue()))
        ret = await init.insert_admin()
        if not ret:
            return response.json(dict(code=5, message='insert_admin', logs=init.logs.getvalue()))
        ret = init.write_config()
        if not ret:
            return response.json(dict(code=6, message='write_config', logs=init.logs.getvalue()))
        ret = init.write_logs()
        if not ret:
            return response.json(dict(code=7, message='write_logs', logs=init.logs.getvalue()))

        return response.json(dict(code=0, message='', logs=init.logs.getvalue()))

    app.run(
        '0.0.0.0',
        8081,
        access_log=True,
        debug=True,
        workers=1,
    )


if __name__ == '__main__':
    create_and_run()
