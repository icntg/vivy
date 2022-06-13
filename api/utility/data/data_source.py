from typing import Dict, Optional
from urllib.parse import quote


class DataSource:
    def __init__(self, cfg: Optional[Dict] = None):
        self.host = ''
        self.port = 0
        self.username = ''
        self.password = ''
        if cfg is not None and 'host' in cfg:
            self.host = cfg['host']
        if cfg is not None and 'port' in cfg:
            self.port = cfg['port']
        if cfg is not None and 'username' in cfg:
            self.username = cfg['username']
        if cfg is not None and 'password' in cfg:
            self.password = cfg['password']


class MySQL(DataSource):
    def __init__(self, cfg: Optional[Dict] = None):
        super().__init__(cfg)
        self.database = ''
        self.option = '?charset=utf8mb4'
        self.max_idle = 10
        self.max_open = 100
        self.show_sql = False
        self.show_exec_time = False
        if cfg is not None and 'database' in cfg:
            self.database = cfg['database']
        if cfg is not None and 'option' in cfg:
            self.option = cfg['option']
        if cfg is not None and 'maxIdle' in cfg:
            self.max_idle = cfg['maxIdle']
        if cfg is not None and 'maxOpen' in cfg:
            self.max_open = cfg['maxOpen']
        if cfg is not None and 'showSql' in cfg:
            self.show_sql = cfg['showSql']
        if cfg is not None and 'showExecTime' in cfg:
            self.show_exec_time = cfg['showExecTime']

    def __repr__(self):
        # SQLALCHEMY_DATABASE_URL = "mysql+aiomysql://root:123456@xx.xx.xx.xx:3306/xx?charset=utf8mb4"
        name = quote(self.username)
        pwd = quote(self.password)
        return f'''
        mysql+aiomysql://{name}:{pwd}@{self.host}:{self.port}/{self.database}{self.option}
        '''.strip()

    def __str__(self):
        return self.__repr__()


class Redis(DataSource):
    # TODO:
    pass


class MongoDB(DataSource):
    # TODO:
    pass
