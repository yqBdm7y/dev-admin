package dadmin

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	d "github.com/yqBdm7y/devtool"
	"gorm.io/gorm"
)

// 应用信息
type App struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	Version   string         // 应用版本
}

const (
	ConfigPathIsDebug = "debug"
)

// 检查APP是否需要更新
func (a App) CheckUpdate(current_version string, update_func map[string]func() error) {
	dbVer := a.CheckDBVersion(current_version)    //检查数据库版本
	i := a.CompareVersion(current_version, dbVer) // 比较数据库版本
	switch i {
	// APP版本和数据库版本相同
	case 0:
		return
		// APP版本小于数据库版本
	case -1:
		panic(fmt.Sprintf("当前APP版本过低，为避免破坏数据结构，请先手动升级APP，再链接数据库！当前APP版本：%v， 数据库版本:%v", current_version, dbVer))
		// APP版本大于数据库版本，APP比数据库新，进行数据库的数据更新
	case 1:
		a.update(dbVer, update_func)
	}
}

// APP更新
func (a App) update(db_version string, update_func map[string]func() error) {
	type version_upgrade struct {
		Version string
		Upgrade func() error
	}

	// 将版本号和升级方法存入Slice
	var vu []version_upgrade
	for k, v := range update_func {
		vu = append(vu, version_upgrade{Version: k, Upgrade: v})
	}

	// 按版本号排序
	sort.Slice(vu, func(i, j int) bool {
		return vu[i].Version < vu[j].Version
	})

	// 按顺序执行升级方法
	for _, v := range vu {
		i := a.CompareVersion(db_version, v.Version)
		if i < 0 {
			err := v.Upgrade()
			if err != nil {
				panic(err)
			}
			db_version = v.Version // 更新完数据库版本为当前升级的版本
			a.SetDBVersion(db_version)
		}
	}
}

// 检查当前APP的数据库版本，如果查不到数据，则把当前版本号写入数据库中，并返回当前版本号
func (a App) CheckDBVersion(current_version string) (version string) {
	var db = d.Database[d.LibraryGorm]{}.Get().DB
	var info App
	//如果查不到数据，则认为当前是最新版本，则把当前版本号写入数据库中，并返回当前版本号
	err := db.First(&info).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		info.Version = current_version
		db.Create(&info)
	}
	return info.Version
}

// 设置数据库中的APP版本
func (a App) SetDBVersion(version string) error {
	var db = d.Database[d.LibraryGorm]{}.Get().DB
	var info App
	err := db.First(&info).Error
	if err != nil {
		return err
	}
	info.Version = version
	db.Debug().Save(&info)
	return nil
}

// 对比当前version1和version2
// 如果version1 ＞ version2，返回1
// 如果version1 = version2，返回0
// 如果version1 ＜ version2，返回-1
func (a App) CompareVersion(version1, version2 string) int {
	// 先去掉版本开头的v字符串，然后分割版本号，按照"主版本号.次版本号.修订号"进行比较
	v1Parts := strings.Split(strings.TrimPrefix(version1, "v"), ".")
	v2Parts := strings.Split(strings.TrimPrefix(version2, "v"), ".")

	for i := 0; i < len(v1Parts) || i < len(v2Parts); i++ {
		v1Part := "0"
		if i < len(v1Parts) {
			v1Part = v1Parts[i]
		}
		v2Part := "0"
		if i < len(v2Parts) {
			v2Part = v2Parts[i]
		}
		var v1Num, v2Num int
		_, err1 := fmt.Sscanf(v1Part, "%d", &v1Num)
		_, err2 := fmt.Sscanf(v2Part, "%d", &v2Num)
		if err1 != nil || err2 != nil {
			// 如果无法转换为整数，则直接比较字符串
			if v1Part < v2Part {
				return -1
			} else if v1Part > v2Part {
				return 1
			}
		} else {
			// 整数比较
			if v1Num < v2Num {
				return -1
			} else if v1Num > v2Num {
				return 1
			}
		}
	}
	return 0
}
