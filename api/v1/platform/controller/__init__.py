from sanic import Blueprint

from api.v1.platform.controller.account import account_bp
from api.v1.platform.controller.platform import platform_bp
from api.v1.platform.controller.security import platform_sec

group = Blueprint.group(strict_slashes=False)

group.extend((
    account_bp,
    platform_bp,
    platform_sec,
))