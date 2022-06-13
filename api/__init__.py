from sanic import Sanic, response
from sanic.handlers import ErrorHandler
from sanic_ext import Extend

from api.utility.context import Context, get_context


def enum_blueprints(path: str = './v1'):
    pass


class CustomErrorHandler(ErrorHandler):
    def default(self, request, exception):
        ''' handles errors that have no error handlers assigned '''
        # You custom error handling logic...
        # return super().default(request, exception)
        return response.html('''
<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>错误</title></head>
<body><h1>错误</h1><p>发生错误！如重试还有问题，请联系管理员！</p></body></html>
        '''.strip())


# app.error_handler = CustomErrorHandler()


def create_app() -> Sanic:
    ctx: Context = get_context()
    app = Sanic(name='VIVY', ctx=ctx)

    ctx.DataSource.init_middleware()

    # 很遗憾，sanic的静态文件好像不支持默认文件名，比如'/' => '/index.html'
    # app.static('/', app.ctx.config.STATIC, stream_large_files=True)
    # app.static('/', str(Path(app.ctx.config.STATIC).joinpath('index.html')), stream_large_files=True)

    # from api.controller.v1 import checkin

    app.static('/', ctx.config.STATIC)
    import api.ssr.index
    import api.ssr.login
    app.blueprint(api.ssr.login.bp)
    Extend(app)

    if not ctx.config.SETTING.DEBUG:
        app.error_handler = CustomErrorHandler()

    return app
