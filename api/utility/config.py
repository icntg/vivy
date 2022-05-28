from typing import Optional

from .constant import Constant


class Config(Constant):
    def __init__(self) -> None:
        super().__init__()
        # config
        self.HTTP_HOST: str = '127.0.0.1'
        self.HTTP_PORT: int = 8999
        self.SECRET_HEX: Optional[str] = None
        # logger
        self.LOGGER_FILE: Optional[str] = None
        self.LOGGER_FORMATTER: str = '[%(levelname)s]%(asctime)s[%(filename)s:%(lineno)s][%(name)s]: %(message)s'
        # setting
        self.DEBUG: Optional[bool] = None

    def read(self, filename='config.yaml'):
        pass

    def write(self, filename='config.yaml'):
        pass
