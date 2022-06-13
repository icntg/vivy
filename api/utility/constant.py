import time
from pathlib import Path


class Constant:
    BASE: Path = Path(__file__).absolute().parent.parent.parent
    INIT: Path = BASE.joinpath('conf', 'initialize.log')
    CONF: Path = BASE.joinpath('conf', 'config.yaml')
    STATIC: Path = BASE.joinpath('resource', 'static')
    TEMPLATE: Path = BASE.joinpath('resource', 'template')
    TIME_FMT: str = '%Y-%m-%d %H:%M:%S'
    TIME0: int = int(time.mktime(time.strptime('2022-01-01 00:00:00', TIME_FMT)))  # 默认应该就是国际标准时间？

    def __init__(self) -> None:
        pass

    class ConstError(TypeError):
        pass

    def __setattr__(self, name, value):
        if name in self.__dict__ and self.__dict__[name] is not None:
            raise self.ConstError("cannot rebind constant (%s)" % name)
        self.__dict__[name] = value
