import functools
import sqlite3
from asyncio import AbstractEventLoop, get_event_loop, Future
from concurrent.futures import Executor, ThreadPoolExecutor
from enum import Enum
from multiprocessing.dummy import Process
from queue import Queue
from typing import Optional, Tuple, Union, List, Any


class CursorState:
    def __init__(self, cursor: sqlite3.Cursor):
        self.last_row_id = cursor.lastrowid
        self.row_count = cursor.rowcount
        self.description = cursor.description

    def __repr__(self):
        return f'last_row_id = {self.last_row_id}, row_count = {self.row_count}, description = {self.description}'

    def __str__(self):
        return self.__repr__()


class _Command(Enum):
    Close = 0b00000000
    Commit = 0b00000001
    RollBack = 0b00000010
    NoMore = 0b00000100
    Execute = 0b00010000
    ExecuteMany = 0b00010001
    SelectOne = 0b00100000
    Select = 0b00100001


class _SQLiteQuery:
    def __init__(
            self,
            command: _Command,
            request: Optional[str] = None,
            arguments: Optional[Tuple] = None,
            response_queue: Optional[Queue] = None,
    ):
        self.cmd = command
        self.sql = request
        self.arg = arguments
        self.res = response_queue

    def __repr__(self):
        return f'cmd = {self.cmd}, sql = {self.sql}, arg = {self.arg}, res = {self.res}'

    def __str__(self):
        return self.__repr__()


class _SQLiteThread(Process):
    def __init__(self, filename: str, autocommit: bool = False, journal_mode: str = "WAL"):
        super(_SQLiteThread, self).__init__()
        self.filename: str = filename
        self.autocommit: bool = autocommit
        self.journal_mode: str = journal_mode
        # use request queue of unlimited size
        self.reqs: Queue[_SQLiteQuery] = Queue()
        # self.setDaemon(True)  # python2.5-compatible
        self.start()

    def run(self):
        if self.autocommit:
            conn: sqlite3.Connection = sqlite3.connect(self.filename, isolation_level=None, check_same_thread=False)
        else:
            conn: sqlite3.Connection = sqlite3.connect(self.filename, check_same_thread=False)
        conn.execute('PRAGMA journal_mode = %s' % self.journal_mode)
        conn.text_factory = str
        cursor = conn.cursor()
        cursor.execute('PRAGMA synchronous=OFF')
        while 1:
            query: _SQLiteQuery = self.reqs.get()
            if query.cmd == _Command.Close:
                break
            elif query.cmd == _Command.Commit:
                conn.commit()
            elif query.sql is not None:
                sql: str = query.sql.strip()
                ret: Optional[Union[CursorState, BaseException]] = None
                try:
                    cursor.execute(sql, query.arg) if query.arg is not None else cursor.execute(sql)
                    if query.cmd in {_Command.Execute, _Command.ExecuteMany}:
                        ret = CursorState(cursor)
                except Exception as e:
                    ret = e
                if query.res is not None:
                    if ret is not None:
                        query.res.put(ret)
                    elif query.cmd == _Command.SelectOne:
                        query.res.put(cursor.fetchone())
                    elif query.cmd == _Command.Select:
                        for rec in cursor:
                            query.res.put(rec)
                        query.res.put(_Command.NoMore)
        conn.close()

    def execute(self, req: str, arg: Optional[Tuple] = None) -> CursorState:
        query = _SQLiteQuery(_Command.Execute, req, arg, Queue())
        self.reqs.put(query)
        ret: Union[CursorState, BaseException] = query.res.get()
        if isinstance(ret, BaseException):
            self.reqs.put(_SQLiteQuery(_Command.RollBack))
            raise ret
        if self.autocommit:
            self.reqs.put(_SQLiteQuery(_Command.Commit))
        return ret

    def execute_many(self, req: str, items: List[Tuple]) -> List[CursorState]:
        results: List[CursorState] = []
        queue = Queue()
        for item in items:
            query = _SQLiteQuery(_Command.ExecuteMany, req, item, queue)
            self.reqs.put(query)
            ret: Union[CursorState, BaseException] = query.res.get()
            if isinstance(ret, BaseException):
                self.reqs.put(_SQLiteQuery(_Command.RollBack))
                raise ret
            results.append(ret)
        if self.autocommit:
            self.reqs.put(_SQLiteQuery(_Command.Commit))
        return results

    def select(self, req: str, arg: Optional[Tuple] = None):
        query: _SQLiteQuery = _SQLiteQuery(_Command.Select, req, arg, Queue())
        self.reqs.put(query)
        while 1:
            ret: Union[Any, _Command, BaseException] = query.res.get()
            if isinstance(ret, BaseException):
                raise ret
            if isinstance(ret, _Command) and ret == _Command.NoMore:
                break
            yield ret

    def select_one(self, req: str, arg: Optional[Tuple] = None):
        query: _SQLiteQuery = _SQLiteQuery(_Command.SelectOne, req, arg, Queue())
        self.reqs.put(query)
        ret: Union[Any, BaseException] = query.res.get()
        if isinstance(ret, BaseException):
            raise ret
        return ret

    def commit(self):
        self.reqs.put(_SQLiteQuery(_Command.Commit))

    def close(self):
        self.reqs.put(_SQLiteQuery(_Command.Close))


class AIOSQlite:
    def __init__(self, filename: str, autocommit: bool = False, journal_mode: str = "WAL"):
        self._thread = _SQLiteThread(filename, autocommit, journal_mode)

    async def execute(self, req: str, arg: Optional[Tuple] = None) -> CursorState:
        loop: AbstractEventLoop = get_event_loop()
        loop.call_soon(functools.partial, self._thread.execute, req, arg)

    async def execute_many(self, req: str, items: List[Tuple]) -> List[CursorState]:
        pass

    async def select(self, req: str, arg: Optional[Tuple] = None):
        pass

    async def select_one(self, req: str, arg: Optional[Tuple] = None):
        pass



def __test__():
    a = AioSQLite(':memory')
    import asyncio
    loop = asyncio.get_event_loop()
    loop.run_until_complete(a.execute('''CREATE TABLE `session`(
  `sid` VARCHAR(255) UNIQUE,
  `start` TIMESTAMP, 
  `expiry` TIMESTAMP, 
  `store` TEXT
);'''))


if __name__ == '__main__':
    __test__()
