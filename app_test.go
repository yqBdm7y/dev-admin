package dadmin

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	d "github.com/yqBdm7y/devtool"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func mock_connect_database(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	// https://github.com/DATA-DOG/go-sqlmock/issues/268
	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	g := d.LibraryGorm{Open: func(dialector gorm.Dialector, opts ...gorm.Option) (db *gorm.DB, err error) {
		return gorm.Open(mysql.New(mysql.Config{
			Conn:                      mockDB,
			SkipInitializeWithVersion: true,
		}), &gorm.Config{})
	}}
	g.Init()
	return mockDB, mock
}

func TestAppInfo_CompareVersion(t *testing.T) {
	type args struct {
		app_version      string
		database_version string
	}
	tests := []struct {
		name string
		a    App
		args args
		want int
	}{
		// 测试大于
		{"测试APP版本大于数据库版本-1", App{}, args{"v1.0.0", "v0.1.0"}, 1},
		{"测试APP版本大于数据库版本-2", App{}, args{"v2.0.0", "v1.0.0"}, 1},
		{"测试APP版本大于数据库版本-3", App{}, args{"v0.2.0", "v0.1.0"}, 1},
		{"测试APP版本大于数据库版本-4", App{}, args{"v0.2.0", "v0.0.1"}, 1},
		{"测试APP版本大于数据库版本-5", App{}, args{"v0.0.2", "v0.0.1"}, 1},
		// 测试等于
		{"测试APP版本等于数据库版本-1", App{}, args{"v1.0.0", "v1.0.0"}, 0},
		{"测试APP版本等于数据库版本-2", App{}, args{"v0.4.0", "v0.4.0"}, 0},
		{"测试APP版本等于数据库版本-3", App{}, args{"v0.0.5", "v0.0.5"}, 0},
		{"测试APP版本等于数据库版本-4", App{}, args{"v3", "v3"}, 0},
		{"测试APP版本等于数据库版本-5", App{}, args{"v7", "v7.0.0"}, 0},
		// 测试小于
		{"测试APP版本小于数据库版本-1", App{}, args{"v0.0.3", "v0.0.5"}, -1},
		{"测试APP版本等于数据库版本-2", App{}, args{"v0.0.5", "v0.1.0"}, -1},
		{"测试APP版本等于数据库版本-3", App{}, args{"v0.5", "v1.1.0"}, -1},
		{"测试APP版本等于数据库版本-4", App{}, args{"v0.8", "v2"}, -1},
		{"测试APP版本等于数据库版本-5", App{}, args{"v6.0.0", "v7.0.0"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.CompareVersion(tt.args.app_version, tt.args.database_version); got != tt.want {
				t.Errorf("App.CompareVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppInfo_update(t *testing.T) {
	mockDB, mock := mock_connect_database(t)
	defer mockDB.Close()
	// https://studygolang.com/articles/27670
	const sql = "SELECT \\* FROM `apps`.+"
	rows := sqlmock.NewRows([]string{"id", "version"}).
		AddRow(1, "v0.0.0").
		AddRow(1, "v0.0.1").
		AddRow(1, "v0.0.3").
		AddRow(1, "v0.2.0")
	mock.ExpectQuery(sql).WillReturnRows(rows)
	mock.ExpectQuery(sql).WillReturnRows(rows)
	mock.ExpectQuery(sql).WillReturnRows(rows)
	mock.ExpectQuery(sql).WillReturnRows(rows)
	// UPDATE `apps` SET `created_at`='0000-00-00 00:00:00',`updated_at`='2024-06-06 23:52:33.156',`deleted_at`=NULL,`version`='v0.0.1' WHERE `apps`.`deleted_at` IS NULL AND `id` = 1
	mock.ExpectBegin() // begin transaction
	mock.ExpectExec("UPDATE `apps`").WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit() // commit transaction

	var history []string

	type args struct {
		db_version  string
		update_func map[string]func() error
	}
	tests := []struct {
		name string
		a    App
		args args
	}{
		{name: "Upgrading to v1.0.0", a: App{}, args: args{"v0.0.0", map[string]func() error{
			"v0.2.0": func() error {
				history = append(history, "v0.2.0")
				return nil
			},
			"v1.0.0": func() error {
				history = append(history, "v1.0.0")
				return nil
			},
			"v0.0.1": func() error {
				history = append(history, "v0.0.1")
				return nil
			},
			"v0.0.3": func() error {
				history = append(history, "v0.0.3")
				return nil
			},
			"v0.0.0": func() error {
				history = append(history, "v0.0.0")
				return nil
			},
		}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.a.update(tt.args.db_version, tt.args.update_func)
			if len(history) != 4 {
				t.Errorf("len(history) got = %v, want %v", len(history), 4)
			}
			if history[0] != "v0.0.1" {
				t.Errorf("history[0] got = %v, want %v", history[0], "v0.0.1")
			}
			if history[1] != "v0.0.3" {
				t.Errorf("history[1] got = %v, want %v", history[1], "v0.0.3")
			}
			if history[2] != "v0.2.0" {
				t.Errorf("history[2] got = %v, want %v", history[2], "v0.2.0")
			}
			if history[3] != "v1.0.0" {
				t.Errorf("history[3] got = %v, want %v", history[3], "v1.0.0")
			}
		})
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
