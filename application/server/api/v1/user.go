package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"application/model"
	"application/pkg/app"
)

type RegisterRequestBody struct {
	ID       string `json:"id"`
	RealName string `json:"realname"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginRequestBody struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

func Register(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		body := new(RegisterRequestBody)

		if err := c.ShouldBind(body); err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", fmt.Sprintf("参数出错%s", err.Error()))
			return
		}
		if body.ID == "" || body.RealName == "" || body.Password == "" || body.Email == "" {
			appG.Response(http.StatusBadRequest, "参数解析失败", "请填写用户名、密码、Email")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("密码加密%s", err.Error()))
			return
		}

		user := model.User{
			ID:       body.ID,
			RealName: body.RealName,
			Password: string(hashedPassword),
			Email:    body.Email,
		}

		if result := db.Create(&user); result.Error != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建账户%s", err.Error()))
			return
		}

		appG.Response(http.StatusOK, "成功", "用户注册成功")
	}
}

func Login(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		body := new(LoginRequestBody)

		if err := c.ShouldBind(body); err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", fmt.Sprintf("参数出错%s", err.Error()))
			return
		}
		if body.ID == "" || body.Password == "" {
			appG.Response(http.StatusBadRequest, "参数解析失败", "请输入用户名和密码")
			return
		}

		var user model.User
		if err := db.Where("id = ?", body.ID).First(&user).Error; err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("用户不存在%s", err.Error()))
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("密码错误%s", err.Error()))
			return
		}

		appG.Response(http.StatusOK, "成功", "用户登录成功")
	}
}

func GetReviewTypes(c *gin.Context) {
	appG := app.Gin{C: c}
	adminTypes := []string{string(model.CET), string(model.CJT), string(model.PHD)}
	appG.Response(http.StatusOK, "成功", adminTypes)
}

type adminRegisterRequestBody struct {
	Password   string `json:"password"`
	ReviewType string `json:"reviewType"`
}

func AdminRegister(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		body := new(adminRegisterRequestBody)

		if err := c.ShouldBind(body); err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", fmt.Sprintf("参数出错%s", err.Error()))
			return
		}
		if body.Password == "" || body.ReviewType == "" {
			appG.Response(http.StatusBadRequest, "参数解析失败", "请填写密码和类型")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("密码加密%s", err.Error()))
			return
		}

		admin := model.Admin{
			Password:   string(hashedPassword),
			ReviewType: model.ReviewType(body.ReviewType),
		}

		if result := db.Create(&admin); result.Error != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建账户%s", err.Error()))
			return
		}

		appG.Response(http.StatusOK, "管理员注册成功", admin.ID)
	}
}

type adminLoginRequestBody struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

func AdminLogin(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		body := new(adminLoginRequestBody)

		if err := c.ShouldBind(body); err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", fmt.Sprintf("参数出错%s", err.Error()))
			return
		}
		if body.ID == 0 || body.Password == "" {
			appG.Response(http.StatusBadRequest, "参数解析失败", "请输入用户名和密码")
			return
		}

		var admin model.Admin
		if err := db.Where("id = ?", body.ID).First(&admin).Error; err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("用户不存在%s", err.Error()))
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.Password)); err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("密码错误%s", err.Error()))
			return
		}

		appG.Response(http.StatusOK, "成功", "管理员登录成功")
	}
}
