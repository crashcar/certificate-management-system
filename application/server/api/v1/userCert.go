package v1

import (
	"application/pkg/app"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// QueryUserCertificate  查询用户自己所在机构的id
func QueryUserCertificate(c *gin.Context) {
	appG := app.Gin{C: c}
	//userID := c.Query("id")

	//查询数据库
	certs := []map[string]interface{}{
		{"id": 1, "name": "证书1", "type": "类型A", "date": "2022-01-01", "issuer": "机构A"},
		{"id": 2, "name": "证书2", "type": "类型B", "date": "2022-02-01", "issuer": "机构B"},
		// 其他证书数据...
	}

	appG.Response(http.StatusOK, "QueryUserCertificate", certs)
}

func UserDownloadCertificate(c *gin.Context) {
	appG := app.Gin{C: c}

	// 设置 PDF 文件路径
	pdfFilePath := "C:\\Users\\Jeremy\\Desktop\\fabric_file\\test.pdf"

	// 打开 PDF 文件
	file, err := os.Open(pdfFilePath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "Failed to open PDF file", nil)
		return
	}
	defer file.Close()

	// 设置响应头，告诉客户端返回的是 PDF 文件
	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "attachment; filename=test.pdf")

	// 将 PDF 文件内容作为响应体发送给客户端
	if _, err := io.Copy(c.Writer, file); err != nil {
		log.Println("Failed to send PDF file")
		appG.Response(http.StatusInternalServerError, "Failed to send PDF file", nil)
		return
	}

	log.Println("Successfully sent PDF file")
	appG.Response(http.StatusOK, "PDF file sent successfully", nil)
}

func UserApplyCertificate(c *gin.Context) {
	appG := app.Gin{C: c}

	// 获取证书颁发机构编号和证书编号
	institutionId := c.PostForm("institutionId")
	certificateId := c.PostForm("certificateId")
	// 打印证书机构和证书编号
	log.Println("Institution ID:", institutionId)
	log.Println("Certificate ID:", certificateId)

	// 保存文件的文件夹路径
	savePath := "C:\\Users\\Jeremy\\Desktop"

	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		appG.Response(http.StatusBadRequest, "UserApplyCertificate", err.Error())
		return
	}

	// 保存文件到指定路径
	filePath := filepath.Join(savePath, file.Filename)
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		appG.Response(http.StatusInternalServerError, "UserApplyCertificate", err.Error())
		return
	}

	// 返回成功响应
	appG.Response(http.StatusOK, "UserApplyCertificate", "ok")
}
