#!.venv/bin/python
# -*- coding: utf-8 -*-

"""
run:

docker run -d --rm --name mariadb \
    -p 3306:3306 \
    -v /home/src/vivy/data:/var/lib/mysql \
    -e MYSQL_ROOT_PASSWORD=root \
    mariadb:latest

source venv/bin/activate
python vivy.py -X pycache_prefix=.  # __pycache__ in one place
"""
import os

from api import create_app
from api.utility.config import Config
from api.utility.context import Context, get_context

ctx: Context = get_context()


def main():
    # 1. read config
    config: Config = Config()
    config.read()
    ctx.init_with_config(config)
    # 2. initial logger
    ctx.init_loggers()
    # 3. connect database
    ctx.init_data_source()
    # 4. start web service
    app = create_app()
    app.run(
        host=ctx.config.HTTP_HOST,
        port=ctx.config.HTTP_PORT,
        access_log=False,
        debug=ctx.config.DEBUG,
        workers=os.cpu_count(),
    )


if __name__ == '__main__':
    main()
