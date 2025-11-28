package main

import (
	"fmt"
	"log"
	"os"
	"time"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
	"github.com/difyz9/gsalary-sdk-go/api"
)

func main() {
	// 1. 配置客户端
	config := gsalary.NewConfig()
	
	// 从环境变量读取配置（推荐）
	config.AppID = os.Getenv("GSALARY_APPID")
	if config.AppID == "" {
		log.Fatal("请设置环境变量 GSALARY_APPID")
	}
	
	config.Endpoint = os.Getenv("GSALARY_ENDPOINT")
	if config.Endpoint == "" {
		config.Endpoint = "https://api-test.gsalary.com" // 默认测试环境
	}
	
	// 2. 加载密钥
	privateKeyFile := os.Getenv("GSALARY_CLIENT_PRIVATE_KEY_FILE")
	if privateKeyFile == "" {
		privateKeyFile = "./private_key.pem"
	}
	if err := config.ConfigClientPrivateKeyPEMFile(privateKeyFile); err != nil {
		log.Fatal("加载私钥失败:", err)
	}
	
	publicKeyFile := os.Getenv("GSALARY_SERVER_PUBLIC_KEY_FILE")
	if publicKeyFile == "" {
		publicKeyFile = "./server_public_key.pem"
	}
	if err := config.ConfigServerPublicKeyPEMFile(publicKeyFile); err != nil {
		log.Fatal("加载服务端公钥失败:", err)
	}
	
	// 3. 创建API客户端
	client := api.NewClient(config)
	
	// 4. 查询持卡人列表
	fmt.Println("=== 查询持卡人列表 ===")
	getCardHolderList(client)
	
	fmt.Println()
	
	// 4. 添加持卡人
	fmt.Println("=== 添加持卡人 ===")
	holderID := addCardHolder(client)
	
	fmt.Println()
	
	// 5. 查询持卡人列表
	fmt.Println("=== 查询持卡人列表 ===")
	getCardHolderList(client)
	
	// 如果成功创建了持卡人，查看其详细信息
	if holderID != "" {
		fmt.Println()
		fmt.Println("=== 查看持卡人详细信息 ===")
		getCardHolderInfo(client, holderID)
	}
	
	fmt.Println()
	
	// 6. 查询卡列表
	fmt.Println("=== 查询卡列表 ===")
	cardID := getCardList(client)
	
	// 如果有卡片，查看第一张卡的详细信息
	if cardID != "" {
		fmt.Println()
		fmt.Println("=== 查看卡详细信息 ===")
		getCardInfo(client, cardID)
	}
	
	fmt.Println()
	
	// 7. 查询可用的卡产品列表
	fmt.Println("=== 查询可用的卡产品列表 ===")
	getProducts(client)
	
	fmt.Println()
	
	// 7. 查询当前汇率
	fmt.Println("=== 查询当前汇率 ===")
	getExchangeRate(client)
	
	fmt.Println()
	
	// 7.1. 请求锁汇报价
	fmt.Println("=== 请求锁汇报价 ===")
	requestExchangeQuote(client)
	
	fmt.Println()
	
	// 8. 查询钱包余额
	fmt.Println("=== 查询钱包余额 ===")
	getWalletBalance(client)
	
	fmt.Println()
	
	// 9. 查询卡可用余额
	fmt.Println("=== 查询卡可用余额 ===")
	getAvailableQuotas(client)
	
	fmt.Println()
	
	// 10. 修改持卡人信息（如果成功创建了持卡人）
	if holderID != "" {
		fmt.Println("=== 修改持卡人信息 ===")
		updateCardHolder(client, holderID)
		
		fmt.Println()
	}
	
	// 11. 申请新卡片（使用之前创建的持卡人ID）
	fmt.Println("=== 申请新卡片 ===")
	requestID := applyCard(client, holderID)
	
	if requestID != "" {
		fmt.Println()
		
		// 12. 查询开卡结果
		fmt.Println("=== 查询开卡结果 ===")
		getCardApplyResult(client, requestID)
	}
}

