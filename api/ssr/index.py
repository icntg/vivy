"""
首页
"""
from sanic import Sanic, Request, response

from api.utility.constant import Constant
from api.utility.context import Context, get_context
from api.v1.platform.controller.__auth__ import need_login

ctx: Context = get_context()
app: Sanic = Sanic.get_app()


@app.route('/', strict_slashes=False)
@app.route('/index.html', strict_slashes=False)
@app.route('/index.php', strict_slashes=False)
@need_login
async def index(_: Request):
    with open(Constant.BASE.joinpath('web', 'dist', 'index.html'), 'rb') as f:
        return response.html(f.read())
