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

// UserRegister  注册
// 输入参数： id, realname, password, email
func UserRegister(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		body := new(RegisterRequestBody)

		if err := c.ShouldBind(body); err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", fmt.Sprintf("参数出错%s", err.Error()))
			return
		}
		if body.ID == "" || body.RealName == "" || body.Password == "" || body.Email == "" {
			appG.Response(http.StatusBadRequest, "Register_info_not_completed", "请填写完整信息(用户名、姓名、密码、Email)")
			return
		}

		// 检查用户是否已注册
		var existingUser model.User
		if err := db.Where("id = ?", body.ID).First(&existingUser).Error; err == nil {
			appG.Response(http.StatusOK, "Register_already", "用户已注册，请直接登录")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("密码加密%s", err.Error()))
			return
		}

		// 创建新用户
		user := model.User{
			ID:       body.ID,
			RealName: body.RealName,
			Password: string(hashedPassword),
			Email:    body.Email,
		}

		if result := db.Create(&user); result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库操作：创建账户失败", result.Error.Error())
			return
		}
		appG.Response(http.StatusOK, "Register_Success", "用户注册成功")
	}
}

// UserLogin  登录
func UserLogin(db *gorm.DB) gin.HandlerFunc {
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
			appG.Response(http.StatusOK, "Login_Failed", "当前用户不存在, 请重新注册")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
			appG.Response(http.StatusOK, "Login_Failed", "用户密码输入错误")
			return
		}

		appG.Response(http.StatusOK, "Login_Success", "user")
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

// AdminRegister  管理员注册
// 管理员注册只输入密码和审查类型，管理员ID由系统（数据库自增主键）分配
// 返回系统分配的管理员ID
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
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建账户%s", result.Error.Error()))
			return
		}

		appG.Response(http.StatusOK, "Admin_Register_Success", admin.ID)
	}
}

type adminLoginRequestBody struct {
	ID       uint   `json:"id"`
	Password string `json:"password"`
}

// AdminLogin  管理员登录
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
			appG.Response(http.StatusInternalServerError, "Login_Failed", "管理员账户不存在")
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.Password)); err != nil {
			appG.Response(http.StatusInternalServerError, "Login_Failed", "管理员密码输入错误")
			return
		}

		appG.Response(http.StatusOK, "Login_Success", "admin")
	}
}
