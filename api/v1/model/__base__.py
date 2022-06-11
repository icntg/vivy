from sqlalchemy import INTEGER, Column, CHAR, TIMESTAMP, Text, func
from sqlalchemy.orm import declarative_base

Base = declarative_base()


class BaseModel(Base):
    __abstract__ = True
    db_id = Column(INTEGER(), primary_key=True, autoincrement=True)
    id = Column(CHAR(20), unique=True, nullable=False)
    create_at = Column(TIMESTAMP, nullable=True, server_default=func.now())
    modify_at = Column(TIMESTAMP, nullable=True, server_default=func.now(), onupdate=func.now())
    delete_at = Column(TIMESTAMP, nullable=True, index=True)
    comment = Column(Text, nullable=True)
