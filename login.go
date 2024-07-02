package dadmin

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	d "github.com/yqBdm7y/devtool"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	jwt2 "github.com/golang-jwt/jwt/v4"
)

type Login struct{}

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

// 处理空路由
func (l Login) HandleNoRoute() func(c *gin.Context) {
	return func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	}
}

type login struct {
	Key        string
	Timeout    time.Duration
	MaxRefresh time.Duration
	Field      struct {
		AccessToken  string
		RefreshToken string
		Expire       string
	}
}

type login_option interface {
	apply(login *login)
}

// Set default key
type login_option_key string

func (l login_option_key) apply(login *login) {
	login.Key = string(l)
}

func LoginWithKey(key string) login_option {
	return login_option_key(key)
}

// Set default Timeout
type login_option_timeout time.Duration

func (l login_option_timeout) apply(login *login) {
	login.Timeout = time.Duration(l)
}

func LoginWithTimeout(value time.Duration) login_option {
	return login_option_timeout(value)
}

// Set default MaxRefresh
type login_option_max_refresh time.Duration

func (l login_option_max_refresh) apply(login *login) {
	login.MaxRefresh = time.Duration(l)
}

func LoginWithMaxRefresh(value time.Duration) login_option {
	return login_option_max_refresh(value)
}

// Set Custom field - Access Token
type login_option_custom_field_access_token string

func (l login_option_custom_field_access_token) apply(login *login) {
	login.Field.AccessToken = string(l)
}

func LoginWithCustomFieldAccessToken(custom_field string) login_option {
	return login_option_custom_field_access_token(custom_field)
}

// Set Custom field - Refresh Token
type login_option_custom_field_refresh_token string

func (l login_option_custom_field_refresh_token) apply(login *login) {
	login.Field.RefreshToken = string(l)
}

func LoginWithCustomFieldRefreshToken(custom_field string) login_option {
	return login_option_custom_field_refresh_token(custom_field)
}

// Set Custom field - Expire
type login_option_custom_field_expire string

func (l login_option_custom_field_expire) apply(login *login) {
	login.Field.Expire = string(l)
}

func LoginWithCustomFieldExpire(custom_field string) login_option {
	return login_option_custom_field_expire(custom_field)
}

func LoginNew(opts ...login_option) *jwt.GinJWTMiddleware {
	// 如果没有设置KEY，则使用随机字符串
	defaultKey, err := d.String{}.GenerateRandomString(64)
	if err != nil {
		panic(err)
	}

	// Set default parameters
	l := &login{
		Key:        defaultKey,
		Timeout:    time.Hour * 2,
		MaxRefresh: time.Hour * 24,
		Field: struct{ AccessToken, RefreshToken, Expire string }{
			AccessToken:  "token",
			RefreshToken: "refresh_token",
			Expire:       "expire",
		},
	}

	for _, opt := range opts {
		opt.apply(l)
	}

	return &jwt.GinJWTMiddleware{
		Realm:           "dev-admin",
		Key:             []byte(l.Key),
		Timeout:         l.Timeout,
		MaxRefresh:      l.MaxRefresh,
		PayloadFunc:     l.PayloadFunc(),
		Authenticator:   l.Authenticator(),
		Authorizator:    l.Authorizator(),
		LoginResponse:   l.LoginResponse(),
		Unauthorized:    l.Unauthorized(),
		LogoutResponse:  l.LogoutResponse(),
		RefreshResponse: l.RefreshResponse(),
		TokenLookup:     "header: Authorization, query: " + l.Field.RefreshToken,
	}
}

// 登录验证
func (l login) Authenticator() func(c *gin.Context) (interface{}, error) {
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
func (l login) PayloadFunc() func(data interface{}) jwt.MapClaims {
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
func (l login) Authorizator() func(data interface{}, c *gin.Context) bool {
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

func (l login) LoginResponse() func(c *gin.Context, code int, message string, time time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		// 解析当前token
		oriToken, err := jwt2.Parse(token, func(t *jwt2.Token) (interface{}, error) {
			return []byte(l.Key), nil
		})
		if err != nil {
			d.Gin{}.Error(c, Err(err))
			return
		}

		// 获取 claims
		oriClaims, ok := oriToken.Claims.(jwt2.MapClaims)
		if !ok {
			d.Gin{}.Error(c, Err(errors.New("invalid claims type")))
			return
		}

		// 获取特定的 claim 值
		userId, exists := oriClaims["id"]
		if !exists {
			d.Gin{}.Error(c, Err(errors.New("id claim not found")))
			return
		}

		// 生成刷新令牌
		refreshToken := jwt2.New(jwt2.GetSigningMethod("HS256"))
		claims := refreshToken.Claims.(jwt2.MapClaims)
		claims["orig_iat"] = time.Now().Unix()
		claims["ip"] = c.ClientIP()
		claims["id"] = userId
		tokenString, _ := refreshToken.SignedString([]byte(l.Key))

		api := d.Api[d.LibraryApi]{}.Get()
		api.Response.Data = gin.H{
			l.Field.AccessToken:  token,
			l.Field.Expire:       expire,
			l.Field.RefreshToken: tokenString,
		}
		d.Gin{}.Success(c, api)
	}
}

func (l login) Unauthorized() func(c *gin.Context, code int, message string) {
	return func(c *gin.Context, code int, message string) {
		api := d.Api[d.LibraryApi]{}.Get()
		api.Response.Code = code
		api.Response.Message = message
		d.Gin{}.Error(c, api)
	}
}

func (l login) LogoutResponse() func(c *gin.Context, code int) {
	return func(c *gin.Context, code int) {
		api := d.Api[d.LibraryApi]{}.Get()
		api.Response.Code = code
		d.Gin{}.Success(c, api)
	}
}

func (l login) RefreshResponse() func(c *gin.Context, code int, token string, expire time.Time) {
	return func(c *gin.Context, code int, token string, expire time.Time) {
		// 判断是否有refresh token
		refToken := c.Query(l.Field.RefreshToken)
		if len(refToken) != 0 {
			// 解析当前refresh token
			oriToken, err := jwt2.Parse(refToken, func(t *jwt2.Token) (interface{}, error) {
				return []byte(l.Key), nil
			})
			if err != nil {
				d.Gin{}.Error(c, Err(err))
				return
			}

			// 获取 claims
			oriClaims, ok := oriToken.Claims.(jwt2.MapClaims)
			if !ok {
				d.Gin{}.Error(c, Err(errors.New("invalid claims type")))
				return
			}

			// 获取特定的 claim 值
			ipAddr, exists := oriClaims["ip"]
			if !exists {
				d.Gin{}.Error(c, Err(errors.New("ip claim not found")))
				return
			}

			// 判断当前IP是否和refresh token中的IP一致
			if c.ClientIP() != ipAddr {
				d.Gin{}.Error(c, Err(errors.New("ip校验不通过")))
				return
			}
		}

		api := d.Api[d.LibraryApi]{}.Get()
		api.Response.Code = http.StatusOK
		api.Response.Data = gin.H{
			l.Field.AccessToken:  token,
			l.Field.Expire:       expire.Format(time.RFC3339),
			l.Field.RefreshToken: c.Query(l.Field.RefreshToken),
		}
		d.Gin{}.Success(c, api)
	}
}
