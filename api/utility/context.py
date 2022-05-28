import sys
from logging import Logger
from types import ModuleType
from typing import Union, Optional

from .config import Config


class Context(ModuleType):
    def __init__(self) -> None:
        super(ModuleType).__init__('context', '')
        self.secret: Optional[Union[bytes, bytearray]] = None

        self.AccessLogger: Optional[Logger] = None
        self.OutputLogger: Optional[Logger] = None
        self.SecureLogger: Optional[Logger] = None

    def init_with_config(self, cfg: Config):
        pass


sys.modules["ApplicationContext"] = Context()


def get_context() -> Context:
    context: Union[Context, ModuleType] = sys.modules["ApplicationContext"]
    return context
