import base64
import hashlib
import hmac
import io
import secrets
import struct
import time

import qrcode
from typing import Tuple


def make_qrcode(b32key: str, service_name: str, account: str, issuer: str) -> bytes:
    """
    otpauth://totp/<供应商>:<账号>?algorithm=SHA1&digits=6&period=30&issuer=<供应商>&secret=<密钥in lower BASE32>
    """
    url = f'''
    otpauth://totp/{service_name}:{account}?algorithm=SHA1&digits=6&period=30&issuer={issuer}&secret={b32key}
    '''.strip()
    qr = qrcode.QRCode(
        version=None,
        error_correction=qrcode.constants.ERROR_CORRECT_H,
        box_size=5,
        border=4,
    )
    qr.add_data(url)
    qr.make(fit=True)
    img = qr.make_image()
    stream = io.BytesIO()
    img.save(stream, format='PNG')
    return stream.getvalue()


def generate_token_and_qrcode(service_name: str, account: str, issuer: str) -> Tuple[str, str]:
    t = secrets.token_bytes(10)
    token = base64.b32encode(t).decode().lower()
    token_show = ' '.join([token[i: i + 4] for i in range(0, len(token), 4)])
    qr_img = 'data:image/png;base64,' + base64.b64encode(make_qrcode(token, service_name, account, issuer)).decode()
    return token_show, qr_img


def _h_otp_token(secret, intervals_no) -> str:
    key = base64.b32decode(secret, True)
    msg = struct.pack(">Q", intervals_no)
    h = hmac.new(key, msg, hashlib.sha1).digest()
    o = h[19] & 15
    h = (struct.unpack(">I", h[o:o+4])[0] & 0x7fffffff) % 1000000
    return "{}".format(h).rjust(6, '0')

def get_totp_token(secret):
    return _h_otp_token(secret, intervals_no=int(time.time())//30)
