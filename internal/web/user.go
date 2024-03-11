package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/liang21/webook/internal/domain"
	"github.com/liang21/webook/internal/service"
	"net/http"
	"time"
)

const (
	emailRegexPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passwordRegexPattern = `^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$`
	userIdKey            = "ssid"
)

type UserHandler struct {
	emailRegexPattern    *regexp.Regexp
	passwordRegexPattern *regexp.Regexp
	svc                  service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexPattern:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRegexPattern: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:                  svc,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	server.POST("/users/signup", u.SignUp)
	//server.POST("/users/login", u.Login)
	server.POST("/users/login", u.LoginJWT)
	server.POST("/users/edit", u.Edit)
	server.POST("/users/profile", u.Profile)

}

type name interface {
}

func (u *UserHandler) SignUp(c *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirm_password"`
	}
	var req SignUpReq
	if err := c.Bind(&req); err != nil {
		return
	}
	isEmail, err := u.emailRegexPattern.MatchString(req.Email)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		c.String(http.StatusOK, "邮箱格式错误")
		return
	}
	if req.Password != req.ConfirmPassword {
		c.String(http.StatusOK, "两次密码不一致")
		return
	}
	isPassword, err := u.passwordRegexPattern.MatchString(req.Password)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		c.String(http.StatusOK, "密码格式错误")
		return
	}
	err = u.svc.SignUp(c, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.JSON(200, gin.H{
		"message": "注册成功!",
	})
}

func (u *UserHandler) Login(c *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	err := c.Bind(&req)
	if err != nil {
		return
	}
	user, err := u.svc.Login(c, req.Email, req.Password)
	if err != nil {
		c.String(http.StatusOK, "登录失败")
		return
	}
	sess := sessions.Default(c)
	sess.Set("userId", user.Id)
	sess.Options(sessions.Options{
		// 十五分钟
		MaxAge: 900,
	})
	err = sess.Save()
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.JSON(200, gin.H{
		"message": "登录成功!",
	})
}

func (u *UserHandler) LoginJWT(c *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	err := c.Bind(&req)
	if err != nil {
		return
	}
	user, err := u.svc.Login(c, req.Email, req.Password)
	if err != nil {
		c.String(http.StatusOK, "登录失败")
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Id: user.Id,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
		UserAgent: c.GetHeader("User-Agent"),
	})
	tokenStr, err := token.SignedString(JWTKey)
	if err != nil {
		c.String(http.StatusOK, "系统错误")
		return
	}
	c.Header("x-jwt-token", tokenStr)
	c.JSON(200, gin.H{
		"message": "登录成功!",
	})
}

func (u *UserHandler) Edit(c *gin.Context) {
	type UserReq struct {
		NikeName string `json:"nike_name"`
		Birthday string `json:"birthday"`
		About    string `json:"about"`
	}
	var req UserReq
	if err := c.Bind(&req); err != nil {
		return
	}
	user, ok := c.MustGet("user").(UserClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		c.String(http.StatusOK, "生日格式不对")
		return
	}
	err = u.svc.Edit(c, domain.User{
		Id:       user.Id,
		NikeName: req.NikeName,
		Birthday: birthday,
		About:    req.About,
	})
	if err != nil {
		c.String(http.StatusOK, "系统异常")
		return
	}
	c.JSON(200, gin.H{
		"message": "修改成功",
	})
}

func (u *UserHandler) Profile(c *gin.Context) {
	uc, ok := c.MustGet("user").(UserClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	user, err := u.svc.Profile(c, uc.Id)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	type UserRes struct {
		Id       int64
		Email    string
		NikeName string
		Birthday string
		About    string
	}
	c.JSON(200, UserRes{
		Id:       user.Id,
		Email:    user.Email,
		NikeName: user.NikeName,
		Birthday: user.Birthday.Format(time.DateOnly),
		About:    user.About,
	})
}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

type UserClaims struct {
	Id int64
	jwt.RegisteredClaims

	UserAgent string
}
