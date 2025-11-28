package api

import (
	"encoding/json"
	"fmt"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// RemittanceAPI 对外付款API接口
type RemittanceAPI struct {
	client *gsalary.GSalaryClient
}

// NewRemittanceAPI 创建对外付款API实例
func NewRemittanceAPI(client *gsalary.GSalaryClient) *RemittanceAPI {
	return &RemittanceAPI{client: client}
}

// GetClearingNetworks 查询可用清算网络
func (api *RemittanceAPI) GetClearingNetworks(req *ClearingNetworkRequest) (*ClearingNetworkResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/remittance/clearing_networks")
	
	// 设置查询参数
	request.QueryArgs["payee_account_id"] = req.PayeeAccountID
	request.QueryArgs["pay_currency"] = req.PayCurrency
	request.QueryArgs["amount"] = fmt.Sprintf("%f", req.Amount)
	request.QueryArgs["amount_type"] = req.AmountType
	request.QueryArgs["receive_currency"] = req.ReceiveCurrency
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get clearing networks failed: %w", err)
	}
	
	// 解析响应
	var networkResp ClearingNetworkResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &networkResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if networkResp.Result.Result != "S" {
		return &networkResp, fmt.Errorf("get clearing networks business error: [%s] %s",
			networkResp.Result.Code, networkResp.Result.Message)
	}
	
	return &networkResp, nil
}

// CreateQuote 申请锁汇
func (api *RemittanceAPI) CreateQuote(req *QuoteRequest) (*QuoteResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/remittance/quotes")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"payee_account_id": req.PayeeAccountID,
		"purpose":          req.Purpose,
		"pay_currency":     req.PayCurrency,
		"amount":           req.Amount,
		"amount_type":      req.AmountType,
	}
	
	// 可选字段
	if req.PayerID != "" {
		request.Body["payer_id"] = req.PayerID
	}
	if req.ReceiveCurrency != "" {
		request.Body["receive_currency"] = req.ReceiveCurrency
	}
	if req.ClearingNetwork != "" {
		request.Body["clearing_network"] = req.ClearingNetwork
	}
	if req.AbaNumber != "" {
		request.Body["aba_number"] = req.AbaNumber
	}
	if req.FpsBankID != "" {
		request.Body["fps_bank_id"] = req.FpsBankID
	}
	if req.IfsCode != "" {
		request.Body["ifs_code"] = req.IfsCode
	}
	if req.IntermediarySwiftCode != "" {
		request.Body["intermediary_swift_code"] = req.IntermediarySwiftCode
	}
	if req.Remark != "" {
		request.Body["remark"] = req.Remark
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("create quote failed: %w", err)
	}
	
	// 解析响应
	var quoteResp QuoteResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &quoteResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if quoteResp.Result.Result != "S" {
		return &quoteResp, fmt.Errorf("create quote business error: [%s] %s",
			quoteResp.Result.Code, quoteResp.Result.Message)
	}
	
	return &quoteResp, nil
}

// SubmitOrder 提交付款订单
func (api *RemittanceAPI) SubmitOrder(req *OrderRequest) (*OrderResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/remittance/orders")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"quote_id":        req.QuoteID,
		"client_order_id": req.ClientOrderID,
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("submit order failed: %w", err)
	}
	
	// 解析响应
	var orderResp OrderResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &orderResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if orderResp.Result.Result != "S" {
		return &orderResp, fmt.Errorf("submit order business error: [%s] %s",
			orderResp.Result.Code, orderResp.Result.Message)
	}
	
	return &orderResp, nil
}

// GetOrderList 查询付款单列表
func (api *RemittanceAPI) GetOrderList(req *OrderListRequest) (*OrderListResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/remittance/orders")
	
	// 设置查询参数
	request.QueryArgs["page"] = fmt.Sprintf("%d", req.Page)
	request.QueryArgs["limit"] = fmt.Sprintf("%d", req.Limit)
	
	if req.PayeeID != "" {
		request.QueryArgs["payee_id"] = req.PayeeID
	}
	if req.PayerID != "" {
		request.QueryArgs["payer_id"] = req.PayerID
	}
	if req.TimeStart != "" {
		request.QueryArgs["time_start"] = req.TimeStart
	}
	if req.TimeEnd != "" {
		request.QueryArgs["time_end"] = req.TimeEnd
	}
	if req.OrderID != "" {
		request.QueryArgs["order_id"] = req.OrderID
	}
	if req.ClientOrderID != "" {
		request.QueryArgs["client_order_id"] = req.ClientOrderID
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get order list failed: %w", err)
	}
	
	// 解析响应
	var listResp OrderListResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &listResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if listResp.Result.Result != "S" {
		return &listResp, fmt.Errorf("get order list business error: [%s] %s",
			listResp.Result.Code, listResp.Result.Message)
	}
	
	return &listResp, nil
}
