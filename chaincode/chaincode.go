package main

import (
	"chaincode/api"
	"chaincode/model"
	"chaincode/pkg/utils"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

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
	HolderID             string               `json:"holderID"`
	HolderName           string               `json:"holderName"`
	CertType             string               `json:"certType"`
	IssueDate            string               `json:"issueDate"`
	ExpiryDate           string               `json:"expiryDate"`
	IssuingAuthority     string               `json:"issuingAuthority"`
	AuthorityContactInfo AuthorityContactInfo `json:"authorityContactInfo"`
	Status               string               `json:"status"`
}

// Init 链码初始化
func (t *Certificate) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("链码初始化")
	var hashfile = string("hashfile")
	var hashpath = string("hashpath")
	var certId = string("certId")
	var holderId = string("holderId")
	var holderName = string("holderName")
	var certType = string("certType")
	var issueDate = string("2020-01-02")
	var expiryDate = string("2020-01-02")
	var issuingAuthority = string("issuingAuthority")
	var phone = string("phone")
	var email = string("email")
	var address = string("address")

	authorityInfo := &model.AuthorityContactInfo{
		Phone:   phone,
		Email:   email,
		Address: address,
	}

	certificate := &model.Certificate{
		HashFile:             hashfile,
		HashPath:             hashpath,
		CertID:               certId,
		HolderID:             holderId,
		HolderName:           holderName,
		CertType:             certType,
		IssueDate:            issueDate,
		ExpiryDate:           expiryDate,
		IssuingAuthority:     issuingAuthority,
		AuthorityContactInfo: *authorityInfo,
	}
	if err := utils.WriteLedger(certificate, stub, model.CertificateKey, []string{certificate.CertID, certificate.HolderID, certificate.HolderName, certificate.IssuingAuthority}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	if err := utils.WriteLedger(certificate, stub, model.AuthorityKey, []string{certificate.IssuingAuthority, certificate.HolderID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	return shim.Success([]byte("Init Success"))
}

// Invoke 实现Invoke接口调用智能合约
func (t *Certificate) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()
	switch funcName {
	case "Init":
		return api.MyInit(stub, args)
	case "hello":
		return api.Hello(stub, args)
		// return api.CreateRealEstate(stub, args)
	case "queryAccountList":
		return api.QueryAccountList(stub, args)
	case "createRealEstate":
		return api.CreateRealEstate(stub, args)
	case "queryRealEstateList":
		return api.QueryRealEstateList(stub, args)
	case "createSelling":
		return api.CreateSelling(stub, args)
	case "createSellingByBuy":
		return api.CreateSellingByBuy(stub, args)
	case "querySellingList":
		return api.QuerySellingList(stub, args)
	case "querySellingListByBuyer":
		return api.QuerySellingListByBuyer(stub, args)
	case "updateSelling":
		return api.UpdateSelling(stub, args)
	case "createDonating":
		return api.CreateDonating(stub, args)
	case "queryDonatingList":
		return api.QueryDonatingList(stub, args)
	case "queryDonatingListByGrantee":
		return api.QueryDonatingListByGrantee(stub, args)
	case "updateDonating":
		return api.UpdateDonating(stub, args)

	case "queryCertByInfos":
		return api.QueryCertByInfos(stub, args)
	case "queryCertByInfosLists":
		return api.QueryCertByInfosLists(stub, args)
	case "queryCertByAuthority":
		return api.QueryCertByAuthority(stub, args)
	case "queryCertByAuthorityLists":
		return api.QueryCertByAuthorityLists(stub, args)
	case "uploadCertOrg":
		return api.UploadCertOrg(stub, args)
	case "deleteCertOrg":
		return api.DeleteCertOrg(stub, args)
	default:
		return shim.Error(fmt.Sprintf("没有该功能: %s", funcName))
	}
}

func main() {
	timeLocal, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	time.Local = timeLocal
	err = shim.Start(new(Certificate))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
