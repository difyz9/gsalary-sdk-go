package api

// WalletBalanceRequest 查询钱包余额请求
type WalletBalanceRequest struct {
	Currency string `json:"currency"` // 币种，请参考ISO-4217币种编码清单
}

// WalletBalanceData 钱包余额数据
type WalletBalanceData struct {
	Currency                 string  `json:"currency"`                    // 币种
	Amount                   float64 `json:"amount"`                      // 总金额
	ShareCardAccountBalance  float64 `json:"share_card_account_balance"`  // 共享卡账户余额
	Available                float64 `json:"available"`                   // 可用余额
	AccountType              string  `json:"account_type"`                // 账户类型 BALANCE
	QueryTime                string  `json:"query_time"`                  // 查询时间
}

// WalletBalanceResponse 查询钱包余额响应
type WalletBalanceResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data WalletBalanceData `json:"data"` // 钱包余额数据
}
