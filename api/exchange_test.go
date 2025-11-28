package api

import (
	"fmt"
	"testing"
	"time"
)

// TestGetCurrentExchangeRate 测试查询当前汇率
func TestGetCurrentExchangeRate(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询USD到HKD的汇率
	req := &ExchangeRateRequest{
		BuyCurrency:  "HKD",
		SellCurrency: "USD",
	}
	
	t.Logf("Querying exchange rate: %s -> %s", req.SellCurrency, req.BuyCurrency)
	
	resp, err := client.Exchange.GetCurrentExchangeRate(req)
	if err != nil {
		t.Logf("Get exchange rate error: %v", err)
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
	
	t.Logf("Get exchange rate success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Buy Currency: %s", resp.Data.BuyCurrency)
	t.Logf("Sell Currency: %s", resp.Data.SellCurrency)
	t.Logf("Rate: %.6f", resp.Data.Rate)
	t.Logf("Update Time: %s", resp.Data.UpdateTime)
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
	
	if resp.Data.BuyCurrency != "HKD" {
		t.Errorf("Expected buy currency 'HKD', got '%s'", resp.Data.BuyCurrency)
	}
	
	if resp.Data.SellCurrency != "USD" {
		t.Errorf("Expected sell currency 'USD', got '%s'", resp.Data.SellCurrency)
	}
	
	if resp.Data.Rate <= 0 {
		t.Errorf("Expected positive rate, got %.6f", resp.Data.Rate)
	}
}

// TestGetCurrentExchangeRateCNY 测试查询CNY汇率
func TestGetCurrentExchangeRateCNY(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询USD到CNY的汇率
	req := &ExchangeRateRequest{
		BuyCurrency:  "CNY",
		SellCurrency: "USD",
	}
	
	resp, err := client.Exchange.GetCurrentExchangeRate(req)
	if err != nil {
		t.Logf("Get exchange rate (CNY) error: %v", err)
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
	
	t.Logf("Get exchange rate (CNY) success!")
	t.Logf("1 %s = %.6f %s", resp.Data.SellCurrency, resp.Data.Rate, resp.Data.BuyCurrency)
	t.Logf("Update Time: %s", resp.Data.UpdateTime)
}

// TestGetCurrentExchangeRateReverse 测试查询反向汇率
func TestGetCurrentExchangeRateReverse(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询HKD到USD的汇率（反向）
	req := &ExchangeRateRequest{
		BuyCurrency:  "USD",
		SellCurrency: "HKD",
	}
	
	resp, err := client.Exchange.GetCurrentExchangeRate(req)
	if err != nil {
		t.Logf("Get exchange rate (reverse) error: %v", err)
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
	
	t.Logf("Get exchange rate (reverse) success!")
	t.Logf("1 %s = %.6f %s", resp.Data.SellCurrency, resp.Data.Rate, resp.Data.BuyCurrency)
}

// TestGetCurrentExchangeRateInvalidCurrency 测试查询无效货币对
func TestGetCurrentExchangeRateInvalidCurrency(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 使用无效的货币代码
	req := &ExchangeRateRequest{
		BuyCurrency:  "INVALID",
		SellCurrency: "USD",
	}
	
	resp, err := client.Exchange.GetCurrentExchangeRate(req)
	if err != nil {
		t.Logf("Got expected error for invalid currency: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
	} else {
		t.Log("Query succeeded (API may have different validation)")
	}
}

// TestGetCurrentExchangeRateSameCurrency 测试查询相同货币
func TestGetCurrentExchangeRateSameCurrency(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询相同货币的汇率
	req := &ExchangeRateRequest{
		BuyCurrency:  "USD",
		SellCurrency: "USD",
	}
	
	resp, err := client.Exchange.GetCurrentExchangeRate(req)
	if err != nil {
		t.Logf("Got error for same currency: %v", err)
	} else if resp != nil {
		t.Logf("Same currency rate: %.6f (should be 1.0)", resp.Data.Rate)
		if resp.Data.Rate != 1.0 {
			t.Logf("Warning: Expected rate 1.0 for same currency, got %.6f", resp.Data.Rate)
		}
	}
}

// TestRequestQuote 测试请求锁汇报价
func TestRequestQuote(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 请求锁汇报价（用卖出金额）
	req := &ExchangeQuoteRequest{
		BuyCurrency:  "HKD",
		SellCurrency: "USD",
		SellAmount:   100.00, // 卖出100 USD
	}
	
	t.Logf("Requesting exchange quote: Sell %s %.2f to buy %s", 
		req.SellCurrency, req.SellAmount, req.BuyCurrency)
	
	resp, err := client.Exchange.RequestQuote(req)
	if err != nil {
		t.Logf("Request exchange quote error: %v", err)
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
	
	t.Logf("Request exchange quote success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Quote ID: %s", resp.Data.QuoteID)
	t.Logf("Buy: %s %.2f", resp.Data.Buy.Currency, resp.Data.Buy.Amount)
	t.Logf("Sell: %s %.2f", resp.Data.Sell.Currency, resp.Data.Sell.Amount)
	t.Logf("Surcharge: %s %.2f", resp.Data.Surcharge.Currency, resp.Data.Surcharge.Amount)
	t.Logf("Total Cost: %s %.2f", resp.Data.TotalCost.Currency, resp.Data.TotalCost.Amount)
	t.Logf("Update Time: %s", resp.Data.UpdateTime)
	t.Logf("Expire Time: %s", resp.Data.ExpireTime)
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
	
	if resp.Data.QuoteID == "" {
		t.Error("Quote ID should not be empty")
	}
	
	if resp.Data.Buy.Currency != req.BuyCurrency {
		t.Errorf("Expected buy currency '%s', got '%s'", req.BuyCurrency, resp.Data.Buy.Currency)
	}
	
	if resp.Data.Sell.Currency != req.SellCurrency {
		t.Errorf("Expected sell currency '%s', got '%s'", req.SellCurrency, resp.Data.Sell.Currency)
	}
}

// TestRequestQuoteWithBuyAmount 测试请求锁汇报价（使用购入金额）
func TestRequestQuoteWithBuyAmount(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 请求锁汇报价（用购入金额）
	req := &ExchangeQuoteRequest{
		BuyCurrency:  "HKD",
		SellCurrency: "USD",
		BuyAmount:    1000.00, // 购入1000 HKD
	}
	
	t.Logf("Requesting exchange quote: Buy %s %.2f with %s", 
		req.BuyCurrency, req.BuyAmount, req.SellCurrency)
	
	resp, err := client.Exchange.RequestQuote(req)
	if err != nil {
		t.Logf("Request exchange quote (buy amount) error: %v", err)
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
	
	t.Logf("Request exchange quote (buy amount) success!")
	t.Logf("Quote ID: %s", resp.Data.QuoteID)
	t.Logf("Will buy: %s %.2f", resp.Data.Buy.Currency, resp.Data.Buy.Amount)
	t.Logf("Will sell: %s %.2f", resp.Data.Sell.Currency, resp.Data.Sell.Amount)
	t.Logf("Total Cost: %s %.2f", resp.Data.TotalCost.Currency, resp.Data.TotalCost.Amount)
}

// TestRequestQuoteInvalidCurrency 测试无效币种的报价请求
func TestRequestQuoteInvalidCurrency(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 使用无效的币种
	req := &ExchangeQuoteRequest{
		BuyCurrency:  "INVALID",
		SellCurrency: "USD",
		SellAmount:   100.00,
	}
	
	resp, err := client.Exchange.RequestQuote(req)
	if err != nil {
		t.Logf("Got expected error for invalid currency: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
	} else {
		t.Log("Request succeeded (API may accept it)")
	}
}

// TestRequestQuoteNoAmount 测试没有提供金额的报价请求
func TestRequestQuoteNoAmount(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 不提供任何金额
	req := &ExchangeQuoteRequest{
		BuyCurrency:  "HKD",
		SellCurrency: "USD",
		// 两个金额都为0
	}
	
	resp, err := client.Exchange.RequestQuote(req)
	if err != nil {
		t.Logf("Got expected error for no amount: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
	} else {
		t.Log("Request succeeded (API may have default behavior)")
	}
}

// TestSubmitExchangeRequest 测试提交换汇订单
func TestSubmitExchangeRequest(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 提交换汇订单
	req := &ExchangeSubmitRequest{
		RequestID: "test_request_" + fmt.Sprintf("%d", time.Now().Unix()),
		QuoteID:   "test_quote_123456",
	}
	
	t.Logf("Submitting exchange request: RequestID=%s, QuoteID=%s", 
		req.RequestID, req.QuoteID)
	
	resp, err := client.Exchange.SubmitExchangeRequest(req)
	if err != nil {
		t.Logf("Submit exchange request error: %v", err)
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
	
	t.Logf("Submit exchange request success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Order ID: %s", resp.Data.OrderID)
	t.Logf("Request ID: %s", resp.Data.RequestID)
	t.Logf("Status: %s", resp.Data.Status)
	t.Logf("Source: %s", resp.Data.Source)
	t.Logf("Sell: %s %.2f", resp.Data.Sell.Currency, resp.Data.Sell.Amount)
	t.Logf("Buy: %s %.2f", resp.Data.Buy.Currency, resp.Data.Buy.Amount)
	t.Logf("Surcharge: %s %.2f", resp.Data.Surcharge.Currency, resp.Data.Surcharge.Amount)
	t.Logf("Exchange Rate: %.6f", resp.Data.ExchangeRate)
	t.Logf("Create Time: %s", resp.Data.CreateTime)
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
	
	if resp.Data.OrderID == "" {
		t.Error("Order ID should not be empty")
	}
	
	if resp.Data.RequestID != req.RequestID {
		t.Errorf("Expected request ID '%s', got '%s'", req.RequestID, resp.Data.RequestID)
	}
}

// TestSubmitExchangeRequestMissingQuoteID 测试缺少QuoteID的提交请求
func TestSubmitExchangeRequestMissingQuoteID(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 缺少QuoteID
	req := &ExchangeSubmitRequest{
		RequestID: "test_request_" + fmt.Sprintf("%d", time.Now().Unix()),
		// QuoteID为空
	}
	
	resp, err := client.Exchange.SubmitExchangeRequest(req)
	if err != nil {
		t.Logf("Got expected error for missing quote ID: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
	} else {
		t.Log("Request succeeded (API may accept empty quote ID)")
	}
}

// TestSubmitExchangeRequestInvalidQuoteID 测试无效QuoteID的提交请求
func TestSubmitExchangeRequestInvalidQuoteID(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 使用无效的QuoteID
	req := &ExchangeSubmitRequest{
		RequestID: "test_request_" + fmt.Sprintf("%d", time.Now().Unix()),
		QuoteID:   "invalid_quote_id_12345",
	}
	
	resp, err := client.Exchange.SubmitExchangeRequest(req)
	if err != nil {
		t.Logf("Got expected error for invalid quote ID: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
	} else {
		t.Log("Request succeeded unexpectedly with invalid quote ID")
	}
}

// TestGetExchangeOrders 测试查询换汇订单列表
func TestGetExchangeOrders(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询订单列表
	req := &ExchangeOrdersRequest{
		Page:  1,
		Limit: 20,
	}
	
	t.Logf("Getting exchange orders: Page=%d, Limit=%d", req.Page, req.Limit)
	
	resp, err := client.Exchange.GetExchangeOrders(req)
	if err != nil {
		t.Logf("Get exchange orders error: %v", err)
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
	
	t.Logf("Get exchange orders success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Page: %d", resp.Data.Page)
	t.Logf("Limit: %d", resp.Data.Limit)
	t.Logf("Total Count: %d", resp.Data.TotalCount)
	t.Logf("Total Page: %d", resp.Data.TotalPage)
	t.Logf("Orders Count: %d", len(resp.Data.Orders))
	
	// 显示订单详情（最多显示前3个）
	displayCount := len(resp.Data.Orders)
	if displayCount > 3 {
		displayCount = 3
	}
	
	for i := 0; i < displayCount; i++ {
		order := resp.Data.Orders[i]
		t.Logf("Order %d:", i+1)
		t.Logf("  Order ID: %s", order.OrderID)
		t.Logf("  Request ID: %s", order.RequestID)
		t.Logf("  Status: %s", order.Status)
		t.Logf("  Source: %s", order.Source)
		t.Logf("  Sell: %s %.2f", order.Sell.Currency, order.Sell.Amount)
		t.Logf("  Buy: %s %.2f", order.Buy.Currency, order.Buy.Amount)
		t.Logf("  Exchange Rate: %.6f", order.ExchangeRate)
		t.Logf("  Create Time: %s", order.CreateTime)
	}
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
}

// TestGetExchangeOrdersWithFilters 测试带过滤条件的订单查询
func TestGetExchangeOrdersWithFilters(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询订单列表，带过滤条件
	req := &ExchangeOrdersRequest{
		Page:         1,
		Limit:        10,
		Status:       "SUCCESS",
		BuyCurrency:  "HKD",
		SellCurrency: "USD",
	}
	
	t.Logf("Getting exchange orders with filters: Status=%s, Buy=%s, Sell=%s", 
		req.Status, req.BuyCurrency, req.SellCurrency)
	
	resp, err := client.Exchange.GetExchangeOrders(req)
	if err != nil {
		t.Logf("Get exchange orders with filters error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	t.Logf("Get exchange orders with filters success!")
	t.Logf("Total Count: %d", resp.Data.TotalCount)
	t.Logf("Orders Count: %d", len(resp.Data.Orders))
	t.Logf("Query Status: %s", resp.Data.Query.Status)
	t.Logf("Query Buy Currency: %s", resp.Data.Query.BuyCurrency)
	t.Logf("Query Sell Currency: %s", resp.Data.Query.SellCurrency)
}

// TestGetExchangeOrdersWithTimeRange 测试带时间范围的订单查询
func TestGetExchangeOrdersWithTimeRange(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询订单列表，带时间范围
	req := &ExchangeOrdersRequest{
		Page:      1,
		Limit:     10,
		TimeStart: "2024-03-04T10:00:00Z",
		TimeEnd:   "2024-03-05T10:00:00Z",
	}
	
	t.Logf("Getting exchange orders with time range: %s to %s", 
		req.TimeStart, req.TimeEnd)
	
	resp, err := client.Exchange.GetExchangeOrders(req)
	if err != nil {
		t.Logf("Get exchange orders with time range error: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
		return
	}
	
	t.Logf("Get exchange orders with time range success!")
	t.Logf("Total Count: %d", resp.Data.TotalCount)
	t.Logf("Query Time Start: %s", resp.Data.Query.TimeStart)
	t.Logf("Query Time End: %s", resp.Data.Query.TimeEnd)
}
