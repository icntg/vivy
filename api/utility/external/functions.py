import sys
from typing import AnyStr


def out_put(s: AnyStr) -> None:
    sys.stdout.write(s)
    sys.stdout.flush()


def err_put(s: AnyStr) -> None:
    sys.stderr.write(s)
    sys.stderr.flush()
