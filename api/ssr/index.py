"""
首页
"""
from sanic import Sanic, Request, response

from api.utility.constant import Constant
from api.utility.context import Context, get_context

ctx: Context = get_context()
app: Sanic = Sanic.get_app()


@app.route('/')
@app.route('/index.html')
async def index(request: Request):
    """
    1. 检查登录凭证。如果未登录则跳转到登录页面。
    2. 如果已登录则渲染主页
    """
    if Constant.SESSION_NAME_CURRENT_ACCOUNT not in request.ctx.web_session:
        return response.redirect('/account/login.html')
    # if verify jwt cookie:
    #     write secure log
    #     return response.redirect('/account/login.html')
    else:
        return response.text('登录成功')
