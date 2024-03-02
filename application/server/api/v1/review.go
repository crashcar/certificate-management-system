package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	bc "application/blockchain"
	"application/model"
	"application/pkg/app"
	"application/pkg/cryptoutils"
	"application/pkg/ipfs"
	"application/ws"
)

type reviewApprovedRequestBody struct {
	CertDBID         uint   `json:"certDBID"`
	AdminID          uint   `json:"adminID"`
	IssuingAuthority string `json:"issuingAuthority"` // 测试不同机构
}

type authorityContactInfo struct {
	Phone   string
	Email   string
	Address string
}

var aci = authorityContactInfo{
	Phone:   "0001-0001",
	Email:   "cet@org.com",
	Address: "wuhan, hubei",
}

// 管理员审核接口
// 审核通过
// 1. 加密； 2. 上传ipfs，获取CID； 4. 删除数据库记录和对应文件； 3. CID和metadata上链； 5. 通知用户
func ApproveCert(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		body := new(reviewApprovedRequestBody)
		err := c.ShouldBind(body)
		if err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", err)
			return
		}
		if body.CertDBID == 0 || body.AdminID == 0 || body.IssuingAuthority == "" {
			appG.Response(http.StatusBadRequest, "失败", "certs表项ID为0或adminID为0或issuingAuthority为空字符串")
			return
		}

		// 根据 body.ID 从数据库取出对应的项
		var cert model.Application
		certDbId := body.CertDBID
		adminId := body.AdminID
		result := db.First(&cert, certDbId)
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库查询错误", result.Error.Error())
			return
		}
		holderID := cert.UploaderID
		holderName := cert.UploaderName
		certType := cert.CertType
		currentDate := time.Now().Format("2006-01-02")
		expiryDate := time.Now().AddDate(10, 0, 0).Format("2006-01-02")

		// 获取原文件的hash
		hashString := cryptoutils.HashFile(*cert.Path, appG)

		// 加密上传ipfs
		ipfsnode := "certman-ipfs:5001"
		cid := ipfs.UploadFileToIPFS(appG, *cert.Path, ipfsnode)

		// 删除文件系统的文件
		err = os.Remove(*cert.Path)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "文件系统删除文件错误", err)
			return
		}

		// 获取存在账本中的 certID
		newID, err := uuid.NewUUID()
		if err != nil {
			appG.Response(http.StatusInternalServerError, "uuid生成错误", err)
			return
		}
		CertID := "cet.com-" + newID.String()

		// 更新数据库存储：path置空，存储certID
		result = db.Model(&cert).Where("id = ?", certDbId).Updates(map[string]interface{}{
			"Path":        nil,
			"IsProcessed": true,
			"ProcessedAt": time.Now(),
			"IsApproved":  true,
			"CertID":      &CertID,
		})

		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库表项更新出错", result.Error.Error())
			return
		}

		// 制作上链数据
		issuingAuthority := body.IssuingAuthority // 测试不同机构
		adminIdStr := strconv.FormatUint(uint64(adminId), 10)
		var bodyBytes [][]byte
		bodyBytes = append(bodyBytes, []byte(hashString))
		bodyBytes = append(bodyBytes, []byte(cid))
		bodyBytes = append(bodyBytes, []byte(CertID))
		bodyBytes = append(bodyBytes, []byte(holderID))
		bodyBytes = append(bodyBytes, []byte(holderName))
		bodyBytes = append(bodyBytes, []byte(certType))
		bodyBytes = append(bodyBytes, []byte(adminIdStr))
		bodyBytes = append(bodyBytes, []byte(currentDate))
		bodyBytes = append(bodyBytes, []byte(expiryDate))
		bodyBytes = append(bodyBytes, []byte(issuingAuthority))
		bodyBytes = append(bodyBytes, []byte(aci.Phone))
		bodyBytes = append(bodyBytes, []byte(aci.Email))
		bodyBytes = append(bodyBytes, []byte(aci.Address))

		//调用智能合约数据上链
		resp, err := bc.ChannelExecute("uploadCertOrg", bodyBytes)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "数据上链失败", err.Error())
			return
		}
		var data map[string]interface{}
		if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
			appG.Response(http.StatusInternalServerError, "账本数据unmarshal失败", err.Error())
			return
		}

		// 通知用户
		notification := model.Notification{
			UserID:       holderID,
			AdminID:      adminId,
			Content:      "证书审核通过",
			IsRead:       false,
			CreatedAt:    time.Now(),
			DenialReason: nil,
		}
		ws.WriteDBAndNotifyUser(db, notification, appG)

		appG.Response(http.StatusOK, "成功", data)
	}
}

type reviewDenialRequestBody struct {
	CertDBID     uint   `json:"certDbID"`
	AdminID      uint   `json:"adminID"`
	DenialReason string `json:"denialReason"`
}

// 管理员审核接口
// 审核不通过
// 1. 删除数据库记录和对应文件； 2. 通知用户
func DenialCert(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		body := new(reviewDenialRequestBody)
		err := c.ShouldBind(body)
		if err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", err)
			return
		}
		if body.CertDBID == 0 || body.AdminID == 0 || body.DenialReason == "" {
			appG.Response(http.StatusBadRequest, "certDbId或adminId为0, 或拒绝理由为空", err)
			return
		}

		var cert model.Application
		certDbId := body.CertDBID
		adminId := body.AdminID

		// 数据库表项更新(审核不通过文件系统文件不删除)
		denialReason := body.DenialReason
		result := db.Model(&cert).Where("id = ?", certDbId).Updates(model.Application{
			IsProcessed:  true,
			ProcessedAt:  time.Now(),
			DenialReason: &denialReason,
		})
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库表项更新错误", result.Error.Error())
			return
		}

		// 通知用户
		notification := model.Notification{
			UserID:       cert.UploaderID,
			AdminID:      adminId,
			Content:      "证书审核未通过",
			IsRead:       false,
			CreatedAt:    time.Now(),
			DenialReason: &denialReason,
		}
		ws.WriteDBAndNotifyUser(db, notification, appG)

		appG.Response(http.StatusOK, "成功", "证书审核不通过")
	}
}

// * 还没有业务逻辑的部分 *
type deleteCertOrgRequestBody struct {
	//HoderID string `json:"holderID"`
	CertID string `json:"certID"`
}

// org：-管理员- 删除证书（可能不放在review.go中，目前还没有业务逻辑来从账本中删除证书）
func DeleteCertOrg(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(deleteCertOrgRequestBody)
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
