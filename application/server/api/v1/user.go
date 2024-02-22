package v1

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"application/model"
	"application/pkg/app"
)

const createUsersTableSQL = `CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    real_name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role INTEGER DEFAULT 0
);`

const (
	maxRetries    = 5                // 最大重试次数
	retryInterval = 10 * time.Second // 重试间隔
)

func connectToDB() (*gorm.DB, error) {
	dsn := "host=certman-db user=admin password=1234 dbname=userdb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var db *gorm.DB
	var err error
	for i := 0; i < maxRetries; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			return db, nil
		}
		log.Printf("Failed to connect to database. Retry %d times: %v", i+1, err)
		time.Sleep(retryInterval)
	}
	return nil, err
}

func InitUserdb() error {
	db, err := connectToDB()
	if err != nil {
		log.Printf("Failed to connect to pg: %s", err)
		return err
	}

	if err := db.Exec(createUsersTableSQL).Error; err != nil {
		return err
	}

	return nil
}

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

func Register(c *gin.Context) {
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

	db, err := connectToDB()
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("数据库连接%s", err.Error()))
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

func Login(c *gin.Context) {
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

	db, err := connectToDB()
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("数据库连接%s", err.Error()))
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
