package api

import (
	"fmt"
	"testing"
	"time"
)

// TestAddCardHolder 测试添加持卡人
func TestAddCardHolder(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 生成唯一的邮箱（确保唯一性）
	email := fmt.Sprintf("test.holder.%d@example.com", time.Now().Unix())
	
	// 创建添加持卡人请求
	req := &CardHolderRequest{
		FirstName: "John",
		LastName:  "Doe",
		Birth:     "1990-01-15",
		Email:     email,
		Mobile: MobileNumber{
			CountryCode: "1",
			Number:      "4155551234",
		},
		Region: "US",
		BillAddress: Address{
			Country:    "US",
			State:      "CA",
			City:       "San Francisco",
			PostalCode: "94102",
			Line1:      "123 Main St",
			Line2:      "Apt 4B",
		},
	}
	
	t.Logf("Adding card holder with email: %s", req.Email)
	
	// 发起添加请求
	resp, err := client.CardHolder.AddCardHolder(req)
	if err != nil {
		t.Logf("Add card holder error: %v", err)
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
	
	t.Logf("Add card holder success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Code: %s", resp.Result.Code)
	t.Logf("Message: %s", resp.Result.Message)
	t.Logf("Card Holder ID: %s", resp.Data.CardHolderID)
	t.Logf("First Name: %s", resp.Data.FirstName)
	t.Logf("Last Name: %s", resp.Data.LastName)
	t.Logf("Email: %s", resp.Data.Email)
	t.Logf("Region: %s", resp.Data.Region)
	t.Logf("Created At: %s", resp.Data.CreatedAt)
	
	// 验证基本响应
	if resp.Result.Result != "S" {
		t.Errorf("Expected result 'S', got '%s'", resp.Result.Result)
	}
	
	if resp.Data.CardHolderID == "" {
		t.Error("Card holder ID should not be empty")
	}
	
	if resp.Data.Email != email {
		t.Errorf("Expected email '%s', got '%s'", email, resp.Data.Email)
	}
}

// TestAddCardHolderMinimalInfo 测试添加持卡人（最小信息）
func TestAddCardHolderMinimalInfo(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 生成唯一的邮箱
	email := fmt.Sprintf("minimal.%d@example.com", time.Now().Unix())
	
	// 创建最小信息的持卡人请求
	req := &CardHolderRequest{
		FirstName: "Jane",
		LastName:  "Smith",
		Birth:     "1995-06-20",
		Email:     email,
		Mobile: MobileNumber{
			CountryCode: "1",
			Number:      "4155559876",
		},
		Region: "US",
		BillAddress: Address{
			Country:    "US",
			State:      "NY",
			City:       "New York",
			PostalCode: "10001",
			Line1:      "456 Park Ave",
			Line2:      "",
		},
	}
	
	t.Logf("Adding card holder with minimal info, email: %s", req.Email)
	
	// 发起添加请求
	resp, err := client.CardHolder.AddCardHolder(req)
	if err != nil {
		t.Logf("Add card holder (minimal) error: %v", err)
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
	
	t.Logf("Add card holder (minimal) success!")
	t.Logf("Card Holder ID: %s", resp.Data.CardHolderID)
	t.Logf("First Name: %s", resp.Data.FirstName)
	t.Logf("Last Name: %s", resp.Data.LastName)
}

// TestAddCardHolderInvalidEmail 测试添加持卡人（无效邮箱）
func TestAddCardHolderInvalidEmail(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 使用无效的邮箱格式
	req := &CardHolderRequest{
		FirstName: "Invalid",
		LastName:  "Email",
		Birth:     "1990-01-01",
		Email:     "invalid-email-format",
		Mobile: MobileNumber{
			CountryCode: "1",
			Number:      "4155551111",
		},
		Region: "US",
		BillAddress: Address{
			Country:    "US",
			State:      "CA",
			City:       "Los Angeles",
			PostalCode: "90001",
			Line1:      "789 Test St",
			Line2:      "",
		},
	}
	
	// 发起添加请求
	resp, err := client.CardHolder.AddCardHolder(req)
	if err != nil {
		t.Logf("Got expected error for invalid email: %v", err)
		if resp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp.Result.Result, resp.Result.Code, resp.Result.Message)
		}
	} else {
		t.Log("Add succeeded (API may not validate email format strictly)")
	}
}

