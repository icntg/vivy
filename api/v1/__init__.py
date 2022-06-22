from sanic import Blueprint

from api.v1.platform.controller import group as platform_group
from api.v1.task.controller import group as task_group

group = Blueprint.group(version='1', version_prefix='/api/v', url_prefix='index.php', strict_slashes=False)

group.extend((
    platform_group,
    task_group,
))
