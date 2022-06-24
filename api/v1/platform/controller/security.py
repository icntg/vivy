from sanic import Blueprint, Request, HTTPResponse

from api.v1.platform.controller.__auth__ import need_login

platform_sec: Blueprint = Blueprint('security', 'security')


@platform_sec.get('totp', strict_slashes=False)
@need_login
async def new_totp(_: Request) -> HTTPResponse:
    raise NotImplementedError
