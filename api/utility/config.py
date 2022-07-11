import js2py


def _check(conf: js2py.base.JsObjectWrapper) -> bool:
    pass


def read(filename='config.js') -> js2py.base.JsObjectWrapper:
    c = js2py.run_file(filename)
    conf = c[0]
    # 检查文件内容？
    return conf

