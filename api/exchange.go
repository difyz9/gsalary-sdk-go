package api

import (
	"encoding/json"
	"fmt"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// ExchangeAPI 换汇API接口
type ExchangeAPI struct {
	client *gsalary.GSalaryClient
}

// NewExchangeAPI 创建换汇API实例
func NewExchangeAPI(client *gsalary.GSalaryClient) *ExchangeAPI {
	return &ExchangeAPI{client: client}
}

// GetCurrentExchangeRate 查询当前汇率
func (api *ExchangeAPI) GetCurrentExchangeRate(req *ExchangeRateRequest) (*ExchangeRateResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/v1/exchange/current_exchange_rate")
	
	// 设置查询参数
	request.QueryArgs["buy_currency"] = req.BuyCurrency
	request.QueryArgs["sell_currency"] = req.SellCurrency
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get exchange rate failed: %w", err)
	}
	
	// 解析响应
	var rateResp ExchangeRateResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &rateResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if rateResp.Result.Result != "S" {
		return &rateResp, fmt.Errorf("get exchange rate business error: [%s] %s",
			rateResp.Result.Code, rateResp.Result.Message)
	}
	
	return &rateResp, nil
}

// RequestQuote 请求锁汇报价
func (api *ExchangeAPI) RequestQuote(req *ExchangeQuoteRequest) (*ExchangeQuoteResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/v1/exchange/quotes")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"buy_currency":  req.BuyCurrency,
		"sell_currency": req.SellCurrency,
	}
	
	// 购入金额和卖出金额不可同时为空，如果同时提供将忽略购入金额
	if req.SellAmount > 0 {
		request.Body["sell_amount"] = req.SellAmount
	} else if req.BuyAmount > 0 {
		request.Body["buy_amount"] = req.BuyAmount
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("request exchange quote failed: %w", err)
	}
	
	// 解析响应
	var quoteResp ExchangeQuoteResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &quoteResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if quoteResp.Result.Result != "S" {
		return &quoteResp, fmt.Errorf("request exchange quote business error: [%s] %s",
			quoteResp.Result.Code, quoteResp.Result.Message)
	}
	
	return &quoteResp, nil
}

// SubmitExchangeRequest 提交换汇订单
func (api *ExchangeAPI) SubmitExchangeRequest(req *ExchangeSubmitRequest) (*ExchangeSubmitResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/v1/exchange/submit_request")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"request_id": req.RequestID,
		"quote_id":   req.QuoteID,
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("submit exchange request failed: %w", err)
	}
	
	// 解析响应
	var submitResp ExchangeSubmitResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &submitResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if submitResp.Result.Result != "S" {
		return &submitResp, fmt.Errorf("submit exchange request business error: [%s] %s",
			submitResp.Result.Code, submitResp.Result.Message)
	}
	
	return &submitResp, nil
}

// GetExchangeOrders 查询换汇订单列表
func (api *ExchangeAPI) GetExchangeOrders(req *ExchangeOrdersRequest) (*ExchangeOrdersResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/v1/exchange/orders")
	
	// 设置查询参数
	if req.Page > 0 {
		request.QueryArgs["page"] = fmt.Sprintf("%d", req.Page)
	}
	if req.Limit > 0 {
		request.QueryArgs["limit"] = fmt.Sprintf("%d", req.Limit)
	}
	if req.TimeStart != "" {
		request.QueryArgs["time_start"] = req.TimeStart
	}
	if req.TimeEnd != "" {
		request.QueryArgs["time_end"] = req.TimeEnd
	}
	if req.Status != "" {
		request.QueryArgs["status"] = req.Status
	}
	if req.BuyCurrency != "" {
		request.QueryArgs["buy_currency"] = req.BuyCurrency
	}
	if req.SellCurrency != "" {
		request.QueryArgs["sell_currency"] = req.SellCurrency
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get exchange orders failed: %w", err)
	}
	
	// 解析响应
	var ordersResp ExchangeOrdersResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &ordersResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if ordersResp.Result.Result != "S" {
		return &ordersResp, fmt.Errorf("get exchange orders business error: [%s] %s",
			ordersResp.Result.Code, ordersResp.Result.Message)
	}
	
	return &ordersResp, nil
}
