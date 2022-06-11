from sqlalchemy import INTEGER, Column, CHAR, TIMESTAMP, Text, func
from sqlalchemy.orm import declarative_base

Base = declarative_base()


class BaseModel(Base):
    __abstract__ = True
    db_id = Column(INTEGER(), primary_key=True, autoincrement=True, doc='数据库ID')
    id = Column(CHAR(20), unique=True, nullable=False, doc='业务逻辑ID')
    create_at = Column(TIMESTAMP, nullable=True, server_default=func.now(), doc='创建时间')
    modify_at = Column(TIMESTAMP, nullable=True, server_default=func.now(), onupdate=func.now(), doc='修改时间')
    # 用于记录删除时间，数据并非真正删除。可以定期清理。
    delete_at = Column(TIMESTAMP, nullable=True, index=True, doc='删除时间')
    comment = Column(Text, nullable=True, doc='备注')
