from sqlalchemy import Column, VARCHAR, Text, CHAR, INTEGER

from api.v1.platform.model.__base__ import BaseModel


class Menu(BaseModel):
    __tablename__ = 'pf_menu'
    path = Column(Text, doc='菜单URL')
    text = Column(VARCHAR(255), doc='显示文本')
    parent_id = Column(CHAR(20), nullable=True, doc='上级菜单ID')
    sort = Column(INTEGER, doc='排序')


class RoleMenu(BaseModel):
    __tablename__ = 'pf_role_menu'
    role_id = Column(CHAR(20), nullable=False, index=True, doc='角色ID')
    menu_id = Column(CHAR(20), nullable=False, doc='菜单ID')
    methods = Column(VARCHAR(255), nullable=True, default='*',
                     doc='可用方法，*表示所有方法均可用。GET,POST表示只能使用GET和POST方法')
    options = Column(INTEGER, nullable=True, default=1, doc='子菜单是否继承，默认可继承')
