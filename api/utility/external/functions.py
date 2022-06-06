import sys
from typing import AnyStr


def std_print(s: AnyStr) -> None:
    sys.stdout.write(s)
    sys.stdout.flush()


def err_print(s: AnyStr) -> None:
    sys.stderr.write(s)
    sys.stderr.flush()
