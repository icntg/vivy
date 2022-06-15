"""
使用pynacl库进行加解密。
XSalsa20 stream cipher / Poly1305 MAC
"""

# from Crypto.Cipher import AES
# from Crypto.Util.Padding import pad

import nacl.secret
import nacl.pwhash


def encrypt(key: bytes, plain: bytes) -> bytes:
    box = nacl.secret.SecretBox(key)
    return box.encrypt(plain)


def decrypt(key: bytes, encrypted: bytes) -> bytes:
    box = nacl.secret.SecretBox(key)
    return box.decrypt(encrypted)


PREFIX = '$argon2id$v=19$m=65536,t=2,p=1$'


def password_hash(password: str) -> str:
    return nacl.pwhash.str(password.encode())[len(PREFIX):].decode()


def password_verify(pwd_hash: str, password: str) -> bool:
    h = (PREFIX + pwd_hash).encode()
    try:
        nacl.pwhash.verify(h, password.encode())
        return True
    except Exception as e:
        _ = e
        return False
