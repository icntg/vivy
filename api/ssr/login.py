"""
@author starvii
为了提升安全性，login模块使用服务端渲染，不使用前后端分离。


"""

from sanic import Sanic, Blueprint, Request, response

from .__ssr__ import jinja2_env
from .. import Context, get_context

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


class LoginUtil:
    # @staticmethod
    # async def find_account_by_code(code: str) -> Account:
    #     session = ctx.DataSource.async_session()
    #     with session.begin():
    #         results = await session.execute(select(Account).where(Account.code==code))
    #         data: Account = results.scalar()
    #         return data

    @staticmethod
    def verify_php_session_id(enc: str) -> int:

        return 0


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
    print(request.form.get('code'))
    print(request.form.get('pass1'))
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
