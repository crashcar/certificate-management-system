package model

import "time"

// Selling 销售要约
// 需要确定ObjectOfSale是否属于Seller
// 买家初始为空
// Seller和ObjectOfSale一起作为复合键,保证可以通过seller查询到名下所有发起的销售
type Selling struct {
	ObjectOfSale  string  `json:"objectOfSale"`  //销售对象(正在出售的车辆RealEstateID)
	Seller        string  `json:"seller"`        //发起销售人、卖家(卖家AccountId)
	Buyer         string  `json:"buyer"`         //参与销售人、买家(买家AccountId)
	Price         float64 `json:"price"`         //价格
	CreateTime    string  `json:"createTime"`    //创建时间
	SalePeriod    int     `json:"salePeriod"`    //智能合约的有效期(单位为天)
	SellingStatus string  `json:"sellingStatus"` //销售状态
}

// SellingStatusConstant 销售状态
var SellingStatusConstant = func() map[string]string {
	return map[string]string{
		"saleStart": "销售中", //正在销售状态,等待买家光顾
		"cancelled": "已取消", //被卖家取消销售或买家退款操作导致取消
		"expired":   "已过期", //销售期限到期
		"delivery":  "交付中", //买家买下并付款,处于等待卖家确认收款状态,如若卖家未能确认收款，买家可以取消并退款
		"done":      "完成",  //卖家确认接收资金，交易完成
	}
}

// Donating 捐赠要约
// 需要确定ObjectOfDonating是否属于Donor
// 需要指定受赠人Grantee，并等待受赠人同意接收
type Donating struct {
	ObjectOfDonating string `json:"objectOfDonating"` //捐赠对象(正在捐赠的车辆RealEstateID)
	Donor            string `json:"donor"`            //捐赠人(捐赠人AccountId)
	Grantee          string `json:"grantee"`          //受赠人(受赠人AccountId)
	CreateTime       string `json:"createTime"`       //创建时间
	DonatingStatus   string `json:"donatingStatus"`   //捐赠状态
}

// DonatingStatusConstant 捐赠状态
var DonatingStatusConstant = func() map[string]string {
	return map[string]string{
		"donatingStart": "捐赠中", //捐赠人发起捐赠合约，等待受赠人确认受赠
		"cancelled":     "已取消", //捐赠人在受赠人确认受赠之前取消捐赠或受赠人取消接收受赠
		"done":          "完成",  //受赠人确认接收，交易完成
	}
}

/*************** CERTIFICATE MANAGEMENT ******************/

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

type Certificate struct {
	// hash
	HashFile string `json:"hashFile"` // 原文件hash
	HashPath string `json:"hashPath"` // ipfs CID hash
	// metadata
	CertID               string               `json:"certID"`   // 认证机构颁发的id：机构名-证书ID(cet.org-0001)
	HolderID             string               `json:"holderID"` // 身份证，机构user表的id，主键
	HolderName           string               `json:"holderName"`
	CertType             string               `json:"certType"` // 类型（但也不知道有哪些类型，可以删了也可以留着吧）
	Reviewer             string               `json:"reviewer"`
	IssueDate            string               `json:"issueDate"`  // 证书上显示的颁发日期
	ExpiryDate           string               `json:"expiryDate"` // 证书本身的过期日期
	IssuingAuthority     string               `json:"issuingAuthority"`
	AuthorityContactInfo AuthorityContactInfo `json:"authorityContactInfo"`
	VerifiedTime         string               `json:"verifiedTime"` // 用户上传后管理员核验的时间
	Status               string               `json:"status"`       // 证书状态（）
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
