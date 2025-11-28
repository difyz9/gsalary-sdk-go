package api

import (
	"encoding/json"
	"fmt"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// PayeeAPI 收款人API接口
type PayeeAPI struct {
	client *gsalary.GSalaryClient
}

// NewPayeeAPI 创建收款人API实例
func NewPayeeAPI(client *gsalary.GSalaryClient) *PayeeAPI {
	return &PayeeAPI{client: client}
}

// AddPayee 新增收款人
func (api *PayeeAPI) AddPayee(req *PayeeRequest) (*PayeeResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", "/remittance/payees")
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"subject_type": req.SubjectType,
		"account_type": req.AccountType,
		"country":      req.Country,
		"currency":     req.Currency,
	}
	
	if req.FirstName != "" {
		request.Body["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		request.Body["last_name"] = req.LastName
	}
	if req.AccountHolder != "" {
		request.Body["account_holder"] = req.AccountHolder
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("add payee failed: %w", err)
	}
	
	// 解析响应
	var payeeResp PayeeResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &payeeResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if payeeResp.Result.Result != "S" {
		return &payeeResp, fmt.Errorf("add payee business error: [%s] %s",
			payeeResp.Result.Code, payeeResp.Result.Message)
	}
	
	return &payeeResp, nil
}

// GetPayeeList 查询收款人列表
func (api *PayeeAPI) GetPayeeList(req *PayeeListRequest) (*PayeeListResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/remittance/payees")
	
	// 设置查询参数
	request.QueryArgs["page"] = fmt.Sprintf("%d", req.Page)
	request.QueryArgs["limit"] = fmt.Sprintf("%d", req.Limit)
	
	if req.PayeeID != "" {
		request.QueryArgs["payee_id"] = req.PayeeID
	}
	if req.Name != "" {
		request.QueryArgs["name"] = req.Name
	}
	if req.Country != "" {
		request.QueryArgs["country"] = req.Country
	}
	if req.Currency != "" {
		request.QueryArgs["currency"] = req.Currency
	}
	if req.Mobile != "" {
		request.QueryArgs["mobile"] = req.Mobile
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get payee list failed: %w", err)
	}
	
	// 解析响应
	var listResp PayeeListResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &listResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if listResp.Result.Result != "S" {
		return &listResp, fmt.Errorf("get payee list business error: [%s] %s",
			listResp.Result.Code, listResp.Result.Message)
	}
	
	return &listResp, nil
}

// UpdatePayee 更新收款人信息
func (api *PayeeAPI) UpdatePayee(payeeID string, req *PayeeRequest) (*PayeeResponse, error) {
	// 创建PUT请求
	request := gsalary.NewRequest("PUT", fmt.Sprintf("/remittance/payees/%s", payeeID))
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"account_type": req.AccountType,
		"country":      req.Country,
		"currency":     req.Currency,
	}
	
	if req.FirstName != "" {
		request.Body["first_name"] = req.FirstName
	}
	if req.LastName != "" {
		request.Body["last_name"] = req.LastName
	}
	if req.AccountHolder != "" {
		request.Body["account_holder"] = req.AccountHolder
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("update payee failed: %w", err)
	}
	
	// 解析响应
	var payeeResp PayeeResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &payeeResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if payeeResp.Result.Result != "S" {
		return &payeeResp, fmt.Errorf("update payee business error: [%s] %s",
			payeeResp.Result.Code, payeeResp.Result.Message)
	}
	
	return &payeeResp, nil
}

// DeletePayee 停用收款人
func (api *PayeeAPI) DeletePayee(payeeID string) error {
	// 创建DELETE请求
	request := gsalary.NewRequest("DELETE", fmt.Sprintf("/remittance/payees/%s", payeeID))
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return fmt.Errorf("delete payee failed: %w", err)
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
		return fmt.Errorf("delete payee business error: [%s] %s",
			result.Result.Code, result.Result.Message)
	}
	
	return nil
}

// AddPayeeAccount 新增收款人收款账户（电子钱包）
func (api *PayeeAPI) AddPayeeAccount(payeeID string, req *PayeeAccountRequest) (*PayeeAccountResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", fmt.Sprintf("/remittance/payees/%s/accounts", payeeID))
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"payment_method": req.PaymentMethod,
		"account_no":     req.AccountNo,
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("add payee account failed: %w", err)
	}
	
	// 解析响应
	var accountResp PayeeAccountResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &accountResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if accountResp.Result.Result != "S" {
		return &accountResp, fmt.Errorf("add payee account business error: [%s] %s",
			accountResp.Result.Code, accountResp.Result.Message)
	}
	
	return &accountResp, nil
}

// GetPayeeAccounts 查看收款人可用收款账户
func (api *PayeeAPI) GetPayeeAccounts(payeeID string, language string) (*PayeeAccountsResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", fmt.Sprintf("/remittance/payees/%s/accounts", payeeID))
	
	// 设置查询参数
	if language != "" {
		request.QueryArgs["language"] = language
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get payee accounts failed: %w", err)
	}
	
	// 解析响应
	var accountsResp PayeeAccountsResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &accountsResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if accountsResp.Result.Result != "S" {
		return &accountsResp, fmt.Errorf("get payee accounts business error: [%s] %s",
			accountsResp.Result.Code, accountsResp.Result.Message)
	}
	
	return &accountsResp, nil
}

