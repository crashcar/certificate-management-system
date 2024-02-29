package ipfs

import (
	"application/pkg/app"
	"application/pkg/cryptoutils"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	shell "github.com/ipfs/go-ipfs-api"
)

func UploadFileToIPFS(appG app.Gin, filepath string, ipfsnode string) (string, error) {
	// 从filepath读取file
	input, err := os.Open(filepath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("打开文件出错: %s", err.Error()))
		return "", err
	}
	defer input.Close()

	// 读取文件内容
	content, err := io.ReadAll(input)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取文件内容时出错: %s", err.Error()))
		return "", err
	}

	// 读取 AES 密钥
	keyFile := "./aesKey.txt"
	key, err := os.ReadFile(keyFile)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取 AES 密钥时出错: %s", err.Error()))
		return "", err
	}

	// 加密file读取的content，返回在buffer中
	buffer, err := cryptoutils.EncryptContent(content, key, appG)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("加密文件内容失败: %s", err.Error()))
	}

	// 连接到IPFS节点
	sh := shell.NewShell(ipfsnode)

	// 加密后的文件上传至ipfs，使用NewReader创建io.Reader
	cid, err := sh.Add(bytes.NewReader(buffer))
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("文件上传IPFS失败: %s", err.Error()))
		return "", err
	}

	//返回CID
	return cid, nil
}

func GetFileFromIPFS(appG app.Gin, cid string, ipfsnode string) ([]byte, error) {
	// 根据cid从ipfs获取文件内容
	sh := shell.NewShell(ipfsnode)
	reader, err := sh.Cat(cid)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("从IPFS读取文件失败: %s", err.Error()))
		return nil, err
	}
	defer reader.Close()

	encryptedData, err := io.ReadAll(reader)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取加密数据失败: %s", err.Error()))
		return nil, err
	}

	// 读取 AES 密钥
	keyFile := "./aesKey.txt"
	key, err := os.ReadFile(keyFile)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取 AES 密钥时出错: %s", err.Error()))
		return nil, err
	}

	// 利用aes密钥解密文件
	buffer, err := cryptoutils.DecryptContent(encryptedData, key, appG)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("解密文件时出错: %s", err.Error()))
		return nil, err
	}

	return buffer, nil
}
