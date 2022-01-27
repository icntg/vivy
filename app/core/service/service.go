package service

type UserService interface {
	GetUserDetail()
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
	GetRoleList()
	GetRoleDetail()
	CreateRole()
	ModifyRole()
	RemoveRole()
}

type UserRoleAdminService interface {
	GetRolesByUser()
	GetUsersByRole()
	ModifyRolesByUser()
	ModifyUsersByRole()
}

type DepartmentService interface {
	GetDepartmentList()
	GetDepartmentDetail()
	CreateDepartment()
	ModifyDepartment()
	RemoveDepartment()
}
