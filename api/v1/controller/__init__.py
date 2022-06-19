from sanic import Blueprint

from api.v1.controller.account import account_bp
from api.v1.controller.platform import platform_bp
from api.v1.controller.utility import platform_util

group = Blueprint.group(version='1', version_prefix='/api/v', url_prefix='index.php', strict_slashes=False)

group.extend((
    account_bp,
    platform_bp,
    platform_util,
))