// UpdatePayeeAccount 更新收款人账户（电子钱包）
func (api *PayeeAPI) UpdatePayeeAccount(payeeID, accountID string, req *PayeeAccountRequest) (*PayeeAccountResponse, error) {
	// 创建PUT请求
	request := gsalary.NewRequest("PUT", fmt.Sprintf("/remittance/payees/%s/payee_accounts/%s", payeeID, accountID))
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"payment_method": req.PaymentMethod,
		"account_no":     req.AccountNo,
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("update payee account failed: %w", err)
	}
	
	// 解析响应
	var accountResp PayeeAccountResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &accountResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if accountResp.Result.Result != "S" {
		return &accountResp, fmt.Errorf("update payee account business error: [%s] %s",
			accountResp.Result.Code, accountResp.Result.Message)
	}
	
	return &accountResp, nil
}

// GetPayeeAccountForm 获取收款人账户表单（银行账户）
func (api *PayeeAPI) GetPayeeAccountForm(payeeID string, req *PayeeAccountFormRequest) (*PayeeAccountFormResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", fmt.Sprintf("/remittance/payees/%s/account_register_format", payeeID))
	
	// 设置查询参数
	request.QueryArgs["payment_method"] = req.PaymentMethod
	if req.Currency != "" {
		request.QueryArgs["currency"] = req.Currency
	}
	if req.Language != "" {
		request.QueryArgs["language"] = req.Language
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get payee account form failed: %w", err)
	}
	
	// 解析响应
	var formResp PayeeAccountFormResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &formResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if formResp.Result.Result != "S" {
		return &formResp, fmt.Errorf("get payee account form business error: [%s] %s",
			formResp.Result.Code, formResp.Result.Message)
	}
	
	return &formResp, nil
}

// AddPayeeAccountBank 新增收款人收款账户（银行账户）
func (api *PayeeAPI) AddPayeeAccountBank(payeeID string, req *PayeeAccountBankRequest) (*PayeeAccountResponse, error) {
	// 创建POST请求
	request := gsalary.NewRequest("POST", fmt.Sprintf("/remittance/payees/%s/account_registry", payeeID))
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"payment_method": req.PaymentMethod,
		"fields":         req.Fields,
	}
	
	if req.Currency != "" {
		request.Body["currency"] = req.Currency
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("add payee account bank failed: %w", err)
	}
	
	// 解析响应
	var accountResp PayeeAccountResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &accountResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if accountResp.Result.Result != "S" {
		return &accountResp, fmt.Errorf("add payee account bank business error: [%s] %s",
			accountResp.Result.Code, accountResp.Result.Message)
	}
	
	return &accountResp, nil
}

// UpdatePayeeAccountBank 更新收款人账户（银行账户）
func (api *PayeeAPI) UpdatePayeeAccountBank(accountID string, req *PayeeAccountBankRequest) (*PayeeAccountResponse, error) {
	// 创建PUT请求
	request := gsalary.NewRequest("PUT", fmt.Sprintf("/remittance/payee_accounts/%s", accountID))
	
	// 设置请求体
	request.Body = map[string]interface{}{
		"fields": req.Fields,
	}
	
	if req.Currency != "" {
		request.Body["currency"] = req.Currency
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("update payee account bank failed: %w", err)
	}
	
	// 解析响应
	var accountResp PayeeAccountResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &accountResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if accountResp.Result.Result != "S" {
		return &accountResp, fmt.Errorf("update payee account bank business error: [%s] %s",
			accountResp.Result.Code, accountResp.Result.Message)
	}
	
	return &accountResp, nil
}

// DeletePayeeAccount 移除收款账户
func (api *PayeeAPI) DeletePayeeAccount(accountID string) error {
	// 创建DELETE请求
	request := gsalary.NewRequest("DELETE", fmt.Sprintf("/remittance/payee_accounts/%s", accountID))
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return fmt.Errorf("delete payee account failed: %w", err)
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
		return fmt.Errorf("delete payee account business error: [%s] %s",
			result.Result.Code, result.Result.Message)
	}
	
	return nil
}

// GetAvailablePaymentMethods 查询可用付款方式
func (api *PayeeAPI) GetAvailablePaymentMethods() (*PaymentMethodsResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/remittance/available_payment_methods")
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get available payment methods failed: %w", err)
	}
	
	// 解析响应
	var methodsResp PaymentMethodsResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &methodsResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if methodsResp.Result.Result != "S" {
		return &methodsResp, fmt.Errorf("get available payment methods business error: [%s] %s",
			methodsResp.Result.Code, methodsResp.Result.Message)
	}
	
	return &methodsResp, nil
}

// GetPayoutCurrencies 查询支持付款国家和币种列表
func (api *PayeeAPI) GetPayoutCurrencies(req *PayoutCurrenciesRequest) (*PayoutCurrenciesResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/remittance/payout_currencies")
	
	// 设置查询参数
	request.QueryArgs["payment_method"] = req.PaymentMethod
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get payout currencies failed: %w", err)
	}
	
	// 解析响应
	var currenciesResp PayoutCurrenciesResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &currenciesResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if currenciesResp.Result.Result != "S" {
		return &currenciesResp, fmt.Errorf("get payout currencies business error: [%s] %s",
			currenciesResp.Result.Code, currenciesResp.Result.Message)
	}
	
	return &currenciesResp, nil
}