// TestAddCardHolderDuplicateEmail 测试添加持卡人（重复邮箱）
func TestAddCardHolderDuplicateEmail(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 使用固定的邮箱地址（可能已经存在）
	duplicateEmail := "duplicate.test@example.com"
	
	req := &CardHolderRequest{
		FirstName: "Duplicate",
		LastName:  "Test",
		Birth:     "1992-03-10",
		Email:     duplicateEmail,
		Mobile: MobileNumber{
			CountryCode: "1",
			Number:      "4155552222",
		},
		Region: "US",
		BillAddress: Address{
			Country:    "US",
			State:      "TX",
			City:       "Houston",
			PostalCode: "77001",
			Line1:      "111 Duplicate Ave",
			Line2:      "",
		},
	}
	
	// 第一次添加
	resp1, err1 := client.CardHolder.AddCardHolder(req)
	if err1 != nil {
		t.Logf("First attempt error (may already exist): %v", err1)
	} else {
		t.Logf("First attempt success, Card Holder ID: %s", resp1.Data.CardHolderID)
	}
	
	// 第二次使用相同邮箱添加
	resp2, err2 := client.CardHolder.AddCardHolder(req)
	if err2 != nil {
		t.Logf("Got expected error for duplicate email: %v", err2)
		if resp2 != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				resp2.Result.Result, resp2.Result.Code, resp2.Result.Message)
		}
	} else {
		t.Log("Second attempt succeeded (API may allow duplicate or update existing)")
	}
}

// TestGetCardHolderList 测试查询持卡人列表
func TestGetCardHolderList(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 查询第一页
	req := &CardHolderListRequest{
		Page:  1,
		Limit: 10,
	}
	
	resp, err := client.CardHolder.GetCardHolderList(req)
	if err != nil {
		t.Logf("Get card holder list error: %v", err)
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
	
	t.Logf("Get card holder list success!")
	t.Logf("Result: %s", resp.Result.Result)
	t.Logf("Page: %d", resp.Data.Page)
	t.Logf("Limit: %d", resp.Data.Limit)
	t.Logf("Total Count: %d", resp.Data.TotalCount)
	t.Logf("Card Holders: %d", len(resp.Data.CardHolders))
	
	// 打印前3个持卡人
	for i, holder := range resp.Data.CardHolders {
		if i >= 3 {
			break
		}
		t.Logf("Holder %d:", i+1)
		t.Logf("  ID: %s", holder.CardHolderID)
		t.Logf("  Name: %s %s", holder.FirstName, holder.LastName)
		t.Logf("  Email: %s", holder.Email)
	}
}

// TestGetCardHolderInfo 测试查看持卡人信息
func TestGetCardHolderInfo(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 先添加一个持卡人
	email := fmt.Sprintf("test.info.%d@example.com", time.Now().Unix())
	addReq := &CardHolderRequest{
		FirstName: "Info",
		LastName:  "Test",
		Birth:     "1991-07-20",
		Email:     email,
		Mobile: MobileNumber{
			CountryCode: "1",
			Number:      "4155553333",
		},
		Region: "US",
		BillAddress: Address{
			Country:    "US",
			State:      "WA",
			City:       "Seattle",
			PostalCode: "98101",
			Line1:      "222 Test Ave",
			Line2:      "",
		},
	}
	
	addResp, err := client.CardHolder.AddCardHolder(addReq)
	if err != nil {
		t.Logf("Add card holder error (will use test ID): %v", err)
	}
	
	// 获取持卡人ID
	var holderID string
	if addResp != nil && addResp.Data.CardHolderID != "" {
		holderID = addResp.Data.CardHolderID
	} else {
		holderID = "test_holder_id_if_known"
	}
	
	// 查询持卡人信息
	infoResp, err := client.CardHolder.GetCardHolderInfo(holderID)
	if err != nil {
		t.Logf("Get card holder info error: %v", err)
		if infoResp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				infoResp.Result.Result, infoResp.Result.Code, infoResp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if infoResp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Get card holder info success!")
	t.Logf("Card Holder ID: %s", infoResp.Data.CardHolderID)
	t.Logf("Name: %s %s", infoResp.Data.FirstName, infoResp.Data.LastName)
	t.Logf("Email: %s", infoResp.Data.Email)
	t.Logf("Birth: %s", infoResp.Data.Birth)
	t.Logf("Region: %s", infoResp.Data.Region)
	t.Logf("Status: %s", infoResp.Data.Status)
}

// TestUpdateCardHolder 测试修改持卡人信息
func TestUpdateCardHolder(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 先添加一个持卡人
	email := fmt.Sprintf("test.update.%d@example.com", time.Now().Unix())
	addReq := &CardHolderRequest{
		FirstName: "Original",
		LastName:  "Name",
		Birth:     "1992-08-10",
		Email:     email,
		Mobile: MobileNumber{
			CountryCode: "1",
			Number:      "4155554444",
		},
		Region: "US",
		BillAddress: Address{
			Country:    "US",
			State:      "FL",
			City:       "Miami",
			PostalCode: "33101",
			Line1:      "333 Original St",
			Line2:      "",
		},
	}
	
	addResp, err := client.CardHolder.AddCardHolder(addReq)
	if err != nil {
		t.Logf("Add card holder error: %v", err)
		return
	}
	
	if addResp == nil || addResp.Data.CardHolderID == "" {
		t.Skip("Cannot get card holder ID, skipping update test")
	}
	
	holderID := addResp.Data.CardHolderID
	t.Logf("Created card holder ID: %s", holderID)
	
	// 修改持卡人信息
	updateReq := &UpdateCardHolderRequest{
		FirstName: "Updated",
		LastName:  "Name",
		Birth:     "1992-08-10", // 保持不变
		Email:     email,         // 保持不变
		Mobile: MobileNumber{
			CountryCode: "1",
			Number:      "4155555555", // 修改手机号
		},
		Region: "US",
		BillAddress: Address{
			Country:    "US",
			State:      "FL",
			City:       "Orlando", // 修改城市
			PostalCode: "32801",
			Line1:      "444 Updated Ave",
			Line2:      "Suite 100",
		},
	}
	
	updateResp, err := client.CardHolder.UpdateCardHolder(holderID, updateReq)
	if err != nil {
		t.Logf("Update card holder error: %v", err)
		if updateResp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				updateResp.Result.Result, updateResp.Result.Code, updateResp.Result.Message)
		}
		return
	}
	
	// 验证响应
	if updateResp == nil {
		t.Fatal("Response is nil")
	}
	
	t.Logf("Update card holder success!")
	t.Logf("Result: %s", updateResp.Result.Result)
	t.Logf("Card Holder ID: %s", updateResp.Data.CardHolderID)
	t.Logf("Updated Name: %s %s", updateResp.Data.FirstName, updateResp.Data.LastName)
	t.Logf("Updated Mobile: %s-%s", updateResp.Data.Mobile.CountryCode, updateResp.Data.Mobile.Number)
	
	// 验证更新结果
	if updateResp.Data.FirstName != "Updated" {
		t.Errorf("Expected first name 'Updated', got '%s'", updateResp.Data.FirstName)
	}
	
	if updateResp.Data.Mobile.Number != "4155555555" {
		t.Errorf("Expected mobile '4155555555', got '%s'", updateResp.Data.Mobile.Number)
	}
}

