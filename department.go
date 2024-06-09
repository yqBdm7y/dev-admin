package dadmin

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	d "github.com/yqBdm7y/devtool"
	"gorm.io/gorm"
)

type Department struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Status    int            `json:"status"`
	Name      string         `json:"name"`
	ParentId  int            `json:"parent_id"`
	Sort      int            `json:"sort"`
	Phone     string         `json:"phone"`
	Principal string         `json:"principal"`
	Email     string         `json:"email"`
	Type      int            `json:"type"`
	Remark    string         `json:"remark"`
}

func (de Department) GetAll(c *gin.Context) {
	var data []Department
	d.Database[d.LibraryGorm]{}.Get().DB.Order("sort asc").Order("created_at desc").Find(&data)

	d.Gin{}.Success(c, Success(data))
}

func (de Department) GetList(c *gin.Context) {
	var query = d.Database[d.LibraryGorm]{}.Get().DB.Model(&Department{}).Order("created_at desc")

	var data []Department
	p, err := d.Gin{}.GetListWithFuzzyQuery(c, query, nil, &data)
	if err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}
	v := p.(d.LibraryPagination)
	v.DataList = data

	d.Gin{}.Success(c, Success(v.ToMap()))
}

func (de Department) Create(c *gin.Context) {
	var form Department
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

func (de Department) Edit(c *gin.Context) {
	var form Department
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

func (de Department) Delete(c *gin.Context) {
	var form Department
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	var total int64
	d.Database[d.LibraryGorm]{}.Get().DB.Model(&Department{}).Where("parent_id = ?", form.ID).Count(&total)
	if total != 0 {
		d.Gin{}.Error(c, Err(errors.New("please delete the sub-department first")))
		return
	}

	result := d.Database[d.LibraryGorm]{}.Get().DB.Delete(&form)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}

	d.Gin{}.Success(c, Success(form.ID))
}
