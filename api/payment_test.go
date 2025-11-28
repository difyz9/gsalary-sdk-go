package api

import (
	"fmt"
	"testing"
	"time"
)

// TestPaymentConsult 测试支付咨询
func TestPaymentConsult(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 支付咨询请求
	req := &PaymentConsultRequest{
		MchAppID:                    "test_mch_app_id",
		PaymentCurrency:             "USD",
		PaymentAmount:               10.00,
		SettlementCurrency:          "USD",
		AllowedPaymentMethodRegions: []string{},
		AllowedPaymentMethods:       []string{},
		UserRegion:                  "US",
		EnvTerminalType:             "APP",
		EnvOsType:                   "ANDROID",
		EnvClientIP:                 "192.168.1.1",
	}
	
	t.Logf("Payment consult: Currency=%s, Amount=%.2f", req.PaymentCurrency, req.PaymentAmount)
	
	resp, err := client.Payment.PaymentConsult(req)
	if err != nil {
		t.Logf("Payment consult error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Payment consult success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Available payment methods: %d", len(resp.Data.PaymentOptions))
	
	for i, option := range resp.Data.PaymentOptions {
		t.Logf("Payment option %d:", i+1)
		t.Logf("  Type: %s", option.PaymentMethodType)
		t.Logf("  Currency: %s", option.Currency)
		t.Logf("  Min: %.2f, Max: %.2f", option.Limit.Min, option.Limit.Max)
		t.Logf("  Country: %s", option.Country)
	}
}

// TestCreatePaymentSession 测试创建支付会话（收银台）
func TestCreatePaymentSession(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 创建支付会话请求
	req := &PaymentSessionRequest{
		MchAppID:          "test_mch_app_id",
		PaymentRequestID:  fmt.Sprintf("PAY_%d", time.Now().Unix()),
		PaymentCurrency:   "USD",
		PaymentAmount:     10.00,
		PaymentMethodType: "ALIPAY_CN",
		PaymentRedirectURL: "https://merchant.com/callback",
		Order: OrderInfo{
			ReferenceOrderID:  fmt.Sprintf("ORD_%d", time.Now().Unix()),
			OrderDescription:  "Test order",
			OrderCurrency:     "USD",
			OrderAmount:       10.00,
			OrderBuyerID:      "BUYER_123",
			OrderBuyerEmail:   "buyer@example.com",
		},
		SettlementCurrency: "USD",
		EnvClientIP:        "192.168.1.1",
		ProductScene:       "CHECKOUT_PAYMENT",
	}
	
	t.Logf("Creating payment session: RequestID=%s, Amount=%.2f %s",
		req.PaymentRequestID, req.PaymentAmount, req.PaymentCurrency)
	
	resp, err := client.Payment.CreatePaymentSession(req)
	if err != nil {
		t.Logf("Create payment session error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Create payment session success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Session ID: %s", resp.Data.PaymentSessionID)
	t.Logf("Session Expiry: %s", resp.Data.PaymentSessionExpiryTime)
	t.Logf("Payment URL: %s", resp.Data.NormalURL)
	
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
}

// TestCreateEasySafePaySession 测试创建钱包授权支付会话
func TestCreateEasySafePaySession(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 创建钱包授权支付会话请求
	req := &EasySafePaySessionRequest{
		PaymentSessionRequest: PaymentSessionRequest{
			MchAppID:          "test_mch_app_id",
			PaymentRequestID:  fmt.Sprintf("PAY_%d", time.Now().Unix()),
			PaymentCurrency:   "USD",
			PaymentAmount:     10.00,
			PaymentMethodType: "ALIPAY_CN",
			PaymentRedirectURL: "https://merchant.com/callback",
			Order: OrderInfo{
				ReferenceOrderID:  fmt.Sprintf("ORD_%d", time.Now().Unix()),
				OrderDescription:  "Test order for easy safe pay",
				OrderCurrency:     "USD",
				OrderAmount:       10.00,
			},
			SettlementCurrency: "USD",
			EnvClientIP:        "192.168.1.1",
			ProductScene:       "CHECKOUT_PAYMENT",
		},
		AuthState:   fmt.Sprintf("AUTH_%d", time.Now().Unix()),
		UserLoginID: "user@example.com",
	}
	
	t.Logf("Creating easy safe pay session: RequestID=%s, AuthState=%s",
		req.PaymentRequestID, req.AuthState)
	
	resp, err := client.Payment.CreateEasySafePaySession(req)
	if err != nil {
		t.Logf("Create easy safe pay session error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Create easy safe pay session success!")
	t.Logf("Session ID: %s", resp.Data.PaymentSessionID)
	t.Logf("Payment URL: %s", resp.Data.NormalURL)
}

// TestEasySafePayPay 测试钱包授权支付（第二次支付）
func TestEasySafePayPay(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 钱包授权支付请求（使用access_token）
	req := &EasySafePayRequest{
		MchAppID:          "test_mch_app_id",
		PaymentRequestID:  fmt.Sprintf("PAY_%d", time.Now().Unix()),
		PaymentCurrency:   "USD",
		PaymentAmount:     10.00,
		PaymentMethodID:   "test_access_token_123456", // 从授权通知中获取
		PaymentMethodType: "ALIPAY_CN",
		PaymentRedirectURL: "https://merchant.com/callback",
		Order: OrderInfo{
			ReferenceOrderID:  fmt.Sprintf("ORD_%d", time.Now().Unix()),
			OrderDescription:  "Test order for second payment",
			OrderCurrency:     "USD",
			OrderAmount:       10.00,
		},
		SettlementCurrency: "USD",
		EnvClientIP:        "192.168.1.1",
		EnvTerminalType:    "WEB",
		EnvOsType:          "IOS",
	}
	
	t.Logf("Easy safe pay payment: RequestID=%s, MethodID=%s",
		req.PaymentRequestID, req.PaymentMethodID)
	
	resp, err := client.Payment.EasySafePayPay(req)
	if err != nil {
		t.Logf("Easy safe pay payment error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Easy safe pay payment success!")
	t.Logf("Payment ID: %s", resp.Data.PaymentID)
	t.Logf("Payment URL: %s", resp.Data.NormalURL)
}

// TestCreateCardAutoDebitSession 测试创建卡授权支付会话
func TestCreateCardAutoDebitSession(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 创建卡授权支付会话请求
	req := &PaymentSessionRequest{
		MchAppID:          "test_mch_app_id",
		PaymentRequestID:  fmt.Sprintf("PAY_%d", time.Now().Unix()),
		PaymentCurrency:   "USD",
		PaymentAmount:     10.00,
		PaymentMethodType: "CARD",
		PaymentRedirectURL: "https://merchant.com/callback",
		Order: OrderInfo{
			ReferenceOrderID:  fmt.Sprintf("ORD_%d", time.Now().Unix()),
			OrderDescription:  "Test order for card auto debit",
			OrderCurrency:     "USD",
			OrderAmount:       10.00,
		},
		SettlementCurrency: "USD",
		EnvClientIP:        "192.168.1.1",
		ProductScene:       "CHECKOUT_PAYMENT",
	}
	
	t.Logf("Creating card auto debit session: RequestID=%s", req.PaymentRequestID)
	
	resp, err := client.Payment.CreateCardAutoDebitSession(req)
	if err != nil {
		t.Logf("Create card auto debit session error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Create card auto debit session success!")
	t.Logf("Session ID: %s", resp.Data.PaymentSessionID)
}

// TestCardAutoDebitPay 测试卡授权支付（第二次支付）
func TestCardAutoDebitPay(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 卡授权支付请求（使用card_token）
	req := &EasySafePayRequest{
		MchAppID:          "test_mch_app_id",
		PaymentRequestID:  fmt.Sprintf("PAY_%d", time.Now().Unix()),
		PaymentCurrency:   "USD",
		PaymentAmount:     10.00,
		PaymentMethodID:   "test_card_token_123456", // 从支付结果通知中获取
		PaymentMethodType: "CARD",
		PaymentRedirectURL: "https://merchant.com/callback",
		Order: OrderInfo{
			ReferenceOrderID:  fmt.Sprintf("ORD_%d", time.Now().Unix()),
			OrderDescription:  "Test order for card second payment",
			OrderCurrency:     "USD",
			OrderAmount:       10.00,
		},
		SettlementCurrency: "USD",
		EnvClientIP:        "192.168.1.1",
		EnvTerminalType:    "WEB",
	}
	
	t.Logf("Card auto debit payment: RequestID=%s, CardToken=%s",
		req.PaymentRequestID, req.PaymentMethodID)
	
	resp, err := client.Payment.CardAutoDebitPay(req)
	if err != nil {
		t.Logf("Card auto debit payment error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Card auto debit payment success!")
	t.Logf("Payment ID: %s", resp.Data.PaymentID)
}

// TestRefreshAuthToken 测试刷新授权令牌
func TestRefreshAuthToken(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 刷新令牌请求
	req := &RefreshTokenRequest{
		MchAppID:       "test_mch_app_id",
		RefreshToken:   "test_refresh_token_123456",
		MerchantRegion: "US",
	}
	
	t.Logf("Refreshing auth token: RefreshToken=%s", req.RefreshToken)
	
	resp, err := client.Payment.RefreshAuthToken(req)
	if err != nil {
		t.Logf("Refresh auth token error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Refresh auth token success!")
	t.Logf("New Access Token: %s", resp.Data.AccessToken)
	t.Logf("Access Token Expiry: %s", resp.Data.AccessTokenExpiryTime)
	t.Logf("New Refresh Token: %s", resp.Data.RefreshToken)
}

// TestRevokeAuthToken 测试取消授权令牌
func TestRevokeAuthToken(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 取消授权请求
	req := &RevokeTokenRequest{
		MchAppID:    "test_mch_app_id",
		AccessToken: "test_access_token_123456",
	}
	
	t.Logf("Revoking auth token: AccessToken=%s", req.AccessToken)
	
	resp, err := client.Payment.RevokeAuthToken(req)
	if err != nil {
		t.Logf("Revoke auth token error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Revoke auth token success!")
	t.Logf("Result: %s", resp.Result.Result)
}

// TestCancelPayment 测试取消支付
func TestCancelPayment(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 取消支付请求
	req := &CancelPaymentRequest{
		MchAppID:         "test_mch_app_id",
		PaymentRequestID: "PAY_TEST_123456",
		PaymentID:        "PAY_ID_123456",
	}
	
	t.Logf("Cancelling payment: PaymentRequestID=%s, PaymentID=%s",
		req.PaymentRequestID, req.PaymentID)
	
	resp, err := client.Payment.CancelPayment(req)
	if err != nil {
		t.Logf("Cancel payment error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Cancel payment success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Cancel Time: %s", resp.Data.CancelTime)
}