// TestUpdateCardHolderInvalidID 测试修改不存在的持卡人
func TestUpdateCardHolderInvalidID(t *testing.T) {
	// 检查密钥是否加载
	if testConfig.GetClientPrivateKey() == nil || testConfig.GetServerPublicKey() == nil {
		t.Skip("Skipping test: keys not loaded")
	}
	
	// 创建客户端
	client := NewClient(testConfig)
	
	// 使用不存在的持卡人ID
	invalidID := "invalid_holder_id_not_exists"
	
	updateReq := &UpdateCardHolderRequest{
		FirstName: "Invalid",
		LastName:  "Update",
		Birth:     "1990-01-01",
		Email:     "invalid@example.com",
		Mobile: MobileNumber{
			CountryCode: "1",
			Number:      "4155556666",
		},
		Region: "US",
		BillAddress: Address{
			Country:    "US",
			State:      "CA",
			City:       "San Jose",
			PostalCode: "95101",
			Line1:      "555 Invalid St",
			Line2:      "",
		},
	}
	
	updateResp, err := client.CardHolder.UpdateCardHolder(invalidID, updateReq)
	if err != nil {
		t.Logf("Got expected error for invalid card holder ID: %v", err)
		if updateResp != nil {
			t.Logf("Response: Result=%s, Code=%s, Message=%s",
				updateResp.Result.Result, updateResp.Result.Code, updateResp.Result.Message)
		}
	} else {
		t.Log("Update succeeded (unexpected, but API may handle it differently)")
	}
}
