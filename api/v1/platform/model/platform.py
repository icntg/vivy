from sqlalchemy import Column, VARCHAR, TIMESTAMP, Text, CHAR, INTEGER

from api.v1.platform.model.__base__ import BaseModel


class Account(BaseModel):
    __tablename__ = 'pf_account'
    code = Column(VARCHAR(20), unique=True, nullable=False, doc='工号')
    name = Column(VARCHAR(10), nullable=False, doc='姓名')
    login_name = Column(VARCHAR(50), unique=True, doc='登录名/昵称')
    email = Column(VARCHAR(255), unique=True, doc='邮箱')
    telephone = Column(VARCHAR(50), doc='座机')
    mobile = Column(VARCHAR(50), doc='手机')
    avatar = Column(Text, doc='头像')
    password = Column(VARCHAR(255), doc='密码')
    token = Column(VARCHAR(50), doc='谷歌令牌')
    department_id = Column(CHAR(20), nullable=True, doc='部门ID')
    last_login = Column(TIMESTAMP, doc='最后登录时间')
    last_address = Column(VARCHAR(50), doc='最后登录IP')

class Department(BaseModel):
    __tablename__ = 'pf_department'
    parent_id = Column(CHAR(20), nullable=True, doc='上级部门ID')
    name = Column(VARCHAR(100), nullable=False, doc='部门名称')


class Role(BaseModel):
    __tablename__ = 'pf_role'
    parent_id = Column(CHAR(20), nullable=True, doc='上级角色ID')
    name = Column(VARCHAR(100), nullable=False, doc='角色名称')
    level = Column(INTEGER, nullable=False, doc='角色等级')


class AccountRole(BaseModel):
    __tablename__ = 'pf_account_role'
    account_id = Column(CHAR(20), nullable=False, index=True, doc='账号ID')
    role_id = Column(CHAR(20), nullable=False, doc='角色ID')
