import hashlib
import os
import socket
import struct
import time
from multiprocessing import Lock
from typing import List
import binascii


__ALL__ = ["ObjectId"]


def _auto_pid() -> bytes:
        pid: int = os.getpid() & 0xffff
        if pid == 1:
            # 某些情况下，docker容器中pid会为1，这样就失去的这个字段的作用。
            # 尝试读取docker id替代。
            try:
                lines: List[bytes] = [
                    x.strip() for x in open('/proc/self/cgroup', 'rb').read().split(b'\n') 
                    if len(x.strip()) > 0 and x.strip().startswith(b'0:')
                ]
                _, docker_id = lines[0].split(b':/docker/')
                return binascii.unhexlify(docker_id[:4])
            except Exception as e:
                _ = e
                # 使用随机pid
                return os.urandom(2)
        return struct.pack('>H', pid)


class ObjectId:
    _index, = struct.unpack('>I', os.urandom(4))
    _mac = hashlib.md5(socket.gethostname().encode('utf-8')).digest()[:3]
    _lock = Lock()
    _pid: bytes = _auto_pid()
        
    @staticmethod
    def generate() -> bytes:
        # | 0 | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 |
        # |   TIMESTAMP   |    MAC    |  PID  |   COUNTER   |
        ObjectId._lock.acquire()
        ObjectId._index = (ObjectId._index + 1) & 0xffffff
        counter: bytes = struct.pack('>I', ObjectId._index)[1:]
        ObjectId._lock.release()
        timestamp: bytes = struct.pack('>I', int(time.time()))
        oid: bytes = timestamp + ObjectId._mac + ObjectId._pid + counter
        return oid
