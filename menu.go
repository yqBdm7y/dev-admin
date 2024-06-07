package dadmin

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	d "github.com/yqBdm7y/devtool"
	"gorm.io/gorm"
)

type Menu struct {
	ID              uint           `gorm:"primarykey" json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ParentId        int            `json:"parent_id"`
	MenuType        int            `json:"menu_type"`
	Title           string         `json:"title"`
	Name            string         `json:"name"`
	Path            string         `json:"path"`
	Component       string         `json:"component"`
	Rank            int            `json:"rank"`
	Sort            int            `json:"sort"`
	Redirect        string         `json:"redirect"`
	Icon            string         `json:"icon"`
	ExtraIcon       string         `json:"extra_icon"`
	EnterTransition string         `json:"enter_transition"`
	LeaveTransition string         `json:"leave_transition"`
	ActivePath      string         `json:"active_path"`
	Auths           string         `json:"auths"`
	FrameSrc        string         `json:"frame_src"`
	FrameLoading    bool           `json:"frame_loading"`
	KeepAlive       bool           `json:"keep_alive"`
	HiddenTag       bool           `json:"hidden_tag"`
	FixedTag        bool           `json:"fixed_tag"`
	ShowLink        bool           `json:"show_link"`
	ShowParent      bool           `json:"show_parent"`
	Roles           []Role         `gorm:"many2many:role_menus;" json:"roles"`
}

const menu_type_api = 4

type menu_tree struct {
	ParentID   int                  `json:"parentId,omitempty"`
	ID         int                  `json:"id,omitempty"`
	MenuType   int                  `json:"menuType"`
	Path       string               `json:"path"`
	Name       string               `json:"name"`
	Redirect   string               `json:"redirect,omitempty"`
	Component  string               `json:"component"`
	Meta       menu_tree_meta       `json:"meta,omitempty"`
	Transition menu_tree_transition `json:"transition,omitempty"`
	HiddenTag  bool                 `json:"hiddenTag,omitempty"`
	ActivePath string               `json:"activePath,omitempty"`
	FixedTag   bool                 `json:"fixedTag,omitempty"`
	Sort       int                  `json:"sort"`
	Children   []menu_tree          `json:"children,omitempty"`
}

type menu_tree_meta struct {
	Title        string `json:"title"`
	Icon         string `json:"icon,omitempty"`
	ExtraIcon    string `json:"extraIcon,omitempty"`
	ShowLink     bool   `json:"showLink,omitempty"`
	Rank         int    `json:"rank,omitempty"`
	ShowParent   bool   `json:"showParent,omitempty"`
	Auths        string `json:"auths,omitempty"`
	KeepAlive    bool   `json:"keepAlive,omitempty"`
	FrameSrc     string `json:"frameSrc,omitempty"`
	FrameLoading bool   `json:"frameLoading"`
}

type menu_tree_transition struct {
	EnterTransition string `json:"enterTransition,omitempty"`
	LeaveTransition string `json:"leaveTransition,omitempty"`
}

// Tree interface function
func (m menu_tree) IsTop() bool {
	return m.ParentID == 0
}
func (m menu_tree) GetId() int {
	return int(m.ID)
}
func (m menu_tree) GetParentId() int {
	return m.ParentID
}
func (m *menu_tree) AppendChildren(child menu_tree) {
	m.Children = append(m.Children, child)
}

// d.InterfaceSort
func (m menu_tree) GetSort() int {
	return m.Sort
}
func (m menu_tree) GetChildren() interface{} {
	return m.Children
}

// 获取全部菜单
func (m Menu) GetAll(c *gin.Context) {
	var data []Menu
	d.Database[d.LibraryGorm]{}.Get().DB.Order("sort asc").Order("created_at desc").Find(&data)

	d.Gin{}.Success(c, Success(data))
}

// 获取树状菜单
func (m Menu) GetTree(c *gin.Context) {
	// 获取用户
	u := Login{}.GetUser(c)

	var dMenus []Menu

	var getAllMenus = func() {
		d.Database[d.LibraryGorm]{}.Get().DB.Where("menu_type != ?", menu_type_api).Find(&dMenus)
	}

	// 如果是超级管理员，则直接获取全部菜单
	if u.ID == 1 {
		getAllMenus()
	} else {
		// 如果不是管理员
		d.Database[d.LibraryGorm]{}.Get().DB.Preload("Roles", "status = ?", 1).Preload("Roles.Menus", "menu_type != ?", menu_type_api).First(&u, u.ID)
		var tmp = make(map[uint]Menu)
		// 遍历期间，顺便判断是不是role.id=1的超级管理员角色，如果是则跳过
		var superRole bool
		for _, v := range u.Roles {
			if v.ID == 1 {
				superRole = true
			}
			for _, sv := range v.Menus {
				tmp[sv.ID] = sv
			}
		}
		// 如果是超级管理员的角色，则获取全部菜单
		if superRole {
			getAllMenus()
		} else {
			for _, v := range tmp {
				dMenus = append(dMenus, v)
			}
		}
	}

	var rtnData []menu_tree
	for _, v := range dMenus {
		rtnData = append(rtnData, menu_tree{
			ParentID:  v.ParentId,
			ID:        int(v.ID),
			MenuType:  v.MenuType,
			Path:      v.Path,
			Name:      v.Name,
			Redirect:  v.Redirect,
			Component: v.Component,
			Meta: menu_tree_meta{
				Title:        v.Title,
				Icon:         v.Icon,
				ExtraIcon:    v.ExtraIcon,
				ShowLink:     v.ShowLink,
				Rank:         v.Rank,
				ShowParent:   v.ShowParent,
				Auths:        v.Auths,
				KeepAlive:    v.KeepAlive,
				FrameSrc:     v.FrameSrc,
				FrameLoading: v.FrameLoading,
			},
			Transition: menu_tree_transition{
				EnterTransition: v.EnterTransition,
				LeaveTransition: v.LeaveTransition,
			},
			HiddenTag:  v.HiddenTag,
			ActivePath: v.ActivePath,
			FixedTag:   v.FixedTag,
			Sort:       v.Sort,
		})
	}

	treeData := d.GenerateTree(rtnData)
	d.SortListWithChildrenBySortField(treeData)

	d.Gin{}.Success(c, Success(treeData))
}

// 创建菜单
func (m Menu) Create(c *gin.Context) {
	var form Menu
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	result := d.Database[d.LibraryGorm]{}.Get().DB.Create(&form)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}
	d.Gin{}.Success(c, Success(form.ID))
}

// 编辑菜单
func (m Menu) Edit(c *gin.Context) {
	var form Menu
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	result := d.Database[d.LibraryGorm]{}.Get().DB.Save(&form)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}

	d.Gin{}.Success(c, Success(form.ID))
}

// 删除菜单
func (m Menu) Delete(c *gin.Context) {
	var form Menu
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}
	// 检查是否有子菜单，若有子菜单则不能删除
	b := m.CheckIfSubmenuExist(form)
	if b {
		d.Gin{}.Error(c, Err(errors.New("please delete the submenu first")))
		return
	}

	err := d.Database[d.LibraryGorm]{}.Get().DB.Transaction(func(tx *gorm.DB) error {
		result := tx.Delete(&form)
		if result.Error != nil {
			return result.Error
		}

		err := tx.Model(&form).Association("Roles").Clear()
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	d.Gin{}.Success(c, Success(form.ID))
}

// 根据角色 id 查对应菜单ID
func (m Menu) GetIdsByRoleId(c *gin.Context) {
	// 获取 URL 中的参数
	idStr := c.Query("role_id")

	// 将字符串格式的 ID 转化为整数
	rid, err := strconv.Atoi(idStr)
	if err != nil || rid <= 0 {
		// 处理转化错误
		d.Gin{}.Error(c, Err(err))
		return
	}

	var r Role
	result := d.Database[d.LibraryGorm]{}.Get().DB.Preload("Menus").First(&r, rid)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}
	var ids []int
	for _, v := range r.Menus {
		ids = append(ids, int(v.ID))
	}

	d.Gin{}.Success(c, Success(ids))
}

// 根据角色 id 查对应叶子菜单ID
func (m Menu) GetLeafIdsByRoleId(c *gin.Context) {
	// 获取 URL 中的参数
	idStr := c.Query("role_id")

	// 将字符串格式的 ID 转化为整数
	rid, err := strconv.Atoi(idStr)
	if err != nil || rid <= 0 {
		// 处理转化错误
		d.Gin{}.Error(c, Err(err))
		return
	}

	var (
		r   Role
		ids []int
	)
	result := d.Database[d.LibraryGorm]{}.Get().DB.Preload("Menus").First(&r, rid)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}

	menus := d.FilterLeafNode(r.Menus, func(dm Menu) int {
		return int(dm.ID)
	}, func(dm Menu) int {
		return int(dm.ParentId)
	})

	for _, v := range menus {
		ids = append(ids, int(v.ID))
	}

	d.Gin{}.Success(c, Success(ids))
}

// 检查子菜单是否存在
func (m Menu) CheckIfSubmenuExist(form Menu) bool {
	var total int64
	d.Database[d.LibraryGorm]{}.Get().DB.Model(&Menu{}).Where("parent_id = ?", form.ID).Count(&total)

	return total > 0
}
