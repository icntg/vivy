# 参考 https://www.jianshu.com/p/86f012e72c24

import casbin
from pathlib import Path

conf = Path(__file__).parent.joinpath('model.conf')

e = casbin.Enforcer(str(conf), ...)  # TODO
