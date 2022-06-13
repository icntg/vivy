from pathlib import Path


class Constant:
    BASE: Path = Path(__file__).absolute().parent.parent.parent
    INIT: Path = BASE.joinpath('conf', 'initialize.log')
    CONF: Path = BASE.joinpath('conf', 'config.yaml')
    STATIC: Path = BASE.joinpath('resource', 'static')
    TEMPLATE: Path = BASE.joinpath('resource', 'template')

    def __init__(self) -> None:
        pass

    class ConstError(TypeError):
        pass

    def __setattr__(self, name, value):
        if name in self.__dict__ and self.__dict__[name] is not None:
            raise self.ConstError("cannot rebind constant (%s)" % name)
        self.__dict__[name] = value
