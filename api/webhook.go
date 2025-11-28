package api

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// WebhookEvent Webhook事件类型
const (
	// 收单服务事件
	EventAcquiringPaymentResult = "ACQUIRING_PAYMENT_RESULT" // 支付结果通知
	EventAcquiringAuthToken     = "ACQUIRING_AUTH_TOKEN"     // 授权Token通知
	
	// 卡相关事件
	EventCardStatusUpdate     = "CARD_STATUS_UPDATE"     // 卡状态变更
	EventCardTransaction      = "CARD_TRANSACTION"       // 卡交易通知
	EventCardAdjustResult     = "CARD_ADJUST_RESULT"     // 卡充值结果
	EventCardApplyResult      = "CARD_APPLY_RESULT"      // 申卡结果
	
	// 换汇相关事件
	EventExchangeOrderResult  = "EXCHANGE_ORDER_RESULT"  // 换汇订单结果
	
	// 对外付款事件
	EventRemittanceOrderResult = "REMITTANCE_ORDER_RESULT" // 付款订单结果
	EventPayeeDeactivated      = "PAYEE_DEACTIVATED"       // 收款人被停用
)

// WebhookHandler Webhook处理器
type WebhookHandler struct {
	config *gsalary.GSalaryConfig
}

// NewWebhookHandler 创建Webhook处理器
func NewWebhookHandler(config *gsalary.GSalaryConfig) *WebhookHandler {
	return &WebhookHandler{config: config}
}

// WebhookRequest Webhook请求数据
type WebhookRequest struct {
	AppID           string          `json:"app_id"`           // 应用ID
	BusinessType    string          `json:"business_type"`    // 业务类型（事件类型）
	Timestamp       int64           `json:"timestamp"`        // 时间戳（毫秒）
	Data            json.RawMessage `json:"data"`             // 业务数据
	SignatureHeader string          `json:"-"`                // 签名头（从HTTP Header获取）
}

// WebhookResponse Webhook响应
type WebhookResponse struct {
	Result  string `json:"result"`  // S-成功 F-失败
	Code    string `json:"code"`    // 结果代码
	Message string `json:"message"` // 结果消息
}

// PaymentResultData 支付结果通知数据
type PaymentResultData struct {
	PaymentRequestID  string                 `json:"payment_request_id"`  // 商户支付请求ID
	PaymentID         string                 `json:"payment_id"`          // GSalary支付ID
	PaymentAmount     float64                `json:"payment_amount"`      // 支付金额
	PaymentCurrency   string                 `json:"payment_currency"`    // 支付币种
	PaymentStatus     string                 `json:"payment_status"`      // 支付状态：SUCCESS/FAIL/PROCESSING
	PaymentResultCode string                 `json:"payment_result_code"` // 支付结果代码
	PaymentResultInfo map[string]interface{} `json:"payment_result_info"` // 支付结果信息（如card_token）
	PaymentTime       string                 `json:"payment_time"`        // 支付时间
}

// AuthTokenData 授权Token通知数据
type AuthTokenData struct {
	AuthState             string `json:"auth_state"`               // 授权状态ID
	AccessToken           string `json:"access_token"`             // 访问令牌
	AccessTokenExpiryTime string `json:"access_token_expiry_time"` // 访问令牌过期时间
	RefreshToken          string `json:"refresh_token"`            // 刷新令牌
	RefreshTokenExpiryTime string `json:"refresh_token_expiry_time"` // 刷新令牌过期时间
	UserLoginID           string `json:"user_login_id"`            // 用户登录ID
	PaymentMethodType     string `json:"payment_method_type"`      // 支付方式类型
	AuthClientID          string `json:"auth_client_id"`           // 授权客户端ID
	Status                string `json:"status"`                   // 状态：ACTIVE/REVOKED
}

// CardStatusUpdateData 卡状态变更通知数据
type CardStatusUpdateData struct {
	CardID     string `json:"card_id"`     // 卡片ID
	Status     string `json:"status"`      // 卡状态
	UpdateTime string `json:"update_time"` // 更新时间
}

// CardTransactionData 卡交易通知数据
type CardTransactionData struct {
	TransactionID      string  `json:"transaction_id"`       // 交易ID
	CardID             string  `json:"card_id"`              // 卡片ID
	TransactionType    string  `json:"transaction_type"`     // 交易类型
	Amount             float64 `json:"amount"`               // 交易金额
	Currency           string  `json:"currency"`             // 币种
	Status             string  `json:"status"`               // 交易状态
	StatusDescription  string  `json:"status_description"`   // 状态描述
	TransactionTime    string  `json:"transaction_time"`     // 交易时间
	MerchantName       string  `json:"merchant_name"`        // 商户名称
	MerchantCountry    string  `json:"merchant_country"`     // 商户国家
}

// HandleWebhook 处理Webhook请求
func (h *WebhookHandler) HandleWebhook(r *http.Request) (*WebhookResponse, error) {
	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return &WebhookResponse{
			Result:  "F",
			Code:    "READ_BODY_FAILED",
			Message: "Failed to read request body",
		}, fmt.Errorf("read body failed: %w", err)
	}
	defer r.Body.Close()
	
	// 获取签名头
	signatureHeader := r.Header.Get("authorization")
	if signatureHeader == "" {
		return &WebhookResponse{
			Result:  "F",
			Code:    "MISSING_SIGNATURE",
			Message: "Missing authorization header",
		}, fmt.Errorf("missing authorization header")
	}
	
	// 验证签名
	if err := h.verifySignature(body, signatureHeader); err != nil {
		return &WebhookResponse{
			Result:  "F",
			Code:    "SIGNATURE_VERIFICATION_FAILED",
			Message: "Signature verification failed",
		}, fmt.Errorf("signature verification failed: %w", err)
	}
	
	// 解析请求
	var req WebhookRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return &WebhookResponse{
			Result:  "F",
			Code:    "INVALID_JSON",
			Message: "Invalid JSON format",
		}, fmt.Errorf("json unmarshal failed: %w", err)
	}
	
	req.SignatureHeader = signatureHeader
	
	// 返回成功响应
	return &WebhookResponse{
		Result:  "S",
		Code:    "SUCCESS",
		Message: "Webhook received successfully",
	}, nil
}

