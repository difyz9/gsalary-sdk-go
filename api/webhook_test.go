package api

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestWebhookHandler_HandleWebhook 测试Webhook处理
func TestWebhookHandler_HandleWebhook(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建Webhook处理器
	handler := NewWebhookHandler(testConfig)
	
	// 构造测试数据
	webhookData := map[string]interface{}{
		"app_id":        testConfig.AppID,
		"business_type": EventAcquiringPaymentResult,
		"timestamp":     time.Now().Unix() * 1000,
		"data": map[string]interface{}{
			"payment_request_id":  "PAY_TEST_123456",
			"payment_id":          "PAY_ID_123456",
			"payment_amount":      100.0,
			"payment_currency":    "USD",
			"payment_status":      "SUCCESS",
			"payment_result_code": "SUCCESS",
			"payment_time":        time.Now().Format(time.RFC3339),
		},
	}
	
	bodyBytes, err := json.Marshal(webhookData)
	if err != nil {
		t.Fatalf("Failed to marshal webhook data: %v", err)
	}
	
	// 使用客户端私钥生成签名（模拟服务端用服务端私钥签名）
	// 注意：在实际场景中，服务端使用服务端私钥签名，客户端使用服务端公钥验证
	// 这里为了测试方便，使用客户端私钥模拟
	hash := sha256.Sum256(bodyBytes)
	signature, err := rsa.SignPKCS1v15(rand.Reader, testConfig.GetClientPrivateKey(), crypto.SHA256, hash[:])
	if err != nil {
		t.Fatalf("Failed to sign: %v", err)
	}
	
	signatureStr := base64.StdEncoding.EncodeToString(signature)
	timestamp := time.Now().Unix()
	authHeader := fmt.Sprintf("algorithm=RSA2,time=%d,signature=%s", timestamp, signatureStr)
	
	// 创建HTTP请求
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", authHeader)
	
	// 处理Webhook
	resp, err := handler.HandleWebhook(req)
	if err != nil {
		t.Logf("Handle webhook error: %v", err)
	}
	
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Webhook response: Result=%s, Code=%s, Message=%s",
		resp.Result, resp.Code, resp.Message)
	
	if resp.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result)
	}
}

// TestWebhookHandler_ParsePaymentResult 测试解析支付结果
func TestWebhookHandler_ParsePaymentResult(t *testing.T) {
	handler := NewWebhookHandler(testConfig)
	
	// 构造支付结果通知
	data := PaymentResultData{
		PaymentRequestID:  "PAY_TEST_123456",
		PaymentID:         "PAY_ID_123456",
		PaymentAmount:     100.0,
		PaymentCurrency:   "USD",
		PaymentStatus:     "SUCCESS",
		PaymentResultCode: "SUCCESS",
		PaymentResultInfo: map[string]interface{}{
			"card_token": "CARD_TOKEN_123456",
		},
		PaymentTime: time.Now().Format(time.RFC3339),
	}
	
	dataBytes, _ := json.Marshal(data)
	
	req := &WebhookRequest{
		AppID:        testConfig.AppID,
		BusinessType: EventAcquiringPaymentResult,
		Timestamp:    time.Now().Unix() * 1000,
		Data:         dataBytes,
	}
	
	// 解析
	result, err := handler.ParsePaymentResult(req)
	if err != nil {
		t.Fatalf("Parse payment result failed: %v", err)
	}
	
	t.Logf("Parsed payment result:")
	t.Logf("  Payment Request ID: %s", result.PaymentRequestID)
	t.Logf("  Payment ID: %s", result.PaymentID)
	t.Logf("  Amount: %.2f %s", result.PaymentAmount, result.PaymentCurrency)
	t.Logf("  Status: %s", result.PaymentStatus)
	t.Logf("  Card Token: %v", result.PaymentResultInfo["card_token"])
	
	if result.PaymentRequestID != data.PaymentRequestID {
		t.Errorf("Expected payment_request_id '%s', got '%s'",
			data.PaymentRequestID, result.PaymentRequestID)
	}
}

// TestWebhookHandler_ParseAuthToken 测试解析授权Token
func TestWebhookHandler_ParseAuthToken(t *testing.T) {
	handler := NewWebhookHandler(testConfig)
	
	// 构造授权Token通知
	data := AuthTokenData{
		AuthState:              "AUTH_STATE_123456",
		AccessToken:            "ACCESS_TOKEN_123456",
		AccessTokenExpiryTime:  time.Now().Add(24 * time.Hour).Format(time.RFC3339),
		RefreshToken:           "REFRESH_TOKEN_123456",
		RefreshTokenExpiryTime: time.Now().Add(30 * 24 * time.Hour).Format(time.RFC3339),
		UserLoginID:            "user@example.com",
		PaymentMethodType:      "ALIPAY_CN",
		AuthClientID:           "CLIENT_123456",
		Status:                 "ACTIVE",
	}
	
	dataBytes, _ := json.Marshal(data)
	
	req := &WebhookRequest{
		AppID:        testConfig.AppID,
		BusinessType: EventAcquiringAuthToken,
		Timestamp:    time.Now().Unix() * 1000,
		Data:         dataBytes,
	}
	
	// 解析
	result, err := handler.ParseAuthToken(req)
	if err != nil {
		t.Fatalf("Parse auth token failed: %v", err)
	}
	
	t.Logf("Parsed auth token:")
	t.Logf("  Auth State: %s", result.AuthState)
	t.Logf("  Access Token: %s", result.AccessToken)
	t.Logf("  Refresh Token: %s", result.RefreshToken)
	t.Logf("  User Login ID: %s", result.UserLoginID)
	t.Logf("  Payment Method Type: %s", result.PaymentMethodType)
	t.Logf("  Status: %s", result.Status)
	
	if result.AuthState != data.AuthState {
		t.Errorf("Expected auth_state '%s', got '%s'", data.AuthState, result.AuthState)
	}
}

