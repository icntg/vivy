from contextvars import ContextVar  # Python>=3.7

from sqlalchemy.engine import Engine
from sqlalchemy.ext.asyncio import create_async_engine


class AsyncSQLAlchemy:
    def __init__(self, data_source_url: str):
        self.engine: Engine = create_async_engine(data_source_url)
        self.base_model_session_ctx = ContextVar("session")
