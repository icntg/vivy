import asyncio
import sqlite3
from asyncio import Queue, AbstractEventLoop, new_event_loop, get_event_loop
from multiprocessing.dummy import Process, Lock
from typing import Optional, Tuple, Any, Union

# import ujson

# from .base import BaseSessionInterface, default_sid_provider
_g_lock = Lock()


class CursorState:
    def __init__(self, cursor: sqlite3.Cursor):
        self.last_row_id = cursor.lastrowid
        self.row_count = cursor.rowcount
        self.description = cursor.description

    def __repr__(self):
        return f'last_row_id = {self.last_row_id}, row_count = {self.row_count}, description = {self.description}'

    def __str__(self):
        return self.__repr__()


class SQLiteLoop(Process):
    def __init__(self, filename: str = ':memory:', autocommit: bool = False, journal_mode: str = "WAL"):
        _g_lock.acquire()  # lock until EventLoop initialized
        super(SQLiteLoop, self).__init__()
        self.filename: str = filename
        self.autocommit: bool = autocommit
        self.journal_mode: str = journal_mode
        self.conn: Optional[sqlite3.Connection] = None
        self.loop: Optional[AbstractEventLoop] = None
        self.setDaemon(True)
        self.start()

    def run(self):
        if self.autocommit:
            self.conn = sqlite3.connect(self.filename, isolation_level=None, check_same_thread=False)
        else:
            self.conn = sqlite3.connect(self.filename, check_same_thread=False)
        self.conn.execute('PRAGMA journal_mode = %s' % self.journal_mode)
        self.conn.text_factory = str
        cursor = self.conn.cursor()
        cursor.execute('PRAGMA synchronous=OFF')
        self.loop = new_event_loop()
        _g_lock.release()  # unlock before EventLoop running
        self.loop.run_forever()

    async def execute(self, sql: str, args: Optional[Tuple] = None) -> Union[CursorState, Exception]:
        main_loop = get_event_loop()

        def _sqlite(_res: Optional[Queue]):
            try:
                cursor = self.conn.cursor()
                cursor.execute(sql, args) if args else cursor.execute(sql)
                _result = CursorState(cursor)
            except Exception as e:
                _result = e
            if _res:
                main_loop.call_soon_threadsafe(_res.put_nowait, _result)
        res: Queue = Queue(maxsize=8)
        self.loop.call_soon_threadsafe(_sqlite, res)
        result: Union[CursorState, Exception] = await res.get()
        return result

    async def select_one(self, sql: str, args: Optional[Tuple] = None) -> Union[Any, Exception]:
        main_loop = get_event_loop()

        def _sqlite(_res: Optional[Queue]):
            cursor = self.conn.cursor()
            cursor.execute(sql, args) if args else cursor.execute(sql)
            if _res:
                _result = cursor.fetchone()
                main_loop.call_soon_threadsafe(_res.put_nowait, _result)
        res: Queue = Queue(maxsize=1)
        self.loop.call_soon_threadsafe(_sqlite, res)
        result = await res.get()
        return result

    def commit(self):
        self.loop.call_soon_threadsafe(self.conn.commit)

    async def close(self):
        self.loop.call_soon_threadsafe(self.conn.close)
        self.loop.stop()

# class MemorySqliteSessionInterface(BaseSessionInterface):
#     SELECT_SQL = '''
# SELECT
#     `sid`,
#     `start`,
#     `expiry`,
#     `store`
# FROM `session`
# WHERE `sid` = ?;
#         '''
#
#     def __init__(
#             self,
#             domain: str = None,
#             expiry: int = 2592000,
#             httponly: bool = True,
#             cookie_name: str = "session",  # PHPSESSID
#             prefix: str = "session:",
#             sessioncookie: bool = False,
#             samesite: str = None,
#             session_name="session",
#             secure: bool = False,
#             sid_provider: Callable[[None], str] = default_sid_provider,
#     ):
#         super().__init__(
#             expiry=expiry,
#             prefix=prefix,
#             cookie_name=cookie_name,
#             domain=domain,
#             httponly=httponly,
#             sessioncookie=sessioncookie,
#             samesite=samesite,
#             session_name=session_name,
#             secure=secure,
#             sid_provider=sid_provider,
#         )
#         self.db_thread = SQLiteThread(':memory')
#         sql = '''
# CREATE TABLE `session`(
#   `sid` VARCHAR(255) UNIQUE,
#   `start` TIMESTAMP,
#   `expiry` TIMESTAMP,
#   `store` TEXT
# );
#         '''
#         self.db_thread.execute(sql)
#         self.executor: Executor = ThreadPoolExecutor(os.cpu_count())
#
#     async def _get_value(self, prefix, sid):
#         session_id = self.prefix + sid
#         rec = await asyncio.get_event_loop().run_in_executor(
#             self.executor,
#             self.db_thread.select_one,
#             self.SELECT_SQL,
#             (session_id,)
#         )
#         if rec and time.time() <= rec[2]:
#             return ujson.loads(rec[3])
#
#     async def _delete_key(self, key):
#         sql = '''DELETE FROM `session` WHERE `sid` = ?;'''
#         await asyncio.get_event_loop().run_in_executor(
#             self.executor,
#             self.db_thread.execute,
#             sql,
#         )
#         self.db_thread.commit()
#
#     async def _set_value(self, key, data):
#         store = ujson.dumps(data)
#         self.db_thread.execute('BEGIN;')
#         rec = await asyncio.get_event_loop().run_in_executor(
#                 self.executor,
#                 self.db_thread.select_one,
#                 self.SELECT_SQL, (key,)
#         )
#         if not rec:
#             start = time.time()
#             expiry = start + self.expiry
#             sql = '''INSERT INTO `session` VALUES(?, ?, ?, ?);'''
#             await asyncio.get_event_loop().run_in_executor(
#                 self.executor,
#                 self.db_thread.execute,
#                 sql, (key, start, expiry, store)
#             )
#         else:
#             expiry = time.time() + self.expiry
#             sql = '''UPDATE `session` SET `expiry`=?, `store`=? WHERE `sid` = ?'''
#             await asyncio.get_event_loop().run_in_executor(
#                 self.executor,
#                 self.db_thread.execute,
#                 sql, (key, expiry, store)
#             )
#         self.db_thread.commit()


async def __test__():
    t = SQLiteLoop()
    _g_lock.acquire()
    print('main线程加锁')
    _g_lock.release()
    print('main线程解锁')
    r0 = await t.execute('''
CREATE TABLE `session`(
  `sid` VARCHAR(255) UNIQUE,
  `start` TIMESTAMP,
  `expiry` TIMESTAMP,
  `store` TEXT
);
    ''')
    print(r0)


if __name__ == '__main__':
    asyncio.run(__test__())
