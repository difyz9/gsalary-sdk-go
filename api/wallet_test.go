package api

import (
	"testing"
)

// TestGetWalletBalance 测试查询钱包余额
func TestGetWalletBalance(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询USD钱包余额
	req := &WalletBalanceRequest{
		Currency: "USD",
	}
	
	resp, err := client.Wallet.GetBalance(req)
	if err != nil {
		t.Logf("Get wallet balance error: %v", err)
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
	
	t.Logf("Get wallet balance success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Currency: %s", resp.Data.Currency)
	t.Logf("Amount: %.2f", resp.Data.Amount)
	t.Logf("Share Card Account Balance: %.2f", resp.Data.ShareCardAccountBalance)
	t.Logf("Available: %.2f", resp.Data.Available)
	t.Logf("Account Type: %s", resp.Data.AccountType)
	t.Logf("Query Time: %s", resp.Data.QueryTime)
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
	
	if resp.Data.Currency != "USD" {
		t.Errorf("Expected currency 'USD', got '%s'", resp.Data.Currency)
	}
}

// TestGetWalletBalanceMultipleCurrencies 测试查询多种货币余额
func TestGetWalletBalanceMultipleCurrencies(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 测试多种货币
	currencies := []string{"USD", "EUR", "GBP"}
	
	for _, currency := range currencies {
		t.Run("Currency_"+currency, func(t *testing.T) {
			req := &WalletBalanceRequest{
				Currency: currency,
			}
			
			resp, err := client.Wallet.GetBalance(req)
			if err != nil {
				t.Logf("Get wallet balance for %s error: %v", currency, err)
				if resp != nil {
					t.Logf("Response: Result=%s, Code=%s, Message=%s",
						resp.Result.Result, resp.Result.Code, resp.Result.Message)
				}
				return
			}
			
			t.Logf("Get wallet balance for %s success!", currency)
			t.Logf("Available: %.2f %s", resp.Data.Available, resp.Data.Currency)
		})
	}
}

// TestGetWalletBalanceInvalidCurrency 测试无效货币
func TestGetWalletBalanceInvalidCurrency(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 使用无效的货币代码
	req := &WalletBalanceRequest{
		Currency: "INVALID",
	}
	
	resp, err := client.Wallet.GetBalance(req)
	if err != nil {
		t.Logf("Got expected error for invalid currency: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
	} else {
		t.Log("Query succeeded (API may accept any currency or return zero balance)")
	}
}

// TestGetWalletBalanceEmpty 测试空货币参数
func TestGetWalletBalanceEmpty(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 不指定货币
	req := &WalletBalanceRequest{
		Currency: "",
	}
	
	resp, err := client.Wallet.GetBalance(req)
	if err != nil {
		t.Logf("Get wallet balance without currency error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
	} else {
		t.Log("Query succeeded (API may return default currency or error)")
		if resp != nil {
			t.Logf("Currency: %s, Available: %.2f", resp.Data.Currency, resp.Data.Available)
		}
	}
}
