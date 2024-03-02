package v1

import (
	"application/model"
	"application/pkg/app"
	"application/pkg/ipfs"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

type uploadRequestBody struct {
	UserID   string `form:"userID"`
	RealName string `form:"realName"`
	CertType string `form:"certType"`
}

// 用户上传接口
// 1. 将文件保存在server的文件系统
// 2. path和其他信息写入数据库
func SaveFile(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}
		file, err := c.FormFile("file")
		if err != nil {
			appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("表单错误：%s", err.Error()))
			return
		}

		body := new(uploadRequestBody)
		if err := c.ShouldBindWith(body, binding.Form); err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", err)
			return
		}
		if body.UserID == "" || body.RealName == "" || body.CertType == "" {
			appG.Response(http.StatusBadRequest, "参数不完整", err)
			return
		}

		// 将文件保存到server
		now := time.Now()
		formattedTime := now.Format("20060102150405")
		fileExt := filepath.Ext(file.Filename)
		if fileExt != ".png" {
			appG.Response(http.StatusBadRequest, "上传文件类型错误", "请上传PNG文件")
			return
		}
		newFileName := fmt.Sprintf("%s-%s%s", body.UserID, formattedTime, fileExt)
		directory := "./uploads/certificates"
		dst := filepath.Join(directory, newFileName)
		if _, err := os.Stat(directory); os.IsNotExist(err) {
			// 如果目录不存在，则创建目录
			err := os.MkdirAll(directory, 0777)
			if err != nil {
				// 处理创建目录的错误
				appG.Response(http.StatusInternalServerError, "创建目录失败", err)
				return
			}
		}
		if err := c.SaveUploadedFile(file, dst); err != nil {
			appG.Response(http.StatusInternalServerError, "文件保存到服务器失败", err)
			return
		}

		// path和其他信息写入数据库
		application := model.Application{
			Path:         &dst,
			CertType:     body.CertType,
			UploaderID:   body.UserID,
			UploaderName: body.RealName,
			CertID:       nil,
		}

		if result := db.Create(&application); result.Error != nil {
			appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("数据库写入失败：%s", result.Error.Error()))
			return
		}

		// 返回申请编号
		appG.Response(http.StatusOK, "成功", application.ID)
	}
}

type downloadRequestBody struct {
	CID string `json:"cid"`
}

func DownloadCertificate(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(downloadRequestBody)
	err := c.ShouldBind(body)
	if err != nil {
		appG.Response(http.StatusBadRequest, "参数解析失败", err)
		return
	}
	if body.CID == "" {
		appG.Response(http.StatusBadRequest, "失败", "cid为空字符串")
		return
	}

	ipfsnode := "certman-ipfs:5001"
	certificateBytes := ipfs.GetFileFromIPFS(appG, body.CID, ipfsnode)
	fileName := body.CID + ".png"

	// 设置响应头以指示内容类型和下载模式
	c.Writer.Header().Set("Content-Type", "application/octet-stream")
	c.Writer.Header().Set("Content-Disposition", "attachment; filename=\""+fileName+"\"")

	// 写入响应体
	if _, err := c.Writer.Write(certificateBytes); err != nil {
		// 如果无法写入响应体，处理错误
		appG.Response(http.StatusInternalServerError, "写入响应体失败", err)
		return
	}
}
