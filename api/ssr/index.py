"""
首页
"""
from sanic import Sanic, Request, response

from api.utility.context import Context, get_context
from .__ssr__ import jinja2_env

ctx: Context = get_context()
app: Sanic = Sanic.get_app()


@app.route('/')
async def index(request: Request):
    """
    1. 检查登录凭证。如果未登录则跳转到登录页面。
    2. 如果已登录则渲染主页
    """
    if ctx.config.COOKIE not in request.cookies:
        return response.redirect('/account/login.html')
    # if verify jwt cookie:
    #     write secure log
    #     return response.redirect('/account/login.html')
    else:
        return dict()
