"""
初始化工具。单独一个文件，用于交互生成配置。
需要在$BASE/conf目录下，生成
"""
import io
from pathlib import Path
from typing import Dict, Union
from urllib.parse import quote

from sanic import Request, response
from sqlalchemy.engine import Engine
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

from api.utility.constant import Constant
from api.utility.external.google_token import generate_token_and_qrcode
from api.v1.model.__base__ import BaseModel

STATIC = Path(Constant.BASE).joinpath('resource', 'static_initialize')
CONF = Path(Constant.BASE).joinpath('conf')


class Initialization:
    def __init__(self, conf_dict: Dict):
        from sanic import Sanic
        self.app = Sanic.get_app()
        self.cfg = conf_dict
        m = self.cfg['mysql']
        self._dsn: str = f'''
mysql+aiomysql://
{quote(m['username'])}:{quote(m['password'])}@
{m['host']}:{m['port']}?charset=utf8mb4
        '''.strip().replace('\n', '')
        self._engine: Engine = create_async_engine(self._dsn)
        self._session: Union[AsyncSession, sessionmaker] = sessionmaker(
            class_=AsyncSession,
            autocommit=False,
            autoflush=False,
            bind=self._engine,
        )
        self.logs = io.StringIO()

    async def do(self):
        self._create_database()
        self._create_tables()

    async def _create_database(self):
        m = self.cfg['mysql']
        if m['createDatabase']:
            self.logs.write('To create database\n')
            session: AsyncSession = yield self._session
            self.logs.write(f'yield session {id(session)}\n')
            sql = f'''CREATE DATABASE `{m['database']}` /*!40100 COLLATE 'utf8mb4_bin' */;'''
            await session.execute(sql)
            self.logs.write(f'create database with sql: {sql}\n')
            sql = f'''CREATE USER '{m['opsUsername']}'@'{m['opsUsernameIP']}' IDENTIFIED BY '{m['opsPassword']}';'''
            await session.execute(sql)
            self.logs.write(f'create user with sql: {sql}\n')
            sql = f'''GRANT ALL PRIVILEGES ON *.`{m['database']}` TO '{m['opsUsername']}'@'{m['opsUsernameIP']}';'''
            await session.execute(sql)
            self.logs.write(f'grant privileges with sql: {sql}\n')

    async def _create_tables(self):
        m = self.cfg['mysql']
        self.logs.write('To create tables\n')
        session: AsyncSession = yield self._session
        self.logs.write(f'yield session {id(session)}\n')
        sql = f'''use {m['database']};'''
        await session.execute(sql)
        self.logs.write(f'switch database with sql: {sql}\n')
        BaseModel.metadata.create_all(self._engine)
        self.logs.write('create tables with sqlalchemy engine\n')

    def _insert_admin(self):
        pass

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
        print(cfg)
        init = Initialization(cfg)
        await app.add_task(init.do())
        return response.json(dict(code=0, message='', logs=init.logs.getvalue()))

    app.run(
        '0.0.0.0',
        8081,
        access_log=True,
        debug=True,
        workers=1,
    )
