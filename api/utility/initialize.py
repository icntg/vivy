"""
初始化工具。单独一个文件，用于交互生成配置。
"""
from pathlib import Path

from api.utility.constant import Constant


class AppContext:
    def __init__(self):
        from jinja2 import Environment, PackageLoader
        self.jinja2 = Environment(
            loader=PackageLoader('initialize', str(Path(Constant.BASE).joinpath('resource', 'template', 'initialize'))))


def create_and_run():
    from sanic import Sanic

    app = Sanic("initialize", ctx=AppContext())
    app.static('/', str(Path(Constant.BASE).joinpath('resource', 'initialize')), stream_large_files=True)

    # @app.get('/')
    # async def main_page(_: Request):
    #     return response.html()

    app.run(
        '0.0.0.0',
        8081,
        access_log=True,
        debug=True,
        workers=1,
    )
