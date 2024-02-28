package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	bc "application/blockchain"
	"application/pkg/app"
)

/**** 管理系统查询部分 ****/

//根据证书ID查询证书
type queryCertByUserSysRequestBody struct {
	CertID string `json:"certID"`
}

type queryCertByFullInfoSysRequestBody struct {
	CertID  string `json:"certID"`
	HoderID string `json:"holderID"`
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
	function := ""
	if body.CertID == "" {
		function = "queryCertByInfosLists"
	} else {
		function = "queryCertByInfos"
	}
	var bodyBytes [][]byte
	if body.CertID != "" {
		bodyBytes = append(bodyBytes, []byte(body.CertID))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery(function, bodyBytes)
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

	// ipfsnode := "certman-ipfs:5001"
	// buffer, err := ipfs.GetFileFromIPFS(appG, body.CID, ipfsnode)
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, "从ipfs获取文件失败", err)
	// 	return
	// }
	// file, err := os.Create("./uploads/certificates/testfile.pdf")
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, "新建pdf失败", err)
	// }
	// defer file.Close()

	// // 写入数据
	// _, err = file.Write(buffer)
	// if err != nil {
	// 	appG.Response(http.StatusInternalServerError, "写入文件时出错", err)
	// 	return
	// }
	// appG.Response(http.StatusOK, "成功", "从ipfs读取文件成功")
}

// 管理系统：通过证书持有人的id和证书id查询该证书，返回{证书ID，持有人ID，证书颁发机构，有效期，状态，修改情况}
func QueryCertByFullInfoSys(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(queryCertByFullInfoSysRequestBody)
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "参数解析失败", err)
		return
	}
	function := ""
	if body.CertID == "" && body.HoderID == "" {
		function = "queryCertByInfosLists"
	} else if body.CertID == "" && body.HoderID != "" {
		appG.Response(http.StatusBadRequest, "CertId为空，HoderId不为空", "参数有误")
		return
	} else {
		function = "queryCertByInfos"
	}

	var bodyBytes [][]byte
	if body.CertID != "" {
		bodyBytes = append(bodyBytes, []byte(body.CertID))
	}

	if body.HoderID != "" {
		bodyBytes = append(bodyBytes, []byte(body.HoderID))
	}

	//调用智能合约
	resp, err := bc.ChannelQuery(function, bodyBytes)
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

/**** 证书机构查询部分 ****/

type queryCertByUserOrgRequestBody struct {
	IssuingAuthority string `json:"issuingAuthority"`
	HoderID          string `json:"holder"`
}

// org：-用户- 通过证书持有人的id查询该人在本机构的所有证书
func QueryCertByUserOrg(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(queryCertByUserOrgRequestBody)
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte

	if body.IssuingAuthority != "" {
		bodyBytes = append(bodyBytes, []byte(body.IssuingAuthority))
	}
	if body.HoderID != "" {
		bodyBytes = append(bodyBytes, []byte(body.HoderID))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryCertByAuthority", bodyBytes)
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

type queryCertOrgRequestBody struct {
	IssuingAuthority string `json:"issuingAuthority"`
}

// org：-管理员- 查询该机构的所有证书以及其持有人
func QueryCertOrg(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(queryCertOrgRequestBody)
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.IssuingAuthority != "" {
		bodyBytes = append(bodyBytes, []byte(body.IssuingAuthority))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryCertByAuthority", bodyBytes)
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
