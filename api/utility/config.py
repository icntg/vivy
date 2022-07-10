import js2py


def read(filename='config.js'):
    with open(filename, 'r') as f:
        config = js2py.eval_js(f.read())
        return config
