from pathlib import Path
from typing import Optional

import yaml

from .constant import Constant
from .external.data_source import DataSource, MySQL
from .external.functions import err_put


class Config(Constant):
    def __init__(self) -> None:
        super().__init__()
        # path
        self.STATIC = str(Path(self.BASE) / Path('resource/static'))
        self.TEMPLATE = str(Path(self.BASE) / Path('resource/template'))
        # config
        self.HTTP_HOST: str = '127.0.0.1'
        self.HTTP_PORT: int = 8999
        self.SECRET_HEX: Optional[str] = None
        # logger
        self.LOGGER_DIRECTORY: str = './logs'
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
                filename = str(Path(self.BASE).joinpath(filename))
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
            err_put(f'[ERROR] cannot read config file [{filename}]: {e}\n')

    def write(self, filename='config.yaml'):
        # TODO:
        pass
