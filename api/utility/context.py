import binascii
import os
import sys
from concurrent.futures import Executor, ThreadPoolExecutor
from logging import Logger
from types import ModuleType
from typing import Union, Optional

import js2py

from .data.data_source import MySQL
from .external import logger
from api.utility.data.async_sqlalchemy import AsyncSQLAlchemy
from .external.functions import err_print


class Context(ModuleType):
    def __init__(self) -> None:
        super().__init__('context')
        self.config: Optional[js2py.base.JsObjectWrapper] = None

        self.secret: Optional[Union[bytes, bytearray]] = None

        self.AccessLogger: Optional[Logger] = None
        self.OutputLogger: Optional[Logger] = None
        self.SecureLogger: Optional[Logger] = None

        self.executor: Executor = ThreadPoolExecutor(os.cpu_count() * 5)

        self.DataSource: Optional[AsyncSQLAlchemy] = None

    def init_with_config(self, cfg: js2py.base.JsObjectWrapper):
        self.config = cfg
        try:
            self.secret = binascii.unhexlify(self.config.session.secret_hex)
            if len(self.secret) != 32:
                import hashlib
                self.secret = hashlib.sha256(self.secret)
        except Exception as e:
            _ = e
            err_print('[WARNING] cannot decode hex secret. to use random instead.\n')
            self.secret = os.urandom(32)


    def init_loggers(self):
        address = (self.config.dependency.rsyslog.host, self.config.dependency.rsyslog.port)
        try:
            self.OutputLogger = logger.init_logger(
                name='output',
                address=address,
                directory=self.config.setting.logger_directory,
                debug=self.config.setting.debug,
                formatter=self.config.setting.logger_formatter,
            )
        except Exception as e:
            err_print(f'[ERROR] cannot initialize OutputLogger, {type(e)}:{e}')
            sys.exit(-1)
        try:
            self.SecureLogger = logger.init_logger(
                name='secure',
                address=address,
                directory=self.config.setting.logger_directory,
                debug=self.config.setting.debug,
                formatter=self.config.setting.logger_formatter,
            )
        except Exception as e:
            err_print(f'[ERROR] cannot initialize SecureLogger, {type(e)}:{e}')
            sys.exit(-1)
        try:
            self.AccessLogger = logger.init_logger(
                name='access',
                address=address,
                directory=self.config.setting.logger_directory,
                debug=self.config.setting.debug,
                formatter=self.config.setting.logger_formatter,
            )
        except Exception as e:
            err_print(f'[ERROR] cannot initialize AccessLogger, {type(e)}:{e}')
            sys.exit(-1)

    def init_data_source(self):
        mysql: MySQL = MySQL(self.config.dependency.mysql)
        self.DataSource = AsyncSQLAlchemy(str(mysql))


sys.modules["ApplicationContext"] = Context()


def get_context() -> Context:
    context: Union[Context, ModuleType] = sys.modules["ApplicationContext"]
    return context
