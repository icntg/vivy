import string
from typing import Optional, List


class BaseX:
    def __init__(self, trans_table: bytes):
        self._enc_table: bytes = trans_table
        self._dec_table: List[Optional[int]] = self.make_dec_table(self._enc_table)

    @staticmethod
    def make_dec_table(table: bytes) -> List[Optional[int]]:
        buf: List[Optional[int]] = [None for _ in range(256)]
        for i, x in enumerate(table):
            buf[x] = i
        return buf

    def encode(self, plain: bytes) -> bytes:
        n = int.from_bytes(plain, 'big', signed=False)
        buf = bytearray()
        while 1:
            if n < len(self._enc_table):
                buf.append(self._enc_table[n])
                break
            m = n % len(self._enc_table)
            n = n // len(self._enc_table)
            buf.append(self._enc_table[m])
        return bytes(buf[::-1])

    def decode(self, enc: bytes) -> bytes:
        buf = []
        s = enc
        for i, b in enumerate(s):
            if self._dec_table[b] is None:
                raise ValueError(f'Non-base{len(self._enc_table)} digit found at {i}: {enc[i]}')
            buf.append(self._dec_table[b])
        z = 0
        for b in buf:
            z = z * len(self._enc_table) + b
        n, rem = divmod(z.bit_length(), 8)
        if rem:
            n += 1
        return z.to_bytes(n, 'big', signed=False)


base255 = BaseX(bytes([x for x in range(1, 256)]))
base36 = BaseX((string.digits + string.ascii_lowercase).encode())


def __test__():
    import string
    base100 = BaseX(string.printable.encode())
    a = base100.encode(b'HelloWorld!')
    print(len(a), repr(a))
    b = base100.decode(a)
    print(b)

    a = base255.encode(b'HelloWorld!')
    print(len(a), repr(a))
    b = base255.decode(a)
    print(b)
    # c = decode('1inndcn1p5hsvnsea7emrduj8t')
    # print(len(c), c)


if __name__ == '__main__':
    __test__()

