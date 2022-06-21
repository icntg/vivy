# 参考 https://www.jianshu.com/p/86f012e72c24

from pathlib import Path

import casbin
import casbin_sqlalchemy_adapter

conf = Path(__file__).parent.joinpath('model.conf')
# adapter = casbin_sqlalchemy_adapter.Adapter('sqlite:///test.db')
adapter = casbin_sqlalchemy_adapter.Adapter('mysql+pymysql://vivy:vivy@127.0.0.1:3306/vivy')
e = casbin.Enforcer(str(conf), adapter)  # TODO


def __test__():
    e.add_policy('g', 'admin1', '/api/admin', 'GET')
    # e.save_policy()


if __name__ == '__main__':
    __test__()
