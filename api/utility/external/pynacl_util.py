"""
使用pynacl库进行加解密。
XSalsa20 stream cipher / Poly1305 MAC
"""

# from Crypto.Cipher import AES
# from Crypto.Util.Padding import pad

import nacl.secret

def encrypt(key: bytes, plain: bytes) -> bytes:
    box = nacl.secret.SecretBox(key)
    return box.encrypt(plain)


def decrypt(key: bytes, encrypted: bytes) -> bytes:
    box = nacl.secret.SecretBox(key)
    return box.decrypt(encrypted)
