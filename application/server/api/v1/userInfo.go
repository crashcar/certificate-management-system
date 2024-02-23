package v1

import (
	"application/pkg/app"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

// TAG 日志定位标签
const TAG = "UserInfo: "

// LoginRequest 定义了登录请求的结构体
type LoginRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

// RegistrationRequest 定义了注册请求的结构体
type RegistrationRequest struct {
	ID       string `json:"id"`       // 用户ID
	Password string `json:"password"` // 密码
	Phone    string `json:"phone"`    // 手机号码
	Email    string `json:"email"`    // 邮箱地址
}

// UserLogin 注册
func UserLogin(c *gin.Context) {
	appG := app.Gin{C: c}
	var loginReq LoginRequest

	// 尝试解析请求体到loginReq结构体
	if err := c.ShouldBindJSON(&loginReq); err != nil {
		// 解析失败，返回错误信息
		log.Println(TAG + "解析失败，返回错误信息")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "用户名或密码格式输入错误", // 移除了 loginReq.ID
		})
		return // 添加了 return 语句
	}

	// 解析成功，打印ID和Password
	log.Println(TAG + "ID: " + loginReq.ID + ", Password: " + loginReq.Password)

	// TODO: 查询身份信息数据库 比对数据
	dbUserID := "1"
	dbUserPassword := "1"

	if loginReq.ID == dbUserID && loginReq.Password == dbUserPassword {
		// login 成功，返回成功消息
		log.Println(TAG + "登录成功")
		appG.Response(http.StatusOK, "登录成功", "Login_Success")
		return
	} else {
		// login 失败，返回失败消息
		log.Println(TAG + "登录失败")
		appG.Response(http.StatusBadRequest, "登录失败", "Login_Fail")
		return
	}
}

// UserRegister 处理用户注册请求
func UserRegister(c *gin.Context) {
	appG := app.Gin{C: c}
	var regReq RegistrationRequest

	// 尝试将请求体解析到RegistrationRequest结构体中
	if err := c.ShouldBindJSON(&regReq); err != nil {
		// 如果解析失败，返回错误信息
		log.Println(TAG + " 解析请求体失败")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "输入数据格式错误",
		})
		return // 发生错误时终止函数执行
	}

	log.Println(TAG + "ID: " + regReq.ID + ", Password: " + regReq.Password + ", Phone: " + regReq.Phone + ", Email: " + regReq.Email)

	// TODO: 此处应添加验证逻辑，例如验证ID是否已存在
	var flag bool
	if regReq.ID == "1" {
		flag = true
	} else {
		flag = false
	}

	if flag {
		// 假设验证通过并且用户成功注册，返回成功消息
		log.Println(TAG + "注册成功")
		appG.Response(http.StatusOK, "注册成功", "Register_Success")
		return
	} else {
		// 假设验证失败，返回失败消息
		log.Println(TAG + "注册失败")
		appG.Response(http.StatusBadRequest, "注册失败", "Register_Fail")
		return
	}

}
