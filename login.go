package dadmin

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	d "github.com/yqBdm7y/devtool"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type Login struct{}

func (l Login) Init() *jwt.GinJWTMiddleware {
	randomKey, err := d.String{}.GenerateRandomString(64)
	if err != nil {
		panic(err)
	}
	b := d.Config[d.LibraryViper]{}.Get().GetBool(ConfigPathIsDebug)
	if b {
		randomKey = ""
	}
	return &jwt.GinJWTMiddleware{
		Realm:           "erp",
		Key:             []byte(randomKey),
		MaxRefresh:      time.Hour,
		PayloadFunc:     l.PayloadFunc(),
		Authenticator:   l.Authenticator(),
		Authorizator:    l.Authorizator(),
		LoginResponse:   l.LoginResponse(),
		Unauthorized:    l.Unauthorized(),
		LogoutResponse:  l.LogoutResponse(),
		RefreshResponse: l.RefreshResponse(),
	}
}

func (l Login) HandleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}

// 登录验证
func (l Login) Authenticator() func(c *gin.Context) (interface{}, error) {
	type login struct {
		Username     string `json:"username" binding:"required"`
		Password     string `json:"password" binding:"required"`
		CaptchaToken string `json:"captcha_token"`
	}

	return func(c *gin.Context) (interface{}, error) {
		var loginVals login
		if err := c.ShouldBindJSON(&loginVals); err != nil {
			return "", jwt.ErrMissingLoginValues
		}

		// 如果是debug模式，无需验证验证码即可登录
		bo := d.Config[d.LibraryViper]{}.Get().GetBool(ConfigPathIsDebug)
		if !bo {
			err := d.Captcha[d.LibraryTurnstile]{}.Get().VerifyToken(loginVals.CaptchaToken)
			if err != nil {
				return "", err
			}
		}

		userName := loginVals.Username
		password := loginVals.Password

		var usr User
		d.Database[d.LibraryGorm]{}.Get().DB.Where("username = ?", userName).First(&usr)
		// 检查用户状态，如果用户状态被禁用，则登录失败
		if usr.Status != 1 {
			return "", jwt.ErrFailedAuthentication
		}

		b := User{}.VerifyPassword(usr.Password, password)
		if b {
			return &User{
				ID:       usr.ID,
				Username: usr.Username,
				Nickname: usr.Nickname,
			}, nil
		}

		return nil, jwt.ErrFailedAuthentication
	}
}

// 设置Payload
func (l Login) PayloadFunc() func(data interface{}) jwt.MapClaims {
	return func(data interface{}) jwt.MapClaims {
		if v, ok := data.(*User); ok {
			return jwt.MapClaims{
				"id":       v.ID,
				"username": v.Username,
				"nickname": v.Nickname,
			}
		}
		return jwt.MapClaims{}
	}
}

// 校验权限
func (l Login) Authorizator() func(data interface{}, c *gin.Context) bool {
	return func(data interface{}, c *gin.Context) bool {

		claims := jwt.ExtractClaims(c)

		id, err := strconv.Atoi(fmt.Sprintf("%v", claims["id"]))
		if err != nil {
			return false
		}

		// 如果是ID为1的则为超级管理员
		if id == 1 {
			return true
		}

		var u User
		result := d.Database[d.LibraryGorm]{}.Get().DB.Debug().Preload("Roles", "status = ?", 1).Preload("Roles.Menus").First(&u, id)
		if result.Error != nil {
			return false
		}

		// 如果角色ID为1的则为超级管理员
		for _, v := range u.Roles {
			if v.ID == 1 {
				return true
			}
		}

		// 检查用户状态，如果用户状态被禁用，则登录失败
		if u.Status != 1 {
			return false
		}

		// 把所有有权限的菜单存到menuList中
		var menuList = make(map[string]bool)
		for _, v := range u.Roles {
			for _, sv := range v.Menus {
				menuList[sv.Path] = true
			}
		}

		// 验证当前api是否在权限范围内
		return menuList[c.Request.URL.Path]
	}
}

func (l Login) LoginResponse() func(c *gin.Context, code int, message string, time time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		api := d.Api[d.LibraryApi]{}.Get()
		api.Response.Data = gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		}
		d.Gin{}.Success(c, api)
	}
}

func (l Login) Unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		api := d.Api[d.LibraryApi]{}.Get()
		api.Response.Code = code
		api.Response.Message = message
		d.Gin{}.Error(c, api)
	}
}

func (l Login) LogoutResponse() func(c *gin.Context, code int) {
	return func(c *gin.Context, code int) {
		api := d.Api[d.LibraryApi]{}.Get()
		api.Response.Code = code
		d.Gin{}.Success(c, api)
	}
}

func (l Login) RefreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		api := d.Api[d.LibraryApi]{}.Get()
		api.Response.Code = http.StatusOK
		api.Response.Data = gin.H{
			"token":  token,
			"expire": expire.Format(time.RFC3339),
		}
		d.Gin{}.Success(c, api)
	}
}

// 获取用户
func (l Login) GetUser(c *gin.Context) (u User) {
	claims := jwt.ExtractClaims(c)
	u = User{
		ID:       uint(claims["id"].(float64)),
		Username: claims["username"].(string),
		Nickname: claims["nickname"].(string),
	}
	return u
}
