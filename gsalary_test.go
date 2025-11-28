package gsalary

import (
	"strings"
	"testing"
)

// TestConfigPEM 测试密钥配置
func TestConfigPEM(t *testing.T) {
	config := NewConfig()
	config.AppID = "test_app_id"

	// 测试配置私钥
	privateKeyPEM := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDIw4HnZKo5K2Az
um/qxI0FStv0XKma6NdItFOPOfNmqU+UJeVaxyr06Mer2LFAvIBudCylJrp29Plt
wvWNBBDkMb/mCTpcV9JrbDskxjGdyKtkK4iN6a2HGHJUt471nX8Vl24+C1yZK8qk
ktu4nlRAU6kmo18qhVgZoSYEmaML7ApE+ncm5cgHrdnOHMsSKJVXpAoGARh9k
Nvksg67Gtknf6+nQNe+tWKALet63i4oxTkUHZWJ4BeIP2g+RmpIYh7PSY7vPqNWj
x/842jo3ThSpJm9yNdIaiqMyFQfICSDdRZGPxaPakng/ufVljmQQDEmxrCwBzRb1
4UqO/pQfrlcb3pwkhGu7GUhCdH4znbDhxYQ5nLECgYEApW/B6MQraHMfDNEYVzG1
hAdKWVI1W3I5KpI5VCogZiR9ta4VILfwTtmlLyYuBnSo38KvfNI+iYg+DfXj4knE
5LbTY+6xBcSEeLTppyzykneLmRn1JcUAr51b5SRrNP+WRWGELXnALGCt5OIrdQM5
KGI8ymPrCm+hIK6pfJ+qcyw=
-----END PRIVATE KEY-----`

	err := config.ConfigClientPrivateKeyPEM(privateKeyPEM)
	if err != nil {
		t.Fatalf("Failed to config private key: %v", err)
	}

	if config.GetClientPrivateKey() == nil {
		t.Error("Private key should not be nil")
	}

	// 测试配置公钥
	publicKeyPEM := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyMOB52SqOStgM7pv6sSN
BUrb9FypmDO51s8OCwA6M92wjEdHhp0ND24Nt0Zix
aQIDAQAB
-----END PUBLIC KEY-----`

	err = config.ConfigServerPublicKeyPEM(publicKeyPEM)
	if err != nil {
		t.Fatalf("Failed to config public key: %v", err)
	}

	if config.GetServerPublicKey() == nil {
		t.Error("Public key should not be nil")
	}
}

// TestAuthorizeHeaderInfo 测试鉴权头部信息
func TestAuthorizeHeaderInfo(t *testing.T) {
	// 测试创建和转换
	auth := NewAuthorizeHeaderInfo("RSA2", "1234567890", "test_signature")
	if !auth.Valid() {
		t.Error("Auth header should be valid")
	}

	headerValue := auth.ToHeaderValue()
	if headerValue == "" {
		t.Error("Header value should not be empty")
	}

	// 测试解析
	parsed := FromHeaderValue(headerValue)
	if !parsed.Valid() {
		t.Error("Parsed auth header should be valid")
	}
	if parsed.Algorithm != "RSA2" {
		t.Errorf("Expected algorithm RSA2, got %s", parsed.Algorithm)
	}
	if parsed.Timestamp != "1234567890" {
		t.Errorf("Expected timestamp 1234567890, got %s", parsed.Timestamp)
	}
}

// TestRequest 测试请求对象
func TestRequest(t *testing.T) {
	// 测试GET请求
	getReq := NewRequest("GET", "/v1/cards")
	getReq.QueryArgs["page"] = "1"
	getReq.QueryArgs["limit"] = "20"

	if err := getReq.Valid(); err != nil {
		t.Errorf("GET request should be valid: %v", err)
	}

	pathWithArgs := getReq.PathWithArgs(false)
	if !strings.Contains(pathWithArgs, "page=1") || !strings.Contains(pathWithArgs, "limit=20") {
		t.Errorf("Path with args incorrect: %s", pathWithArgs)
	}

	// 测试POST请求
	postReq := NewRequest("POST", "/v1/exchange/quotes")
	postReq.Body["sell_currency"] = "USD"
	postReq.Body["buy_currency"] = "CNY"

	if !postReq.HasBody() {
		t.Error("POST request should have body")
	}

	// 测试无效请求
	invalidReq := NewRequest("", "")
	if err := invalidReq.Valid(); err == nil {
		t.Error("Invalid request should return error")
	}
}

// TestSignature 测试签名和验证
func TestSignature(t *testing.T) {
	config := NewConfig()
	config.AppID = "test_app"

	privateKeyPEM := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDIw4HnZKo5K2Az
um/qxI0FStv0XKma6NdItFOPOfNmqU+UJeVaxyr06Mer2LFAvIBudCylJrp29Plt
wvWNBBDkMb/mCTpc7GUhCdH4znbDhxYQ5nLECgYEApW/B6MQraHMfDNEYVzG1
hAdKWVI1W3I5KpI5VCogZiR9ta4VILfwTtmlLyYuBnSo38KvfNI+iYg+DfXj4knE
5LbTY+6xBcSEeLTppyzykneLmRn1JcUAr51b5SRrNP+WRWGELXnALGCt5OIrdQM5
KGI8ymPrCm+hIK6pfJ+qcyw=
-----END PRIVATE KEY-----`

	publicKeyPEM := `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAyMOB52SqOStgM7pv6sSN
BUrb9FypmujXUjbSfnzCqYGQ3idAq/hItd8EXFo6Cx+BHTnG0pk6VkAKkR
QFshA2ETnDbHNBAF+5a1lPitWV2h46m8DO51s8OCwA6M92wjEdHhp0ND24Nt0Zix
aQIDAQAB
-----END PUBLIC KEY-----`

	if err := config.ConfigClientPrivateKeyPEM(privateKeyPEM); err != nil {
		t.Fatalf("Failed to config private key: %v", err)
	}
	if err := config.ConfigServerPublicKeyPEM(publicKeyPEM); err != nil {
		t.Fatalf("Failed to config public key: %v", err)
	}

	// 测试签名
	req := NewRequest("POST", "/v1/test")
	req.Body["test"] = "data"

	authHeader, err := req.SignRequest(config)
	if err != nil {
		t.Fatalf("Failed to sign request: %v", err)
	}

	if !authHeader.Valid() {
		t.Error("Auth header should be valid after signing")
	}

	if authHeader.Algorithm != "RSA2" {
		t.Errorf("Expected algorithm RSA2, got %s", authHeader.Algorithm)
	}

	if authHeader.Signature == "" {
		t.Error("Signature should not be empty")
	}
}
