import hashlib
import os
import socket
import struct
import time
from multiprocessing import Lock


__ALL__ = ["ObjectId"]


class ObjectId:
    _index, = struct.unpack('>I', os.urandom(4))
    _mac = hashlib.md5(socket.gethostname().encode('utf-8')).digest()[:3]
    _lock = Lock()

    @staticmethod
    def generate() -> bytes:
        # | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 |
        # |   TIMESTAMP   |    MAC    |  PID  |   COUNTER   |
        ObjectId._lock.acquire()
        ObjectId._index = (ObjectId._index + 1) & 0xffffff
        counter: bytes = struct.pack('>I', ObjectId._index)[1:]
        ObjectId._lock.release()
        timestamp: bytes = struct.pack('>I', int(time.time()))
        pid: bytes = struct.pack('>H', os.getpid())
        oid: bytes = timestamp + ObjectId._mac + pid + counter
        return oid

# TODO: 如果在docker下，pid使用docker id替代。
