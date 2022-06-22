from sanic import Blueprint, Request, HTTPResponse

from api.v1.platform.controller.__auth__ import need_login

platform_util: Blueprint = Blueprint('utility', 'utility')


@platform_util.get('totp', strict_slashes=False)
@need_login
async def new_totp(_: Request) -> HTTPResponse:
    raise NotImplementedError
