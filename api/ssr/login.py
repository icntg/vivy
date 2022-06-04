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


@bp.route('/login-with-code.html', methods=['GET'])
@bp.route('/login.html', methods=['GET'])
async def login_with_code_page(_: Request):
    """
    使用账号、密码登录
    TODO: 获取图形验证码，并将校验信息加密后保存到cookie中。
    """
    template = jinja2_env.get_template('account/login0code.html')
    return response.html(template.render())


@bp.route('/login-with-portal.html', methods=['GET'])
async def login_with_portal_page(_: Request):
    template = jinja2_env.get_template('account/login0portal.html')
    return response.html(template.render())


@bp.route('/register.html', methods=['GET'])
async def register_page(_: Request):
    template = jinja2_env.get_template('account/register.html')
    return response.html(template.render())


@bp.route('/login.php', methods=['POST'])
async def login(request: Request):
    """
    TODO: 校验图形验证码；
    获取用户信息；
    如果不需要TOTP校验，则校验用户名密码；
    如果需要TOTP校验，将用户信息保存到cookie，跳转到TOTP页面。
    """
    if ctx.config.COOKIE in request.cookies:
        pass


def login_with_auth_code(_: Request):
    pass


def logout(_: Request):
    pass
