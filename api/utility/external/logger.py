import logging
from logging import Logger, StreamHandler, Handler, Formatter
from logging.handlers import TimedRotatingFileHandler
from pathlib import Path
import socket
from typing import List

from api.utility.constant import Constant


def init_logger(
        name: str,
        debug: bool = False,
        formatter: str = '[%(levelname)s]%(asctime)s[%(filename)s:%(lineno)s][%(name)s]: %(message)s',
) -> Logger:
    handlers: List[Handler] = []
    fmt: Formatter = Formatter(formatter)
    # init TCP SysLogHandler
    tcp_sys_handler = logging.handlers.SysLogHandler(
        address=('localhost', 514),
        facility=logging.handlers.SysLogHandler.LOG_UUCP,
        # socktype=socket.SOCK_DGRAM,
        socktype=socket.SOCK_STREAM,
    )
    handlers.append(tcp_sys_handler)
    # logger to stdout if debug
    if debug:
        handlers.append(StreamHandler())

    logger: Logger = logging.getLogger(name)
    for handler in handlers:
        handler.setLevel(logging.DEBUG if debug else logging.INFO)
        handler.setFormatter(fmt)
        logger.addHandler(handler)
    return logger


# def init_logger(
#         name: str,
#         directory: str,
#         debug: bool = False,
#         formatter: str = '[%(levelname)s]%(asctime)s[%(filename)s:%(lineno)s][%(name)s]: %(message)s',
#         when: str = 'M',
#         interval: int = 1,
#         backup_count: int = 12,
# ) -> Logger:
#     if Path(directory).is_absolute():
#         filepath = Path(directory)
#     else:
#         filepath = Path(Constant.BASE).joinpath(directory)
#     if not filepath.exists():
#         filepath.mkdir(0o755, True, True)
#     filename = str(filepath.joinpath(name + '.log'))
#     handlers: List[Handler] = [TimedRotatingFileHandler(
#         filename,
#         when=when,
#         interval=interval,
#         backupCount=backup_count
#     )]
#     fmt: Formatter = Formatter(formatter)
#     logger: Logger = logging.getLogger(name)
#     if debug:
#         handlers.append(StreamHandler())
#     for handler in handlers:
#         handler.setLevel(logging.DEBUG if debug else logging.INFO)
#         handler.setFormatter(fmt)
#         logger.addHandler(handler)
#     return logger


def __test__():
    # 测试rsyslog
    handler_udp = logging.handlers.SysLogHandler(
        address=('localhost', 514),
        facility=logging.handlers.SysLogHandler.LOG_UUCP,
        socktype=socket.SOCK_DGRAM,
        # socktype = socket.SOCK_STREAM,
    )
    handler_tcp = logging.handlers.SysLogHandler(
        address=('localhost', 514),
        facility=logging.handlers.SysLogHandler.LOG_UUCP,
        # socktype=socket.SOCK_DGRAM,
        socktype=socket.SOCK_STREAM,
    )
    logger = logging.getLogger()
    logger.level = logging.DEBUG
    # logger.handlers.clear()
    logger.handlers.append(handler_udp)
    logger.handlers.append(handler_tcp)
    import time

    logger.info(f'test info {time.time()}')
    logger.debug(f'test debug {time.time()}')
    logger.critical(f'test critical {time.time()}')


if __name__ == '__main__':
    __test__()
