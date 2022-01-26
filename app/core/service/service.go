package service

type UserService interface {
	GetUser()
	ModifyUser()
	ChangePassword()
	NewGoogleToken()
}

type UserAdminService interface {
	UserService
	GetUserList()
	CreateUser()
	RemoveUser()
	DisableUser()
	EnableUser()
}

type RoleAdminService interface {
}

type UserRoleAdminService interface {
}

type DepartmentService interface {
}
