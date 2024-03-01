package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
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

type showCertListRequestBody struct {
	AdminID uint `json:"adminID"`
}

type CertListDisplay struct {
	ID           uint      `json:"id"`
	UploaderID   string    `json:"uploaderId"`
	UploaderName string    `json:"uploaderName"`
	CreatedAt    time.Time `json:"createdAt"`
}

// 管理员审核接口
// 显示未处理证书列表
func ShowCertList(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		body := new(showCertListRequestBody)
		err := c.ShouldBind(body)
		if err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", err)
			return
		}
		log.Print(body)
		if body.AdminID == 0 {
			appG.Response(http.StatusBadRequest, "参数有误", "adminID为0")
			return
		}

		// 查询adminID对应的reviewType
		var admin model.Admin
		result := db.First(&admin, body.AdminID)
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库查询错误", result.Error.Error())
			return
		}
		reviewType := admin.ReviewType

		// 从数据库查询certtype为body.type类型的项
		var certs []model.Cert
		if err := db.Where("cert_type = ?", reviewType).Order("created_at asc").Find(&certs).Error; err != nil {
			appG.Response(http.StatusInternalServerError, "数据库查询错误", result.Error.Error())
			return
		}

		// 生成用于前端展示的信息
		var displayData []CertListDisplay
		for _, cert := range certs {
			displayData = append(displayData, CertListDisplay{
				ID:           cert.ID,
				UploaderID:   cert.UploaderID,
				UploaderName: cert.UploaderName,
				CreatedAt:    cert.CreatedAt,
			})
		}

		appG.Response(http.StatusOK, "成功", displayData)
	}
}

type showCertRequestBody struct {
	ID uint `json:"id"`
}

type CertDisplay struct {
	UploaderID   string    `json:"uploaderId"`
	UploaderName string    `json:"uploaderName"`
	CreatedAt    time.Time `json:"createdAt"`
	ImageURL     string    `json:"imageURL"`
}

// 管理员审核接口
// 点击列表中的一项，显示证书文件
func ShowProcessedCert(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		appG := app.Gin{C: c}

		body := new(showCertRequestBody)
		err := c.ShouldBind(body)
		if err != nil {
			appG.Response(http.StatusBadRequest, "参数解析失败", err)
			return
		}
		if body.ID == 0 {
			appG.Response(http.StatusBadRequest, "参数错误", "db证书ID有误")
			return
		}

		// 从数据库查询
		var cert model.Cert
		result := db.First(&cert, body.ID)
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库查询错误", result.Error.Error())
			return
		}

		// 生成用于前端展示的信息
		url := "http://localhost:8000/" + cert.Path
		displayData := CertDisplay{
			UploaderID:   cert.UploaderID,
			UploaderName: cert.UploaderName,
			CreatedAt:    cert.CreatedAt,
			ImageURL:     url,
		}

		appG.Response(http.StatusOK, "成功", displayData)
	}
}

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
		var cert model.Cert
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
		hashString := cryptoutils.HashFile(cert.Path, appG)

		// 加密上传ipfs
		ipfsnode := "certman-ipfs:5001"
		cid := ipfs.UploadFileToIPFS(appG, cert.Path, ipfsnode)

		// 删除文件系统的文件以及数据库记录
		err = os.Remove(cert.Path)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "文件系统删除文件错误", err)
			return
		}
		result = db.Delete(&cert, certDbId)
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库删除记录错误", err)
			return
		}

		// 获取存在账本中的 certID
		newID, err := uuid.NewUUID()
		if err != nil {
			appG.Response(http.StatusInternalServerError, "uuid生成错误", err)
			return
		}
		CertID := "cet.com-" + newID.String()

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

		// 根据 body.ID 从数据库取出对应的项
		var cert model.Cert
		certDbId := body.CertDBID
		adminId := body.AdminID
		result := db.First(&cert, certDbId)
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库查询错误", result.Error.Error())
			return
		}

		// 删除文件系统的文件以及数据库记录
		err = os.Remove(cert.Path)
		if err != nil {
			appG.Response(http.StatusInternalServerError, "文件系统删除文件错误", err)
			return
		}
		result = db.Delete(&cert, certDbId)
		if result.Error != nil {
			appG.Response(http.StatusInternalServerError, "数据库删除记录错误", result.Error.Error())
			return
		}

		// 通知用户
		denialReason := body.DenialReason
		notification := model.Notification{
			UserID:       cert.UploaderID,
			AdminID:      adminId,
			Content:      "证书审核未通过",
			IsRead:       false,
			CreatedAt:    time.Now(),
			DenialReason: &denialReason,
		}
		ws.WriteDBAndNotifyUser(db, notification, appG)

		appG.Response(http.StatusOK, "成功", "证书审核不通过，成功删除证书数据")
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
