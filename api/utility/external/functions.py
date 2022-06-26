import sys
from typing import AnyStr

from api.utility.external.base32 import encode_for_id
from api.utility.external.identity import ObjectId


def std_print(s: AnyStr) -> None:
    sys.stdout.write(s)
    sys.stdout.flush()


def err_print(s: AnyStr) -> None:
    sys.stderr.write(s)
    sys.stderr.flush()


def object_id() -> str:
    return encode_for_id(ObjectId.generate()).decode()


def __test__():
    for _ in range(20):
        print(object_id())


if __name__ == '__main__':
    __test__()
