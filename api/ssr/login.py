"""
@author starvii
为了提升安全性，login模块使用服务端渲染，不使用前后端分离。


"""
import base64
import os
import struct
import time
from typing import Optional, Union
from urllib.parse import urlparse

from sanic import Sanic, Blueprint, Request, response, HTTPResponse
from sqlalchemy import select

from .__ssr__ import jinja2_env
from .. import Context, get_context
from ..utility.captcha import get_captcha
from ..utility.constant import Constant
from ..utility.error import CaptchaError
from ..utility.external.functions import err_print
from ..utility.external.pynacl_util import password_verify, encrypt, decrypt
from ..v1.model.platform import Account

ctx: Context = get_context()
app: Sanic = Sanic.get_app()
bp: Blueprint = Blueprint('account', 'account')


@bp.route('/login-with-code.html')
@bp.route('/login.html')
async def login_with_code_page(_: Request):
    """
    使用账号、密码登录
    使用图形验证码，并将校验信息加密后保存到cookie中。
    """
    cap, url = get_captcha()
    now = struct.pack('>I', int(time.time()) - Constant.TIME0)
    enc_cap = base64.urlsafe_b64encode(encrypt(ctx.secret, now + cap.encode())).decode()
    template = jinja2_env.get_template('account/login0code.html')
    res = response.html(template.render(url=url))
    res.cookies['c'] = enc_cap
    res.cookies['c']['httponly'] = True
    res.cookies['c']['samesite'] = 'Strict'
    return res


@bp.route('/login-with-portal.html')
async def login_with_portal_page(_: Request):
    template = jinja2_env.get_template('account/login0portal.html')
    return response.html(template.render())


@bp.route('/register.html')
async def register_page(_: Request):
    template = jinja2_env.get_template('account/register.html')
    return response.html(template.render())


class LoginVO:
    def __init__(self) -> None:
        self.timestamp: int = int(time.time() - Constant.TIME0)  # 时间戳，为了范围更加合理，从2022年1月1日开始计时的秒数。
        self.form_password: str = ''
        self.form_token: str = ''

        self._ctrl_byte: int = 0
        self.id: int = 0
        self.password_hash: str = ''
        self.token: Optional[str] = None
        self.captcha: Optional[str] = None  # TODO: 图形校验码暂时未启用。

    @property
    def fake(self) -> bool:
        return self._ctrl_byte & 1 == 0

    @fake.setter
    def fake(self, f: bool) -> None:
        if f:
            self._ctrl_byte &= 0b11111110
        else:
            self._ctrl_byte |= 0b00000001

    @property
    def has_token(self) -> bool:
        return (self._ctrl_byte >> 1) & 1 == 1

    @has_token.setter
    def has_token(self, t: bool) -> None:
        if t:
            self._ctrl_byte |= 0b00000010
        else:
            self._ctrl_byte &= 0b11111101

    def serialize(self) -> bytes:
        if not self.fake:
            buf = bytearray()
            buf.extend(struct.pack('>I', self.id))
            try:
                s_salt, s_hash = self.password_hash.split('$')
                b_salt = base64.b64decode(s_salt + '==')
                b_hash = base64.b64decode(s_hash + '=')
                buf.extend(b_salt)
                buf.extend(b_hash)
                if self.token is not None and len(self.token) != 16:
                    self.has_token = True
                    buf.extend(base64.b32decode(self.token.upper()))
                else:
                    self.has_token = False
                    buf.extend(os.urandom(10))
                buf.insert(0, self._ctrl_byte)
                return buf
            except Exception as e:
                err_print(f'error in LoginVO.serialize: {e}\n')
        return b'\x02' + os.urandom(16 + 32 + 10)  # 0b00000010 第二位代表token存在，第一位代表fake

    def deserialize(self, stream: bytes) -> bool:
        result = True
        if len(stream) != 16 + 32 + 10 + 1:
            result = False
            s = b'\x02' + os.urandom(16 + 32 + 10)
        else:
            s = stream
        self._ctrl_byte = s[0]
        self.id = struct.unpack('>I', s[1:][:4])[0]
        b_salt = s[1:][4:][:16]
        b_hash = s[1:][4:][16:][:32]
        self.password_hash = '{}${}'.format(
            base64.b64encode(b_salt).decode().replace('=', ''),
            base64.b64encode(b_hash).decode().replace('=', ''),
        )
        b_token = s[1:][4:][16:][32:]
        self.token = base64.b32encode(b_token).decode().lower()
        return result


