from sanic import Sanic
from sanic_ext import Extend


def enum_blueprints(path: str = './v1'):
    pass


def create_app() -> Sanic:
    app = Sanic(name='VIVY')

    # 很遗憾，sanic的静态文件好像不支持默认文件名，比如'/' => '/index.html'
    # app.static('/', app.ctx.config.STATIC, stream_large_files=True)
    # app.static('/', str(Path(app.ctx.config.STATIC).joinpath('index.html')), stream_large_files=True)

    # from api.controller.v1 import checkin
    # app.blueprint(checkin.bp)

    Extend(app)
    return app