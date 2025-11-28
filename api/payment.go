package api

import (
	"encoding/json"
	"fmt"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// PaymentAPI 支付API接口
type PaymentAPI struct {
	client *gsalary.GSalaryClient
}

// NewPaymentAPI 创建支付API实例
func NewPaymentAPI(client *gsalary.GSalaryClient) *PaymentAPI {
	return &PaymentAPI{client: client}
}

// PaymentConsult 支付咨询 - 查询可用支付方式、限额、国家/货币支持等信息
func (api *PaymentAPI) PaymentConsult(req *PaymentConsultRequest) (*PaymentConsultResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/gateway/v1/acquiring/pay_consult")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"mch_app_id":                       req.MchAppID,
		"payment_currency":                 req.PaymentCurrency,
		"payment_amount":                   req.PaymentAmount,
		"settlement_currency":              req.SettlementCurrency,
		"allowed_payment_method_regions":   req.AllowedPaymentMethodRegions,
		"allowed_payment_methods":          req.AllowedPaymentMethods,
		"env_terminal_type":                req.EnvTerminalType,
	}
	
	// 可选字段
	if req.UserRegion != "" {
		request.Body["user_region"] = req.UserRegion
	}
	if req.EnvOsType != "" {
		request.Body["env_os_type"] = req.EnvOsType
	}
	if req.EnvClientIP != "" {
		request.Body["env_client_ip"] = req.EnvClientIP
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("payment consult failed: %w", err)
	}
	
	// 解析响应
	var consultResp PaymentConsultResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &consultResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if consultResp.Result.Result != "S" {
		return &consultResp, fmt.Errorf("payment consult business error: [%s] %s",
			consultResp.Result.Code, consultResp.Result.Message)
	}
	
	return &consultResp, nil
}

// CreatePaymentSession 创建支付会话（收银台）
func (api *PaymentAPI) CreatePaymentSession(req *PaymentSessionRequest) (*PaymentSessionResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/gateway/v1/acquiring/pay_session")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"mch_app_id":           req.MchAppID,
		"payment_request_id":   req.PaymentRequestID,
		"payment_currency":     req.PaymentCurrency,
		"payment_amount":       req.PaymentAmount,
		"payment_method_type":  req.PaymentMethodType,
		"payment_redirect_url": req.PaymentRedirectURL,
		"order":                req.Order,
		"settlement_currency":  req.SettlementCurrency,
		"product_scene":        req.ProductScene,
	}
	
	// 可选字段
	if req.PaymentSessionExpiryTime != "" {
		request.Body["payment_session_expiry_time"] = req.PaymentSessionExpiryTime
	}
	if req.EnvClientIP != "" {
		request.Body["env_client_ip"] = req.EnvClientIP
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("create payment session failed: %w", err)
	}
	
	// 解析响应
	var sessionResp PaymentSessionResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &sessionResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if sessionResp.Result.Result != "S" {
		return &sessionResp, fmt.Errorf("create payment session business error: [%s] %s",
			sessionResp.Result.Code, sessionResp.Result.Message)
	}
	
	return &sessionResp, nil
}

// CreateEasySafePaySession 创建钱包授权支付会话（第一次支付）
func (api *PaymentAPI) CreateEasySafePaySession(req *EasySafePaySessionRequest) (*PaymentSessionResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/gateway/v1/acquiring/easy_safe_pay/pay_session")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"mch_app_id":           req.MchAppID,
		"payment_request_id":   req.PaymentRequestID,
		"payment_currency":     req.PaymentCurrency,
		"payment_amount":       req.PaymentAmount,
		"payment_method_type":  req.PaymentMethodType,
		"payment_redirect_url": req.PaymentRedirectURL,
		"order":                req.Order,
		"settlement_currency":  req.SettlementCurrency,
		"product_scene":        req.ProductScene,
		"auth_state":           req.AuthState,
	}
	
	// 可选字段
	if req.PaymentSessionExpiryTime != "" {
		request.Body["payment_session_expiry_time"] = req.PaymentSessionExpiryTime
	}
	if req.EnvClientIP != "" {
		request.Body["env_client_ip"] = req.EnvClientIP
	}
	if req.UserLoginID != "" {
		request.Body["user_login_id"] = req.UserLoginID
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("create easy safe pay session failed: %w", err)
	}
	
	// 解析响应
	var sessionResp PaymentSessionResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &sessionResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if sessionResp.Result.Result != "S" {
		return &sessionResp, fmt.Errorf("create easy safe pay session business error: [%s] %s",
			sessionResp.Result.Code, sessionResp.Result.Message)
	}
	
	return &sessionResp, nil
}

