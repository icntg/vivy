from sqlalchemy.engine import Engine
from sqlalchemy.ext.asyncio import create_async_engine, AsyncSession
from sqlalchemy.orm import sessionmaker


class AsyncSQLAlchemy:
    def __init__(self, data_source_url: str):
        self.engine: Engine = create_async_engine(data_source_url)
        self.session_local = sessionmaker(
            class_=AsyncSession,
            autocommit=False,
            autoflush=False,
            bind=self.engine,
        )

    async def async_session(self) -> AsyncSession:
        async with self.session_local() as session:
            yield session
