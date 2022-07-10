import time
from pathlib import Path


class Constant:
    def __init__(self) -> None:
        self.BASE: Path = Path(__file__).absolute().parent.parent.parent
        self.LOG_DIR: Path = self.BASE.joinpath('logs')
        self.INIT: Path = self.BASE.joinpath('conf', 'initialize.log')
        self.CONF: Path = self.BASE.joinpath('conf', 'config.js')
        self.STATIC: Path = self.BASE.joinpath('resource', 'static')
        self.TEMPLATE: Path = self.BASE.joinpath('resource', 'template')
        self.TIME_FMT: str = '%Y-%m-%d %H:%M:%S'
        # 默认应该就是国际标准时间？
        self.TIME0: int = int(time.mktime(time.strptime('2022-01-01 00:00:00 UTC', '%Y-%m-%d %H:%M:%S %Z')))
        self.SESSION_NAME_LOGIN: str = 'SessionNameLogin'
        self.SESSION_NAME_CURRENT_ACCOUNT: str = 'SessionNameCurrentAccount'

    class ConstError(TypeError):
        pass

    def __setattr__(self, name, value):
        if name in self.__dict__ and self.__dict__[name] is not None:
            raise self.ConstError("cannot rebind constant (%s)" % name)
        self.__dict__[name] = value


# 模块默认就是单例模式。
constant: Constant = Constant()