// EasySafePayPay 钱包授权支付（第二次支付 - 使用access_token）
func (api *PaymentAPI) EasySafePayPay(req *EasySafePayRequest) (*PaymentResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/gateway/v1/acquiring/easy_safe_pay/pay")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"mch_app_id":           req.MchAppID,
		"payment_request_id":   req.PaymentRequestID,
		"payment_currency":     req.PaymentCurrency,
		"payment_amount":       req.PaymentAmount,
		"payment_method_id":    req.PaymentMethodID,
		"payment_method_type":  req.PaymentMethodType,
		"payment_redirect_url": req.PaymentRedirectURL,
		"order":                req.Order,
		"settlement_currency":  req.SettlementCurrency,
		"env_terminal_type":    req.EnvTerminalType,
	}
	
	// 可选字段
	if req.EnvClientIP != "" {
		request.Body["env_client_ip"] = req.EnvClientIP
	}
	if req.PaymentExpiryTime != "" {
		request.Body["payment_expiry_time"] = req.PaymentExpiryTime
	}
	if req.EnvOsType != "" {
		request.Body["env_os_type"] = req.EnvOsType
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("easy safe pay payment failed: %w", err)
	}
	
	// 解析响应
	var payResp PaymentResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &payResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if payResp.Result.Result != "S" {
		return &payResp, fmt.Errorf("easy safe pay payment business error: [%s] %s",
			payResp.Result.Code, payResp.Result.Message)
	}
	
	return &payResp, nil
}

// CreateCardAutoDebitSession 创建卡授权支付会话（第一次支付）
func (api *PaymentAPI) CreateCardAutoDebitSession(req *PaymentSessionRequest) (*PaymentSessionResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/gateway/v1/acquiring/card_auto_debit/pay_session")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"mch_app_id":           req.MchAppID,
		"payment_request_id":   req.PaymentRequestID,
		"payment_currency":     req.PaymentCurrency,
		"payment_amount":       req.PaymentAmount,
		"payment_method_type":  req.PaymentMethodType,
		"payment_redirect_url": req.PaymentRedirectURL,
		"order":                req.Order,
		"settlement_currency":  req.SettlementCurrency,
		"product_scene":        req.ProductScene,
	}
	
	// 可选字段
	if req.PaymentSessionExpiryTime != "" {
		request.Body["payment_session_expiry_time"] = req.PaymentSessionExpiryTime
	}
	if req.EnvClientIP != "" {
		request.Body["env_client_ip"] = req.EnvClientIP
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("create card auto debit session failed: %w", err)
	}
	
	// 解析响应
	var sessionResp PaymentSessionResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &sessionResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if sessionResp.Result.Result != "S" {
		return &sessionResp, fmt.Errorf("create card auto debit session business error: [%s] %s",
			sessionResp.Result.Code, sessionResp.Result.Message)
	}
	
	return &sessionResp, nil
}

