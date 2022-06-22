from functools import wraps

from sanic import Request, HTTPResponse, response

from api.utility.constant import Constant


def check_login_status(request: Request) -> bool:
    web_session = request.ctx.web_session
    if Constant.SESSION_NAME_CURRENT_ACCOUNT not in web_session:
        return False
    uid = web_session[Constant.SESSION_NAME_CURRENT_ACCOUNT]
    if isinstance(uid, int) and uid > 0:
        return True
    return False


def need_login(wrapped):
    def decorator(f):
        @wraps(f)
        async def decorated_function(request: Request, *args, **kwargs):
            is_authenticated = check_login_status(request)

            if is_authenticated:
                res: HTTPResponse = await f(request, *args, **kwargs)
                return res
            else:
                if request.path.startswith('/api'):
                    return response.json(dict(
                        code=401,
                        message='not login',
                        url='/account/login.html',
                    ))
                else:
                    return response.redirect('/account/login.html')
        return decorated_function
    return decorator(wrapped)
