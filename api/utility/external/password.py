import base64
import hashlib
import hmac
import secrets
from typing import List

hash_algorithms = {
    1: hashlib.md5,
    5: hashlib.sha256,
    6: hashlib.sha512,
}

hash_trans = str.maketrans('ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=',
                           './0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz')
de_hash_trans = str.maketrans('./0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz',
                              'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=')


def _hash_base64(a: bytes) -> str:
    enc = base64.b64encode(a).decode()
    return enc.translate(hash_trans)


def password_hash(password: str, algo: int = 6) -> str:
    nonce: bytes = secrets.token_bytes(12)
    hash_func = hash_algorithms[algo]
    enc = hmac.new(nonce, password.encode(), hash_func).digest()
    return f'''${algo}${_hash_base64(nonce)}${_hash_base64(enc)}'''


def password_verify(password: str, hashed: str) -> bool:
    if hashed[0] != '$':
        raise ValueError('hashed_format')
    a: List[str] = hash[1:].split('$')
    if len(a) != 3:
        raise ValueError('hashed_format')
    if not a[0].isdigit():
        raise ValueError('hashed_format')
    if int(a[0]) not in hash_algorithms:
        raise ValueError('hashed_format')
    hash_func = hash_algorithms[int(a[0])]
    nonce = base64.b64decode(a[1].translate(de_hash_trans))
    hashed_enc = base64.b64decode(a[2].translate(de_hash_trans))
    enc = hmac.new(nonce, password.encode(), hash_func).digest()
    return secrets.compare_digest(enc, hashed_enc)
