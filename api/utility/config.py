import binascii
import os
from abc import ABCMeta, abstractmethod
from pathlib import Path
from typing import Optional, Dict, List

import yaml

from .constant import Constant
from api.utility.data.data_source import DataSource, MySQL
from .external.functions import err_print


class AbstractConfig(metaclass=ABCMeta):
    @abstractmethod
    def to_json(self) -> Dict:
        pass

    @abstractmethod
    def use_default_values(self):
        pass


class Http(Constant, AbstractConfig):
    def __init__(self) -> None:
        super(AbstractConfig).__init__()
        super(Constant).__init__()
        self.HOST: Optional[str] = None
        self.PORT: Optional[int] = None

    def to_json(self) -> Dict:
        return self.__dict__

    def use_default_values(self):
        self.HOST = '127.0.0.1'
        self.PORT = 8999


class Session(Constant, AbstractConfig):
    def __init__(self) -> None:
        super(AbstractConfig).__init__()
        super(Constant).__init__()
        self.SECRET_HEX: Optional[str] = None
        self.COOKIE: Optional[str] = None
        self.LOGIN_COOKIE_TIMEOUT: Optional[int] = None
        self.SESSION_TIMEOUT: Optional[int] = None

    def to_json(self) -> Dict:
        return self.__dict__

    def use_default_values(self):
        self.SECRET_HEX: Optional[str] = binascii.hexlify(os.urandom(32)).decode()
        self.COOKIE = 'PHPSESSID'
        self.LOGIN_COOKIE_TIMEOUT: int = 60 * 5  # 登录过程时间5分钟
        self.SESSION_TIMEOUT: int = 60 * 60 * 2  # 会话时间2个小时


class Setting(Constant, AbstractConfig):
    def __init__(self) -> None:
        super(AbstractConfig).__init__()
        super(Constant).__init__()
        # logger
        self.LOGGER_DIRECTORY: Optional[Path] = None
        self.LOGGER_FORMATTER: Optional[str] = None
        # setting
        self.DEBUG: Optional[bool] = None

    def to_json(self) -> Dict:
        ret = {}
        for k, v in self.__dict__.items():
            if isinstance(v, Path):
                ret[k] = str(v)
            else:
                ret[k] = v
        return ret

    def use_default_values(self):
        self.DEBUG = False
        self.LOGGER_DIRECTORY = self.BASE.joinpath('logs')
        self.LOGGER_FORMATTER = '[%(levelname)s]%(asctime)s[%(filename)s:%(lineno)s][%(name)s]: %(message)s'


class Config(Constant, AbstractConfig):
    def __init__(self) -> None:
        super(AbstractConfig).__init__()
        super(Constant).__init__()
        self.HTTP: Http = Http()
        self.SESSION: Session = Session()
        self.SETTING: Setting = Setting()
        self.DATA_SOURCES: Optional[List[DataSource]] = None

    def use_default_values(self):
        self.HTTP.use_default_values()
        self.SESSION.use_default_values()
        self.SETTING.use_default_values()
        self.DATA_SOURCES = [MySQL(dict(
            host='127.0.0.1',
            port=3306,
            username='vivy',
            password='vivy',
            database='vivy',
            option='?parseTime=true&charset=utf8mb4&loc=Local',
            maxIdle=10,
            maxOpen=100,
            showSql=False,
            showExecTime=False,
        ))]

    def to_json(self) -> Dict:
        ret = {}
        for k, v in self.__dict__.items():
            if isinstance(v, AbstractConfig):
                x: AbstractConfig = v
                ret[k] = x.to_json()
        ret['DATA_SOURCES'] = {}
        for ds in self.DATA_SOURCES:
            k = str(ds.__class__)[8:-2].split('.')[-1]
            ret['DATA_SOURCES'][k] = ds.__dict__
        return ret

    def read(self, filename='config.yaml'):
        try:
            if not Path(filename).is_absolute():
                filename = str(Path(self.BASE).joinpath('conf', filename))
            cfg = yaml.load(open(filename, 'r', encoding='utf-8').read(), Loader=yaml.Loader)
            if 'HTTP' in cfg and 'HOST' in cfg['HTTP']:
                self.HTTP.__dict__['HOST'] = cfg['HTTP']['HOST']
            if 'HTTP' in cfg and 'PORT' in cfg['HTTP']:
                self.HTTP.__dict__['PORT'] = cfg['HTTP']['PORT']
            if 'SESSION' in cfg and 'SECRET_HEX' in cfg['SESSION']:
                self.SESSION.__dict__['SECRET_HEX'] = cfg['SESSION']['SECRET_HEX']
            if 'SETTING' in cfg and 'LOGGER_DIRECTORY' in cfg['SETTING']:
                self.SETTING.__dict__['LOGGER_DIRECTORY'] = cfg['SETTING']['LOGGER_DIRECTORY']
            if 'SETTING' in cfg and 'LOGGER_FORMATTER' in cfg['SETTING']:
                self.SETTING.__dict__['LOGGER_FORMATTER'] = cfg['SETTING']['LOGGER_FORMATTER']
            if 'SETTING' in cfg and 'DEBUG' in cfg['SETTING']:
                self.SETTING.__dict__['DEBUG'] = cfg['SETTING']['DEBUG']
            if 'DATA_SOURCES' in cfg and 'MySQL' in cfg['DATA_SOURCES']:
                self.__dict__['DATA_SOURCES'] = []
                self.DATA_SOURCES.append(MySQL(cfg['DATA_SOURCES']['MySQL']))
        except Exception as e:
            err_print(f'[ERROR] cannot read config file [{filename}]: {e}\n')
            self.use_default_values()

    def write(self, filename='config.yaml'):
        data = self.to_json()
        with open(filename, "w", encoding="utf-8") as f:
            yaml.dump(data, f)
