import aioredis
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


def create_app() -> Sanic:
    ctx: Context = get_context()
    app = Sanic(name='VIVY', ctx=ctx)

    from api.utility.session import Session, AIORedisSessionInterface, PHP_PROVIDER
    session = Session()

    @app.listener('before_server_start')
    async def aioredis_server_init(app, loop):
        app.ctx.redis = aioredis.from_url('redis://localhost', decode_responses=True)
        # init extensions fabrics
        session_instance = AIORedisSessionInterface(
            app.ctx.redis,
            expiry=ctx.config.SESSION.SESSION_TIMEOUT,
            cookie_name=ctx.config.SESSION.COOKIE,
            session_name='web_session',
        )
        session_instance.sid_provider = PHP_PROVIDER
        session_instance.__dict__['sid_provider'] = PHP_PROVIDER
        session.init_app(app, interface=session_instance)


    ctx.DataSource.init_middleware()

    # 很遗憾，sanic的静态文件好像不支持默认文件名，比如'/' => '/index.html'
    # 但是可以使用'/'渲染到'index.html'
    # app.static('/', app.ctx.config.STATIC, stream_large_files=True)
    # app.static('/', str(Path(app.ctx.config.STATIC).joinpath('index.html')), stream_large_files=True)

    # from api.controller.v1 import checkin

    app.static('/favicon.ico', ctx.config.STATIC.joinpath('favicon.ico'))
    app.static('/css/style.css', ctx.config.STATIC.joinpath('css', 'style.css'))
    app.static('/', ctx.config.BASE.joinpath('web', 'dist'))
    import api.ssr.index
    import api.ssr.login
    app.blueprint(api.ssr.login.bp)
    Extend(app)

    if not ctx.config.SETTING.DEBUG:
        app.error_handler = CustomErrorHandler()

    return app
