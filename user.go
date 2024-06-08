package dadmin

import (
	"errors"
	"time"

	"github.com/dlclark/regexp2"

	"github.com/gin-gonic/gin"
	d "github.com/yqBdm7y/devtool"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Status       int            `json:"status"`
	Avatar       string         `json:"avatar"`
	Username     string         `json:"username"`
	Password     string         `json:"-"`
	Nickname     string         `json:"nickname"`
	Phone        string         `json:"phone"`
	Email        string         `json:"email"`
	Sex          int            `json:"sex"`
	Remark       string         `json:"remark"`
	Roles        []Role         `gorm:"many2many:user_roles;" json:"roles"`
	DepartmentID uint           `json:"-"`
	Department   Department     `json:"department"`
}

// Create User
func (u User) Create(c *gin.Context) {
	type user struct {
		User
		Password string `json:"password"`
	}

	var form user
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	var checkDuplicate int64
	d.Database[d.LibraryGorm]{}.Get().DB.Model(&User{}).Where("username = ?", form.Username).Count(&checkDuplicate)
	if checkDuplicate > 0 {
		d.Gin{}.Error(c, Err(errors.New("duplicate username")))
		return
	}

	// Encryption password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}
	form.Password = string(hashedPassword)

	result := d.Database[d.LibraryGorm]{}.Get().DB.Create(&form)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}

	d.Gin{}.Success(c, Success(form.ID))
}

// Edit user
func (u User) Edit(c *gin.Context) {
	var form User
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	var checkDuplicate int64
	d.Database[d.LibraryGorm]{}.Get().DB.Model(&User{}).Where("id != ?", form.ID).Where("username = ?", form.Username).Count(&checkDuplicate)
	if checkDuplicate > 0 {
		d.Gin{}.Error(c, Err(errors.New("duplicate username")))
		return
	}

	result := d.Database[d.LibraryGorm]{}.Get().DB.Omit("password", "created_at").Save(&form)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}

	d.Gin{}.Success(c, Success(form.ID))
}

// Edit user status
func (u User) EditStatus(c *gin.Context) {
	var form User
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	if form.Status != 0 && form.Status != 1 {
		d.Gin{}.Error(c, Err(errors.New("status format is incorrect")))
		return
	}

	// If the user is a super administrator, the status cannot be modified
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

// Edit user password
func (u User) EditPassword(c *gin.Context) {
	type user struct {
		User
		Password string `json:"password"`
	}

	var form user
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	b := u.ValidatePassword(form.Password)
	if !b {
		d.Gin{}.Error(c, Err(errors.New("incorrect password format")))
		return
	}
	// Create a new password hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(form.Password), bcrypt.DefaultCost)
	if err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}
	form.Password = string(hashedPassword)

	result := d.Database[d.LibraryGorm]{}.Get().DB.Model(&form).Update("password", form.Password)
	if result.Error != nil {
		d.Gin{}.Error(c, Err(result.Error))
		return
	}
	d.Gin{}.Success(c, Success(form.ID))
}

// Delete user
func (u User) Delete(c *gin.Context) {
	var form User
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	// If the user is a super administrator, it cannot be deleted
	if form.ID == 1 {
		d.Gin{}.Error(c, Err(errors.New("unable to delete super administrator")))
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

// Delete Multiple Users
func (u User) DeleteMultiple(c *gin.Context) {
	var form []User
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	// If the user is a super administrator, it cannot be deleted
	for _, v := range form {
		if v.ID == 1 {
			d.Gin{}.Error(c, Err(errors.New("unable to delete super administrator")))
			return
		}
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

	var idList []int
	for _, v := range form {
		idList = append(idList, int(v.ID))
	}

	d.Gin{}.Success(c, Success(idList))
}

func (u User) GetList(c *gin.Context) {
	var query = d.Database[d.LibraryGorm]{}.Get().DB.Model(&User{}).Preload("Department").Order("created_at desc")
	var data []User
	p, err := d.Gin{}.GetListWithFuzzyQuery(c, query, []string{"username", "phone", "status", "department_id"}, &data)
	if err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}
	v := p.(d.LibraryPagination)
	v.DataList = data

	d.Gin{}.Success(c, Success(v.ToMap()))
}

func (u User) AssociateRole(c *gin.Context) {
	var form User
	if err := c.ShouldBindJSON(&form); err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	err := d.Database[d.LibraryGorm]{}.Get().DB.Model(&form).Association("Roles").Replace(form.Roles)
	if err != nil {
		d.Gin{}.Error(c, Err(err))
		return
	}

	d.Gin{}.Success(c, Success(form.ID))
}

// Verify user password
func (u User) VerifyPassword(stored_hash, input_password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(stored_hash), []byte(input_password))
	return err == nil
}

func (u User) ValidatePassword(password string) bool {
	// Regular expression matching rules: 8-18 characters, including at least two of numbers, letters, and symbols
	re := regexp2.MustCompile(`^(?=.*[0-9])(?=.*[a-zA-Z!@#$%^&*()_+])[0-9a-zA-Z!@#$%^&*()_+]{8,18}$`, 0)

	if b, _ := re.MatchString(password); b {
		return true
	}
	return false
}
