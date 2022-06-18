import asyncio
import sqlite3
import time
from asyncio import Queue, AbstractEventLoop, new_event_loop, get_event_loop
from multiprocessing.dummy import Process, Lock
from typing import Optional, Tuple, Union, Callable, List

import ujson

from .base import BaseSessionInterface

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
    def __init__(
            self,
            filename: str = ':memory:',  # default to use memory database
            autocommit: bool = False,
            journal_mode: str = "WAL",
            init_sql: Optional[str] = None
    ):
        _g_lock.acquire()  # lock until EventLoop initialized
        super(SQLiteLoop, self).__init__()
        self.filename: str = filename
        self.autocommit: bool = autocommit
        self.journal_mode: str = journal_mode
        self.init_sql = init_sql
        self.conn: Optional[sqlite3.Connection] = None
        self.thread_loop: Optional[AbstractEventLoop] = None
        self.setDaemon(True)
        self.start()

    def __del__(self):
        if self.thread_loop is not None and self.thread_loop.is_running():
            self.thread_loop.stop()
        if self.conn is not None:
            self.conn.close()

    def run(self):
        if self.autocommit:
            self.conn = sqlite3.connect(self.filename, isolation_level=None, check_same_thread=False)
        else:
            self.conn = sqlite3.connect(self.filename, check_same_thread=False)
        self.conn.execute('PRAGMA journal_mode = %s;' % self.journal_mode)
        self.conn.text_factory = str
        cursor = self.conn.cursor()
        cursor.execute('PRAGMA synchronous=OFF;')
        if self.init_sql is not None and len(self.init_sql.strip()) > 0:
            sql_s = [x.strip() for x in self.init_sql.strip().split(';') if len(x.strip()) > 0]
            for sql in sql_s:
                cursor.execute(sql)
        self.thread_loop = new_event_loop()
        _g_lock.release()  # unlock before EventLoop running
        self.thread_loop.run_forever()

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
        self.thread_loop.call_soon_threadsafe(_sqlite, res)
        result: Union[CursorState, Exception] = await res.get()
        return result

    async def select_one(self, sql: str, args: Optional[Tuple] = None) -> Union[List, Exception]:
        main_loop = get_event_loop()

        def _sqlite(_res: Optional[Queue]):
            cursor = self.conn.cursor()
            cursor.execute(sql, args) if args else cursor.execute(sql)
            if _res:
                _result = cursor.fetchone()
                main_loop.call_soon_threadsafe(_res.put_nowait, _result)

        res: Queue = Queue(maxsize=1)
        self.thread_loop.call_soon_threadsafe(_sqlite, res)
        result = await res.get()
        return result

    def commit(self):
        self.thread_loop.call_soon_threadsafe(self.conn.commit)

    def close(self):
        self.thread_loop.call_soon_threadsafe(self.conn.close)
        self.thread_loop.stop()


class AIOSqliteSessionInterface(BaseSessionInterface):
    SQL_SELECT = 'SELECT `sid`, `start`, `expiry`, `store` FROM `session` WHERE `sid` = ?;'
    SQL_CREATE = '''
    CREATE TABLE IF NOT EXISTS `session`
    (`sid` VARCHAR(255) UNIQUE, `start` INTEGER,`expiry` INTEGER,`store` TEXT);'''
    SQL_INSERT = 'INSERT INTO `session` VALUES(?, ?, ?, ?);'
    SQL_UPDATE = 'UPDATE `session` SET `expiry`=?, `store`=? WHERE `sid` = ?'
    SQL_UPDATE_EXP = 'UPDATE `session` SET `expiry`=? WHERE `sid` = ?'
    SQL_DELETE = 'DELETE FROM `session` WHERE `sid` = ?;'
    SQL_CLEAN = 'DELETE FROM `session` WHERE `expiry` < ?;'

    def __init__(
            self,
            domain: str = None,
            expiry: int = 2592000,
            httponly: bool = True,
            cookie_name: str = "session",  # you could use 'PHPSESSID' to puzzle crackers
            prefix: str = "session:",
            sessioncookie: bool = False,
            samesite: str = None,
            session_name="session",  # you could use 'web_session' to distinct from 'db_session' of SQLAlchemy
            secure: bool = False,
            sid_provider: Callable[[None], str] = lambda: __import__('uuid').uuid4().hex,
            # sid_provider: Callable[[None], str] = lambda: ''.join([(__import__('string').
            #                                                         digits + __import__('string').ascii_lowercase)[
            #                                                            x % len(__import__('string').digits + __import__(
            #                                                                'string').ascii_lowercase)]
            #                                                        for x in __import__('os').urandom(26)]),
    ):
        super().__init__(
            expiry=expiry,
            prefix=prefix,
            cookie_name=cookie_name,
            domain=domain,
            httponly=httponly,
            sessioncookie=sessioncookie,
            samesite=samesite,
            session_name=session_name,
            secure=secure,
            sid_provider=sid_provider,
        )
        self.db = SQLiteLoop(init_sql=self.SQL_CREATE)
        async def _clean_timeout_session_loop():
            while 1:
                print('to clean')
                await self._clean_timeout()
                await asyncio.sleep(1, get_event_loop())
        get_event_loop().create_task(_clean_timeout_session_loop())

    async def _get_value(self, prefix, sid):
        session_id = self.prefix + sid

        rec = await self.db.select_one(self.SQL_SELECT, (session_id,))
        if not rec:
            return None
        store = ujson.loads(rec[3])
        new_expiry = int(time.time() + self.expiry)
        _ = await self.db.execute('BEGIN;')
        _ = await self.db.execute(self.SQL_UPDATE_EXP, (new_expiry, session_id))
        self.db.commit()
        return store

    async def _delete_key(self, key):
        _ = await self.db.execute('BEGIN;')
        _ = await self.db.execute(self.SQL_DELETE, (key,))
        self.db.commit()

    async def _set_value(self, key, data):
        store = ujson.dumps(data)
        _ = await self.db.execute('BEGIN;')
        rec = await self.db.select_one(self.SQL_SELECT, (key,))
        if not rec:
            start = int(time.time())
            expiry = int(start + self.expiry)
            _ = await self.db.execute(self.SQL_INSERT, (key, start, expiry, store))
        else:
            expiry = int(time.time() + self.expiry)
            _ = await self.db.execute(self.SQL_UPDATE, (expiry, store, key))
        self.db.commit()

    async def _clean_timeout(self):
        _ = await self.db.execute('BEGIN;')
        _ = await self._debug_()
        _ = await self.db.execute(self.SQL_CLEAN, (int(time.time()),))
        _ = await self._debug_()
        self.db.commit()

    async def _debug_(self):
        sql = 'SELECT * FROM `session`;'
        rec = await self.db.select_one(sql)
        print(rec)