// verifySignature 验证Webhook签名
func (h *WebhookHandler) verifySignature(body []byte, signatureHeader string) error {
	// 解析签名头: algorithm=RSA2,time=1234567890,signature=base64_signature
	parts := strings.Split(signatureHeader, ",")
	signatureMap := make(map[string]string)
	for _, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			signatureMap[kv[0]] = kv[1]
		}
	}
	
	algorithm := signatureMap["algorithm"]
	timestamp := signatureMap["time"]
	signature := signatureMap["signature"]
	
	if algorithm != "RSA2" {
		return fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
	
	if signature == "" {
		return fmt.Errorf("missing signature")
	}
	
	// 检查时间戳（防重放攻击，允许5分钟误差）
	if timestamp != "" {
		// 这里可以根据需要检查时间戳有效性
		// 实际生产环境建议检查时间戳防止重放攻击
	}
	
	// 解码签名
	signatureBytes, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		return fmt.Errorf("decode signature failed: %w", err)
	}
	
	// 计算消息摘要
	hash := sha256.Sum256(body)
	
	// 验证签名
	publicKey := h.config.GetServerPublicKey()
	if publicKey == nil {
		return fmt.Errorf("server public key not configured")
	}
	
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hash[:], signatureBytes)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}
	
	return nil
}

// ParsePaymentResult 解析支付结果通知
func (h *WebhookHandler) ParsePaymentResult(req *WebhookRequest) (*PaymentResultData, error) {
	if req.BusinessType != EventAcquiringPaymentResult {
		return nil, fmt.Errorf("invalid business type: %s", req.BusinessType)
	}
	
	var data PaymentResultData
	if err := json.Unmarshal(req.Data, &data); err != nil {
		return nil, fmt.Errorf("parse payment result failed: %w", err)
	}
	
	return &data, nil
}

// ParseAuthToken 解析授权Token通知
func (h *WebhookHandler) ParseAuthToken(req *WebhookRequest) (*AuthTokenData, error) {
	if req.BusinessType != EventAcquiringAuthToken {
		return nil, fmt.Errorf("invalid business type: %s", req.BusinessType)
	}
	
	var data AuthTokenData
	if err := json.Unmarshal(req.Data, &data); err != nil {
		return nil, fmt.Errorf("parse auth token failed: %w", err)
	}
	
	return &data, nil
}

// ParseCardStatusUpdate 解析卡状态变更通知
func (h *WebhookHandler) ParseCardStatusUpdate(req *WebhookRequest) (*CardStatusUpdateData, error) {
	if req.BusinessType != EventCardStatusUpdate {
		return nil, fmt.Errorf("invalid business type: %s", req.BusinessType)
	}
	
	var data CardStatusUpdateData
	if err := json.Unmarshal(req.Data, &data); err != nil {
		return nil, fmt.Errorf("parse card status update failed: %w", err)
	}
	
	return &data, nil
}

// ParseCardTransaction 解析卡交易通知
func (h *WebhookHandler) ParseCardTransaction(req *WebhookRequest) (*CardTransactionData, error) {
	if req.BusinessType != EventCardTransaction {
		return nil, fmt.Errorf("invalid business type: %s", req.BusinessType)
	}
	
	var data CardTransactionData
	if err := json.Unmarshal(req.Data, &data); err != nil {
		return nil, fmt.Errorf("parse card transaction failed: %w", err)
	}
	
	return &data, nil
}

// WebhookServer Webhook服务器示例（可选）
type WebhookServer struct {
	handler *WebhookHandler
	port    string
}

// NewWebhookServer 创建Webhook服务器
func NewWebhookServer(config *gsalary.GSalaryConfig, port string) *WebhookServer {
	return &WebhookServer{
		handler: NewWebhookHandler(config),
		port:    port,
	}
}

// Start 启动Webhook服务器
func (s *WebhookServer) Start() error {
	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(WebhookResponse{
				Result:  "F",
				Code:    "METHOD_NOT_ALLOWED",
				Message: "Only POST method is allowed",
			})
			return
		}
		
		// 处理Webhook
		resp, err := s.handler.HandleWebhook(r)
		if err != nil {
			fmt.Printf("Webhook error: %v\n", err)
		}
		
		// 读取请求体用于日志
		r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB max
		body, _ := io.ReadAll(r.Body)
		
		// 解析业务类型
		var webhookReq WebhookRequest
		json.Unmarshal(body, &webhookReq)
		
		fmt.Printf("[%s] Webhook received: business_type=%s, timestamp=%d\n",
			time.Now().Format("2006-01-02 15:04:05"),
			webhookReq.BusinessType,
			webhookReq.Timestamp)
		
		// 返回响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
	
	fmt.Printf("Webhook server started on port %s\n", s.port)
	return http.ListenAndServe(":"+s.port, nil)
}
