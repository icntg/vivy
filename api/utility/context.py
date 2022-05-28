import binascii
import os
import sys
from concurrent.futures import Executor, ThreadPoolExecutor
from logging import Logger
from types import ModuleType
from typing import Union, Optional

from .config import Config
from .external import logger
from .external.async_sqlalchemy import AsyncSQLAlchemy
from .external.functions import err_put


class Context(ModuleType):
    def __init__(self) -> None:
        super(ModuleType).__init__('context', '')
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
            self.secret = binascii.unhexlify(self.config.SECRET_HEX)
        except Exception as e:
            _ = e
            err_put('[WARNING] cannot decode hex secret. to use random instead.\n')
            self.secret = os.urandom(64)

    def init_loggers(self):
        try:
            self.OutputLogger = logger.init_logger(
                'output',
                self.config.LOGGER_DIRECTORY,
                self.config.DEBUG,
                self.config.LOGGER_FORMATTER,
            )
        except Exception as e:
            err_put(f'[ERROR] cannot initialize OutputLogger, {type(e)}:{e}')
            sys.exit(-1)
        try:
            self.SecureLogger = logger.init_logger(
                'secure',
                self.config.LOGGER_DIRECTORY,
                self.config.DEBUG,
                self.config.LOGGER_FORMATTER,
            )
        except Exception as e:
            err_put(f'[ERROR] cannot initialize SecureLogger, {type(e)}:{e}')
            sys.exit(-1)
        try:
            self.AccessLogger = logger.init_logger(
                'access',
                self.config.LOGGER_DIRECTORY,
                self.config.DEBUG,
                self.config.LOGGER_FORMATTER,
            )
        except Exception as e:
            err_put(f'[ERROR] cannot initialize AccessLogger, {type(e)}:{e}')
            sys.exit(-1)

    def init_data_source(self):
        self.DataSource = AsyncSQLAlchemy(str(self.config.DATASOURCE))


sys.modules["ApplicationContext"] = Context()


def get_context() -> Context:
    context: Union[Context, ModuleType] = sys.modules["ApplicationContext"]
    return context
