"""
@author starvii
为了提升安全性，login模块使用服务端渲染，不使用前后端分离。


"""
import base64
import os
import struct
from typing import Tuple, Optional

from sanic import Sanic, Blueprint, Request, response
from sqlalchemy import select

from .__ssr__ import jinja2_env
from .. import Context, get_context
from ..utility.external.functions import err_print
from ..utility.external.pynacl_util import password_verify
from ..v1.model.platform import Account

ctx: Context = get_context()
app: Sanic = Sanic.get_app()
bp: Blueprint = Blueprint('account', 'account')


@bp.route('/login-with-code.html')
@bp.route('/login.html')
async def login_with_code_page(_: Request):
    """
    使用账号、密码登录
    TODO: 获取图形验证码，并将校验信息加密后保存到cookie中。
    """
    template = jinja2_env.get_template('account/login0code.html')
    return response.html(template.render())


@bp.route('/login-with-portal.html')
async def login_with_portal_page(_: Request):
    template = jinja2_env.get_template('account/login0portal.html')
    return response.html(template.render())


@bp.route('/register.html')
async def register_page(_: Request):
    template = jinja2_env.get_template('account/register.html')
    return response.html(template.render())


class AccountVO:
    def __init__(self) -> None:
        self.fake: bool = False
        self.id: int = 0
        self.password: str = ''
        self.has_token: bool = True
        self.token: Optional[str] = None

    def serialize(self) -> bytes:
        if not self.fake:
            buf = bytearray()
            ctrl_byte = 1
            buf.extend(struct.pack('>I', self.id))
            try:
                s_salt, s_hash = self.password.split('$')
                b_salt = base64.b64decode(s_salt + '==')
                b_hash = base64.b64decode(s_hash + '=')
                buf.extend(b_salt)
                buf.extend(b_hash)
                if self.token is None or len(self.token) != 16:
                    buf.extend(os.urandom(10))
                    self.has_token = False
                    ctrl_byte &= 0b11111101
                else:
                    buf.extend(base64.b32decode(self.token.upper()))
                    self.has_token = True
                    ctrl_byte |= 0b00000010
                buf.insert(0, ctrl_byte)
                return buf
            except Exception as e:
                err_print(f'error in AccountVO.serialize: {e}\n')
        return b'\x02' + os.urandom(16 + 32 + 10)  # 0b00000010 第二位代表token存在，第一位代表fake

    def deserialize(self, stream: bytes) -> bool:
        result = True
        if len(stream) != 16 + 32 + 10 + 1:
            result = False
            s = b'\x02' + os.urandom(16 + 32 + 10)
        else:
            s = stream
        cb = s[0]
        if cb & 1 == 1:
            self.fake = False
        else:
            self.fake = True
        if (cb >> 1) & 1 == 1:
            self.has_token = True
        else:
            self.has_token = False
        self.id = struct.unpack('>I', s[1:][:4])[0]
        b_salt = s[1:][4:][:16]
        b_hash = s[1:][4:][16:][:32]
        self.password = (base64.b64encode(b_salt).decode() + '$' + base64.b64encode(b_hash).decode()).replace('=', '')
        b_token = s[1:][4:][16:][32:]
        self.token = base64.b32encode(b_token).decode().lower()
        return result


class LoginUtil:
    @staticmethod
    def verify_php_session_id(enc: str) -> int:
        return 0

    @staticmethod
    async def step_1_query_account(db_session, code: str) -> AccountVO:
        async with db_session.begin():
            stmt = select(
                Account.db_id,
                Account.password,
                Account.token,
            ).where(Account.code == code)
            result = await db_session.execute(stmt)
            account = result.first()
        a = AccountVO()
        if account is None:
            a.fake = True
        else:
            a.id = account[0]
            a.password = account[1]
            a.token = account[2]
            a.has_token = a.token is not None and len(a.token) > 0
        return a


@bp.route('/login.php', methods=['POST'])
async def login(request: Request):
    """
    TODO: 校验图形验证码；
    获取用户信息；
    如果不需要TOTP校验，则校验用户名密码；
    如果需要TOTP校验，将用户信息保存到cookie，跳转到TOTP页面。
    """
    cn = ctx.config.SESSION.COOKIE
    if cn in request.cookies and LoginUtil.verify_php_session_id(request.cookies[cn]) > 0:
        return response.redirect('/')  # 已登录

    # if ctx.config.COOKIE in request.cookies:
    #     pass
    # if LoginUtil.verify_php_session_id(request.cookies[ctx.config.COOKIE]) <= 0:
    #     pass
    code = request.form.get('code')
    pwd = request.form.get('pass')
    avo = await LoginUtil.step_1_query_account(request.ctx.session, code)
    if not avo.fake:
        if not avo.has_token and password_verify(avo.password, pwd):  # 直接比较密码
            # 登录成功，设置jwt，显示提示页面。
            pass
        elif avo.has_token:
            # cookie中加密保存密码和avo，显示token页面
            pass
    else:
        # cookie中加密保存密码和avo，显示token页面
        pass
    res = response.text('aaa')
    res.cookies[]
    return response.text(str(request.form))


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
