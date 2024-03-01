package api

import (
	"chaincode/model"
	"chaincode/pkg/utils"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// 管理系统：单次查询调用接口，证书作为查询主键
func QueryCertByInfos(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var certificateList []model.Certificate
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.CertificateKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var certificate model.Certificate
			err := json.Unmarshal(v, &certificate)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryCertByInfos-反序列化出错: %s", err))
			}
			certificateList = append(certificateList, certificate)
		}
	}
	certificateListByte, err := json.Marshal(certificateList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryCertByInfos-序列化出错: %s", err))
	}
	return shim.Success(certificateListByte)
}

// 管理系统：多次查询调用接口，证书作为查询主键
func QueryCertByInfosLists(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var certificateList []model.Certificate
	results, err := utils.GetStateByPartialCompositeKeys(stub, model.CertificateKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var certificate model.Certificate
			err := json.Unmarshal(v, &certificate)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryCertByInfosLists-反序列化出错: %s", err))
			}
			certificateList = append(certificateList, certificate)
		}
	}
	certificateListByte, err := json.Marshal(certificateList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryCertByInfosLists-序列化出错: %s", err))
	}
	return shim.Success(certificateListByte)
}

// 管理系统：单次查询调用接口, 机构作为查询主键
func QueryCertByAuthority(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var certificateList []model.Certificate
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.AuthorityKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var certificate model.Certificate
			err := json.Unmarshal(v, &certificate)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryCertByAuthority-反序列化出错: %s", err))
			}
			certificateList = append(certificateList, certificate)
		}
	}
	certificateListByte, err := json.Marshal(certificateList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryCertQueryCertByAuthorityInfos-序列化出错: %s", err))
	}
	return shim.Success(certificateListByte)
}

// 管理系统：多次查询调用接口, 机构作为查询主键
func QueryCertByAuthorityLists(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	var certificateList []model.Certificate
	results, err := utils.GetStateByPartialCompositeKeys(stub, model.AuthorityKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v != nil {
			var certificate model.Certificate
			err := json.Unmarshal(v, &certificate)
			if err != nil {
				return shim.Error(fmt.Sprintf("QueryCertByAuthorityLists-反序列化出错: %s", err))
			}
			certificateList = append(certificateList, certificate)
		}
	}
	certificateListByte, err := json.Marshal(certificateList)
	if err != nil {
		return shim.Error(fmt.Sprintf("QueryCertByAuthorityLists-序列化出错: %s", err))
	}
	return shim.Success(certificateListByte)
}

// org：-管理员- 上传证书
func UploadCertOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 13 {
		return shim.Error(fmt.Sprintf("参数个数不满足, 输入%d, 应为13\n", len(args)))
	}

	//检测参数是否为空
	for _, v := range args {
		if v == "" {
			return shim.Error("参数存在空值")
		}
	}

	hashFile := args[0]
	hashPath := args[1]
	certID := args[2]
	holderID := args[3]
	holderName := args[4]
	certType := args[5]
	reviewer := args[6]
	issueDate := args[7]
	expiryDate := args[8]
	issuingAuthority := args[9]
	authorityPhone := args[10]
	authorityEmail := args[11]
	authorityAddress := args[12]
	//输入参数的逻辑检查：待添加

	authorityContactInfo := &model.AuthorityContactInfo{
		Phone:   authorityPhone,
		Email:   authorityEmail,
		Address: authorityAddress,
	}

	//判断一下是否过期
	expireTime, err := time.Parse("2006-01-02", expiryDate)
	if err != nil {
		return shim.Error(fmt.Sprintf("过期日期格式解析出错:%s", err))
	}

	now := time.Now()
	todayDate, err := time.Parse("2006-01-02", now.Format("2006-01-02"))
	if err != nil {
		return shim.Error(fmt.Sprintf("当天日期格式解析出错:%s", err))
	}

	certificate := &model.Certificate{
		HashFile:             hashFile,
		HashPath:             hashPath,
		CertID:               certID,
		HolderID:             holderID,
		HolderName:           holderName,
		Reviewer:             reviewer,
		CertType:             certType,
		IssueDate:            issueDate,
		ExpiryDate:           expiryDate,
		IssuingAuthority:     issuingAuthority,
		AuthorityContactInfo: *authorityContactInfo,
	}

	if expireTime.Before(todayDate) {
		certificate.Status = model.CertStatusConstant()["expired"]
	} else {
		certificate.Status = model.CertStatusConstant()["valid"]
	}
	//写入账本

	if err := utils.WriteLedger(certificate, stub, model.CertificateKey, []string{certificate.CertID, certificate.HolderID, certificate.HolderName, certificate.IssuingAuthority}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	if err := utils.WriteLedger(certificate, stub, model.AuthorityKey, []string{certificate.IssuingAuthority, certificate.HolderID, certificate.CertID}); err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	//将成功创建的信息返回
	certificateByte, err := json.Marshal(certificate)
	if err != nil {
		return shim.Error(fmt.Sprintf("序列化成功创建的信息出错: %s", err))
	}
	// 成功返回
	return shim.Success(certificateByte)
}

// org：-管理员- 删除证书(或许不需要？model/cert结构的字段status，用invalid、expired表示证书已无效)
func DeleteCertOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	if len(args) != 1 {
		return shim.Error("参数个数不满足")
	}
	//certID := args[0]

	//查询证书是否存在
	var certificateList []model.Certificate
	results, err := utils.GetStateByPartialCompositeKeys2(stub, model.CertificateKey, args)
	if err != nil {
		return shim.Error(fmt.Sprintf("%s", err))
	}
	for _, v := range results {
		if v == nil {
			return shim.Error(fmt.Sprintf("待删除的证书不存在"))
		}
		var certificate model.Certificate
		err := json.Unmarshal(v, &certificate)
		if err != nil {
			return shim.Error(fmt.Sprintf("DeleteCertOrg-反序列化出错: %s", err))
		}
		certificateList = append(certificateList, certificate)
		//删除证书
		if err := utils.DelLedger(stub, model.CertificateKey, []string{certificate.CertID, certificate.HolderID, certificate.HolderName, certificate.IssuingAuthority}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
		if err := utils.DelLedger(stub, model.AuthorityKey, []string{certificate.IssuingAuthority, certificate.HolderID, certificate.CertID}); err != nil {
			return shim.Error(fmt.Sprintf("%s", err))
		}
	}

	certificateByte, err := json.Marshal(certificateList)
	//fmt.Sprintf("删除完毕")
	return shim.Success(certificateByte)
}
