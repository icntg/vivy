"""
@author starvii
为了提升安全性，login模块使用服务端渲染，不使用前后端分离。


"""
import base64
import os
import pickle
import random
import struct
import time
from typing import Optional, Union
from urllib.parse import urlparse

from sanic import Sanic, Blueprint, Request, response, HTTPResponse
from sqlalchemy import select

from .__ssr__ import jinja2_env, render
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
    """
    return response.html(render('account/login0code.html'))


@bp.route('/login-with-portal.html')
async def login_with_portal_page(_: Request):
    return response.html(render('account/login0portal.html'))


@bp.route('/register.html')
async def register_page(_: Request):
    return response.html(render('account/register.html'))


class LoginVO:
    def __init__(self) -> None:
        self.form_user: str = ''
        self.form_pass: str = ''
        self.form_token: str = ''
        self.form_captcha: str = ''

        self._ctrl_byte: int = 0
        self.db_id: int = 0
        self.db_password_hash: str = ''
        self.db_token: Optional[str] = None

    def init_with_account(self, db_id: int, ph: str, token: Optional[str]):
        self.fake = False
        self.db_id = db_id
        self.db_password_hash = ph
        self.db_token = token
        self.has_token = self.db_token is not None and len(self.db_token) > 0

    def init_fake(self):
        self.fake = True
        self.db_id = struct.unpack('>I', os.urandom(4))[0]
        self.db_password_hash = '{}${}'.format(
            base64.b64encode(os.urandom(16)).decode().replace('=', ''),
            base64.b64encode(os.urandom(32)).decode().replace('=', ''),
        )
        self.db_token = base64.b32encode(os.urandom(10)).decode().lower()
        self.has_token = True

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


class LoginUtil:
    @staticmethod
    def random_code(code: str) -> int:
        c = code.encode()
        n = int.from_bytes(c, 'big', signed=False)
        random.seed(n)
        return random.randint(0, 65536)

    @staticmethod
    def render_captcha_page(req: Request) -> HTTPResponse:
        cap, url = get_captcha()
        req.ctx.web_session['captcha'] = cap
        return response.html(render('account/login1captcha.html', url=url))

    @staticmethod
    def render_token_page() -> HTTPResponse:
        return response.html(render('account/login1token.html'))

    @staticmethod
    def step_1_verify_captcha(request: Request) -> Optional[Exception]:
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
    async def step_0_query_account(db_session, code: str) -> LoginVO:
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
            db_id, ph, token = result
            a.init_with_account(db_id, ph, token)
        else:
            a.init_fake()
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
    if parsed.path in {_login_with_code_url, '/account/login.html'}:
        # 使用用户名密码登陆
        username = request.form.get('code')
        password = request.form.get('pass')
        lvo = await LoginUtil.step_0_query_account(request.ctx.db_session, username)
        lvo.form_user = username
        lvo.form_pass = password
        request.ctx.web_session['login'] = base64.b64encode(pickle.dumps(lvo)).decode()
        if lvo.fake:  # 未查到用户，根据用户名随机使用校验方式。
            if LoginUtil.random_code(username) % 2 == 0:
                return LoginUtil.render_captcha_page(request)
            else:
                return LoginUtil.render_token_page()
        elif lvo.has_token:
            return LoginUtil.render_token_page()
        else:
            return LoginUtil.render_captcha_page(request)
    elif parsed.path == '/account/login.php':
        pass

        #
        # ret = LoginUtil.step_0_verify_captcha(request)
        # if ret is not None:
        #     template = jinja2_env.get_template('message/uni-message.html')
        #     return response.html(template.render(
        #         title='错误',
        #         panel_title='错误',
        #         panel_message='验证码错误',
        #         back_url=_login_with_code_url,
        #     ))
        # # TODO: 获取验证信息
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
