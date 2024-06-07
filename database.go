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
func (da Database) InsertInitData() {
	// 获取初始化用户数据
	init_user_data, err := get_init_user_data()
	if err != nil {
		panic(err)
	}

	err = d.Database[d.LibraryGorm]{}.Get().InsertInitializationData(
		init_department_data,
		init_role_data,
		init_menu_data,
		init_user_data,
	)
	if err != nil {
		panic(err)
	}
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
		{ID: 6, Path: "/purchase", Title: "采购", Name: "Purchase", Sort: 1, Icon: "ep:shopping-cart-full"},
		{ID: 7, ParentId: 6, Path: "/supplier/list", Name: "SupplierList", Component: "supplier/list", Title: "供应商管理", Sort: 1},
		{ID: 8, Path: "/sale", Title: "销售", Name: "Sale", Sort: 2, Icon: "ep:suitcase"},
		{ID: 9, ParentId: 8, Path: "/product/list", Name: "ProductList", Component: "product/list", Title: "产品管理", Sort: 1},
		{ID: 10, ParentId: 8, Path: "/platform/list", Name: "PlatformList", Component: "platform/list", Title: "平台管理", Sort: 2},
		{ID: 11, ParentId: 8, Path: "/model/list", Name: "ModelList", Component: "model/list", Title: "型号管理", Sort: 3},
		{ID: 12, ParentId: 8, Path: "/quotation/list", Name: "QuotationList", Component: "quotation/list", Title: "报价单管理", Sort: 4},
		{ID: 13, ParentId: 8, Path: "/price/list", Name: "PriceList", Component: "price/list", Title: "价格管理", Sort: 5},
		{ID: 14, ParentId: 8, Path: "/profit/list", Name: "ProfitList", Component: "profit/list", Title: "利润管理", Sort: 6},

		{ParentId: 2, MenuType: 4, Path: "/api/user/create", Name: "UserCreate", Title: "创建", Sort: 1},
		{ParentId: 2, MenuType: 4, Path: "/api/user/edit", Name: "UserEdit", Title: "编辑", Sort: 2},
		{ParentId: 2, MenuType: 4, Path: "/api/user/editStatus", Name: "UserEditStatus", Title: "修改状态", Sort: 3},
		{ParentId: 2, MenuType: 4, Path: "/api/user/delete", Name: "UserDelete", Title: "删除", Sort: 4},
		{ParentId: 2, MenuType: 4, Path: "/api/user/deleteMultiple", Name: "UserDeleteMultiple", Title: "批量删除", Sort: 5},
		{ParentId: 2, MenuType: 4, Path: "/api/user/getList", Name: "UserGetList", Title: "获取列表", Sort: 6},
		{ParentId: 2, MenuType: 4, Path: "/api/user/associateRole", Name: "UserAssociateRole", Title: "分配角色", Sort: 7},

		{ParentId: 3, MenuType: 4, Path: "/api/role/create", Name: "RoleCreate", Title: "创建", Sort: 1},
		{ParentId: 3, MenuType: 4, Path: "/api/role/edit", Name: "RoleEdit", Title: "编辑", Sort: 2},
		{ParentId: 3, MenuType: 4, Path: "/api/role/editStatus", Name: "RoleEditStatus", Title: "修改状态", Sort: 3},
		{ParentId: 3, MenuType: 4, Path: "/api/role/delete", Name: "RoleDelete", Title: "删除", Sort: 4},
		{ParentId: 3, MenuType: 4, Path: "/api/role/getList", Name: "RoleGetList", Title: "获取列表", Sort: 5},
		{ParentId: 3, MenuType: 4, Path: "/api/role/getAll", Name: "RoleGetAll", Title: "获取全部", Sort: 6},
		{ParentId: 3, MenuType: 4, Path: "/api/role/getIdsByUserId", Name: "RoleGetIdsByUserId", Title: "根据用户ID获取ID", Sort: 7},
		{ParentId: 3, MenuType: 4, Path: "/api/role/associateMenu", Name: "RoleAssociateMenu", Title: "分配菜单", Sort: 8},

		{ParentId: 4, MenuType: 4, Path: "/api/department/getList", Name: "DepartmentGetList", Title: "获取列表", Sort: 1},
		{ParentId: 4, MenuType: 4, Path: "/api/department/getAll", Name: "DepartmentGetAll", Title: "获取全部", Sort: 2},
		{ParentId: 4, MenuType: 4, Path: "/api/department/create", Name: "DepartmentCreate", Title: "创建", Sort: 3},
		{ParentId: 4, MenuType: 4, Path: "/api/department/edit", Name: "DepartmentEdit", Title: "编辑", Sort: 4},
		{ParentId: 4, MenuType: 4, Path: "/api/department/delete", Name: "DepartmentDelete", Title: "删除", Sort: 5},

		{ParentId: 5, MenuType: 4, Path: "/api/menu/create", Name: "MenuCreate", Title: "创建", Sort: 1},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/edit", Name: "MenuEdit", Title: "编辑", Sort: 2},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/delete", Name: "MenuDelete", Title: "删除", Sort: 3},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/getTree", Name: "MenuGetTree", Title: "获取树级", Sort: 4},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/getAll", Name: "MenuGetAll", Title: "获取全部", Sort: 5},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/getIdsByRoleId", Name: "MenuGetIdsByRoleId", Title: "通过角色ID获取列表ID", Sort: 6},
		{ParentId: 5, MenuType: 4, Path: "/api/menu/getLeafIdsByRoleId", Name: "MenuGetLeafIdsByRoleId", Title: "通过角色ID获取最底级ID", Sort: 7},

		{ParentId: 7, MenuType: 4, Path: "/api/supplier/getList", Name: "SupplierGetList", Title: "获取列表", Sort: 1},
		{ParentId: 7, MenuType: 4, Path: "/api/supplier/getListByProductId", Name: "SupplierGetListByProductId", Title: "通过产品ID获取列表", Sort: 2},
		{ParentId: 7, MenuType: 4, Path: "/api/supplier/get", Name: "SupplierGet", Title: "获取", Sort: 3},
		{ParentId: 7, MenuType: 4, Path: "/api/supplier/create", Name: "SupplierCreate", Title: "创建", Sort: 4},
		{ParentId: 7, MenuType: 4, Path: "/api/supplier/edit", Name: "SupplierEdit", Title: "编辑", Sort: 5},

		{ParentId: 9, MenuType: 4, Path: "/api/product/get", Name: "ProductGet", Title: "获取", Sort: 1},
		{ParentId: 9, MenuType: 4, Path: "/api/product/getList", Name: "ProductGetList", Title: "获取列表", Sort: 2},
		{ParentId: 9, MenuType: 4, Path: "/api/product/getListBySupplierId", Name: "ProductGetListBySupplierId", Title: "通过供应商ID获取列表", Sort: 3},
		{ParentId: 9, MenuType: 4, Path: "/api/product/create", Name: "ProductCreate", Title: "创建", Sort: 4},
		{ParentId: 9, MenuType: 4, Path: "/api/product/edit", Name: "ProductEdit", Title: "编辑", Sort: 5},

		{ParentId: 10, MenuType: 4, Path: "/api/platform/getList", Name: "PlatformGetList", Title: "获取列表", Sort: 1},
		{ParentId: 10, MenuType: 4, Path: "/api/platform/get", Name: "PlatformGet", Title: "获取", Sort: 2},
		{ParentId: 10, MenuType: 4, Path: "/api/platform/create", Name: "PlatformCreate", Title: "创建", Sort: 3},
		{ParentId: 10, MenuType: 4, Path: "/api/platform/edit", Name: "PlatformEdit", Title: "编辑", Sort: 4},

		{ParentId: 11, MenuType: 4, Path: "/api/model/get", Name: "ModelGet", Title: "获取", Sort: 1},
		{ParentId: 11, MenuType: 4, Path: "/api/model/getListByProductId", Name: "ModelGetListByProductId", Title: "通过产品ID获取列表", Sort: 2},
		{ParentId: 11, MenuType: 4, Path: "/api/model/getList", Name: "ModelGetList", Title: "获取列表", Sort: 3},
		{ParentId: 11, MenuType: 4, Path: "/api/model/create", Name: "ModelCreate", Title: "创建", Sort: 4},
		{ParentId: 11, MenuType: 4, Path: "/api/model/edit", Name: "ModelEdit", Title: "编辑", Sort: 5},

		{ParentId: 12, MenuType: 4, Path: "/api/quotation/getList", Name: "QuotationGetList", Title: "获取列表", Sort: 1},
		{ParentId: 12, MenuType: 4, Path: "/api/quotation/get", Name: "QuotationGet", Title: "获取", Sort: 2},
		{ParentId: 12, MenuType: 4, Path: "/api/quotation/getListByProductId", Name: "QuotationGetListByProductId", Title: "通过产品ID获取列表", Sort: 3},
		{ParentId: 12, MenuType: 4, Path: "/api/quotation/getListByModelId", Name: "QuotationGetListByModelId", Title: "根据型号ID获取列表", Sort: 4},
		{ParentId: 12, MenuType: 4, Path: "/api/quotation/create", Name: "QuotationCreate", Title: "创建", Sort: 5},
		{ParentId: 12, MenuType: 4, Path: "/api/quotation/edit", Name: "QuotationEdit", Title: "编辑", Sort: 6},

		{ParentId: 13, MenuType: 4, Path: "/api/price/getList", Name: "PriceGetList", Title: "获取列表", Sort: 1},
		{ParentId: 13, MenuType: 4, Path: "/api/price/get", Name: "PriceGet", Title: "获取", Sort: 2},
		{ParentId: 13, MenuType: 4, Path: "/api/price/getListByProductId", Name: "PriceGetListByProductId", Title: "通过产品ID获取列表", Sort: 3},
		{ParentId: 13, MenuType: 4, Path: "/api/price/getListByModelId", Name: "PriceGetListByModelId", Title: "根据型号ID获取列表", Sort: 4},
		{ParentId: 13, MenuType: 4, Path: "/api/price/create", Name: "PriceCreate", Title: "创建", Sort: 5},
		{ParentId: 13, MenuType: 4, Path: "/api/price/edit", Name: "PriceEdit", Title: "编辑", Sort: 6},

		{ParentId: 14, MenuType: 4, Path: "/api/profit/getList", Name: "ProfitGetList", Title: "获取列表", Sort: 1},
		{ParentId: 14, MenuType: 4, Path: "/api/profit/get", Name: "ProfitGet", Title: "获取", Sort: 2},
		{ParentId: 14, MenuType: 4, Path: "/api/profit/getListByProductId", Name: "ProfitGetListByProductId", Title: "通过产品ID获取列表", Sort: 3},
		{ParentId: 14, MenuType: 4, Path: "/api/profit/create", Name: "ProfitCreate", Title: "创建", Sort: 4},
		{ParentId: 14, MenuType: 4, Path: "/api/profit/edit", Name: "ProfitEdit", Title: "编辑", Sort: 5},

		{ParentId: 14, MenuType: 4, Path: "/api/profit/item/get", Name: "ProfitItemGet", Title: "获取明细", Sort: 6},
		{ParentId: 14, MenuType: 4, Path: "/api/profit/item/getListByProfitId", Name: "ProfitItemGetListByProfitId", Title: "根据利润ID获取明细列表", Sort: 7},
		{ParentId: 14, MenuType: 4, Path: "/api/profit/item/create", Name: "ProfitItemCreate", Title: "创建", Sort: 8},
		{ParentId: 14, MenuType: 4, Path: "/api/profit/item/edit", Name: "ProfitItemEdit", Title: "编辑", Sort: 9},
		{ParentId: 14, MenuType: 4, Path: "/api/profit/item/delete", Name: "ProfitItemDelete", Title: "删除", Sort: 10},
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
