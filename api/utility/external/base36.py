import string

TABLE = (string.digits + string.ascii_lowercase).encode()


def make_dec_table() -> bytes:
    buf = [255 for _ in range(256)]
    for i, x in enumerate(TABLE):
        buf[x] = i
    return bytes(buf)


DEC_TABLE = make_dec_table()


def encode(data: bytes) -> str:
    n = int.from_bytes(data, 'big', signed=False)
    buf = bytearray()
    while 1:
        if n < len(TABLE):
            buf.append(TABLE[n])
            break
        m = n % len(TABLE)
        n = n // len(TABLE)
        buf.append(TABLE[m])
    buf = buf[::-1]
    return buf.decode()


def decode(data: str) -> bytes:
    buf = []
    s = data.encode()
    for i, b in enumerate(s):
        if DEC_TABLE[b] >= 255:
            raise ValueError(f'Non-base36 digit found at {i}: {data[i]}')
        buf.append(DEC_TABLE[b])
    z = 0
    for b in buf:
        z = z * len(TABLE) + b
    n, rem = divmod(z.bit_length(), 8)
    if rem:
        n += 1
    return z.to_bytes(n, 'big', signed=False)


def __test__():
    a = encode(b'HelloWorld!')
    print(a)
    b = decode(a)
    print(b)
    c = decode('1inndcn1p5hsvnsea7emrduj8t')
    print(len(c), c)


if __name__ == '__main__':
    __test__()

