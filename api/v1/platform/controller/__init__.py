from sanic import Blueprint

from api.v1.platform.controller.account import account_bp
from api.v1.platform.controller.platform import platform_bp
from api.v1.platform.controller.utility import platform_util

group = Blueprint.group(strict_slashes=False)

group.extend((
    account_bp,
    platform_bp,
    platform_util,
))