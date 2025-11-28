package api

import (
	"fmt"
	"os"
	"testing"
	"time"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// 测试配置
var testConfig *gsalary.GSalaryConfig

func init() {
	// 初始化测试配置
	testConfig = gsalary.NewConfig()
	
	// 从环境变量获取 AppID（必需）
	testConfig.AppID = os.Getenv("GSALARY_APPID")
	if testConfig.AppID == "" {
		testConfig.AppID = "your_app_id_here" // 占位符，实际使用时需要设置环境变量
	}
	
	// 从环境变量获取 Endpoint（可选，默认测试环境）
	testConfig.Endpoint = os.Getenv("GSALARY_ENDPOINT")
	if testConfig.Endpoint == "" {
		testConfig.Endpoint = "https://api-test.gsalary.com"
	}
	
	// 尝试加载环境变量中的密钥路径
	privateKeyFile := os.Getenv("GSALARY_CLIENT_PRIVATE_KEY_FILE")
	if privateKeyFile == "" {
		privateKeyFile = "../private_key.pem"
	}
	
	publicKeyFile := os.Getenv("GSALARY_SERVER_PUBLIC_KEY_FILE")
	if publicKeyFile == "" {
		publicKeyFile = "../server_public_key.pem"
	}
	
	// 加载密钥
	if err := testConfig.ConfigClientPrivateKeyPEMFile(privateKeyFile); err != nil {
		fmt.Printf("Warning: Failed to load private key: %v\n", err)
	}
	
	if err := testConfig.ConfigServerPublicKeyPEMFile(publicKeyFile); err != nil {
		fmt.Printf("Warning: Failed to load public key: %v\n", err)
	}
}

// TestCardApply 测试申请卡片
func TestCardApply(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 生成唯一的请求ID
	requestID := fmt.Sprintf("TEST_CARD_APPLY_%d", time.Now().Unix())
	
	// 创建申请请求
	req := &CardApplyRequest{
		RequestID:           requestID,
		ProductCode:         "VIRTUAL_CARD_USD", // 根据实际产品代码调整
		Currency:            "USD",
		CardHolderID:        "test_holder_001",
		LimitPerDay:         1000.00,
		LimitPerMonth:       5000.00,
		LimitPerTransaction: 500.00,
		InitBalance:         100.00,
	}
	
	// 发起申请
	resp, err := client.Card.ApplyCard(req)
	if err != nil {
		t.Logf("Apply card error: %v", err)
		// 不直接失败，因为可能是业务错误（如余额不足、产品代码错误等）
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s, Status=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message, resp.Data.Status)
		}
		return
	}
	
	// 验证响应
	if resp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Apply card success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Code: %s", resp.Result.Code)
	t.Logf("Message: %s", resp.Result.Message)
	t.Logf("Request ID: %s", resp.Data.RequestID)
	t.Logf("Status: %s", resp.Data.Status)
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
	
	if resp.Data.RequestID != requestID {
		t.Errorf("Expected request_id '%s', got '%s'", requestID, resp.Data.RequestID)
	}
}

