package routers

import (
	v1 "application/api/v1"
	ws "application/ws"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由信息
func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{

		AllowOrigins:     []string{"*"}, // 允许的域名列表
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允许带凭证的请求，如Cookies
		MaxAge:           12 * time.Hour,
	}))

	apiV1 := r.Group("/api/v1")
	{
		// user - 注册登录接口
		// 1. 用户登录注册
		apiV1.POST("/register", v1.UserRegister(db))
		apiV1.POST("/login", v1.UserLogin(db))
		// 2. 管理员登录注册
		apiV1.POST("/adminRegister", v1.AdminRegister(db))
		apiV1.POST("/adminLogin", v1.AdminLogin(db))
		// 3. 获取管理员的审查类型/该机构的证书类型（管理员注册下拉菜单/上传证书类型下拉菜单）
		apiV1.GET("/reviewTypes", v1.GetReviewTypes)

		// file - 用户证书文件类型接口
		// 1. 用户上传证书（需要确定上传证书类型GETreviewTypes）
		apiV1.POST("/upload", v1.SaveFile(db))
		// 2. 用户下载证书文件（用户查看审查通过证书列表单击表项）
		apiV1.POST("/downloadCertificate", v1.DownloadCertificate)

		// review - 管理员审查接口
		// 1. 获取所有待审查证书列表
		apiV1.POST("/showCertList", v1.ShowCertList(db))
		// 2. 获取某个证书详情含图片（单击列表中的证书表项）
		apiV1.POST("/showProcessedCert", v1.ShowProcessedCert(db))
		// 3. 审查通过证书
		apiV1.POST("/approveCert", v1.ApproveCert(db))
		// 4. 审查不通过证书
		apiV1.POST("/denialCert", v1.DenialCert(db))

		// query - 查看链上证书接口
		// 1. 查看用户在所有机构的所有证书（没有用上）
		apiV1.POST("/queryCertByUserSys", v1.QueryCertByUserSys)
		// 2. 证书验证 通过证书id、userID、userName验证查询证书，验证真实性
		apiV1.POST("/queryCertByFullInfoSys", v1.QueryCertByFullInfoSys)
		// 3. 查看用户在某一机构的所有证书
		apiV1.POST("/queryCertByUserOrg", v1.QueryCertByUserOrg)
		// 4. 查看某机构的所有证书（管理员调用）
		apiV1.POST("/queryCertOrg", v1.QueryCertOrg)

		// motification - 通知接口
		// 1. 用户已读通知，用户已读通知列表中的通知，数据库中表项改为已读
		apiV1.POST("/readNotification", v1.ReadNotification(db))
	}

	// 静态文件路由
	r.StaticFS("/web", http.Dir("./dist/"))
	r.StaticFS("/uploads", http.Dir("./uploads"))

	// websocket路由
	r.GET("/ws", ws.HandleConnections(db))

	return r
}
