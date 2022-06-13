import hashlib
import hmac
import os
from typing import Optional, Union, List


class Crypto:
    IV_SIZE = 8  # HASH_VALUE[16] | IV[8] | ENCRYPTED
    HASH_SIZE = 16

    @staticmethod
    def _rc4(key: Union[bytes, bytearray], data: Union[bytes, bytearray]) -> Union[bytes, bytearray]:
        def get_int(ch):
            if isinstance(ch, int):
                return ch
            elif isinstance(ch, str or bytes or bytearray):
                return ord(ch)
            raise ValueError(ch)

        s: List[int] = [i for i in range(256)]
        j: int = 0
        for i in range(256):
            x = key[i % len(key)]
            x: int = get_int(x)
            j: int = (j + s[i] + x) & 0xff
            s[i], s[j] = s[j], s[i]
        j: int = 0
        y: int = 0
        out: bytearray = bytearray()
        for char in data:
            j = (j + 1) & 0xff
            y = (y + s[j]) & 0xff
            s[j], s[y] = s[y], s[j]
            x: int = get_int(char)
            out.append(x ^ s[(s[j] + s[y]) & 0xff])
        return bytes(out)

    @staticmethod
    def encrypt(secret: Union[bytes, bytearray], message: Union[bytes, bytearray],
                iv: Optional[Union[bytes, bytearray]] = None) -> Union[bytes, bytearray]:
        if iv is not None and len(iv) != Crypto.IV_SIZE:
            raise ValueError(f"size of iv must equal {Crypto.IV_SIZE}")
        if iv is None:
            iv = os.urandom(Crypto.IV_SIZE)
        enc_key = Crypto._calc_enc_key(secret, iv)
        hash_key = Crypto._calc_hash_key(secret, iv)
        encrypted: bytes = Crypto._rc4(enc_key, message)
        iv_enc = iv + encrypted
        hash_value: bytes = hmac.new(hash_key, iv_enc, hashlib.sha256).digest()[:Crypto.HASH_SIZE]
        return hash_value + iv_enc

    @staticmethod
    def decrypt(secret: Union[bytes, bytearray], encrypted: Union[bytes, bytearray]) -> Union[bytes, bytearray]:
        if len(encrypted) <= Crypto.HASH_SIZE + Crypto.IV_SIZE:
            raise ValueError("size of encrypted is too short")
        expect_hash: bytes = encrypted[:Crypto.HASH_SIZE]
        iv: bytes = encrypted[Crypto.HASH_SIZE:][:Crypto.IV_SIZE]
        iv_enc: bytes = encrypted[Crypto.HASH_SIZE:]
        enc: bytes = iv_enc[Crypto.IV_SIZE:]
        hash_key = Crypto._calc_hash_key(secret, iv)
        hash_value: bytes = hmac.new(hash_key, iv_enc, hashlib.sha256).digest()[:Crypto.HASH_SIZE]
        if expect_hash != hash_value:
            raise ValueError("hash verify failed")
        enc_key = Crypto._calc_enc_key(secret, iv)
        message: bytes = Crypto._rc4(enc_key, enc)
        return message

    @staticmethod
    def _calc_enc_key(secret: bytes, iv: bytes) -> bytes:
        return hmac.new(iv, secret, hashlib.sha256).digest()

    @staticmethod
    def _calc_hash_key(secret: bytes, iv: bytes) -> bytes:
        return hmac.new(secret, iv, hashlib.sha256).digest()
