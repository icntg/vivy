"""
@author starvii
为了提升安全性，login模块使用服务端渲染，不使用前后端分离。


"""


def login_page():
    """
    获取图形验证码，并将校验信息加密后保存到cookie中。
    """
    pass


def login_with_password():
    """
    校验图形验证码；
    获取用户信息；
    如果不需要TOTP校验，则校验用户名密码；
    如果需要TOTP校验，将用户信息保存到cookie，跳转到TOTP页面。
    """
    pass


def login_with_auth_code():
    pass


def logout():
    pass
