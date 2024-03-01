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

func UploadFileToIPFS(appG app.Gin, filepath string, ipfsnode string) string {
	// 从filepath读取file
	input, err := os.Open(filepath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("打开文件出错: %s", err.Error()))
		return ""
	}
	defer input.Close()

	// 读取文件内容
	content, err := io.ReadAll(input)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取文件内容时出错: %s", err.Error()))
		return ""
	}

	// 读取 AES 密钥
	keyFile := "./aesKey.txt"
	key, err := os.ReadFile(keyFile)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取 AES 密钥时出错: %s", err.Error()))
		return ""
	}

	// 加密file读取的content，返回在buffer中
	buffer := cryptoutils.EncryptContent(content, key, appG)

	// 连接到IPFS节点
	sh := shell.NewShell(ipfsnode)

	// 加密后的文件上传至ipfs，使用NewReader创建io.Reader
	cid, err := sh.Add(bytes.NewReader(buffer))
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("文件上传IPFS失败: %s", err.Error()))
		return ""
	}

	//返回CID
	return cid
}

func GetFileFromIPFS(appG app.Gin, cid string, ipfsnode string) []byte {
	// 根据cid从ipfs获取文件内容
	sh := shell.NewShell(ipfsnode)
	reader, err := sh.Cat(cid)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("从IPFS读取文件失败: %s", err.Error()))
		return nil
	}
	defer reader.Close()

	encryptedData, err := io.ReadAll(reader)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取加密数据失败: %s", err.Error()))
		return nil
	}

	// 读取 AES 密钥
	keyFile := "./aesKey.txt"
	key, err := os.ReadFile(keyFile)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("读取 AES 密钥时出错: %s", err.Error()))
		return nil
	}

	// 利用aes密钥解密文件
	buffer := cryptoutils.DecryptContent(encryptedData, key, appG)

	return buffer
}
