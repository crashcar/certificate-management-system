package api

import (
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// 管理系统：查看用户所有证书，通过用户id查询用户在所有机构的所有证书
func QueryCertByUserSys(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return nil
}

// 管理系统：通过证书持有人的id和证书id查询该证书，返回{证书ID，持有人ID，证书颁发机构，有效期，状态，修改情况}
func QueryCertByFullInfoSys(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return nil
}

// org：-用户- 通过证书持有人的id查询该人在本机构的所有证书
func QueryAllCertByUserOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return nil
}

// org：-管理员- 查询该机构的所有证书以及其持有人
func QueryCertOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return nil
}

// org：-管理员- 上传证书
func UploadCertOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return nil
}

// org：-管理员- 删除证书(或许不需要？model/cert结构的字段status，用invalid、expired表示证书已无效)
func DeleteCertOrg(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	return nil
}
