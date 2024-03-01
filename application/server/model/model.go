package model

import "time"

// server端用户信息结构体，用于注册登录，数据库表格
type User struct {
	ID       string `gorm:"primary_key"` // 身份证号，主键
	RealName string `gorm:"not null"`    // 实名
	Password string `gorm:"not null"`
	Email    string `gorm:"not null"`
}

type ReviewType string

const (
	CET ReviewType = "CET"
	CJT ReviewType = "CJT"
	PHD ReviewType = "PHD"
)

type Admin struct {
	ID         uint       `gorm:"primary_key"` // 主键
	Password   string     `gorm:"not null"`
	ReviewType ReviewType `gorm:"not null"` // 管理员的类型，用于处理对应类型的证书
}

// server端暂存用户证书的结构体，数据库表格
type Cert struct {
	ID           uint      `gorm:"primary_key"` // 主键
	Path         string    `gorm:"not null"`    // 证书存储url
	CertType     string    `gorm:"not null"`    // 证书类型，需要对应类型的管理员进行处理
	CreatedAt    time.Time // 证书上传时间，按照时间升序排列
	UploaderID   string    `gorm:"not null"`      // 上传者ID
	UploaderName string    `gorm:"not null"`      // 上传者姓名
	IsProcessed  bool      `gorm:"default:false"` //是否处理，true的定时删除
}

type AuthorityContactInfo struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type RespondCertificate struct {
	// hash
	HashFile string `json:"hashFile"` // 原文件hash
	HashPath string `json:"hashPath"` // ipfs CID
	// metadata
	CertID               string               `json:"certID"`   // 认证机构颁发的id：机构名-证书ID(cet.org-0001)
	HolderID             string               `json:"holderID"` // 身份证，机构user表的id，主键
	HolderName           string               `json:"holderName"`
	CertType             string               `json:"certType"` // 类型（该证书颁发机构包含的证书类型，也是管理员的审核类型）
	Reviewer             string               `json:"reviewer"`
	IssueDate            string               `json:"issueDate"`  // 证书上显示的颁发日期
	ExpiryDate           string               `json:"expiryDate"` // 证书本身的过期日期
	IssuingAuthority     string               `json:"issuingAuthority"`
	AuthorityContactInfo AuthorityContactInfo `json:"authorityContactInfo"`
	VerifiedTime         string               `json:"verifiedTime"` // 用户上传后管理员核验的时间
	Status               string               `json:"status"`       // 证书状态
	// appendded data
	RetrievedHash string `json:"retrievedHash"` // 从ipfs取出来的文件的哈希，用于对比
}

type LedgerCertificate struct {
	// hash
	HashFile string `json:"hashFile"` // 原文件hash
	HashPath string `json:"hashPath"` // ipfs CID
	// metadata
	CertID               string               `json:"certID"`   // 认证机构颁发的id：机构名-证书ID(cet.org-0001)
	HolderID             string               `json:"holderID"` // 身份证，机构user表的id，主键
	HolderName           string               `json:"holderName"`
	CertType             string               `json:"certType"` // 类型（该证书颁发机构包含的证书类型，也是管理员的审核类型）
	Reviewer             string               `json:"reviewer"`
	IssueDate            string               `json:"issueDate"`  // 证书上显示的颁发日期
	ExpiryDate           string               `json:"expiryDate"` // 证书本身的过期日期
	IssuingAuthority     string               `json:"issuingAuthority"`
	AuthorityContactInfo AuthorityContactInfo `json:"authorityContactInfo"`
	VerifiedTime         string               `json:"verifiedTime"` // 用户上传后管理员核验的时间
	Status               string               `json:"status"`       // 证书状态
}

var CertStatusConstant = func() map[string]string {
	return map[string]string{
		"valid":   "有效",  //正常的有效证书
		"expired": "已过期", //针对有实效的证书类型，过了时效显示已过期状态
		"invaild": "无效",  //已撤销（删除）证书
	}
}

type Notification struct {
	ID           uint      `gorm:"primary_key"`   // 主键
	UserID       string    `gorm:"not null"`      // 接收用户的ID
	AdminID      uint      `gorm:"not null"`      // 发送管理员的ID
	Content      string    `gorm:"not null"`      // 消息内容
	IsRead       bool      `gorm:"default:false"` // 是否已读
	CreatedAt    time.Time // 发送时间
	DenialReason *string
}
