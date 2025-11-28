package api

import (
	"encoding/json"
	"fmt"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// CardHolderAPI 持卡人API接口
type CardHolderAPI struct {
	client *gsalary.GSalaryClient
}

// NewCardHolderAPI 创建持卡人API实例
func NewCardHolderAPI(client *gsalary.GSalaryClient) *CardHolderAPI {
	return &CardHolderAPI{client: client}
}

// AddCardHolder 添加持卡人
func (api *CardHolderAPI) AddCardHolder(req *CardHolderRequest) (*CardHolderResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/v1/card_holders")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"birth":      req.Birth,
		"email":      req.Email,
		"mobile": map[string]interface{}{
			"country_code": req.Mobile.CountryCode,
			"number":       req.Mobile.Number,
		},
		"region": req.Region,
		"bill_address": map[string]interface{}{
			"country":     req.BillAddress.Country,
			"state":       req.BillAddress.State,
			"city":        req.BillAddress.City,
			"postal_code": req.BillAddress.PostalCode,
			"line1":       req.BillAddress.Line1,
			"line2":       req.BillAddress.Line2,
		},
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("add card holder failed: %w", err)
	}
	
	// 解析响应
	var holderResp CardHolderResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &holderResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if holderResp.Result.Result != "S" {
		return &holderResp, fmt.Errorf("add card holder business error: [%s] %s",
			holderResp.Result.Code, holderResp.Result.Message)
	}
	
	return &holderResp, nil
}

// GetCardHolderList 查询持卡人列表
func (api *CardHolderAPI) GetCardHolderList(req *CardHolderListRequest) (*CardHolderListResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/v1/card_holders")
	
	// 设置分页参数
	if req.Page > 0 {
		request.QueryArgs["page"] = fmt.Sprintf("%d", req.Page)
	}
	if req.Limit > 0 {
		request.QueryArgs["limit"] = fmt.Sprintf("%d", req.Limit)
	}
	
	// 设置时间过滤参数
	if req.TimeStart != "" {
		request.QueryArgs["time_start"] = req.TimeStart
	}
	if req.TimeEnd != "" {
		request.QueryArgs["time_end"] = req.TimeEnd
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get card holder list failed: %w", err)
	}
	
	// 解析响应
	var listResp CardHolderListResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &listResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if listResp.Result.Result != "S" {
		return &listResp, fmt.Errorf("get card holder list business error: [%s] %s",
			listResp.Result.Code, listResp.Result.Message)
	}
	
	return &listResp, nil
}

// GetCardHolderInfo 查看持卡人信息
func (api *CardHolderAPI) GetCardHolderInfo(cardHolderID string) (*CardHolderDetailResponse, error) {
	if cardHolderID == "" {
		return nil, fmt.Errorf("card_holder_id is required")
	}
	
	// 创建GET请求，路径中包含card_holder_id
	path := fmt.Sprintf("/v1/card_holders/%s", cardHolderID)
	request := gsalary.NewRequest("GET", path)
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get card holder info failed: %w", err)
	}
	
	// 解析响应
	var infoResp CardHolderDetailResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &infoResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if infoResp.Result.Result != "S" {
		return &infoResp, fmt.Errorf("get card holder info business error: [%s] %s",
			infoResp.Result.Code, infoResp.Result.Message)
	}
	
	return &infoResp, nil
}

// UpdateCardHolder 修改持卡人信息
func (api *CardHolderAPI) UpdateCardHolder(cardHolderID string, req *UpdateCardHolderRequest) (*UpdateCardHolderResponse, error) {
	if cardHolderID == "" {
		return nil, fmt.Errorf("card_holder_id is required")
	}
	
	// 创建PUT请求，路径中包含card_holder_id
	path := fmt.Sprintf("/v1/card_holders/%s", cardHolderID)
	request := gsalary.NewRequest("PUT", path)
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"first_name": req.FirstName,
		"last_name":  req.LastName,
		"birth":      req.Birth,
		"email":      req.Email,
		"mobile": map[string]interface{}{
			"country_code": req.Mobile.CountryCode,
			"number":       req.Mobile.Number,
		},
		"region": req.Region,
		"bill_address": map[string]interface{}{
			"country":     req.BillAddress.Country,
			"state":       req.BillAddress.State,
			"city":        req.BillAddress.City,
			"postal_code": req.BillAddress.PostalCode,
			"line1":       req.BillAddress.Line1,
			"line2":       req.BillAddress.Line2,
		},
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("update card holder failed: %w", err)
	}
	
	// 解析响应
	var updateResp UpdateCardHolderResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &updateResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if updateResp.Result.Result != "S" {
		return &updateResp, fmt.Errorf("update card holder business error: [%s] %s",
			updateResp.Result.Code, updateResp.Result.Message)
	}
	
	return &updateResp, nil
}
