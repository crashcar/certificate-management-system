package model

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

type User struct {
	ID       string `gorm:"primary_key"` // 身份证号
	RealName string `gorm:"not null"`    // 实名
	Password string `gorm:"not null"`
	Email    string `gorm:"not null"`
	Role     int    `gorm:"default:0"` // 0--普通用户 1--管理员（机构中的角色，系统中全是用户）
}

type AuthorityContactInfo struct {
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Address string `json:"address"`
}

type Certificate struct {
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

var CertStatusConstant = func() map[string]string {
	return map[string]string{
		"valid":   "有效",  //正常的有效证书
		"expired": "已过期", //针对有实效的证书类型，过了时效显示已过期状态
		"invaild": "无效",  //已撤销（删除）证书
	}
}
