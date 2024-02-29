package cryptoutils

import (
	"application/pkg/app"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
)

// 使用 AES 密钥加密文件
func EncryptContent(content []byte, key []byte, appG app.Gin) ([]byte, error) {
	// 创建cipher.Block实例
	block, err := aes.NewCipher(key)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建cipher.Block实例时出错: %s", err.Error()))
		return nil, err
	}

	// 如果AES加密需要IV，可以在这里生成，这里假设key足够长且为了简化直接使用key前16字节作为IV
	iv := key[:aes.BlockSize]

	// 设置加密模式，这里使用CBC
	mode := cipher.NewCBCEncrypter(block, iv)

	// PKCS#7填充
	padding := aes.BlockSize - len(content)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	content = append(content, padText...)

	// 加密文件内容
	encrypted := make([]byte, len(content))
	mode.CryptBlocks(encrypted, content)

	// 创建bytes.Buffer并写入加密后的数据
	buffer := bytes.NewBuffer(encrypted)
	return buffer.Bytes(), nil
}

// 使用 AES 密钥解密文件
func DecryptContent(encryptContent []byte, key []byte, appG app.Gin) ([]byte, error) {
	// 创建cipher.Block实例
	block, err := aes.NewCipher(key)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", fmt.Sprintf("创建cipher.Block实例时出错: %s", err.Error()))
		return nil, err
	}

	// IV, 加密时使用的是key的前16字节
	iv := key[:aes.BlockSize]

	// 设置解密模式 CBC
	mode := cipher.NewCBCDecrypter(block, iv)

	// 解密
	decrypted := make([]byte, len(encryptContent))
	mode.CryptBlocks(decrypted, encryptContent)

	// 去除PKCS#7填充
	padding := decrypted[len(decrypted)-1]
	padLen := int(padding)
	if padLen > aes.BlockSize || padLen == 0 {
		appG.Response(http.StatusInternalServerError, "失败", "解密失败，填充错误")
		return nil, err
	}

	// 验证并去除填充
	for _, v := range decrypted[len(decrypted)-padLen:] {
		if v != padding {
			appG.Response(http.StatusInternalServerError, "失败", "解密失败，填充错误")
			return nil, err
		}
	}
	decrypted = decrypted[:len(decrypted)-padLen]

	return decrypted, nil
}

func HashFile(filePath string, appG app.Gin) string {
	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "HashFile打开文件失败", err)
		return ""
	}
	defer file.Close()

	// 创建一个新的SHA256哈希对象
	hash := sha256.New()

	// 将文件内容复制到哈希对象中
	if _, err := io.Copy(hash, file); err != nil {
		appG.Response(http.StatusInternalServerError, "计算hash失败", err)
		return ""
	}

	// 计算最终的哈希值并将其格式化为16进制字符串
	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)
	return hashString
}
