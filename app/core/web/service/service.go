package service

type IUserService interface {
	GetUserDetail()
	ModifyUser()
	ChangePassword()
	NewGoogleToken()
}

type IUserAdminService interface {
	IUserService
	GetUserList()
	CreateUser()
	RemoveUser()
	DisableUser()
	EnableUser()
}

type IRoleAdminService interface {
	GetRoleList()
	GetRoleDetail()
	CreateRole()
	ModifyRole()
	RemoveRole()
}

type IUserRoleAdminService interface {
	GetRolesByUser()
	GetUsersByRole()
	ModifyRolesByUser()
	ModifyUsersByRole()
}

type IDepartmentService interface {
	GetDepartmentList()
	GetDepartmentDetail()
	CreateDepartment()
	ModifyDepartment()
	RemoveDepartment()
}