// CardAutoDebitPay 卡授权支付（第二次支付 - 使用card_token）
func (api *PaymentAPI) CardAutoDebitPay(req *EasySafePayRequest) (*PaymentResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/gateway/v1/acquiring/card_auto_debit/pay")
	
	// 设置请求体（与EasySafePayPay相同，只是payment_method_id传入card_token）
	request.Body = map[string]interface{}{
		"mch_app_id":           req.MchAppID,
		"payment_request_id":   req.PaymentRequestID,
		"payment_currency":     req.PaymentCurrency,
		"payment_amount":       req.PaymentAmount,
		"payment_method_id":    req.PaymentMethodID, // 这里是card_token
		"payment_method_type":  req.PaymentMethodType,
		"payment_redirect_url": req.PaymentRedirectURL,
		"order":                req.Order,
		"settlement_currency":  req.SettlementCurrency,
		"env_terminal_type":    req.EnvTerminalType,
	}
	
	// 可选字段
	if req.EnvClientIP != "" {
		request.Body["env_client_ip"] = req.EnvClientIP
	}
	if req.PaymentExpiryTime != "" {
		request.Body["payment_expiry_time"] = req.PaymentExpiryTime
	}
	if req.EnvOsType != "" {
		request.Body["env_os_type"] = req.EnvOsType
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("card auto debit payment failed: %w", err)
	}
	
	// 解析响应
	var payResp PaymentResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &payResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if payResp.Result.Result != "S" {
		return &payResp, fmt.Errorf("card auto debit payment business error: [%s] %s",
			payResp.Result.Code, payResp.Result.Message)
	}
	
	return &payResp, nil
}

// RefreshAuthToken 刷新授权令牌
func (api *PaymentAPI) RefreshAuthToken(req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/gateway/v1/acquiring/auth_refresh_token")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"mch_app_id":    req.MchAppID,
		"refresh_token": req.RefreshToken,
	}
	
	// 可选字段
	if req.MerchantRegion != "" {
		request.Body["merchant_region"] = req.MerchantRegion
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("refresh auth token failed: %w", err)
	}
	
	// 解析响应
	var refreshResp RefreshTokenResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &refreshResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if refreshResp.Result.Result != "S" {
		return &refreshResp, fmt.Errorf("refresh auth token business error: [%s] %s",
			refreshResp.Result.Code, refreshResp.Result.Message)
	}
	
	return &refreshResp, nil
}

// RevokeAuthToken 取消授权令牌
func (api *PaymentAPI) RevokeAuthToken(req *RevokeTokenRequest) (*RevokeTokenResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/gateway/v1/acquiring/auth_revoke_token")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"mch_app_id":    req.MchAppID,
		"access_token":  req.AccessToken,
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("revoke auth token failed: %w", err)
	}
	
	// 解析响应
	var revokeResp RevokeTokenResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &revokeResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果（注意：取消授权可能返回S/F/U）
	if revokeResp.Result.Result == "F" {
		return &revokeResp, fmt.Errorf("revoke auth token business error: [%s] %s",
			revokeResp.Result.Code, revokeResp.Result.Message)
	}
	
	return &revokeResp, nil
}

// CancelPayment 取消支付
func (api *PaymentAPI) CancelPayment(req *CancelPaymentRequest) (*CancelPaymentResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/gateway/v1/acquiring/cancel")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"mch_app_id":         req.MchAppID,
		"payment_request_id": req.PaymentRequestID,
		"payment_id":         req.PaymentID,
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("cancel payment failed: %w", err)
	}
	
	// 解析响应
	var cancelResp CancelPaymentResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &cancelResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果（注意：取消支付可能返回S/F/U）
	if cancelResp.Result.Result == "F" {
		return &cancelResp, fmt.Errorf("cancel payment business error: [%s] %s",
			cancelResp.Result.Code, cancelResp.Result.Message)
	}
	
	return &cancelResp, nil
}

// QueryPayment 查询支付状态
func (api *PaymentAPI) QueryPayment(req *QueryPaymentRequest) (*QueryPaymentResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/gateway/v1/acquiring/query")
	
	// 设置查询参数
	request.QueryArgs["mch_app_id"] = req.MchAppID
	if req.PaymentRequestID != "" {
		request.QueryArgs["payment_request_id"] = req.PaymentRequestID
	}
	if req.PaymentID != "" {
		request.QueryArgs["payment_id"] = req.PaymentID
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("query payment failed: %w", err)
	}
	
	// 解析响应
	var queryResp QueryPaymentResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &queryResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if queryResp.Result.Result != "S" {
		return &queryResp, fmt.Errorf("query payment business error: [%s] %s",
			queryResp.Result.Code, queryResp.Result.Message)
	}
	
	return &queryResp, nil
}
