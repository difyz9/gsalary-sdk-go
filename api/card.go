package api

import (
	"encoding/json"
	"fmt"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// CardAPI 卡片API接口
type CardAPI struct {
	client *gsalary.GSalaryClient
}

// NewCardAPI 创建卡片API实例
func NewCardAPI(client *gsalary.GSalaryClient) *CardAPI {
	return &CardAPI{client: client}
}

// ApplyCard 申请新卡片
func (api *CardAPI) ApplyCard(req *CardApplyRequest) (*CardApplyResponse, error) {
	// 创建请求
	request := gsalary.NewRequest("POST", "/v1/card_applies")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"request_id":      req.RequestID,
		"product_code":    req.ProductCode,
		"currency":        req.Currency,
		"card_holder_id":  req.CardHolderID,
		"init_balance":    req.InitBalance,
	}
	
	// 可选参数
	if req.LimitPerDay > 0 {
		request.Body["limit_per_day"] = req.LimitPerDay
	}
	if req.LimitPerMonth > 0 {
		request.Body["limit_per_month"] = req.LimitPerMonth
	}
	if req.LimitPerTransaction > 0 {
		request.Body["limit_per_transaction"] = req.LimitPerTransaction
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("apply card failed: %w", err)
	}
	
	// 解析响应
	var cardResp CardApplyResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &cardResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if cardResp.Result.Result != "S" {
		return &cardResp, fmt.Errorf("apply card business error: [%s] %s", 
			cardResp.Result.Code, cardResp.Result.Message)
	}
	
	return &cardResp, nil
}

// GetAvailableQuotas 查询卡可用余额
func (api *CardAPI) GetAvailableQuotas(req *CardAvailableQuotasRequest) (*CardAvailableQuotasResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/v1/cards/available_quotas")
	
	// 设置查询参数
	request.QueryArgs["currency"] = req.Currency
	
	// 如果指定了卡账务类型，添加到查询参数
	if req.AccountingCardType != "" {
		request.QueryArgs["accounting_card_type"] = req.AccountingCardType
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get available quotas failed: %w", err)
	}
	
	// 解析响应
	var quotasResp CardAvailableQuotasResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &quotasResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if quotasResp.Result.Result != "S" {
		return &quotasResp, fmt.Errorf("get available quotas business error: [%s] %s",
			quotasResp.Result.Code, quotasResp.Result.Message)
	}
	
	return &quotasResp, nil
}

// GetProducts 查询可用的卡产品列表
func (api *CardAPI) GetProducts(req *CardProductsRequest) (*CardProductsResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/v1/card_support/products")
	
	// 设置查询参数
	if req.CardType != "" {
		request.QueryArgs["card_type"] = req.CardType
	}
	if req.BrandCode != "" {
		request.QueryArgs["brand_code"] = req.BrandCode
	}
	if req.Currency != "" {
		request.QueryArgs["currency"] = req.Currency
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get card products failed: %w", err)
	}
	
	// 解析响应
	var productsResp CardProductsResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &productsResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if productsResp.Result.Result != "S" {
		return &productsResp, fmt.Errorf("get card products business error: [%s] %s",
			productsResp.Result.Code, productsResp.Result.Message)
	}
	
	return &productsResp, nil
}

// GetCardApplyResult 查询开卡结果
func (api *CardAPI) GetCardApplyResult(requestID string) (*CardApplyResultResponse, error) {
	if requestID == "" {
		return nil, fmt.Errorf("request_id is required")
	}
	
	// 创建GET请求，路径中包含request_id
	path := fmt.Sprintf("/v1/card_applies/%s", requestID)
	request := gsalary.NewRequest("GET", path)
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get card apply result failed: %w", err)
	}
	
	// 解析响应
	var resultResp CardApplyResultResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &resultResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if resultResp.Result.Result != "S" {
		return &resultResp, fmt.Errorf("get card apply result business error: [%s] %s",
			resultResp.Result.Code, resultResp.Result.Message)
	}
	
	return &resultResp, nil
}

// GetCardList 查询卡列表
func (api *CardAPI) GetCardList(req *CardListRequest) (*CardListResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/v1/cards")
	
	// 设置分页参数
	if req.Page > 0 {
		request.QueryArgs["page"] = fmt.Sprintf("%d", req.Page)
	}
	if req.Limit > 0 {
		request.QueryArgs["limit"] = fmt.Sprintf("%d", req.Limit)
	}
	
	// 设置过滤参数
	if req.ProductCode != "" {
		request.QueryArgs["product_code"] = req.ProductCode
	}
	if req.BrandCode != "" {
		request.QueryArgs["brand_code"] = req.BrandCode
	}
	if req.CardHolderID != "" {
		request.QueryArgs["card_holder_id"] = req.CardHolderID
	}
	if req.CreateStart != "" {
		request.QueryArgs["create_start"] = req.CreateStart
	}
	if req.CreateEnd != "" {
		request.QueryArgs["create_end"] = req.CreateEnd
	}
	if req.Status != "" {
		request.QueryArgs["status"] = req.Status
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get card list failed: %w", err)
	}
	
	// 解析响应
	var listResp CardListResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &listResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if listResp.Result.Result != "S" {
		return &listResp, fmt.Errorf("get card list business error: [%s] %s",
			listResp.Result.Code, listResp.Result.Message)
	}
	
	return &listResp, nil
}