class LoginUtil:
    @staticmethod
    def verify_php_session_id(enc: str) -> int:
        return 0

    @staticmethod
    def step_0_verify_captcha(request: Request) -> Optional[Exception]:
        capt = request.form.get('capt').lower()
        c = request.cookies['c']
        c = base64.urlsafe_b64decode(c)
        c = decrypt(ctx.secret, c)
        t = struct.unpack('>I', c[:4])[0]
        if t - time.time() + Constant.TIME0 > ctx.config.SESSION.LOGIN_COOKIE_TIMEOUT:
            return TimeoutError('captcha timeout')
        c = c[4:].decode().lower()
        if capt != c:
            return CaptchaError('captcha verify failed')

    @staticmethod
    async def step_1_query_account(db_session, code: str) -> LoginVO:
        """
        数据库ID与是否存在TOKEN
        """
        async with db_session.begin():
            stmt = select(
                Account.db_id,
                Account.password,
                Account.token,
            ).where(Account.code == code)
            cur = await db_session.execute(stmt)
            result = cur.first()
        a = LoginVO()
        if result:
            a.id, a.password, a.token = result
            a.ctrl_token = a.token is not None and len(a.token) > 0
        else:
            a.ctrl_fake = True
        return a

    @staticmethod
    async def step_3_verify_password(lvo: LoginVO) -> bool:
        pass

    @staticmethod
    async def step_2_verify_token(lvo: LoginVO) -> bool:
        pass


@bp.route('/login.php', methods=['POST'])
async def login(request: Request):
    """
    校验图形验证码；
    cookie 分成两部分，

    获取用户信息；
    如果不需要TOTP校验，则校验用户名密码；
    如果需要TOTP校验，将用户信息保存到cookie，跳转到TOTP页面。
    """

    referer = request.headers['referer']
    parsed = urlparse(referer)
    _login_with_code_url = '/account/login-with-code.html'
    if parsed.path == _login_with_code_url:
        ret = LoginUtil.step_0_verify_captcha(request)
        if ret is not None:
            template = jinja2_env.get_template('message/uni-message.html')
            return response.html(template.render(
                title='错误',
                panel_title='错误',
                panel_message='验证码错误',
                back_url=_login_with_code_url,
            ))
        # TODO: 获取验证信息
    return response.text(parsed.path)
    #
    # return response.json(request.headers)
    #
    # cn = ctx.config.SESSION.COOKIE
    # if cn in request.cookies and LoginUtil.verify_php_session_id(request.cookies[cn]) > 0:
    #     return response.redirect('/')  # 已登录
    #
    # # if ctx.config.COOKIE in request.cookies:
    # #     pass
    # # if LoginUtil.verify_php_session_id(request.cookies[ctx.config.COOKIE]) <= 0:
    # #     pass
    # code = request.form.get('code')
    # pwd = request.form.get('pass')
    # avo = await LoginUtil.step_1_query_account(request.ctx.session, code)
    # if not avo.fake:
    #     if not avo.ctrl_token and password_verify(avo.password_hash, pwd):  # 直接比较密码
    #         # 登录成功，设置jwt，显示提示页面。
    #         pass
    #     elif avo.has_token:
    #         # cookie中加密保存密码和avo，显示token页面
    #         pass
    # else:
    #     # cookie中加密保存密码和avo，显示token页面
    #     pass
    # res = response.text('aaa')
    # # res.cookies[]
    # return response.text(str(request.form))


def login_with_auth_code(_: Request):
    pass


@bp.route('/logout.php', methods=['GET'])
def logout(request: Request):
    if ctx.config.SESSION.COOKIE in request.cookies:
        template = jinja2_env.get_template('account/logout.html')
        res = response.html(template.render())
        del res.cookies[ctx.config.SESSION.COOKIE]
    else:
        return response.redirect('/')
