package v1

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"application/pkg/app"
	"application/pkg/ipfs"
)

/**** 管理系统后端代码 ****/

type queryCertByUserSysRequestBody struct {
	CID string `json:"cid"`
}

type QueryCertByFullInfoSysRequestBody struct {
	HoderID string `json:"holderID"`
	CertID  string `json:"certID"`
}

// 管理系统：查看用户所有证书，通过用户id查询用户在所有机构的所有证书
func QueryCertByUserSys(c *gin.Context) {
	appG := app.Gin{C: c}

	body := new(queryCertByUserSysRequestBody)
	err := c.ShouldBind(body)
	if err != nil {
		appG.Response(http.StatusBadRequest, "参数解析失败", err)
		return
	}
	if body.CID == "" {
		appG.Response(http.StatusBadRequest, "参数不完整", err)
		return
	}

	ipfsnode := "certman-ipfs:5001"
	buffer, err := ipfs.GetFileFromIPFS(appG, body.CID, ipfsnode)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "从ipfs获取文件失败", err)
		return
	}
	file, err := os.Create("./uploads/certificates/testfile.pdf")
	if err != nil {
		appG.Response(http.StatusInternalServerError, "新建pdf失败", err)
	}
	defer file.Close()

	// 写入数据
	_, err = file.Write(buffer)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "写入文件时出错", err)
		return
	}
	appG.Response(http.StatusOK, "成功", "从ipfs读取文件成功")
}

// 管理系统：通过证书持有人的id和证书id查询该证书，返回{证书ID，持有人ID，证书颁发机构，有效期，状态，修改情况}
func QueryCertByFullInfoSys(c *gin.Context) {

}
