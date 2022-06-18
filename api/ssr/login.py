"""
@author starvii
为了提升安全性，login模块使用服务端渲染，不使用前后端分离。


"""
import base64
import os
import random
import secrets
import string
import struct
from typing import Optional, Tuple
from urllib.parse import urlparse

from sanic import Sanic, Blueprint, Request, response, HTTPResponse
from sqlalchemy import select

from .__ssr__ import render
from .. import Context, get_context
from ..utility.captcha import get_captcha
from ..utility.constant import Constant
from ..utility.external.google_token import get_totp_token
from ..utility.external.pynacl_util import password_verify
from ..v1.model.platform import Account

ctx: Context = get_context()
app: Sanic = Sanic.get_app()
bp: Blueprint = Blueprint('account', 'account')


@bp.route('/login-with-code.html')
@bp.route('/login.html')
async def login_with_code_page(_: Request):
    """ 使用账号、密码登录 """
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

        self.session_captcha: str = ''

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
        fid = bytearray(os.urandom(4))
        fid[0] |= 0x80  # 保证是一个负数
        self.db_id = struct.unpack('>i', fid)[0]
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
        """ 根据username产生固定的随机数 """
        c = code.encode()
        n = int.from_bytes(c, 'big', signed=False)
        random.seed(n)
        return random.randint(0, 65536)

    @staticmethod
    async def step_0_query_account(request: Request, code: str) -> LoginVO:
        """
        数据库ID与是否存在TOKEN
        """
        db_session = request.ctx.db_session
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
    async def step_1_render_enhance_pages(request: Request) -> HTTPResponse:
        # 使用用户名密码登陆
        username = request.form.get('code')
        password = request.form.get('pass')
        lvo = await LoginUtil.step_0_query_account(request, username)
        lvo.form_user = username
        lvo.form_pass = password
        if lvo.fake:  # 未查到用户，根据用户名随机使用校验方式。
            if not LoginUtil.is_username_legal(username):
                next_action = 0
            else:
                if LoginUtil.random_code(username) % 10 == 0:  # 合法用户名的话，有1/10的概率需要TOTP验证。
                    next_action = 1
                else:
                    next_action = 0
        elif lvo.has_token:  # 用户存在，而且需要TOTP验证
            next_action = 1
        else:  # 用户存在，无需TOTP验证。进行图形人机验证。
            next_action = 0
        if next_action == 0:
            cap, url = get_captcha()
            lvo.session_captcha = cap.lower()
            res = response.html(render('account/login1captcha.html', url=url))
        else:
            res = response.html(render('account/login1token.html'))
        session = request.ctx.web_session
        if Constant.SESSION_NAME_LOGIN in session:
            del session[Constant.SESSION_NAME_LOGIN]
        session[Constant.SESSION_NAME_LOGIN] = lvo.__dict__
        return res

    @staticmethod
    async def step_2_verify_factors(request: Request) -> Tuple[int, HTTPResponse]:
        session = request.ctx.web_session
        lvo = LoginVO()
        lvo.__dict__ = session[Constant.SESSION_NAME_LOGIN]
        del session[Constant.SESSION_NAME_LOGIN]

        capt: str = request.form.get('capt')
        totp: str = request.form.get('auth')
        if capt:
            ret = secrets.compare_digest(capt.lower(), lvo.session_captcha)
            if not ret:
                # 图形验证码错误
                return -1, response.html(render(
                    'message/uni-message.html',
                    title='captcha',
                    panel_title='错误',
                    panel_message='图形验证码错误！请重新登录！',
                    back_url='/',
                ))
        elif totp:
            ret = secrets.compare_digest(totp, get_totp_token(lvo.db_token))
            if not ret:
                # 谷歌时间令牌错误
                return -2, response.html(render(
                    'message/uni-message.html',
                    title='google time-based one-time password',
                    panel_title='错误',
                    panel_message='TOTP验证码错误！请重新登录！',
                    back_url='/',
                ))
        else:
            return -3, response.html(render(
                'message/uni-message.html',
                title='unknown',
                panel_title='错误',
                panel_message='未知错误！请重新登录！',
                back_url='/',
            ))
        # 校验密码
        ret = password_verify(lvo.db_password_hash, lvo.form_pass)
        f0 = 1 if ret else 0
        f1 = 0 if lvo.fake else 1
        ff = f0 & f1
        if ff == 0:
            return -4, response.html(render(
                'message/uni-message.html',
                title='credential',
                panel_title='错误',
                panel_message='用户名或密码错误！请重新登录！',
                back_url='/',
            ))
        return lvo.db_id, response.html(render(
            'message/uni-message.html',
            title='welcome',
            panel_title=f'欢迎你，{lvo.form_user}！',
            panel_message='登录成功，请继续。',
            back_url='/',
        ))

    @staticmethod
    def is_username_legal(login_name: str) -> bool:
        if len(login_name) < 4 or len(login_name) > 20:
            return False
        s = set(string.digits + string.ascii_lowercase + '_')
        ls = set(login_name.lower())
        if len(ls.difference(s)) > 0:
            return False
        return True


@bp.route('/login.php', methods=['POST'])
async def login(request: Request):
    """
    校验图形验证码；
    cookie 分成两部分，

    获取用户信息；
    如果不需要TOTP校验，则人机校验，再校验用户名密码；
    如果需要TOTP校验，将用户信息保存到session，跳转到TOTP页面。
    """
    referer = request.headers['referer']
    parsed = urlparse(referer)
    if parsed.path.lower().endswith('.html'):
        # 来自用户名、密码登录页面
        return await LoginUtil.step_1_render_enhance_pages(request)
    elif parsed.path.lower().endswith('.php'):
        # 来自第二步校验。校验验证码，通过后继续校验密码。
        uid, res = await LoginUtil.step_2_verify_factors(request)
        if uid > 0:
            # account: Account = await LoginUtil.step_3_query_account(request, uid)
            # print(account.__dict__)
            session = request.ctx.web_session
            session[Constant.SESSION_NAME_CURRENT_ACCOUNT] = uid
        return res


@bp.route('/logout.php', methods=['GET'])
def logout(request: Request):
    session = request.ctx.web_session
    del session[Constant.SESSION_NAME_CURRENT_ACCOUNT]
    if ctx.config.SESSION.COOKIE in request.cookies:
        return response.html(render('account/logout.html'))
    else:
        return response.redirect('/')
