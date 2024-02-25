package v1

import (
	bc "application/blockchain"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**** 管理系统后端代码 ****/

// type QueryCertByUserSysRequestBody struct {
// 	HoderID string `json:"holder"`
// }

// type QueryCertByFullInfoSysRequestBody struct {
// 	HoderID string `json:"holderID"`
// 	CertID  string `json:"certID"`
// }

// // 管理系统：查看用户所有证书，通过用户id查询用户在所有机构的所有证书
// func QueryCertByUserSys(c *gin.Context) {

// }
// 管理系统：查看用户所有证书，通过用户id查询用户在所有机构的所有证书
func QueryCertByUserSys(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(QueryCertByUserSysRequestBody)
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.HoderID != "" {
		bodyBytes = append(bodyBytes, []byte(body.HoderID))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryCertByUserSys", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

// 管理系统：通过证书持有人的id和证书id查询该证书，返回{证书ID，持有人ID，证书颁发机构，有效期，状态，修改情况}
func QueryCertByFullInfoSys(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(QueryCertByFullInfoSysRequestBody)
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.HoderID != "" {
		bodyBytes = append(bodyBytes, []byte(body.HoderID))
	}
	if body.CertID != "" {
		bodyBytes = append(bodyBytes, []byte(body.CertID))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryCertByFullInfoSys", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

/**** 证书机构后端代码 ****/

type QueryCertByUserOrgRequestBody struct {
	HoderID          string `json:"holder"`
	IssuingAuthority string `json:"issuingAuthority"`
}

type QueryCertOrgRequestBody struct {
	IssuingAuthority string `json:"issuingAuthority"`
}

type AuthorityContactInfo struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type UploadCertOrgRequestBody struct {
	// hash
	HashFile string `json:"hashFile"`
	HashPath string `json:"hashPath"`
	// metadata
	CertID               string               `json:"certID"`
	HoderID              string               `json:"hoderID"`
	HoderName            string               `json:"hoderName"`
	CertType             string               `json:"certType"`
	IssueDate            string               `json:"issueDate"`
	ExpiryDate           string               `json:"expiryDate"`
	IssuingAuthority     string               `json:"issuingAuthority"`
	AuthorityContactInfo AuthorityContactInfo `json:"authorityContactInfo"`
	//Status               string               `json:"status"`
}

type DeleteCertOrgRequestBody struct {
	//HoderID string `json:"holderID"`
	CertID string `json:"certID"`
}

// org：-用户- 通过证书持有人的id查询该人在本机构的所有证书
func QueryCertByUserOrg(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(QueryCertByUserOrgRequestBody)
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.HoderID != "" {
		bodyBytes = append(bodyBytes, []byte(body.HoderID))
	}
	if body.IssuingAuthority != "" {
		bodyBytes = append(bodyBytes, []byte(body.IssuingAuthority))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryCertByInfos", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

// org：-管理员- 查询该机构的所有证书以及其持有人
func QueryCertOrg(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(QueryCertOrgRequestBody)
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.IssuingAuthority != "" {
		bodyBytes = append(bodyBytes, []byte(body.IssuingAuthority))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryCertByInfos", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

// org：-管理员- 上传证书
func UploadCertOrg(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(UploadCertOrgRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.HashFile == "" || body.HashPath == "" {
		appG.Response(http.StatusBadRequest, "失败", "HashFile或HashPath不能为空")
		return
	}
	if body.CertID == "" || body.HoderID == "" || body.HoderName == "" || body.IssuingAuthority == "" {
		appG.Response(http.StatusBadRequest, "失败", "主键不能为空")
		return
	}

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.HashFile))
	bodyBytes = append(bodyBytes, []byte(body.HashPath))
	bodyBytes = append(bodyBytes, []byte(body.CertID))
	bodyBytes = append(bodyBytes, []byte(body.HoderID))
	bodyBytes = append(bodyBytes, []byte(body.HoderName))
	bodyBytes = append(bodyBytes, []byte(body.CertType))
	bodyBytes = append(bodyBytes, []byte(body.IssueDate))
	bodyBytes = append(bodyBytes, []byte(body.ExpiryDate))
	bodyBytes = append(bodyBytes, []byte(body.AuthorityContactInfo.Phone))
	bodyBytes = append(bodyBytes, []byte(body.AuthorityContactInfo.Email))
	bodyBytes = append(bodyBytes, []byte(body.AuthorityContactInfo.Address))
	//调用智能合约
	resp, err := bc.ChannelExecute("uploadCertOrg", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

// org：-管理员- 删除证书(或许不需要？model/cert结构的字段status，用invalid、expired表示证书已无效)
func DeleteCertOrg(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(DeleteCertOrgRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.CertID != "" {
		bodyBytes = append(bodyBytes, []byte(body.CertID))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("deleteCertOrg", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}
