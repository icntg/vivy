from sqlalchemy import Column, VARCHAR, TIMESTAMP, Text, CHAR, INTEGER

from .__base__ import BaseModel


class Account(BaseModel):
    __tablename__ = 'vv_task'
