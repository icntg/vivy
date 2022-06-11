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

from api.utility.constant import Constant
from api.utility.external.functions import err_print, std_print


def check_init_state() -> bool:
    """
    检查初始化状态。
    True = 已初始化
    False = 未初始化
    """
    use_service_mode = True
    if not Constant.INIT.exists():
        err_print(f'[{Constant.INIT}] does not exist\n')
        use_service_mode = False
    if not Constant.CONF.exists():
        err_print(f'[{Constant.CONF}] does not exist\n')
        use_service_mode = False
    if not use_service_mode:
        std_print('to start static_initialize mode ...\n')
    return use_service_mode


def initialize_mode():
    from api.v1.initialize import create_and_run
    create_and_run()


def service_mode():
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
        host=ctx.config.HTTP.HOST,
        port=ctx.config.HTTP.PORT,
        access_log=ctx.config.SETTING.DEBUG,
        debug=ctx.config.SETTING.DEBUG,
        # workers=os.cpu_count(),
        workers=1,  # 目前由于采用文件log原因，只能一个进程。等采用rsyslog服务之后再改进。
    )


def main():
    if check_init_state():
        service_mode()
    else:
        initialize_mode()


if __name__ == '__main__':
    main()
