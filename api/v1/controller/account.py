from typing import List

from sanic import Blueprint, Request, HTTPResponse, response
from sqlalchemy import select

from api.utility.constant import Constant
from api.v1.controller.__auth__ import need_login
from api.v1.model.platform import Account

account_bp: Blueprint = Blueprint('account', 'account')


class AccountVO:
    def __init__(self, query: List):
        self.id = query[0]
        self.code = query[1]
        self.name = query[2]
        self.login_name = query[3]
        self.email = query[4]
        self.telephone = query[5]
        self.mobile = query[6]
        self.avatar = query[7]
        self.department_id = query[8]


@account_bp.get('current', strict_slashes=False)  # /api/v1/index.php/account/current
@need_login
async def current_account(request: Request) -> HTTPResponse:
    """ 当前用户与可用菜单？ """
    web_session = request.ctx.web_session
    uid: int = web_session[Constant.SESSION_NAME_CURRENT_ACCOUNT]

    db_session = request.ctx.db_session
    async with db_session.begin():
        stmt = select(
            Account.id,
            Account.code,
            Account.name,
            Account.login_name,
            Account.email,
            Account.telephone,
            Account.mobile,
            Account.avatar,
            Account.department_id,
        ).where(Account.db_id == uid)
        cur = await db_session.execute(stmt)
        result = cur.first()
    if result:
        a = AccountVO(result)
        return response.json(dict(
            code=0,
            message='ok',
            data=dict(
                account=a.__dict__,
                menu=dict(),
            ),
        ))
    else:
        return response.json(dict(
            code=500,
            message='error of querying current user',
        ))

