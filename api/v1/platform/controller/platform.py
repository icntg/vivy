from sanic import Blueprint, Request, HTTPResponse

from api.v1.platform.controller.__auth__ import need_login

platform_bp: Blueprint = Blueprint('platform', 'platform')


class AccountController:
    @staticmethod
    @platform_bp.get('account', strict_slashes=False)
    @need_login
    async def get_account_list(request: Request) -> HTTPResponse:
        raise NotImplementedError

    @staticmethod
    @platform_bp.post('account', strict_slashes=False)
    @need_login
    async def create_account(request: Request) -> HTTPResponse:
        raise NotImplementedError

    @staticmethod
    @platform_bp.get('account/<account_id>', strict_slashes=False)
    @need_login
    async def get_account_detail(request: Request, account_id: str) -> HTTPResponse:
        raise NotImplementedError

    @staticmethod
    @platform_bp.post('account/<account_id>', strict_slashes=False)
    @need_login
    async def modify_account(request: Request, account_id: str) -> HTTPResponse:
        raise NotImplementedError

    @staticmethod
    @platform_bp.delete('account/<account_id>', strict_slashes=False)
    @need_login
    async def delete_account(request: Request, account_id: str) -> HTTPResponse:
        raise NotImplementedError


class DepartmentController:
    pass


class RoleController:
    pass


class MenuController:
    pass
