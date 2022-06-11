from sqlalchemy import Column, VARCHAR, TIMESTAMP, Text

from .__base__ import BaseModel


class Account(BaseModel):
    __tablename__ = 'account'
    code = Column(VARCHAR(20), unique=True, nullable=False, doc='工号')
    name = Column(VARCHAR(10), nullable=False, doc='姓名')
    login_name = Column(VARCHAR(50), unique=True, doc='登录名/昵称')
    email = Column(VARCHAR(255), unique=True, doc='邮箱')
    telephone = Column(VARCHAR(50), doc='座机')
    mobile = Column(VARCHAR(50), doc='手机')
    avatar = Column(Text, doc='头像')
    password = Column(VARCHAR(255), doc='密码')
    token = Column(VARCHAR(50), doc='谷歌令牌')
    last_login = Column(TIMESTAMP, doc='最后登录时间')
    last_address = Column(VARCHAR(50), doc='最后登录IP')

# class Department(BaseModel):
#     pass
#
#
# class Role(BaseModel):
#     pass
#
#
# class AccountRole(BaseModel):
#     pass
#
#
# class Resource(BaseModel):
#     pass
#
#
# class RoleResource(BaseModel):
#     pass
