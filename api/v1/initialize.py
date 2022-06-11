"""
初始化工具。单独一个文件，用于交互生成配置。
需要在$BASE/conf目录下，生成
"""
import asyncio
import io
from pathlib import Path
from typing import Dict, Optional
from urllib.parse import quote

import aiomysql
import nacl.pwhash
from sanic import Request, response
from sqlalchemy.engine import Engine
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

from api.utility.constant import Constant
from api.utility.external import base32, identity
from api.utility.external.google_token import generate_token_and_qrcode

STATIC = Path(Constant.BASE).joinpath('resource', 'static_initialize')
CONF = Path(Constant.BASE).joinpath('conf')


class Initialization:
    def __init__(self, conf_dict: Dict):
        self.cfg: Dict = conf_dict
        self.logs: io.StringIO = io.StringIO()
        self.loop = asyncio.get_event_loop()
        self.pool: Optional[aiomysql.Pool] = None
        self.engine: Optional[Engine] = None

    async def aio_mysql_connect(self):
        """
        创建数据库连接池（异步）
        由于需要异步创建，无法写在__init__函数中。需要单独进行。
        """
        m = self.cfg['mysql']
        self.pool = await aiomysql.create_pool(
            host=m['host'],
            port=m['port'], user=m['username'],
            password=m['password'],
            loop=self.loop,
            cursorclass=aiomysql.DictCursor,
            # cursorclass=aiomysql.SSCursor,
            echo=True,
        )

    async def alchemy_connect(self):
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
        self.engine = create_async_engine(dsn, echo=True)

    def __del__(self):
        """
        关闭数据库连接池
        """
        if self.pool is not None:
            self.pool.close()

    def check_config(self) -> Optional[str]:
        """
        检查传入的配置文件，是否存在业务错误。
        """
        pass

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
            self.logs.write(f'create user with sql: {sql}\n')
            ret &= await run_sql(sql)

            sql = f'''GRANT ALL PRIVILEGES ON `{m['database']}`.* TO '{m['opsUsername']}'@'{m['opsUsernameIP']}';'''
            self.logs.write(f'grant privileges with sql: {sql}\n')
            ret &= await run_sql(sql)
            return ret

    async def create_tables(self) -> bool:
        self.logs.write('To create tables\n')

        import api.v1.model.platform as _
        import api.v1.model.__base__
        base = api.v1.model.__base__.Base
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
        admin = api.v1.model.platform.Account()
        admin.id = base32.encode_for_id(identity.ObjectId.generate()).decode()
        admin.code = u['loginName']
        admin.login_name = u['loginName']
        admin.password = nacl.pwhash.str(u['password'].encode()).decode()
        admin.token = u['token'].replace(' ', '')
        admin.comment = '系统管理员'
        admin.name = '系统管理员'
        # now = int(time.time())
        # admin.create_at = now
        # admin.modify_at = now

        try:
            async_session = sessionmaker(self.engine, expire_on_commit=False, class_=AsyncSession)
            async with async_session() as session:
                async with session.begin():
                    session.add(admin)
            self.logs.write(f'success\n')
            return True
        except Exception as e:
            self.logs.write(f'error: {e}\n')
            return False

    def _write_logs(self):
        pass

    def _write_config(self):
        pass


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
        await init.aio_mysql_connect()
        await init.create_database()

        await init.alchemy_connect()
        await init.create_tables()
        await init.insert_admin()

        return response.json(dict(code=0, message='', logs=init.logs.getvalue()))

    app.run(
        '0.0.0.0',
        8081,
        access_log=True,
        debug=True,
        workers=1,
    )
