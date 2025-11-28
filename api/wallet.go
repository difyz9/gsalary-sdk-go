package api

import (
	"encoding/json"
	"fmt"
	
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// WalletAPI 钱包API接口
type WalletAPI struct {
	client *gsalary.GSalaryClient
}

// NewWalletAPI 创建钱包API实例
func NewWalletAPI(client *gsalary.GSalaryClient) *WalletAPI {
	return &WalletAPI{client: client}
}

// GetBalance 查询钱包余额
func (api *WalletAPI) GetBalance(req *WalletBalanceRequest) (*WalletBalanceResponse, error) {
	// 创建GET请求
	request := gsalary.NewRequest("GET", "/v1/wallets/balance")
	
	// 设置查询参数
	if req.Currency != "" {
		request.QueryArgs["currency"] = req.Currency
	}
	
	// 发送请求
	resp, err := api.client.Request(request)
	if err != nil {
		return nil, fmt.Errorf("get wallet balance failed: %w", err)
	}
	
	// 解析响应
	var balanceResp WalletBalanceResponse
	respBytes, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("marshal response failed: %w", err)
	}
	
	if err := json.Unmarshal(respBytes, &balanceResp); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}
	
	// 检查业务结果
	if balanceResp.Result.Result != "S" {
		return &balanceResp, fmt.Errorf("get wallet balance business error: [%s] %s",
			balanceResp.Result.Code, balanceResp.Result.Message)
	}
	
	return &balanceResp, nil
}