func getExchangeRate(client *api.Client) {
	// 查询USD到HKD的汇率
	req := &api.ExchangeRateRequest{
		BuyCurrency:  "HKD",
		SellCurrency: "USD",
	}
	
	fmt.Printf("查询汇率: %s -> %s\n", req.SellCurrency, req.BuyCurrency)
	fmt.Println()
	
	// 发起查询
	resp, err := client.Exchange.GetCurrentExchangeRate(req)
	if err != nil {
		log.Printf("❌ 查询汇率失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 查询汇率成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   卖出币种: %s\n", resp.Data.SellCurrency)
	fmt.Printf("   买入币种: %s\n", resp.Data.BuyCurrency)
	fmt.Printf("   汇率: %.6f\n", resp.Data.Rate)
	fmt.Printf("   说明: 1 %s = %.6f %s\n", resp.Data.SellCurrency, resp.Data.Rate, resp.Data.BuyCurrency)
	fmt.Printf("   更新时间: %s\n", resp.Data.UpdateTime)
}

func requestExchangeQuote(client *api.Client) {
	// 请求锁汇报价（卖出100 USD，换取HKD）
	req := &api.ExchangeQuoteRequest{
		BuyCurrency:  "HKD",
		SellCurrency: "USD",
		SellAmount:   100.00,
	}
	
	fmt.Printf("请求锁汇报价: 卖出 %s %.2f，换取 %s\n", 
		req.SellCurrency, req.SellAmount, req.BuyCurrency)
	fmt.Println()
	
	// 发起请求
	resp, err := client.Exchange.RequestQuote(req)
	if err != nil {
		log.Printf("❌ 请求锁汇报价失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 请求锁汇报价成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   报价ID: %s\n", resp.Data.QuoteID)
	fmt.Printf("   买入: %s %.2f\n", resp.Data.Buy.Currency, resp.Data.Buy.Amount)
	fmt.Printf("   卖出: %s %.2f\n", resp.Data.Sell.Currency, resp.Data.Sell.Amount)
	fmt.Printf("   手续费: %s %.2f\n", resp.Data.Surcharge.Currency, resp.Data.Surcharge.Amount)
	fmt.Printf("   总成本: %s %.2f\n", resp.Data.TotalCost.Currency, resp.Data.TotalCost.Amount)
	fmt.Printf("   更新时间: %s\n", resp.Data.UpdateTime)
	fmt.Printf("   过期时间: %s\n", resp.Data.ExpireTime)
	fmt.Println()
	fmt.Printf("💡 提示: 该报价将在 %s 过期，请在此之前完成换汇操作\n", resp.Data.ExpireTime)
}

func getWalletBalance(client *api.Client) {
	// 查询USD钱包余额
	req := &api.WalletBalanceRequest{
		Currency: "USD",
	}
	
	fmt.Printf("查询钱包余额，币种: %s\n", req.Currency)
	fmt.Println()
	
	// 发起查询
	resp, err := client.Wallet.GetBalance(req)
	if err != nil {
		log.Printf("❌ 查询钱包余额失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 查询钱包余额成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   币种: %s\n", resp.Data.Currency)
	fmt.Printf("   总金额: %.2f\n", resp.Data.Amount)
	fmt.Printf("   共享卡账户余额: %.2f\n", resp.Data.ShareCardAccountBalance)
	fmt.Printf("   可用余额: %.2f\n", resp.Data.Available)
	fmt.Printf("   账户类型: %s\n", resp.Data.AccountType)
	fmt.Printf("   查询时间: %s\n", resp.Data.QueryTime)
}

func getCardHolderList(client *api.Client) {
	// 查询第一页，每页10条
	req := &api.CardHolderListRequest{
		Page:  1,
		Limit: 10,
	}
	
	fmt.Printf("查询持卡人列表，页码: %d，每页数量: %d\n", req.Page, req.Limit)
	fmt.Println()
	
	// 发起查询
	resp, err := client.CardHolder.GetCardHolderList(req)
	if err != nil {
		log.Printf("❌ 查询持卡人列表失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 查询持卡人列表成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   当前页: %d\n", resp.Data.Page)
	fmt.Printf("   每页数量: %d\n", resp.Data.Limit)
	fmt.Printf("   总记录数: %d\n", resp.Data.TotalCount)
	fmt.Printf("   总页数: %d\n", resp.Data.TotalPage)
	fmt.Printf("   本页持卡人数: %d\n", len(resp.Data.CardHolders))
	fmt.Println()
	
	// 显示持卡人详情（最多显示前3个）
	displayCount := len(resp.Data.CardHolders)
	if displayCount > 3 {
		displayCount = 3
	}
	
	for i := 0; i < displayCount; i++ {
		holder := resp.Data.CardHolders[i]
		fmt.Printf("持卡人 %d:\n", i+1)
		fmt.Printf("  持卡人ID: %s\n", holder.CardHolderID)
		fmt.Printf("  姓名: %s %s\n", holder.FirstName, holder.LastName)
		fmt.Printf("  邮箱: %s\n", holder.Email)
		fmt.Printf("  地区: %s\n", holder.Region)
		fmt.Printf("  生日: %s\n", holder.Birth)
		fmt.Printf("  手机: +%s %s\n", holder.Mobile.CountryCode, holder.Mobile.Number)
		fmt.Printf("  创建时间: %s\n", holder.CreatedAt)
		fmt.Println()
	}
	
	if len(resp.Data.CardHolders) > 3 {
		fmt.Printf("... 还有 %d 个持卡人\n", len(resp.Data.CardHolders)-3)
		fmt.Println()
	}
}

func getCardHolderInfo(client *api.Client, holderID string) {
	fmt.Printf("查询持卡人详细信息，持卡人ID: %s\n", holderID)
	fmt.Println()
	
	// 发起查询
	resp, err := client.CardHolder.GetCardHolderInfo(holderID)
	if err != nil {
		log.Printf("❌ 查询持卡人信息失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 查询持卡人信息成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   持卡人ID: %s\n", resp.Data.CardHolderID)
	fmt.Printf("   姓名: %s %s\n", resp.Data.FirstName, resp.Data.LastName)
	fmt.Printf("   邮箱: %s\n", resp.Data.Email)
	fmt.Printf("   生日: %s\n", resp.Data.Birth)
	fmt.Printf("   地区: %s\n", resp.Data.Region)
	fmt.Printf("   状态: %s\n", resp.Data.Status)
	fmt.Printf("   创建时间: %s\n", resp.Data.CreateTime)
	
	// 显示账单地址
	if resp.Data.BillAddress != nil {
		fmt.Println("   账单地址:")
		if addrMap, ok := resp.Data.BillAddress.(map[string]interface{}); ok {
			for key, value := range addrMap {
				fmt.Printf("     %s: %v\n", key, value)
			}
		}
	}
}

func updateCardHolder(client *api.Client, holderID string) {
	fmt.Printf("修改持卡人信息，持卡人ID: %s\n", holderID)
	fmt.Println()
	
	// 创建修改请求
	req := &api.UpdateCardHolderRequest{
		FirstName: "John",
		LastName:  "Smith", // 修改姓氏
		Birth:     "1990-05-15",
		Email:     fmt.Sprintf("demo.holder.updated.%d@example.com", time.Now().Unix()), // 新邮箱
		Mobile: api.MobileNumber{
			CountryCode: "1",
			Number:      "4155559999", // 修改手机号
		},
		Region: "US",
		BillAddress: api.Address{
			Country:    "US",
			State:      "NY",         // 修改州
			City:       "New York",   // 修改城市
			PostalCode: "10001",
			Line1:      "456 Updated Street", // 修改地址
			Line2:      "Floor 5",
		},
	}
	
	fmt.Println("修改内容:")
	fmt.Printf("  新姓氏: %s\n", req.LastName)
	fmt.Printf("  新邮箱: %s\n", req.Email)
	fmt.Printf("  新手机号: %s-%s\n", req.Mobile.CountryCode, req.Mobile.Number)
	fmt.Printf("  新地址: %s, %s, %s\n", req.BillAddress.City, req.BillAddress.State, req.BillAddress.Country)
	fmt.Println()
	
	// 发起修改请求
	resp, err := client.CardHolder.UpdateCardHolder(holderID, req)
	if err != nil {
		log.Printf("❌ 修改持卡人信息失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 修改持卡人信息成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   持卡人ID: %s\n", resp.Data.CardHolderID)
	fmt.Printf("   更新后姓名: %s %s\n", resp.Data.FirstName, resp.Data.LastName)
	fmt.Printf("   更新后邮箱: %s\n", resp.Data.Email)
	fmt.Printf("   更新后手机号: %s-%s\n", resp.Data.Mobile.CountryCode, resp.Data.Mobile.Number)
	fmt.Printf("   更新时间: %s\n", resp.Data.UpdatedAt)
}

func addCardHolder(client *api.Client) string {
	// 生成唯一的邮箱（实际使用时应该使用真实的用户信息）
	email := fmt.Sprintf("demo.holder.%d@example.com", time.Now().Unix())
	
	// 创建添加持卡人请求
	req := &api.CardHolderRequest{
		FirstName: "John",
		LastName:  "Doe",
		Birth:     "1990-05-15",
		Email:     email,
		Mobile: api.MobileNumber{
			CountryCode: "1",
			Number:      "4155551234",
		},
		Region: "US",
		BillAddress: api.Address{
			Country:    "US",
			State:      "CA",
			City:       "San Francisco",
			PostalCode: "94102",
			Line1:      "123 Main Street",
			Line2:      "Apt 4B",
		},
	}
	
	fmt.Printf("添加持卡人:\n")
	fmt.Printf("  名字: %s %s\n", req.FirstName, req.LastName)
	fmt.Printf("  邮箱: %s\n", req.Email)
	fmt.Printf("  生日: %s\n", req.Birth)
	fmt.Printf("  地区: %s\n", req.Region)
	fmt.Println()
	
	// 发起添加请求
	resp, err := client.CardHolder.AddCardHolder(req)
	if err != nil {
		log.Printf("❌ 添加持卡人失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return ""
	}
	
	// 显示结果
	fmt.Println("✅ 添加持卡人成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   代码: %s\n", resp.Result.Code)
	fmt.Printf("   消息: %s\n", resp.Result.Message)
	fmt.Printf("   持卡人ID: %s\n", resp.Data.CardHolderID)
	fmt.Printf("   名字: %s %s\n", resp.Data.FirstName, resp.Data.LastName)
	fmt.Printf("   邮箱: %s\n", resp.Data.Email)
	fmt.Printf("   创建时间: %s\n", resp.Data.CreatedAt)
	
	return resp.Data.CardHolderID
}

func applyCard(client *api.Client, holderID string) string {
	// 生成唯一的请求ID（实际使用时应该使用UUID或其他唯一标识）
	requestID := fmt.Sprintf("DEMO_CARD_APPLY_%d", time.Now().Unix())
	
	// 如果没有提供持卡人ID，使用默认值
	if holderID == "" {
		holderID = "holder_demo_001"
	}
	
	// 创建申请请求
	req := &api.CardApplyRequest{
		RequestID:           requestID,
		ProductCode:         "VIRTUAL_CARD_USD", // 虚拟卡产品代码
		Currency:            "USD",               // 美元
		CardHolderID:        holderID,            // 持卡人ID
		LimitPerDay:         1000.00,             // 每日限额
		LimitPerMonth:       5000.00,             // 每月限额
		LimitPerTransaction: 500.00,              // 单笔限额
		InitBalance:         100.00,              // 初始余额100美元
	}
	
	fmt.Printf("请求ID: %s\n", req.RequestID)
	fmt.Printf("产品代码: %s\n", req.ProductCode)
	fmt.Printf("货币: %s\n", req.Currency)
	fmt.Printf("持卡人ID: %s\n", req.CardHolderID)
	fmt.Printf("初始余额: %.2f\n", req.InitBalance)
	fmt.Println()
	
	// 发起申请
	resp, err := client.Card.ApplyCard(req)
	if err != nil {
		log.Printf("❌ 申请卡片失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return ""
	}
	
	// 显示结果
	fmt.Println("✅ 申请卡片成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   代码: %s\n", resp.Result.Code)
	fmt.Printf("   消息: %s\n", resp.Result.Message)
	fmt.Printf("   请求ID: %s\n", resp.Data.RequestID)
	fmt.Printf("   状态: %s\n", resp.Data.Status)
	
	return requestID
}

func getAvailableQuotas(client *api.Client) {
	// 查询USD币种的SHARE类型余额
	req := &api.CardAvailableQuotasRequest{
		Currency:           "USD",
		AccountingCardType: "SHARE", // 可选: "SHARE" 或 "RECHARGE"，不填默认为SHARE
	}
	
	fmt.Printf("货币: %s\n", req.Currency)
	fmt.Printf("卡账务类型: %s\n", req.AccountingCardType)
	fmt.Println()
	
	// 发起查询
	resp, err := client.Card.GetAvailableQuotas(req)
	if err != nil {
		log.Printf("❌ 查询余额失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 查询余额成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   代码: %s\n", resp.Result.Code)
	fmt.Printf("   消息: %s\n", resp.Result.Message)
	fmt.Printf("   货币: %s\n", resp.Data.Currency)
	fmt.Printf("   卡账务类型: %s\n", resp.Data.AccountingCardType)
	fmt.Printf("   可用余额: %.2f\n", resp.Data.AvailableQuota)
}

func getProducts(client *api.Client) {
	// 查询所有可用产品
	req := &api.CardProductsRequest{
		// 可以添加过滤条件，不填则查询所有产品
		// CardType:  "VIRTUAL",
		// BrandCode: "VISA",
		// Currency:  "USD",
	}
	
	fmt.Println("查询所有可用的卡产品...")
	fmt.Println()
	
	// 发起查询
	resp, err := client.Card.GetProducts(req)
	if err != nil {
		log.Printf("❌ 查询产品列表失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 查询产品列表成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   产品数量: %d\n", len(resp.Data.Products))
	fmt.Println()
	
	// 显示产品详情
	for i, product := range resp.Data.Products {
		fmt.Printf("产品 %d:\n", i+1)
		fmt.Printf("  代码: %s\n", product.ProductCode)
		fmt.Printf("  名称: %s\n", product.ProductName)
		fmt.Printf("  类型: %s\n", product.CardType)
		fmt.Printf("  品牌: %s\n", product.BrandCode)
		fmt.Printf("  币种: %s\n", product.Currency)
		if product.Description != "" {
			fmt.Printf("  描述: %s\n", product.Description)
		}
		fmt.Println()
	}
}

func getCardApplyResult(client *api.Client, requestID string) {
	fmt.Printf("查询开卡结果，Request ID: %s\n", requestID)
	fmt.Println()
	
	// 发起查询
	resp, err := client.Card.GetCardApplyResult(requestID)
	if err != nil {
		log.Printf("❌ 查询开卡结果失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 查询开卡结果成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   代码: %s\n", resp.Result.Code)
	fmt.Printf("   消息: %s\n", resp.Result.Message)
	
	// 显示data内容
	if len(resp.Data) > 0 {
		fmt.Println("   数据:")
		for key, value := range resp.Data {
			fmt.Printf("     %s: %v\n", key, value)
		}
	}
}

func getCardInfo(client *api.Client, cardID string) {
	fmt.Printf("查询卡详细信息，Card ID: %s\n", cardID)
	fmt.Println()
	
	// 发起查询
	resp, err := client.Card.GetCardInfo(cardID)
	if err != nil {
		log.Printf("❌ 查询卡信息失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return
	}
	
	// 显示结果
	fmt.Println("✅ 查询卡信息成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Println()
	
	// 显示卡片详细信息
	fmt.Println("卡片详细信息:")
	fmt.Printf("  卡片ID: %s\n", resp.Data.CardID)
	fmt.Printf("  卡片名称: %s\n", resp.Data.CardName)
	fmt.Printf("  掩码卡号: %s\n", resp.Data.MaskCardNumber)
	fmt.Printf("  币种: %s\n", resp.Data.CardCurrency)
	fmt.Printf("  可用余额: %.2f\n", resp.Data.AvailableBalance)
	fmt.Printf("  品牌: %s\n", resp.Data.BrandCode)
	fmt.Printf("  状态: %s\n", resp.Data.Status)
	fmt.Printf("  卡类型: %s\n", resp.Data.CardType)
	fmt.Printf("  账务类型: %s\n", resp.Data.AccountingType)
	fmt.Printf("  卡地区: %s\n", resp.Data.CardRegion)
	fmt.Printf("  持卡人ID: %s\n", resp.Data.CardHolderID)
	
	if resp.Data.FirstName != "" || resp.Data.LastName != "" {
		fmt.Printf("  持卡人姓名: %s %s\n", resp.Data.FirstName, resp.Data.LastName)
	}
	
	if resp.Data.Email != "" {
		fmt.Printf("  邮箱: %s\n", resp.Data.Email)
	}
	
	if len(resp.Data.Mobile) > 0 {
		fmt.Printf("  手机号: %v\n", resp.Data.Mobile)
	}
	
	fmt.Printf("  每日限额: %.2f\n", resp.Data.LimitPerDay)
	fmt.Printf("  每月限额: %.2f\n", resp.Data.LimitPerMonth)
	fmt.Printf("  单笔限额: %.2f\n", resp.Data.LimitPerTransaction)
	fmt.Printf("  支持3DS交易: %v\n", resp.Data.SupportTdsTrans)
	fmt.Printf("  创建时间: %s\n", resp.Data.CreateTime)
	
	if len(resp.Data.BillAddress) > 0 {
		fmt.Println("  账单地址:")
		for key, value := range resp.Data.BillAddress {
			fmt.Printf("    %s: %v\n", key, value)
		}
	}
}

func getCardList(client *api.Client) string {
	// 查询第一页，每页10条
	req := &api.CardListRequest{
		Page:  1,
		Limit: 10,
		// 可以添加过滤条件
		// Status: "ACTIVE",
		// BrandCode: "VISA",
	}
	
	fmt.Printf("查询卡列表，页码: %d，每页数量: %d\n", req.Page, req.Limit)
	fmt.Println()
	
	// 发起查询
	resp, err := client.Card.GetCardList(req)
	if err != nil {
		log.Printf("❌ 查询卡列表失败: %v\n", err)
		if resp != nil {
			fmt.Printf("   结果: %s\n", resp.Result.Result)
			fmt.Printf("   代码: %s\n", resp.Result.Code)
			fmt.Printf("   消息: %s\n", resp.Result.Message)
		}
		return ""
	}
	
	// 显示结果
	fmt.Println("✅ 查询卡列表成功!")
	fmt.Printf("   结果: %s\n", resp.Result.Result)
	fmt.Printf("   当前页: %d\n", resp.Data.Page)
	fmt.Printf("   每页数量: %d\n", resp.Data.Limit)
	fmt.Printf("   总记录数: %d\n", resp.Data.TotalCount)
	fmt.Printf("   总页数: %d\n", resp.Data.TotalPage)
	fmt.Printf("   本页卡片数: %d\n", len(resp.Data.Cards))
	fmt.Println()
	
	// 显示卡片详情（最多显示前3张）
	displayCount := len(resp.Data.Cards)
	if displayCount > 3 {
		displayCount = 3
	}
	
	for i := 0; i < displayCount; i++ {
		card := resp.Data.Cards[i]
		fmt.Printf("卡片 %d:\n", i+1)
		fmt.Printf("  卡片ID: %s\n", card.CardID)
		fmt.Printf("  产品代码: %s\n", card.ProductCode)
		fmt.Printf("  品牌: %s\n", card.BrandCode)
		fmt.Printf("  持卡人ID: %s\n", card.CardHolderID)
		fmt.Printf("  状态: %s\n", card.Status)
		fmt.Printf("  创建时间: %s\n", card.CreatedAt)
		fmt.Printf("  更新时间: %s\n", card.UpdatedAt)
		fmt.Println()
	}
	
	if len(resp.Data.Cards) > 3 {
		fmt.Printf("... 还有 %d 张卡片\n", len(resp.Data.Cards)-3)
		fmt.Println()
	}
	
	// 返回第一张卡的ID（如果有）用于后续查询详情
	if len(resp.Data.Cards) > 0 {
		return resp.Data.Cards[0].CardID
	}
	return ""
}
