package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	bc "application/blockchain"
	"application/pkg/app"
	"application/pkg/cryptoutils"
	"application/pkg/ipfs"
)

/**** 管理系统查询部分 ****/

// 根据证书ID查询证书
type queryCertByUserSysRequestBody struct {
	UserID string `json:"userID"`
}

// 此功能暂时不使用，因为不符合现实情境
// 管理系统：查看用户所有证书，通过用户id查询用户在所有机构的所有证书
func QueryCertByUserSys(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(queryCertByUserSysRequestBody)
	err := c.ShouldBind(body)
	if err != nil {
		appG.Response(http.StatusBadRequest, "参数解析失败", err)
		return
	}
	if body.UserID == "" {
		appG.Response(http.StatusBadRequest, "参数有误", "userID为空字符串")
	}

	var bodyBytes [][]byte

	//调用智能合约
	resp, err := bc.ChannelQuery("queryCertByInfosLists", bodyBytes)
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

	processDataAndRespond(data, appG)
}

type queryCertByFullInfoSysRequestBody struct {
	CertID     string `json:"certID"`
	HolderID   string `json:"holderID"`
	HolderName string `json:"holderName"`
}

// 管理系统 - 验证：通过证书持有人的id和证书id查询该证书，返回证书完整信息
func QueryCertByFullInfoSys(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(queryCertByFullInfoSysRequestBody)
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "参数解析失败", err)
		return
	}
	if body.HolderID == "" || body.HolderName == "" || body.CertID == "" {
		appG.Response(http.StatusBadRequest, "参数有误", "certID或holderID或holderName为空字符串")
	}

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.CertID))
	bodyBytes = append(bodyBytes, []byte(body.HolderID))
	bodyBytes = append(bodyBytes, []byte(body.HolderName))

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

	processDataAndRespond(data, appG)
}

/**** 证书机构查询部分 ****/

type queryCertByUserOrgRequestBody struct {
	IssuingAuthority string `json:"issuingAuthority"`
	HolderID         string `json:"holderID"`
}

// org：-用户- 通过证书持有人的id查询该人在本机构的所有证书
func QueryCertByUserOrg(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(queryCertByUserOrgRequestBody)
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.IssuingAuthority == "" || body.HolderID == "" {
		appG.Response(http.StatusBadRequest, "参数有误", "issuingAuthority或holderID为空")
		return
	}

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.IssuingAuthority))
	bodyBytes = append(bodyBytes, []byte(body.HolderID))

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

	processDataAndRespond(data, appG)
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
	if body.IssuingAuthority == "" {
		appG.Response(http.StatusBadRequest, "参数有误", "issuingAuthority为空字符串")
		return
	}

	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.IssuingAuthority))

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

	processDataAndRespond(data, appG)
}

func processDataAndRespond(data []map[string]interface{}, appG app.Gin) {
	for i, item := range data {
		// 获取cid
		var cid string
		if hashPath, ok := item["hashPath"].(string); ok {
			cid = hashPath
		}

		// 获取从ipfs获取文件
		ipfsnode := "certman-ipfs:5001"
		buffer := ipfs.GetFileFromIPFS(appG, cid, ipfsnode)

		// 哈希buffer，返回前端显示当前cid文件的哈希以便对比
		retrievedHash := cryptoutils.HashBuffer(buffer)

		// 将retrievedHash加入到item中
		item["retrievedHash"] = retrievedHash

		// 更新data
		data[i] = item
	}

	appG.Response(http.StatusOK, "链上查询成功", data)
}