// GetCardInfo 查看卡信息
func (api *CardAPI) GetCardInfo(cardID string) (*CardInfoResponse, error) {
	if cardID == "" {
		return nil, fmt.Errorf("card_id is required")
	}
	
	// 创建GET请求，路径中包含card_id
	path := fmt.Sprintf("/v1/cards/%s", cardID)
	request := gsalary.NewRequest("GET", path)
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get card info failed: %w", err)
	}
	
	// 解析响应
	var infoResp CardInfoResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &infoResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if infoResp.Result.Result != "S" {
		return &infoResp, fmt.Errorf("get card info business error: [%s] %s",
			infoResp.Result.Code, infoResp.Result.Message)
	}
	
	return &infoResp, nil
}

// UpdateCard 修改卡信息
func (api *CardAPI) UpdateCard(req *UpdateCardRequest) (*UpdateCardResponse, error) {
	// 创建PUT请求
	request := gsalary.NewRequest("PUT", fmt.Sprintf("/v1/cards/%s", req.CardID))
	
	// 设置请求体
	body := make(map[string]interface{})
	if req.CardName != "" {
		body["card_name"] = req.CardName
	}
	if req.LimitPerDay > 0 {
		body["limit_per_day"] = req.LimitPerDay
	}
	if req.LimitPerMonth > 0 {
		body["limit_per_month"] = req.LimitPerMonth
	}
	if req.LimitPerTransaction > 0 {
		body["limit_per_transaction"] = req.LimitPerTransaction
	}
	request.Body = body
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("update card failed: %w", err)
	}
	
	// 解析响应
	var updateResp UpdateCardResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &updateResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if updateResp.Result.Result != "S" {
		return &updateResp, fmt.Errorf("update card business error: [%s] %s",
			updateResp.Result.Code, updateResp.Result.Message)
	}
	
	return &updateResp, nil
}

// DeleteCard 销卡
func (api *CardAPI) DeleteCard(cardID string) (*DeleteCardResponse, error) {
	// 创建DELETE请求
	request := gsalary.NewRequest("DELETE", fmt.Sprintf("/v1/cards/%s", cardID))
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("delete card failed: %w", err)
	}
	
	// 解析响应
	var deleteResp DeleteCardResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &deleteResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果（注意：可能返回S/F/U）
	if deleteResp.Result.Result == "F" {
		return &deleteResp, fmt.Errorf("delete card business error: [%s] %s",
			deleteResp.Result.Code, deleteResp.Result.Message)
	}
	
	return &deleteResp, nil
}

// GetCardSecureInfo 获取卡机密信息（PAN、CVV、有效期）
func (api *CardAPI) GetCardSecureInfo(cardID string) (*CardSecureInfoResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", fmt.Sprintf("/v1/cards/%s/secure_info", cardID))
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get card secure info failed: %w", err)
	}
	
	// 解析响应
	var secureResp CardSecureInfoResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &secureResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if secureResp.Result.Result != "S" {
		return &secureResp, fmt.Errorf("get card secure info business error: [%s] %s",
			secureResp.Result.Code, secureResp.Result.Message)
	}
	
	return &secureResp, nil
}

// AdjustCardBalance 卡片调额（增加或减少余额）
func (api *CardAPI) AdjustCardBalance(req *AdjustCardBalanceRequest) (*AdjustCardBalanceResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/v1/cards/balance_modifies")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"card_id":    req.CardID,
		"amount":     req.Amount,
		"type":       req.Type,
		"request_id": req.RequestID,
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("adjust card balance failed: %w", err)
	}
	
	// 解析响应
	var adjustResp AdjustCardBalanceResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &adjustResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if adjustResp.Result.Result != "S" {
		return &adjustResp, fmt.Errorf("adjust card balance business error: [%s] %s",
			adjustResp.Result.Code, adjustResp.Result.Message)
	}
	
	return &adjustResp, nil
}

