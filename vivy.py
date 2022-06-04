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


def main():
    # 1. read config
    from api.utility.config import Config
    config: Config = Config()
    config.read()

    from api.utility.context import Context, get_context
    ctx: Context = get_context()

    ctx.init_with_config(config)
    # 2. initial logger
    ctx.init_loggers()
    # 3. connect database
    ctx.init_data_source()
    # 4. start web service
    from api import create_app
    app = create_app()
    app.run(
        host=ctx.config.HTTP_HOST,
        port=ctx.config.HTTP_PORT,
        access_log=False,
        debug=ctx.config.DEBUG,
        # workers=os.cpu_count(),
        workers=1,  # 目前由于采用文件log原因，只能一个进程。等采用rsyslog服务之后再改进。
    )


if __name__ == '__main__':
    main()
