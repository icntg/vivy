from api.v1.platform.model.__base__ import BaseModel


class Task(BaseModel):  # 任务表
    """
    名称
    详情
    备注
    重要程度
    计划开始时间
    计划结束时间
    实际开始时间
    实际结束时间
    是否在日历上显示？
    上级任务
    前序任务
    来源任务（汇总时使用？）
    """
    __tablename__ = 'vv_task'


class Label(BaseModel):  # 标签表
    """
    必选标签：
    权限：全局可见、组内可见（分级：同级可见、上级可见、上上级可见？）、仅自己可见、
    可选标签：
    可分包
    可分解
    任务进度
    """
    __tablename__ = 'vv_label'


class EntityLabel(BaseModel):  # 实体标签表
    __tablename__ = 'vv_entity_label'

