package api

import (
	"encoding/json"
	"fmt"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// PayerAPI 付款人API接口
type PayerAPI struct {
	client *gsalary.GSalaryClient
}

// NewPayerAPI 创建付款人API实例
func NewPayerAPI(client *gsalary.GSalaryClient) *PayerAPI {
	return &PayerAPI{client: client}
}

// UploadAttachment 上传附件
func (api *PayerAPI) UploadAttachment(req *UploadAttachmentRequest) (*UploadAttachmentResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/attachments")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"type":     req.Type,
		"filename": req.Filename,
		"base64":   req.Base64,
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("upload attachment failed: %w", err)
	}
	
	// 解析响应
	var uploadResp UploadAttachmentResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &uploadResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if uploadResp.Result.Result != "S" {
		return &uploadResp, fmt.Errorf("upload attachment business error: [%s] %s",
			uploadResp.Result.Code, uploadResp.Result.Message)
	}
	
	return &uploadResp, nil
}

// AddPayer 新增付款人
func (api *PayerAPI) AddPayer(req *PayerRequest) (*PayerResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/remittance/payers")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"subject_type": req.SubjectType,
		"cert_type":    req.CertType,
		"cert_number":  req.CertNumber,
		"cert_files":   req.CertFiles,
		"region":       req.Region,
		"address":      req.Address,
	}
	
	// 个人类型必填字段
	if req.FirstName != "" {
		request.Body["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		request.Body["last_name"] = req.LastName
	}
	if req.Birthday != "" {
		request.Body["birthday"] = req.Birthday
	}
	
	// 企业类型必填字段
	if req.CompanyName != "" {
		request.Body["company_name"] = req.CompanyName
	}
	if req.RegisterNumber != "" {
		request.Body["register_number"] = req.RegisterNumber
	}
	if len(req.BusinessScopes) > 0 {
		request.Body["business_scopes"] = req.BusinessScopes
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("add payer failed: %w", err)
	}
	
	// 解析响应
	var payerResp PayerResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &payerResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if payerResp.Result.Result != "S" {
		return &payerResp, fmt.Errorf("add payer business error: [%s] %s",
			payerResp.Result.Code, payerResp.Result.Message)
	}
	
	return &payerResp, nil
}

// GetPayerList 列出付款人列表
func (api *PayerAPI) GetPayerList() (*PayerListResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/remittance/payers")
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get payer list failed: %w", err)
	}
	
	// 解析响应
	var listResp PayerListResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &listResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if listResp.Result.Result != "S" {
		return &listResp, fmt.Errorf("get payer list business error: [%s] %s",
			listResp.Result.Code, listResp.Result.Message)
	}
	
	return &listResp, nil
}

// GetPayer 查看付款人详情
func (api *PayerAPI) GetPayer(payerID string) (*PayerResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", fmt.Sprintf("/remittance/payers/%s", payerID))
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get payer failed: %w", err)
	}
	
	// 解析响应
	var payerResp PayerResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &payerResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if payerResp.Result.Result != "S" {
		return &payerResp, fmt.Errorf("get payer business error: [%s] %s",
			payerResp.Result.Code, payerResp.Result.Message)
	}
	
	return &payerResp, nil
}

// UpdatePayer 更新付款人信息
func (api *PayerAPI) UpdatePayer(payerID string, req *PayerRequest) (*PayerResponse, error) {
	// 创建PUT请求
	request := gsalary.NewRequest("PUT", fmt.Sprintf("/remittance/payers/%s", payerID))
	
	// 设置请求体（全量更新）
	request.Body = map[string]interface{}{
		"subject_type": req.SubjectType,
		"cert_type":    req.CertType,
		"cert_number":  req.CertNumber,
		"cert_files":   req.CertFiles,
		"region":       req.Region,
		"address":      req.Address,
	}
	
	// 个人类型必填字段
	if req.FirstName != "" {
		request.Body["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		request.Body["last_name"] = req.LastName
	}
	if req.Birthday != "" {
		request.Body["birthday"] = req.Birthday
	}
	
	// 企业类型必填字段
	if req.CompanyName != "" {
		request.Body["company_name"] = req.CompanyName
	}
	if req.RegisterNumber != "" {
		request.Body["register_number"] = req.RegisterNumber
	}
	if len(req.BusinessScopes) > 0 {
		request.Body["business_scopes"] = req.BusinessScopes
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("update payer failed: %w", err)
	}
	
	// 解析响应
	var payerResp PayerResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &payerResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if payerResp.Result.Result != "S" {
		return &payerResp, fmt.Errorf("update payer business error: [%s] %s",
			payerResp.Result.Code, payerResp.Result.Message)
	}
	
	return &payerResp, nil
}

// DeletePayer 移除付款人信息
func (api *PayerAPI) DeletePayer(payerID string) error {
	// 创建DELETE请求
	request := gsalary.NewRequest("DELETE", fmt.Sprintf("/remittance/payers/%s", payerID))
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return fmt.Errorf("delete payer failed: %w", err)
	}
	
	// 解析响应
	var result struct {
		Result struct {
			Result  string `json:"result"`
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"result"`
	}
	
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果（S=成功，F=失败，U=未知需后续查询）
	if result.Result.Result == "F" {
		return fmt.Errorf("delete payer business error: [%s] %s",
			result.Result.Code, result.Result.Message)
	}
	
	return nil
}
