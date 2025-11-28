package gsalary

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"strings"
)

// GSalaryConfig 配置信息
type GSalaryConfig struct {
	Endpoint         string
	AppID            string
	clientPrivateKey *rsa.PrivateKey
	serverPublicKey  *rsa.PublicKey
}

// NewConfig 创建新的配置
func NewConfig() *GSalaryConfig {
	return &GSalaryConfig{
		Endpoint: "https://api-test.gsalary.com",
	}
}

// GetClientPrivateKey 获取客户端私钥
func (c *GSalaryConfig) GetClientPrivateKey() *rsa.PrivateKey {
	return c.clientPrivateKey
}

// GetServerPublicKey 获取服务端公钥
func (c *GSalaryConfig) GetServerPublicKey() *rsa.PublicKey {
	return c.serverPublicKey
}

// ConcatPath 拼接完整URL路径
func (c *GSalaryConfig) ConcatPath(path string) string {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	if strings.HasSuffix(c.Endpoint, "/") {
		return c.Endpoint + path
	}
	return c.Endpoint + "/" + path
}

// insertNewLines 在PEM字符串中插入换行符
func insertNewLines(value string) string {
	value = strings.TrimSpace(value)
	// 去除已有的头尾
	if strings.HasPrefix(value, "-----BEGIN PUBLIC KEY-----") {
		value = strings.TrimSpace(value[len("-----BEGIN PUBLIC KEY-----"):])
	}
	if strings.HasSuffix(value, "-----END PUBLIC KEY-----") {
		value = strings.TrimSpace(value[:len(value)-len("-----END PUBLIC KEY-----")])
	}
	if strings.HasPrefix(value, "-----BEGIN PRIVATE KEY-----") {
		value = strings.TrimSpace(value[len("-----BEGIN PRIVATE KEY-----"):])
	}
	if strings.HasSuffix(value, "-----END PRIVATE KEY-----") {
		value = strings.TrimSpace(value[:len(value)-len("-----END PRIVATE KEY-----")])
	}
	
	// 如果已经有换行符，无需再次插入
	if strings.Contains(value, "\n") {
		return value
	}
	
	// 按每64个字符插入一个换行符
	var result strings.Builder
	for i := 0; i < len(value); i += 64 {
		end := i + 64
		if end > len(value) {
			end = len(value)
		}
		result.WriteString(value[i:end])
		if end < len(value) {
			result.WriteString("\n")
		}
	}
	return result.String()
}

// ConfigClientPrivateKeyPEM 配置客户端私钥（PEM格式）
func (c *GSalaryConfig) ConfigClientPrivateKeyPEM(pemStr string) error {
	if !strings.HasPrefix(pemStr, "-----BEGIN PRIVATE KEY-----") {
		if !strings.Contains(strings.TrimSpace(pemStr), "\n") {
			pemStr = insertNewLines(pemStr)
		}
		pemStr = "-----BEGIN PRIVATE KEY-----\n" + pemStr + "\n-----END PRIVATE KEY-----"
	}
	
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return errors.New("failed to parse PEM block containing the private key")
	}
	
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return err
	}
	
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return errors.New("not an RSA private key")
	}
	
	c.clientPrivateKey = rsaPrivateKey
	return nil
}

// ConfigClientPrivateKeyPEMFile 从文件加载客户端私钥
func (c *GSalaryConfig) ConfigClientPrivateKeyPEMFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return c.ConfigClientPrivateKeyPEM(string(content))
}

// ConfigServerPublicKeyPEM 配置服务端公钥（PEM格式）
func (c *GSalaryConfig) ConfigServerPublicKeyPEM(pemStr string) error {
	if !strings.HasPrefix(pemStr, "-----BEGIN PUBLIC KEY-----") {
		if !strings.Contains(strings.TrimSpace(pemStr), "\n") {
			pemStr = insertNewLines(pemStr)
		}
		pemStr = "-----BEGIN PUBLIC KEY-----\n" + pemStr + "\n-----END PUBLIC KEY-----"
	}
	
	block, _ := pem.Decode([]byte(pemStr))
	if block == nil {
		return errors.New("failed to parse PEM block containing the public key")
	}
	
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}
	
	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("not an RSA public key")
	}
	
	c.serverPublicKey = rsaPublicKey
	return nil
}

// ConfigServerPublicKeyPEMFile 从文件加载服务端公钥
func (c *GSalaryConfig) ConfigServerPublicKeyPEMFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return c.ConfigServerPublicKeyPEM(string(content))
}