// GetBalanceModifyResult 查询卡片调额结果
func (api *CardAPI) GetBalanceModifyResult(requestID string) (*GetBalanceModifyResultResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", fmt.Sprintf("/v1/cards/balance_modifies/%s", requestID))
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get balance modify result failed: %w", err)
	}
	
	// 解析响应
	var resultResp GetBalanceModifyResultResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &resultResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if resultResp.Result.Result != "S" {
		return &resultResp, fmt.Errorf("get balance modify result business error: [%s] %s",
			resultResp.Result.Code, resultResp.Result.Message)
	}
	
	return &resultResp, nil
}

// SetCardFreezeStatus 冻结/解冻卡
func (api *CardAPI) SetCardFreezeStatus(req *SetCardFreezeStatusRequest) (*SetCardFreezeStatusResponse, error) {
	// 创建PUT请求
	request := gsalary.NewRequest("PUT", fmt.Sprintf("/v1/cards/%s/freeze_status", req.CardID))
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"freeze": req.Freeze,
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("set card freeze status failed: %w", err)
	}
	
	// 解析响应
	var freezeResp SetCardFreezeStatusResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &freezeResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果（注意：可能返回S/F/U）
	if freezeResp.Result.Result == "F" {
		return &freezeResp, fmt.Errorf("set card freeze status business error: [%s] %s",
			freezeResp.Result.Code, freezeResp.Result.Message)
	}
	
	return &freezeResp, nil
}

// GetCardTransactions 查询卡交易列表
func (api *CardAPI) GetCardTransactions(req *CardTransactionsRequest) (*CardTransactionsResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/v1/card_bill/card_transactions")
	
	// 设置查询参数
	request.QueryArgs["page"] = fmt.Sprintf("%d", req.Page)
	request.QueryArgs["limit"] = fmt.Sprintf("%d", req.Limit)
	
	if req.TransactionID != "" {
		request.QueryArgs["transaction_id"] = req.TransactionID
	}
	if req.MchRequestID != "" {
		request.QueryArgs["mch_request_id"] = req.MchRequestID
	}
	if req.TimeStart != "" {
		request.QueryArgs["time_start"] = req.TimeStart
	}
	if req.TimeEnd != "" {
		request.QueryArgs["time_end"] = req.TimeEnd
	}
	if req.CardID != "" {
		request.QueryArgs["card_id"] = req.CardID
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get card transactions failed: %w", err)
	}
	
	// 解析响应
	var txResp CardTransactionsResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &txResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if txResp.Result.Result != "S" {
		return &txResp, fmt.Errorf("get card transactions business error: [%s] %s",
			txResp.Result.Code, txResp.Result.Message)
	}
	
	return &txResp, nil
}

// GetBalanceHistory 查询卡余额变更记录
func (api *CardAPI) GetBalanceHistory(req *BalanceHistoryRequest) (*BalanceHistoryResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/v1/card_bill/balance_history")
	
	// 设置查询参数
	request.QueryArgs["page"] = fmt.Sprintf("%d", req.Page)
	request.QueryArgs["limit"] = fmt.Sprintf("%d", req.Limit)
	
	if req.TransactionID != "" {
		request.QueryArgs["transaction_id"] = req.TransactionID
	}
	if req.LogID != "" {
		request.QueryArgs["log_id"] = req.LogID
	}
	if req.TimeStart != "" {
		request.QueryArgs["time_start"] = req.TimeStart
	}
	if req.TimeEnd != "" {
		request.QueryArgs["time_end"] = req.TimeEnd
	}
	if req.CardID != "" {
		request.QueryArgs["card_id"] = req.CardID
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get balance history failed: %w", err)
	}
	
	// 解析响应
	var historyResp BalanceHistoryResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &historyResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if historyResp.Result.Result != "S" {
		return &historyResp, fmt.Errorf("get balance history business error: [%s] %s",
			historyResp.Result.Code, historyResp.Result.Message)
	}
	
	return &historyResp, nil
}

// UpdateCardContact 修改卡联系信息（email用于ApplePay绑卡验证）
func (api *CardAPI) UpdateCardContact(req *UpdateCardContactRequest) (*UpdateCardContactResponse, error) {
	// 创建PUT请求
	request := gsalary.NewRequest("PUT", fmt.Sprintf("/v1/cards/%s/contact", req.CardID))
	
	// 设置请求体
	body := make(map[string]interface{})
	if req.Email != "" {
		body["email"] = req.Email
	}
	request.Body = body
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("update card contact failed: %w", err)
	}
	
	// 解析响应
	var contactResp UpdateCardContactResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &contactResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if contactResp.Result.Result != "S" {
		return &contactResp, fmt.Errorf("update card contact business error: [%s] %s",
			contactResp.Result.Code, contactResp.Result.Message)
	}
	
	return &contactResp, nil
}
