from pathlib import Path
from typing import Optional


class Constant:
    BASE: str = str(Path(__file__).absolute().parent.parent.parent)

    def __init__(self) -> None:
        pass

    class ConstError(TypeError):
        pass

    def __setattr__(self, name, value):
        if name in self.__dict__ and self.__dict__[name] is not None:
            raise self.ConstError("cannot rebind constant (%s)" % name)
        self.__dict__[name] = value
