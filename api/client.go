package api

import (
	gsalary "github.com/difyz9/gsalary-sdk-go"
)

// Client 统一的API客户端
type Client struct {
	client     *gsalary.GSalaryClient
	Card       *CardAPI
	CardHolder *CardHolderAPI
	Wallet     *WalletAPI
	Exchange   *ExchangeAPI
	Payment    *PaymentAPI
	Payee      *PayeeAPI       // 收款人管理
	Payer      *PayerAPI       // 付款人管理
	Remittance *RemittanceAPI  // 对外付款订单
}

// NewClient 创建新的API客户端
func NewClient(config *gsalary.GSalaryConfig) *Client {
	gsalaryClient := gsalary.NewClient(config)
	
	return &Client{
		client:     gsalaryClient,
		Card:       NewCardAPI(gsalaryClient),
		CardHolder: NewCardHolderAPI(gsalaryClient),
		Wallet:     NewWalletAPI(gsalaryClient),
		Exchange:   NewExchangeAPI(gsalaryClient),
		Payment:    NewPaymentAPI(gsalaryClient),
		Payee:      NewPayeeAPI(gsalaryClient),
		Payer:      NewPayerAPI(gsalaryClient),
		Remittance: NewRemittanceAPI(gsalaryClient),
	}
}

// GetRawClient 获取底层的GSalary客户端
func (c *Client) GetRawClient() *gsalary.GSalaryClient {
	return c.client
}
