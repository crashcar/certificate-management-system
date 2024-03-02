package v1

import (
	"application/model"
	"application/pkg/app"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type showApplicationListRequestBody struct {
	QueryType string `json:"queryType"`
	UserID    string `json:"userID"`
	AdminID   uint   `json:"adminID"`
}

type applicationDisplay struct {
	ApplicationID uint      `json:"applicationID"`
	UploaderID    string    `json:"uploaderID"`
	UploaderName  string    `json:"uploaderName"`
	CreatedAt     time.Time `json:"createdAt"`
	CertType      string    `json:"certType"`
	IsProcessed   bool      `json:"isProcessed"`
	IsApproved    bool      `json:"isApproved"`
}

// 获取证书列表
func ShowApplicationList(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		body := new(showApplicationListRequestBody)
		err := c.ShouldBind(body)
		if err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", err)
			return
		}

		var applications []model.Application

		// 用户查询已申请列表
		if body.QueryType == "user" {
			if body.UserID == "" {
				appG.Response(http.StatusBadRequest, "参数有误", "userID为空字符串")
			} else {
				result := db.Where("uploader_id = ?", body.UserID).Order("is_processed").Find(&applications)
				if result.Error != nil {
					appG.Response(http.StatusInternalServerError, "数据库查询错误", result.Error.Error())
					return
				}
			}
			// 管理员查询对应审查类型表项
		} else if body.QueryType == "admin" {
			if body.AdminID == 0 {
				appG.Response(http.StatusBadRequest, "参数有误", "adminID为0")
			} else {
				// 查询adminID对应的reviewType
				var admin model.Admin
				result := db.First(&admin, body.AdminID)
				if result.Error != nil {
					appG.Response(http.StatusInternalServerError, "数据库查询错误", result.Error.Error())
					return
				}
				reviewType := admin.ReviewType

				// 从数据库查询certtype为body.type类型的项
				result = db.Where("cert_type = ?", reviewType).Order("created_at asc").Find(&applications)
				if result.Error != nil {
					appG.Response(http.StatusInternalServerError, "数据库查询错误", result.Error.Error())
					return
				}
			}
		} else {
			appG.Response(http.StatusBadRequest, "参数有误", "不支持的queryType")
			return
		}

		// 生成用于前端展示的信息
		var displayData []applicationDisplay
		for _, cert := range applications {
			displayData = append(displayData, applicationDisplay{
				ApplicationID: cert.ID,           // user & admin
				UploaderID:    cert.UploaderID,   // admin
				UploaderName:  cert.UploaderName, // admin
				CreatedAt:     cert.CreatedAt,    // user & admin
				CertType:      cert.CertType,     // user
				IsProcessed:   cert.IsProcessed,  // user
				IsApproved:    cert.IsApproved,   // user
			})
		}

		appG.Response(http.StatusOK, "成功", displayData)
	}
}

type deleteRecordRequestBody struct {
	ApplicationID uint `json:"applicationId"`
}

// 用户撤销未审核的证书 和 删除记录
func DeleteRecord(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		body := new(deleteRecordRequestBody)
		err := c.ShouldBind(body)
		if err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", err)
			return
		}
		if body.ApplicationID == 0 {
			appG.Response(http.StatusBadRequest, "参数有误", "applicationId为0")
			return
		}

		// 删除数据库记录
		result := db.Delete(&model.Application{}, body.ApplicationID)
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库删除表项错误", result.Error.Error())
			return
		}

		appG.Response(http.StatusOK, "成功", "删除记录成功")
	}
}

type applicationDetailRequestBody struct {
	ApplicationID uint `json:"applicationID"`
}

type applicationDetailDisplay struct {
	UploaderID   string    `json:"uploaderID"`
	UploaderName string    `json:"uploaderName"`
	ProcessedAt  time.Time `json:"processedAt"`
	ImageURL     string    `json:"imageURL"`
	CertID       string    `json:"certID"`
	DenialReason string    `json:"denialReason"`
}

// 查看申请详情
func ApplicationDetail(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		body := new(applicationDetailRequestBody)
		err := c.ShouldBind(body)
		if err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", err)
			return
		}

		// 从数据库查询
		var application model.Application
		result := db.First(&application, body.ApplicationID)
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库查询错误", result.Error.Error())
			return
		}

		// 生成用于前端展示的信息
		certID := ""
		url := ""
		denialReason := ""
		if application.CertID != nil {
			certID = *application.CertID
		}
		if application.Path != nil {
			url = "http://localhost:8000/" + *application.Path
		}
		if application.DenialReason != nil {
			denialReason = *application.DenialReason
		}
		displayData := applicationDetailDisplay{
			UploaderID:   application.UploaderID,   // admin
			UploaderName: application.UploaderName, // admin
			ProcessedAt:  application.ProcessedAt,  // user-processed-approved, user-processed-denialed
			ImageURL:     url,                      // user-notProcessed, admin
			CertID:       certID,                   // user-processed-approved
			DenialReason: denialReason,             // user-processed-denialed
		}

		appG.Response(http.StatusOK, "成功", displayData)
	}

}
