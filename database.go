package dadmin

import (
	d "github.com/yqBdm7y/devtool"
	"golang.org/x/crypto/bcrypt"
)

type Database struct{}

func (da Database) Migrate() {
	d.Database[d.LibraryGorm]{}.Get().DB.AutoMigrate(
		&App{},
		&User{},
		&Role{},
		&Menu{},
		&Department{},
	)
}

// 插入初始化数据
func (da Database) InsertInitData() (initialized bool, err error) {
	// 获取初始化用户数据
	init_user_data, err := get_init_user_data()
	if err != nil {
		panic(err)
	}

	initialized, err = d.Database[d.LibraryGorm]{}.Get().InsertInitializationData(
		init_department_data,
		init_role_data,
		init_menu_data,
		init_user_data,
	)
	if err != nil {
		panic(err)
	}

	return initialized, nil
}

var (
	init_department_data = []Department{
		{ID: 1, Status: 1, Name: "管理组"},
	}

	init_role_data = []Role{
		{ID: 1, Status: 1, Name: "超级管理员"},
	}

	init_menu_data = []Menu{
		{ID: 1, Path: "/system", Title: "系统", Name: "System", Sort: 100, Icon: "ep:setting"},
		{ID: 2, ParentId: 1, Path: "/system/user", Name: "UserList", Component: "system/user/index", Title: "用户管理", Sort: 1},
		{ID: 3, ParentId: 1, Path: "/system/role", Name: "RoleList", Component: "system/role/index", Title: "角色管理", Sort: 2},
		{ID: 4, ParentId: 1, Path: "/system/department", Name: "DepartmentList", Component: "system/dept/index", Title: "部门管理", Sort: 3},
		{ID: 5, ParentId: 1, Path: "/system/menu", Name: "MenuList", Component: "system/menu/index", Title: "菜单管理", Sort: 4},

		{ParentId: 2, MenuType: 4, Path: "/api/user/getList", Name: "UserGetList", Title: "获取列表", Sort: 1},
		{ParentId: 2, MenuType: 4, Path: "/api/user/create", Name: "UserCreate", Title: "创建用户", Sort: 2},
		{ParentId: 2, MenuType: 4, Path: "/api/user/edit", Name: "UserEdit", Title: "编辑用户", Sort: 3},
		{ParentId: 2, MenuType: 4, Path: "/api/user/editStatus", Name: "UserEditStatus", Title: "修改状态", Sort: 4},
		{ParentId: 2, MenuType: 4, Path: "/api/user/editPassword", Name: "UserEditPassword", Title: "修改密码", Sort: 5},
		{ParentId: 2, MenuType: 4, Path: "/api/user/associateRole", Name: "UserAssociateRole", Title: "分配角色", Sort: 6},
		{ParentId: 2, MenuType: 4, Path: "/api/user/delete", Name: "UserDelete", Title: "删除用户", Sort: 7},
		{ParentId: 2, MenuType: 4, Path: "/api/user/deleteMultiple", Name: "UserDeleteMultiple", Title: "批量删除", Sort: 8},

		{ParentId: 3, MenuType: 4, Path: "/api/role/getAll", Name: "RoleGetAll", Title: "获取全部", Sort: 1},
		{ParentId: 3, MenuType: 4, Path: "/api/role/getList", Name: "RoleGetList", Title: "获取列表", Sort: 2},
		{ParentId: 3, MenuType: 4, Path: "/api/role/getIdsByUserId", Name: "RoleGetIdsByUserId", Title: "获取ID列表(根据用户ID)", Sort: 3},
		{ParentId: 3, MenuType: 4, Path: "/api/role/create", Name: "RoleCreate", Title: "创建角色", Sort: 4},
		{ParentId: 3, MenuType: 4, Path: "/api/role/edit", Name: "RoleEdit", Title: "编辑角色", Sort: 5},
		{ParentId: 3, MenuType: 4, Path: "/api/role/editStatus", Name: "RoleEditStatus", Title: "修改状态", Sort: 6},
		{ParentId: 3, MenuType: 4, Path: "/api/role/associateMenu", Name: "RoleAssociateMenu", Title: "分配菜单", Sort: 7},
		{ParentId: 3, MenuType: 4, Path: "/api/role/delete", Name: "RoleDelete", Title: "删除角色", Sort: 8},

		{ParentId: 4, MenuType: 4, Path: "/api/department/getAll", Name: "DepartmentGetAll", Title: "获取全部", Sort: 1},
		{ParentId: 4, MenuType: 4, Path: "/api/department/getList", Name: "DepartmentGetList", Title: "获取列表", Sort: 2},
		{ParentId: 4, MenuType: 4, Path: "/api/department/create", Name: "DepartmentCreate", Title: "创建部门", Sort: 3},
		{ParentId: 4, MenuType: 4, Path: "/api/department/edit", Name: "DepartmentEdit", Title: "编辑部门", Sort: 4},
		{ParentId: 4, MenuType: 4, Path: "/api/department/delete", Name: "DepartmentDelete", Title: "删除部门", Sort: 5},

		{ParentId: 5, MenuType: 4, Path: "/api/menu/getAll", Name: "MenuGetAll", Title: "获取全部", Sort: 1},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/getTree", Name: "MenuGetTree", Title: "获取树级", Sort: 2},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/getIdsByRoleId", Name: "MenuGetIdsByRoleId", Title: "获取ID列表(根据角色ID)", Sort: 3},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/create", Name: "MenuCreate", Title: "创建菜单", Sort: 4},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/edit", Name: "MenuEdit", Title: "编辑菜单", Sort: 5},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/delete", Name: "MenuDelete", Title: "删除菜单", Sort: 6},
	}
)

func get_init_user_data() ([]User, error) {
	// 创建初始用户
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	var init_user_data = []User{
		{ID: 1, Status: 1, Username: "admin", Password: string(hashedPassword), Nickname: "Administrator", DepartmentID: 1, Roles: []Role{{ID: 1}}},
	}
	return init_user_data, nil
}
