from sqlalchemy import INTEGER, Column, CHAR, TIMESTAMP, Text
from sqlalchemy.orm import declarative_base

Base = declarative_base()


class BaseModel(Base):
    __abstract__ = True
    db_id = Column(INTEGER(), primary_key=True, autoincrement=True)
    id = Column(CHAR(20), unique=True, nullable=False)
    create_at = Column(TIMESTAMP, nullable=False)
    modify_at = Column(TIMESTAMP, nullable=False)
    delete_at = Column(TIMESTAMP, nullable=True, index=True)
    comment = Column(Text, nullable=True)
