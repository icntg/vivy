from contextvars import ContextVar  # Python>=3.7
from typing import Union

from sqlalchemy.engine import Engine
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker

from api.utility.data.data_source import DataSource


class AsyncSQLAlchemy:
    def __init__(self, data_source_url: str):
        from api import Context, get_context
        ctx: Context = get_context()
        from api.utility.data.data_source import MySQL
        mysql_cfg: Union[MySQL, DataSource] = ctx.config.DATA_SOURCES[0]
        self.engine: Engine = create_async_engine(data_source_url, echo=mysql_cfg.show_sql)
        self.base_model_session_ctx = ContextVar("session")

    def init_middleware(self):
        from sanic import Sanic
        app: Sanic = Sanic.get_app()

        @app.middleware("request")
        async def inject_session(request):
            request.ctx.session = sessionmaker(self.engine, AsyncSession, expire_on_commit=False)()
            request.ctx.session_ctx_token = self.base_model_session_ctx.set(request.ctx.session)

        @app.middleware("response")
        async def close_session(request, response):
            if hasattr(request.ctx, "session_ctx_token"):
                self.base_model_session_ctx.reset(request.ctx.session_ctx_token)
                await request.ctx.session.close()
