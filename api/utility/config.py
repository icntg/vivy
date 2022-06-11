import binascii
import os
from pathlib import Path
from typing import Optional

import yaml

from .constant import Constant
from .external.data_source import DataSource, MySQL
from .external.functions import err_print


class Config(Constant):
    def __init__(self) -> None:
        super().__init__()
        # web
        self.COOKIE: Optional[str] = None
        # path
        self.STATIC: Optional[str] = None
        self.TEMPLATE: Optional[str] = None
        # http
        self.HTTP_HOST: Optional[str] = None
        self.HTTP_PORT: Optional[int] = None
        self.SECRET_HEX: Optional[str] = None
        # session
        self.LOGIN_COOKIE_TIMEOUT: Optional[int] = None
        self.SESSION_TIMEOUT: Optional[int] = None
        # logger
        self.LOGGER_DIRECTORY: Optional[str] = None
        self.LOGGER_FORMATTER: Optional[str] = None
        # setting
        self.DEBUG: Optional[bool] = None
        self.DATASOURCE: Optional[DataSource] = None

    def use_default_value(self):
        # web
        self.COOKIE = 'PHPSESSID'
        # path
        self.STATIC = str(Path(self.BASE).joinpath('resource', 'static'))
        self.TEMPLATE = str(Path(self.BASE).joinpath('resource', 'template'))
        # http
        self.HTTP_HOST = '127.0.0.1'
        self.HTTP_PORT = 8999
        self.SECRET_HEX: Optional[str] = binascii.hexlify(os.urandom(32)).decode()
        # session
        self.LOGIN_COOKIE_TIMEOUT: int = 60 * 5  # 登录过程时间5分钟
        self.SESSION_TIMEOUT: int = 60 * 60 * 2  # 会话时间2个小时
        # logger
        self.LOGGER_DIRECTORY: str = str(Path(self.BASE).joinpath('logs'))
        self.LOGGER_FORMATTER: str = '[%(levelname)s]%(asctime)s[%(filename)s:%(lineno)s][%(name)s]: %(message)s'
        # setting
        self.DEBUG: bool = False
        self.DATASOURCE: DataSource = MySQL(dict(
            host='127.0.0.1',
            port=3306,
            username='root',
            password='root',
            database='vivy',
            option='?parseTime=true&charset=utf8mb4&loc=Local',
            maxIdle=10,
            maxOpen=100,
            showSql=False,
            showExecTime=False,
        ))


    def read(self, filename='config.yaml'):
        try:
            if not Path(filename).is_absolute():
                filename = str(Path(self.BASE).joinpath('conf', filename))
            cfg = yaml.load(open(filename, 'r', encoding='utf-8').read(), Loader=yaml.Loader)
            if 'http' in cfg and 'host' in cfg['http']:
                # self.HTTP_HOST = cfg['http']['host']
                self.__dict__['HTTP_HOST'] = cfg['http']['host']
            if 'http' in cfg and 'port' in cfg['http']:
                # self.HTTP_PORT = cfg['http']['port']
                self.__dict__['HTTP_PORT'] = cfg['http']['port']
            if 'http' in cfg and 'secret' in cfg['http']:
                # self.SECRET_HEX = cfg['http']['secret']
                self.__dict__['SECRET_HEX'] = cfg['http']['secret']
            if 'logger' in cfg and 'directory' in cfg['logger']:
                # self.LOGGER_DIRECTORY = cfg['logger']['directory']
                self.__dict__['LOGGER_DIRECTORY'] = cfg['logger']['directory']
            if 'logger' in cfg and 'formatter' in cfg['logger']:
                # self.LOGGER_FORMATTER = cfg['logger']['formatter']
                self.__dict__['LOGGER_FORMATTER'] = cfg['logger']['formatter']
            if 'debug' in cfg:
                # self.DEBUG = cfg['debug']
                self.__dict__['DEBUG'] = cfg['debug']
            if 'dataSource' in cfg and 'MySQL' in cfg['dataSource']:
                # self.DATASOURCE = MySQL(cfg['dataSource'])
                self.__dict__['DATASOURCE'] = MySQL(cfg['dataSource']['MySQL'])
        except Exception as e:
            err_print(f'[ERROR] cannot read config file [{filename}]: {e}\n')
            self.use_default_value()

    def write(self, filename='config.yaml'):
        from copy import deepcopy
        data = deepcopy(self.__dict__)
        # todo: del data some p
        with open(filename, "w", encoding="utf-8") as f:
            yaml.dump(data, f)
