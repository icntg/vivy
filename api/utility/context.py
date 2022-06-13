import binascii
import os
import sys
from concurrent.futures import Executor, ThreadPoolExecutor
from logging import Logger
from types import ModuleType
from typing import Union, Optional

from .config import Config
from .external import logger
from api.utility.data.async_sqlalchemy import AsyncSQLAlchemy
from .external.functions import err_print


class Context(ModuleType):
    def __init__(self) -> None:
        super().__init__('context')
        self.config: Optional[Config] = None

        self.secret: Optional[Union[bytes, bytearray]] = None

        self.AccessLogger: Optional[Logger] = None
        self.OutputLogger: Optional[Logger] = None
        self.SecureLogger: Optional[Logger] = None

        self.executor: Executor = ThreadPoolExecutor(os.cpu_count() * 5)

        self.DataSource: Optional[AsyncSQLAlchemy] = None

    def init_with_config(self, cfg: Config):
        self.config = cfg
        try:
            self.secret = binascii.unhexlify(self.config.SESSION.SECRET_HEX)
            if len(self.secret) != 32:
                import hashlib
                self.secret = hashlib.sha256(self.secret)
        except Exception as e:
            _ = e
            err_print('[WARNING] cannot decode hex secret. to use random instead.\n')
            self.secret = os.urandom(32)

    def init_loggers(self):
        try:
            self.OutputLogger = logger.init_logger(
                'output',
                str(self.config.SETTING.LOGGER_DIRECTORY),
                self.config.SETTING.DEBUG,
                self.config.SETTING.LOGGER_FORMATTER,
            )
        except Exception as e:
            err_print(f'[ERROR] cannot initialize OutputLogger, {type(e)}:{e}')
            sys.exit(-1)
        try:
            self.SecureLogger = logger.init_logger(
                'secure',
                str(self.config.SETTING.LOGGER_DIRECTORY),
                self.config.SETTING.DEBUG,
                self.config.SETTING.LOGGER_FORMATTER,
            )
        except Exception as e:
            err_print(f'[ERROR] cannot initialize SecureLogger, {type(e)}:{e}')
            sys.exit(-1)
        try:
            self.AccessLogger = logger.init_logger(
                'access',
                str(self.config.SETTING.LOGGER_DIRECTORY),
                self.config.SETTING.DEBUG,
                self.config.SETTING.LOGGER_FORMATTER,
            )
        except Exception as e:
            err_print(f'[ERROR] cannot initialize AccessLogger, {type(e)}:{e}')
            sys.exit(-1)

    def init_data_source(self):
        self.DataSource = AsyncSQLAlchemy(str(self.config.DATA_SOURCES[0]))


sys.modules["ApplicationContext"] = Context()


def get_context() -> Context:
    context: Union[Context, ModuleType] = sys.modules["ApplicationContext"]
    return context
