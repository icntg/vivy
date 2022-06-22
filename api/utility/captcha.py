import base64
import os
from typing import Tuple

from captcha.image import ImageCaptcha

TABLE = '23467acdefhkmnqrtwxyACEFHKMNQRTWXY'  # 去掉了容易混淆的字符
N = 5
# CAP = ImageCaptcha(fonts=[str(__import__('api.utility.constant.Constant').BASE.joinpath('resource', 'static', 'font',
#                                                                                         'wkzyzt_tty_v4.ttf'))])
CAP = ImageCaptcha()


def get_captcha(n: int = N) -> Tuple[str, str]:
    rands = os.urandom(n)
    chars = ''.join([TABLE[x % len(TABLE)] for x in rands])
    img_bytes = CAP.generate(chars).getvalue()
    img_url = 'data:image/png;base64,' + base64.b64encode(img_bytes).decode()
    return chars, img_url
