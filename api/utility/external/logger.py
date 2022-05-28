import logging
from logging import Logger, StreamHandler, Handler, Formatter
from logging.handlers import TimedRotatingFileHandler
from pathlib import Path
from typing import List

from api.utility.constant import Constant


def init_logger(
        name: str,
        directory: str,
        debug: bool = False,
        formatter: str = '[%(levelname)s]%(asctime)s[%(filename)s:%(lineno)s][%(name)s]: %(message)s',
        when: str = 'M',
        interval: int = 1,
        backup_count: int = 12,
) -> Logger:
    if Path(directory).is_absolute():
        filepath = Path(directory)
    else:
        filepath = Path(Constant.BASE) / Path(directory)
    filename = str(filepath / Path(name + '.log'))
    handlers: List[Handler] = [TimedRotatingFileHandler(
        filename,
        when=when,
        interval=interval,
        backupCount=backup_count
    )]
    fmt: Formatter = Formatter(formatter)
    logger: Logger = logging.getLogger(name)
    if debug:
        handlers.append(StreamHandler())
    for handler in handlers:
        handler.setLevel(logging.DEBUG if debug else logging.INFO)
        handler.setFormatter(fmt)
        logger.addHandler(handler)
    return logger
