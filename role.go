package dadmin

import (
	"errors"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	d "github.com/yqBdm7y/devtool"
	"gorm.io/gorm"
)

type Role struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Status    int            `json:"status"`
	Name      string         `json:"name"`
	Code      string         `json:"code"`
	Remark    string         `json:"remark"`
	Users     []User         `gorm:"many2many:user_roles;" json:"users"`
	Menus     []Menu         `gorm:"many2many:role_menus;" json:"menus"`
}

func (r Role) Create(c *gin.Context) {
	var form Role
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

func (r Role) Edit(c *gin.Context) {
	var form Role
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	result := d.Database[d.LibraryGorm]{}.Get().DB.Omit("created_at").Save(&form)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}

	d.Gin{}.Success(c, Success(form.ID))
}

// Edit role status
func (r Role) EditStatus(c *gin.Context) {
	var form Role
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	if form.Status != 0 && form.Status != 1 {
		d.Gin{}.Error(c, Err(errors.New("status format is incorrect")))
		return
	}

	// If the role is a super administrator, the status cannot be modified
	if form.ID == 1 {
		d.Gin{}.Error(c, Err(errors.New("super administrator cannot modify status")))
		return
	}

	result := d.Database[d.LibraryGorm]{}.Get().DB.Model(&form).Update("status", form.Status)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}

	d.Gin{}.Success(c, Success(form.ID))
}

// Delete role
func (r Role) Delete(c *gin.Context) {
	var form Role
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	// If the role is a super administrator, it cannot be deleted
	if form.ID == 1 {
		d.Gin{}.Error(c, Err(errors.New("unable to delete super administrator")))
		return
	}

	result := d.Database[d.LibraryGorm]{}.Get().DB.Delete(&form)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}

	d.Gin{}.Success(c, Success(form.ID))
}

func (r Role) GetList(c *gin.Context) {
	var query = d.Database[d.LibraryGorm]{}.Get().DB.Model(&Role{}).Order("created_at desc")
	var data []Role
	pr, err := d.Gin{}.GetListWithFuzzyQuery(c, query, []string{"name", "code", "status"}, &data)
	if err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}
	v := pr.(d.LibraryPagination)
	v.DataList = data
	d.Gin{}.Success(c, Success(v.ToMap()))
}

func (r Role) GetAll(c *gin.Context) {
	var data []Role
	d.Database[d.LibraryGorm]{}.Get().DB.Order("created_at desc").Find(&data)
	d.Gin{}.Success(c, Success(data))
}

func (r Role) GetIdsByUserId(c *gin.Context) {
	// 获取 URL 中的参数
	idStr := c.Query("user_id")

	// 将字符串格式的 ID 转化为整数
	uid, err := strconv.Atoi(idStr)
	if err != nil || uid <= 0 {
		// 处理转化错误
		d.Gin{}.Error(c, Err(errors.New("invalid ID format")))
		return
	}

	var usr User
	d.Database[d.LibraryGorm]{}.Get().DB.Preload("Roles").First(&usr, uid)
	var ids []uint
	for _, v := range usr.Roles {
		ids = append(ids, v.ID)
	}
	d.Gin{}.Success(c, Success(ids))
}

// 角色关联菜单
func (r Role) AssociateMenu(c *gin.Context) {
	var form Role
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	err := d.Database[d.LibraryGorm]{}.Get().DB.Model(&form).Association("Menus").Replace(form.Menus)
	if err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}
	d.Gin{}.Success(c, Success(form.Menus))
}
