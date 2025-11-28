package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	
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
		config.Endpoint = "https://api-test.gsalary.com"
	}
	
	// 2. 加载密钥
	privateKeyFile := os.Getenv("GSALARY_CLIENT_PRIVATE_KEY_FILE")
	if privateKeyFile == "" {
		privateKeyFile = "../../private_key.pem"
	}
	if err := config.ConfigClientPrivateKeyPEMFile(privateKeyFile); err != nil {
		log.Fatal("加载私钥失败:", err)
	}
	
	publicKeyFile := os.Getenv("GSALARY_SERVER_PUBLIC_KEY_FILE")
	if publicKeyFile == "" {
		publicKeyFile = "../../server_public_key.pem"
	}
	if err := config.ConfigServerPublicKeyPEMFile(publicKeyFile); err != nil {
		log.Fatal("加载服务端公钥失败:", err)
	}
	
	// 3. 创建Webhook处理器
	webhookHandler := api.NewWebhookHandler(config)
	
	// 4. 设置HTTP路由
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		// 只接受POST请求
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		
		// 处理Webhook
		resp, err := webhookHandler.HandleWebhook(r)
		if err != nil {
			log.Printf("❌ Webhook处理失败: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK) // 仍返回200，避免重复推送
			json.NewEncoder(w).Encode(resp)
			return
		}
		
		// 读取原始请求以便处理业务逻辑
		var webhookReq api.WebhookRequest
		// 注意：这里需要重新读取body，实际使用时应该在HandleWebhook中保存
		// 为了演示，这里假设已经验证通过
		
		// 根据业务类型处理不同的事件
		// 这里演示如何使用，实际需要从请求中解析
		fmt.Printf("✅ Webhook接收成功\n")
		fmt.Printf("   业务类型: %s\n", webhookReq.BusinessType)
		
		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	
	// 5. 添加具体业务处理示例
	http.HandleFunc("/webhook/payment", handlePaymentWebhook(webhookHandler))
	http.HandleFunc("/webhook/auth", handleAuthWebhook(webhookHandler))
	http.HandleFunc("/webhook/card", handleCardWebhook(webhookHandler))
	
	// 6. 启动服务器
	port := "8080"
	fmt.Printf("🚀 Webhook服务器启动在端口 %s\n", port)
	fmt.Printf("   支付结果通知: http://localhost:%s/webhook/payment\n", port)
	fmt.Printf("   授权Token通知: http://localhost:%s/webhook/auth\n", port)
	fmt.Printf("   卡事件通知:   http://localhost:%s/webhook/card\n", port)
	fmt.Println()
	
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("启动服务器失败:", err)
	}
}

// handlePaymentWebhook 处理支付结果通知
func handlePaymentWebhook(handler *api.WebhookHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 验证并解析Webhook
		resp, err := handler.HandleWebhook(r)
		if err != nil {
			log.Printf("❌ 支付Webhook验证失败: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		
		// 这里需要重新解析请求体获取业务数据
		// 实际使用时应该在HandleWebhook返回解析后的数据
		// 这里为演示目的，假设已经有了WebhookRequest对象
		
		fmt.Println("=== 收到支付结果通知 ===")
		fmt.Println("✅ 签名验证通过")
		
		// 处理支付结果业务逻辑
		// 例如：
		// 1. 更新订单状态
		// 2. 发送通知给用户
		// 3. 记录日志
		// 4. 如果是首次支付且返回card_token，保存card_token用于后续代扣
		
		fmt.Println("💡 提示: 请在这里添加您的业务逻辑")
		fmt.Println()
		
		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// handleAuthWebhook 处理授权Token通知
func handleAuthWebhook(handler *api.WebhookHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 验证并解析Webhook
		resp, err := handler.HandleWebhook(r)
		if err != nil {
			log.Printf("❌ 授权Webhook验证失败: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		
		fmt.Println("=== 收到授权Token通知 ===")
		fmt.Println("✅ 签名验证通过")
		
		// 处理授权Token业务逻辑
		// 例如：
		// 1. 保存access_token和refresh_token
		// 2. 记录token过期时间
		// 3. 如果status=REVOKED，删除保存的token
		// 4. 关联用户账号
		
		fmt.Println("💡 提示: 请保存access_token用于后续代扣支付")
		fmt.Println("💡 提示: 请监控access_token过期时间，及时调用刷新接口")
		fmt.Println()
		
		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// handleCardWebhook 处理卡事件通知
func handleCardWebhook(handler *api.WebhookHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 验证并解析Webhook
		resp, err := handler.HandleWebhook(r)
		if err != nil {
			log.Printf("❌ 卡Webhook验证失败: %v\n", err)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
		
		fmt.Println("=== 收到卡事件通知 ===")
		fmt.Println("✅ 签名验证通过")
		
		// 处理卡事件业务逻辑
		// 例如：
		// 1. 卡状态变更：更新本地卡状态缓存
		// 2. 卡交易通知：记录交易历史，发送通知给用户
		// 3. 卡充值结果：更新充值订单状态
		// 4. 申卡结果：通知用户申卡结果
		
		fmt.Println("💡 提示: 请根据事件类型处理相应的业务逻辑")
		fmt.Println()
		
		// 返回成功响应
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
