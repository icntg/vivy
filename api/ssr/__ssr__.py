from pathlib import Path

from jinja2 import Environment, PackageLoader

from api.utility.context import Context, get_context

ctx: Context = get_context()

jinja2_env = Environment(loader=PackageLoader('vivy', str(ctx.config.TEMPLATE)))
# jinja2_dyn_env = Environment(loader=PackageLoader('vivy_dyn', str(Path(ctx.config.BASE).joinpath('web', 'dist'))))


def render(template_name: str, *args, **kwargs) -> str:
    t = jinja2_env.get_template(template_name)
    return t.render(*args, **kwargs)