// TestWebhookHandler_ParseCardStatusUpdate 测试解析卡状态变更
func TestWebhookHandler_ParseCardStatusUpdate(t *testing.T) {
	handler := NewWebhookHandler(testConfig)
	
	// 构造卡状态变更通知
	data := CardStatusUpdateData{
		CardID:     "CARD_123456",
		Status:     "ACTIVE",
		UpdateTime: time.Now().Format(time.RFC3339),
	}
	
	dataBytes, _ := json.Marshal(data)
	
	req := &WebhookRequest{
		AppID:        testConfig.AppID,
		BusinessType: EventCardStatusUpdate,
		Timestamp:    time.Now().Unix() * 1000,
		Data:         dataBytes,
	}
	
	// 解析
	result, err := handler.ParseCardStatusUpdate(req)
	if err != nil {
		t.Fatalf("Parse card status update failed: %v", err)
	}
	
	t.Logf("Parsed card status update:")
	t.Logf("  Card ID: %s", result.CardID)
	t.Logf("  Status: %s", result.Status)
	t.Logf("  Update Time: %s", result.UpdateTime)
	
	if result.CardID != data.CardID {
		t.Errorf("Expected card_id '%s', got '%s'", data.CardID, result.CardID)
	}
}

// TestWebhookHandler_ParseCardTransaction 测试解析卡交易通知
func TestWebhookHandler_ParseCardTransaction(t *testing.T) {
	handler := NewWebhookHandler(testConfig)
	
	// 构造卡交易通知
	data := CardTransactionData{
		TransactionID:     "TXN_123456",
		CardID:            "CARD_123456",
		TransactionType:   "PURCHASE",
		Amount:            50.0,
		Currency:          "USD",
		Status:            "SUCCESS",
		StatusDescription: "Transaction successful",
		TransactionTime:   time.Now().Format(time.RFC3339),
		MerchantName:      "Test Merchant",
		MerchantCountry:   "US",
	}
	
	dataBytes, _ := json.Marshal(data)
	
	req := &WebhookRequest{
		AppID:        testConfig.AppID,
		BusinessType: EventCardTransaction,
		Timestamp:    time.Now().Unix() * 1000,
		Data:         dataBytes,
	}
	
	// 解析
	result, err := handler.ParseCardTransaction(req)
	if err != nil {
		t.Fatalf("Parse card transaction failed: %v", err)
	}
	
	t.Logf("Parsed card transaction:")
	t.Logf("  Transaction ID: %s", result.TransactionID)
	t.Logf("  Card ID: %s", result.CardID)
	t.Logf("  Type: %s", result.TransactionType)
	t.Logf("  Amount: %.2f %s", result.Amount, result.Currency)
	t.Logf("  Status: %s - %s", result.Status, result.StatusDescription)
	t.Logf("  Merchant: %s (%s)", result.MerchantName, result.MerchantCountry)
	
	if result.TransactionID != data.TransactionID {
		t.Errorf("Expected transaction_id '%s', got '%s'",
			data.TransactionID, result.TransactionID)
	}
}

// TestWebhookHandler_InvalidSignature 测试无效签名
func TestWebhookHandler_InvalidSignature(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	handler := NewWebhookHandler(testConfig)
	
	// 构造测试数据
	webhookData := map[string]interface{}{
		"app_id":        testConfig.AppID,
		"business_type": EventAcquiringPaymentResult,
		"timestamp":     time.Now().Unix() * 1000,
		"data": map[string]interface{}{
			"payment_request_id": "PAY_TEST_123456",
			"payment_status":     "SUCCESS",
		},
	}
	
	bodyBytes, _ := json.Marshal(webhookData)
	
	// 使用错误的签名
	authHeader := "algorithm=RSA2,time=1234567890,signature=INVALID_SIGNATURE"
	
	// 创建HTTP请求
	req := httptest.NewRequest(http.MethodPost, "/webhook", bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", authHeader)
	
	// 处理Webhook（应该失败）
	resp, err := handler.HandleWebhook(req)
	if err == nil {
		t.Error("Expected error for invalid signature, got nil")
	}
	
	if resp != nil && resp.Result == "S" {
		t.Error("Expected failure for invalid signature")
	}
	
	t.Logf("Got expected error for invalid signature: %v", err)
}
