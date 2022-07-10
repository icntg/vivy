import logging
from logging import Logger, StreamHandler, Handler, Formatter
from logging.handlers import TimedRotatingFileHandler, SysLogHandler
from pathlib import Path
import socket
from typing import List, Optional, Tuple, Union

from api.utility.constant import constant
from api.utility.external.functions import err_print, std_print


def init_logger(
        name: str,
        debug: bool = False,
        formatter: str = '[%(levelname)s]%(asctime)s[%(filename)s:%(lineno)s][%(name)s]: %(message)s',
        address: Tuple[str, int] = ('localhost', 514),
        sock_type: socket.SocketKind = socket.SOCK_STREAM,
        directory: str = constant.LOG_DIR,
        when: str = 'M',
        interval: int = 1,
        backup_count: int = 12,
) -> Logger:
    handlers: List[Handler] = []
    fmt: Formatter = Formatter(formatter)
    # init TCP SysLogHandler
    tcp_sys_handler: Optional[SysLogHandler] = None
    try:
        tcp_sys_handler = SysLogHandler(
            address=address,
            facility=logging.handlers.SysLogHandler.LOG_UUCP,
            # socktype=socket.SOCK_DGRAM,
            # socktype=socket.SOCK_STREAM,
            socktype=sock_type,
        )
        handlers.append(tcp_sys_handler)
    except Exception as e:
        err_print(f'cannot create tcp_sys_handler: {e}\n')
    # use file_handler
    if debug or tcp_sys_handler is None:
        if Path(directory).is_absolute():
            filepath = Path(directory)
        else:
            filepath = Path(constant.BASE).joinpath(directory)
        std_print(f'to use {filepath} as {name} logger instead.\n')
        try:
            if not filepath.exists():
                filepath.mkdir(0o755, True, True)
            filename = str(filepath.joinpath(name + '.log'))
            file_handler = TimedRotatingFileHandler(
                filename,
                when=when,
                interval=interval,
                backupCount=backup_count,
            )
            handlers.append(file_handler)
        except Exception as e:
            err_print(f'cannot create file_handler: {e}\n')
    # logger to stdout if debug
    if debug:
        handlers.append(StreamHandler())

    logger: Logger = logging.getLogger(name)
    for handler in handlers:
        handler.setLevel(logging.DEBUG if debug else logging.INFO)
        handler.setFormatter(fmt)
        logger.addHandler(handler)
    return logger


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
