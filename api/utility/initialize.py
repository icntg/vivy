"""
初始化工具。单独一个文件，用于交互生成配置。
"""
from pathlib import Path

from sanic import Request, response

from api.utility.constant import Constant
from api.utility.external.google_token import generate_token_and_qrcode

STATIC = Path(Constant.BASE).joinpath('resource', 'static_initialize')


def initialize():
    pass


def create_and_run():
    from sanic import Sanic

    app = Sanic("initialize")
    app.static('/', str(STATIC), stream_large_files=True)

    @app.get('/')
    async def main_page(_: Request):
        return response.html(open(str(STATIC.joinpath('index.html')), 'rb').read())

    @app.post('/api/token')
    async def generate_token(req: Request):
        login_name = req.body.decode()
        token, img = generate_token_and_qrcode('日晷', login_name, '温州信通')
        return response.json(dict(token=token, qrcode=img))

    @app.post('/api/initialization')
    async def initialize(req: Request):
        print(req.json)
        return response.json(dict(code=0, message=''))

    app.run(
        '0.0.0.0',
        8081,
        access_log=True,
        debug=True,
        workers=1,
    )
