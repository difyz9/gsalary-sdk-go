package gsalary

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// GSalaryRequest 请求对象
type GSalaryRequest struct {
	Method    string
	Path      string
	QueryArgs map[string]string
	Body      map[string]interface{}
}

// NewRequest 创建新的请求
func NewRequest(method, path string) *GSalaryRequest {
	return &GSalaryRequest{
		Method:    method,
		Path:      path,
		QueryArgs: make(map[string]string),
		Body:      make(map[string]interface{}),
	}
}

// Valid 验证请求参数
func (r *GSalaryRequest) Valid() error {
	if r.Method == "" || r.Path == "" {
		return errors.New("invalid request: method and path are required")
	}
	validMethods := map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}
	if !validMethods[r.Method] {
		return errors.New("invalid method: must be GET, POST, PUT, or DELETE")
	}
	return nil
}

// HasBody 检查请求是否有body
func (r *GSalaryRequest) HasBody() bool {
	if r.Method == "POST" || r.Method == "PUT" {
		return len(r.Body) > 0
	}
	return false
}

// PathWithArgs 返回带查询参数的路径
func (r *GSalaryRequest) PathWithArgs(escape bool) string {
	if len(r.QueryArgs) == 0 {
		return r.Path
	}

	var pairs []string
	for k, v := range r.QueryArgs {
		if escape {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, url.QueryEscape(v)))
		} else {
			pairs = append(pairs, fmt.Sprintf("%s=%s", k, v))
		}
	}
	return r.Path + "?" + strings.Join(pairs, "&")
}

// getBodyHash 计算body的SHA256哈希值
func (r *GSalaryRequest) getBodyHash() (string, error) {
	if !r.HasBody() {
		return "", nil
	}

	bodyBytes, err := json.Marshal(r.Body)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(bodyBytes)
	return base64.StdEncoding.EncodeToString(hash[:]), nil
}

// SignRequest 对请求进行签名
func (r *GSalaryRequest) SignRequest(config *GSalaryConfig) (*AuthorizeHeaderInfo, error) {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	
	bodyHash, err := r.getBodyHash()
	if err != nil {
		return nil, err
	}

	// 构造签名基础字符串
	signBase := fmt.Sprintf("%s %s\n%s\n%s\n%s\n",
		r.Method,
		r.PathWithArgs(false),
		config.AppID,
		timestamp,
		bodyHash)

	// 生成RSA签名
	signature, err := generateRSASignature(config.GetClientPrivateKey(), signBase)
	if err != nil {
		return nil, err
	}

	return NewAuthorizeHeaderInfo("RSA2", timestamp, signature), nil
}

// VerifySignature 验证响应签名
func (r *GSalaryRequest) VerifySignature(config *GSalaryConfig, headerInfo *AuthorizeHeaderInfo, responseBody string) bool {
	// 计算响应body的哈希
	hash := sha256.Sum256([]byte(responseBody))
	bodyHash := base64.StdEncoding.EncodeToString(hash[:])

	// 构造签名基础字符串
	signBase := fmt.Sprintf("%s %s\n%s\n%s\n%s\n",
		r.Method,
		r.PathWithArgs(false),
		config.AppID,
		headerInfo.Timestamp,
		bodyHash)

	// 验证签名
	return verifyRSASignature(config.GetServerPublicKey(), signBase, headerInfo.Signature)
}

// generateRSASignature 生成RSA签名
func generateRSASignature(privateKey *rsa.PrivateKey, signBase string) (string, error) {
	hash := sha256.Sum256([]byte(signBase))
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hash[:])
	if err != nil {
		return "", fmt.Errorf("failed to generate RSA signature: %w", err)
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

// verifyRSASignature 验证RSA签名
func verifyRSASignature(publicKey *rsa.PublicKey, signBase, signatureBase64 string) bool {
	signature, err := base64.StdEncoding.DecodeString(signatureBase64)
	if err != nil {
		return false
	}

	hash := sha256.Sum256([]byte(signBase))
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signature)
	return err == nil
}