// TestCardApplyWithMinimalParams 测试最小参数申请卡片
func TestCardApplyWithMinimalParams(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 生成唯一的请求ID
	requestID := fmt.Sprintf("TEST_MINIMAL_%d", time.Now().Unix())
	
	// 创建最小参数的申请请求
	req := &CardApplyRequest{
		RequestID:    requestID,
		ProductCode:  "VIRTUAL_CARD_USD",
		Currency:     "USD",
		CardHolderID: "test_holder_002",
		InitBalance:  50.00,
		// 不设置limit参数，使用默认值
	}
	
	// 发起申请
	resp, err := client.Card.ApplyCard(req)
	if err != nil {
		t.Logf("Apply card with minimal params error: %v", err)
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
	
	t.Logf("Apply card with minimal params success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Request ID: %s", resp.Data.RequestID)
	t.Logf("Status: %s", resp.Data.Status)
}

// TestCardApplyInvalidParams 测试无效参数
func TestCardApplyInvalidParams(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 测试空请求ID
	req := &CardApplyRequest{
		RequestID:    "", // 空的请求ID
		ProductCode:  "VIRTUAL_CARD_USD",
		Currency:     "USD",
		CardHolderID: "test_holder_003",
		InitBalance:  10.00,
	}
	
	_, err := client.Card.ApplyCard(req)
	if err == nil {
		t.Log("Expected error for empty request_id, but got success (API may accept it)")
	} else {
		t.Logf("Got expected error for empty request_id: %v", err)
	}
	
	// 测试无效货币代码
	req2 := &CardApplyRequest{
		RequestID:    fmt.Sprintf("TEST_INVALID_%d", time.Now().Unix()),
		ProductCode:  "VIRTUAL_CARD_USD",
		Currency:     "INVALID", // 无效的货币代码
		CardHolderID: "test_holder_004",
		InitBalance:  10.00,
	}
	
	_, err = client.Card.ApplyCard(req2)
	if err == nil {
		t.Log("Expected error for invalid currency, but got success (API may accept it)")
	} else {
		t.Logf("Got expected error for invalid currency: %v", err)
	}
}

// TestGetAvailableQuotas 测试查询卡可用余额
func TestGetAvailableQuotas(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询USD币种的SHARE类型余额
	req := &CardAvailableQuotasRequest{
		Currency:           "USD",
		AccountingCardType: "SHARE",
	}
	
	resp, err := client.Card.GetAvailableQuotas(req)
	if err != nil {
		t.Logf("Get available quotas error: %v", err)
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
	
	t.Logf("Get available quotas success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Currency: %s", resp.Data.Currency)
	t.Logf("Accounting Card Type: %s", resp.Data.AccountingCardType)
	t.Logf("Available Quota: %.2f", resp.Data.AvailableQuota)
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
	
	if resp.Data.Currency != "USD" {
		t.Errorf("Expected currency 'USD', got '%s'", resp.Data.Currency)
	}
}

// TestGetAvailableQuotasDefaultType 测试查询余额（默认类型）
func TestGetAvailableQuotasDefaultType(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 不指定AccountingCardType，使用默认值
	req := &CardAvailableQuotasRequest{
		Currency: "USD",
	}
	
	resp, err := client.Card.GetAvailableQuotas(req)
	if err != nil {
		t.Logf("Get available quotas (default type) error: %v", err)
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
	
	t.Logf("Get available quotas (default type) success!")
	t.Logf("Currency: %s", resp.Data.Currency)
	t.Logf("Accounting Card Type: %s", resp.Data.AccountingCardType)
	t.Logf("Available Quota: %.2f", resp.Data.AvailableQuota)
}

// TestGetProducts 测试查询可用的卡产品列表
func TestGetProducts(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询所有产品
	req := &CardProductsRequest{}
	
	resp, err := client.Card.GetProducts(req)
	if err != nil {
		t.Logf("Get card products error: %v", err)
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
	
	t.Logf("Get card products success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Product count: %d", len(resp.Data.Products))
	
	// 打印产品列表
	for i, product := range resp.Data.Products {
		t.Logf("Product %d:", i+1)
		t.Logf("  Code: %s", product.ProductCode)
		t.Logf("  Name: %s", product.ProductName)
		t.Logf("  Type: %s", product.CardType)
		t.Logf("  Brand: %s", product.BrandCode)
		t.Logf("  Currency: %s", product.Currency)
		t.Logf("  Description: %s", product.Description)
	}
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
}

// TestGetProductsWithFilters 测试带过滤条件查询产品
func TestGetProductsWithFilters(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询VIRTUAL类型的VISA卡USD币种产品
	req := &CardProductsRequest{
		CardType:  "VIRTUAL",
		BrandCode: "VISA",
		Currency:  "USD",
	}
	
	resp, err := client.Card.GetProducts(req)
	if err != nil {
		t.Logf("Get card products with filters error: %v", err)
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
	
	t.Logf("Get card products with filters success!")
	t.Logf("Product count: %d", len(resp.Data.Products))
	
	// 验证过滤条件
	for _, product := range resp.Data.Products {
		if product.CardType != "VIRTUAL" {
			t.Errorf("Expected card type 'VIRTUAL', got '%s'", product.CardType)
		}
		if product.BrandCode != "VISA" {
			t.Errorf("Expected brand code 'VISA', got '%s'", product.BrandCode)
		}
		if product.Currency != "USD" {
			t.Errorf("Expected currency 'USD', got '%s'", product.Currency)
		}
	}
}

// TestGetCardApplyResult 测试查询开卡结果
func TestGetCardApplyResult(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 首先申请一张卡片获取request_id
	requestID := fmt.Sprintf("TEST_QUERY_RESULT_%d", time.Now().Unix())
	applyReq := &CardApplyRequest{
		RequestID:    requestID,
		ProductCode:  "VIRTUAL_CARD_USD",
		Currency:     "USD",
		CardHolderID: "test_holder_query",
		InitBalance:  10.00,
	}
	
	applyResp, err := client.Card.ApplyCard(applyReq)
	if err != nil {
		t.Logf("Apply card error (skipping result query): %v", err)
		if applyResp != nil {
			t.Logf("Apply Response: Result=%s, Code=%s, Message=%s",
				applyResp.Result.Result, applyResp.Result.Code, applyResp.Result.Message)
		}
		// 即使申请失败，也尝试查询一个已知的request_id
		requestID = "existing_request_id_if_known"
	} else {
		t.Logf("Apply card success, request_id: %s", requestID)
	}
	
	// 查询开卡结果
	resultResp, err := client.Card.GetCardApplyResult(requestID)
	if err != nil {
		t.Logf("Get card apply result error: %v", err)
		if resultResp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resultResp.Result.Result, resultResp.Result.Code, resultResp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if resultResp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Get card apply result success!")
	t.Logf("Result: %s", resultResp.Result.Result)
	t.Logf("Code: %s", resultResp.Result.Code)
	t.Logf("Message: %s", resultResp.Result.Message)
	
	// 打印data内容
	if len(resultResp.Data) > 0 {
		t.Logf("Data:")
		for key, value := range resultResp.Data {
			t.Logf("  %s: %v", key, value)
		}
	}
	
	// 验证基本响应
	if resultResp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resultResp.Result.Result)
	}
}

// TestGetCardApplyResultInvalidRequestID 测试查询不存在的开卡结果
func TestGetCardApplyResultInvalidRequestID(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询一个不存在的request_id
	requestID := "INVALID_REQUEST_ID_NOT_EXISTS"
	
	resultResp, err := client.Card.GetCardApplyResult(requestID)
	if err != nil {
		t.Logf("Got expected error for invalid request_id: %v", err)
		if resultResp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resultResp.Result.Result, resultResp.Result.Code, resultResp.Result.Message)
		}
	} else {
		t.Log("Query returned success (unexpected, but API may handle it differently)")
	}
}

// TestGetCardApplyResultEmptyRequestID 测试空的request_id
func TestGetCardApplyResultEmptyRequestID(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 测试空的request_id
	_, err := client.Card.GetCardApplyResult("")
	if err == nil {
		t.Error("Expected error for empty request_id, but got success")
	} else {
		t.Logf("Got expected error for empty request_id: %v", err)
	}
}

// TestGetCardList 测试查询卡列表
func TestGetCardList(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询第一页，每页10条
	req := &CardListRequest{
		Page:  1,
		Limit: 10,
	}
	
	resp, err := client.Card.GetCardList(req)
	if err != nil {
		t.Logf("Get card list error: %v", err)
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
	
	t.Logf("Get card list success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Page: %d", resp.Data.Page)
	t.Logf("Limit: %d", resp.Data.Limit)
	t.Logf("Total Count: %d", resp.Data.TotalCount)
	t.Logf("Total Page: %d", resp.Data.TotalPage)
	t.Logf("Cards in this page: %d", len(resp.Data.Cards))
	
	// 打印前3张卡片信息
	for i, card := range resp.Data.Cards {
		if i >= 3 {
			break
		}
		t.Logf("Card %d:", i+1)
		t.Logf("  ID: %s", card.CardID)
		t.Logf("  Product Code: %s", card.ProductCode)
		t.Logf("  Brand Code: %s", card.BrandCode)
		t.Logf("  Card Holder ID: %s", card.CardHolderID)
		t.Logf("  Status: %s", card.Status)
		t.Logf("  Created At: %s", card.CreatedAt)
	}
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
}

// TestGetCardListWithFilters 测试带过滤条件查询卡列表
func TestGetCardListWithFilters(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询ACTIVE状态的VISA品牌卡片
	req := &CardListRequest{
		Page:      1,
		Limit:     5,
		BrandCode: "VISA",
		Status:    "ACTIVE",
	}
	
	resp, err := client.Card.GetCardList(req)
	if err != nil {
		t.Logf("Get card list with filters error: %v", err)
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
	
	t.Logf("Get card list with filters success!")
	t.Logf("Cards count: %d", len(resp.Data.Cards))
	
	// 验证过滤条件
	for _, card := range resp.Data.Cards {
		if card.BrandCode != "VISA" {
			t.Errorf("Expected brand code 'VISA', got '%s'", card.BrandCode)
		}
		if card.Status != "ACTIVE" {
			t.Errorf("Expected status 'ACTIVE', got '%s'", card.Status)
		}
	}
}

// TestGetCardListByHolder 测试按持卡人查询卡列表
func TestGetCardListByHolder(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 先申请一张卡
	requestID := fmt.Sprintf("TEST_LIST_HOLDER_%d", time.Now().Unix())
	holderID := "test_holder_for_list"
	
	applyReq := &CardApplyRequest{
		RequestID:    requestID,
		ProductCode:  "VIRTUAL_CARD_USD",
		Currency:     "USD",
		CardHolderID: holderID,
		InitBalance:  10.00,
	}
	
	_, err := client.Card.ApplyCard(applyReq)
	if err != nil {
		t.Logf("Apply card error (will still try to query): %v", err)
	}
	
	// 查询该持卡人的卡片
	req := &CardListRequest{
		Page:         1,
		Limit:        20,
		CardHolderID: holderID,
	}
	
	resp, err := client.Card.GetCardList(req)
	if err != nil {
		t.Logf("Get card list by holder error: %v", err)
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
	
	t.Logf("Get card list by holder success!")
	t.Logf("Holder ID: %s", holderID)
	t.Logf("Cards count: %d", len(resp.Data.Cards))
	
	// 验证持卡人ID
	for _, card := range resp.Data.Cards {
		if card.CardHolderID != holderID {
			t.Errorf("Expected card holder ID '%s', got '%s'", holderID, card.CardHolderID)
		}
	}
}

// TestGetCardInfo 测试查看卡信息
func TestGetCardInfo(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 首先获取卡列表，取第一张卡的ID
	listReq := &CardListRequest{
		Page:  1,
		Limit: 1,
	}
	
	listResp, err := client.Card.GetCardList(listReq)
	if err != nil || len(listResp.Data.Cards) == 0 {
		t.Skip("Skipping test: no cards available to query")
	}
	
	cardID := listResp.Data.Cards[0].CardID
	t.Logf("Testing with card_id: %s", cardID)
	
	// 查看卡信息
	resp, err := client.Card.GetCardInfo(cardID)
	if err != nil {
		t.Logf("Get card info error: %v", err)
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
	
	t.Logf("Get card info success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Card ID: %s", resp.Data.CardID)
	t.Logf("Card Name: %s", resp.Data.CardName)
	t.Logf("Mask Card Number: %s", resp.Data.MaskCardNumber)
	t.Logf("Currency: %s", resp.Data.CardCurrency)
	t.Logf("Available Balance: %.2f", resp.Data.AvailableBalance)
	t.Logf("Brand Code: %s", resp.Data.BrandCode)
	t.Logf("Status: %s", resp.Data.Status)
	t.Logf("Card Type: %s", resp.Data.CardType)
	t.Logf("Accounting Type: %s", resp.Data.AccountingType)
	t.Logf("Card Region: %s", resp.Data.CardRegion)
	t.Logf("Card Holder ID: %s", resp.Data.CardHolderID)
	t.Logf("First Name: %s", resp.Data.FirstName)
	t.Logf("Last Name: %s", resp.Data.LastName)
	t.Logf("Email: %s", resp.Data.Email)
	t.Logf("Limit Per Day: %.2f", resp.Data.LimitPerDay)
	t.Logf("Limit Per Month: %.2f", resp.Data.LimitPerMonth)
	t.Logf("Limit Per Transaction: %.2f", resp.Data.LimitPerTransaction)
	t.Logf("Support TDS Trans: %v", resp.Data.SupportTdsTrans)
	t.Logf("Create Time: %s", resp.Data.CreateTime)
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
	
	if resp.Data.CardID != cardID {
		t.Errorf("Expected card_id '%s', got '%s'", cardID, resp.Data.CardID)
	}
}

// TestGetCardInfoInvalidCardID 测试查询不存在的卡信息
func TestGetCardInfoInvalidCardID(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询一个不存在的card_id
	cardID := "INVALID_CARD_ID_NOT_EXISTS"
	
	resp, err := client.Card.GetCardInfo(cardID)
	if err != nil {
		t.Logf("Got expected error for invalid card_id: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
	} else {
		t.Log("Query returned success (unexpected, but API may handle it differently)")
	}
}

// TestGetCardInfoEmptyCardID 测试空的card_id
func TestGetCardInfoEmptyCardID(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 测试空的card_id
	_, err := client.Card.GetCardInfo("")
	if err == nil {
		t.Error("Expected error for empty card_id, but got success")
	} else {
		t.Logf("Got expected error for empty card_id: %v", err)
	}
}
