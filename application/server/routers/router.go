package routers

import (
	v1 "application/api/v1"
	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// InitRouter 初始化路由信息
func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	apiV1 := r.Group("/api/v1")
	{
		apiV1.GET("/hello", v1.Hello)
		apiV1.POST("/queryAccountList", v1.QueryAccountList)
		apiV1.POST("/createRealEstate", v1.CreateRealEstate)
		apiV1.POST("/queryRealEstateList", v1.QueryRealEstateList)
		apiV1.POST("/createSelling", v1.CreateSelling)
		apiV1.POST("/createSellingByBuy", v1.CreateSellingByBuy)
		apiV1.POST("/querySellingList", v1.QuerySellingList)
		apiV1.POST("/querySellingListByBuyer", v1.QuerySellingListByBuyer)
		apiV1.POST("/updateSelling", v1.UpdateSelling)
		apiV1.POST("/createDonating", v1.CreateDonating)
		apiV1.POST("/queryDonatingList", v1.QueryDonatingList)
		apiV1.POST("/queryDonatingListByGrantee", v1.QueryDonatingListByGrantee)
		apiV1.POST("/updateDonating", v1.UpdateDonating)
		apiV1.POST("/register", v1.Register(db))
		apiV1.POST("/login", v1.Login(db))
		apiV1.GET("/reviewTypes", v1.GetReviewTypes)
		apiV1.POST("/adminRegister", v1.AdminRegister(db))
		apiV1.POST("/adminLogin", v1.AdminLogin(db))
		apiV1.POST("/upload", v1.SaveFile(db))
		apiV1.POST("/showCertList", v1.ShowCertList(db))
		apiV1.POST("/showProcessedCert", v1.ShowProcessedCert(db))
		apiV1.POST("/approveCert", v1.ApproveCert(db))
		apiV1.POST("/testIPFS", v1.QueryCertByUserSys)
	}
	// 静态文件路由
	r.StaticFS("/web", http.Dir("./dist/"))
	r.StaticFS("/uploads", http.Dir("./uploads"))
	return r
}
