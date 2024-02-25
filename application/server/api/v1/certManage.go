package v1

import (
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

// // 管理系统：通过证书持有人的id和证书id查询该证书，返回{证书ID，持有人ID，证书颁发机构，有效期，状态，修改情况}
// func QueryCertByFullInfoSys(c *gin.Context) {

// }

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
	Status               string               `json:"status"`
}

type DeleteCertOrgRequestBody struct {
	HoderID string `json:"holderID"`
	CertID  string `json:"certID"`
}

// org：-用户- 通过证书持有人的id查询该人在本机构的所有证书
func QueryCertByUserOrg(c *gin.Context) {

}

// org：-管理员- 查询该机构的所有证书以及其持有人
func QueryCertOrg(c *gin.Context) {

}

// org：-管理员- 上传证书
func UploadCertOrg(c *gin.Context) {

}

// org：-管理员- 删除证书(或许不需要？model/cert结构的字段status，用invalid、expired表示证书已无效)
func DeleteCertOrg(c *gin.Context) {

}
