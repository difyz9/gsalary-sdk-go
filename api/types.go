package api

// CardApplyRequest 申请新卡片请求
type CardApplyRequest struct {
	RequestID            string  `json:"request_id"`              // 唯一请求ID，最大50字符
	ProductCode          string  `json:"product_code"`            // 卡产品代码
	Currency             string  `json:"currency"`                // 卡币种，ISO-4217货币代码
	CardHolderID         string  `json:"card_holder_id"`          // 持卡人ID
	LimitPerDay          float64 `json:"limit_per_day,omitempty"` // 每日交易限额
	LimitPerMonth        float64 `json:"limit_per_month,omitempty"` // 每月交易限额
	LimitPerTransaction  float64 `json:"limit_per_transaction,omitempty"` // 单笔交易限额
	InitBalance          float64 `json:"init_balance"`            // 初始余额
}

// CardApplyResponse 申请卡片响应
type CardApplyResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果状态，S表示成功
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		RequestID string `json:"request_id"` // 请求ID
		Status    string `json:"status"`     // 状态，如PENDING
	} `json:"data"`
}

// Response 通用响应结构
type Response struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data interface{} `json:"data"`
}

// CardAvailableQuotasRequest 查询卡可用余额请求
type CardAvailableQuotasRequest struct {
	Currency           string `json:"currency"`             // 卡币种，参考ISO-4217币种清单(支持的币种：USD)
	AccountingCardType string `json:"accounting_card_type"` // 卡账务类型，不填默认认为SHARE。Enum: "SHARE" "RECHARGE"
}

// CardAvailableQuotasResponse 查询卡可用余额响应
type CardAvailableQuotasResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data struct {
		Currency           string  `json:"currency"`             // 币种
		AccountingCardType string  `json:"accounting_card_type"` // 卡账务类型
		AvailableQuota     float64 `json:"available_quota"`      // 可用余额
	} `json:"data"`
}

// CardProductsRequest 查询可用卡产品列表请求
type CardProductsRequest struct {
	CardType   string `json:"card_type"`   // 卡类型 Enum: "PHYSICAL" "VIRTUAL"
	BrandCode  string `json:"brand_code"`  // 卡品牌 Enum: "VISA" "MASTER"
	Currency   string `json:"currency"`    // 卡币种，参考ISO-4217币种清单
}

// CardProduct 卡产品信息
type CardProduct struct {
	ProductCode string `json:"product_code"` // 产品代码
	ProductName string `json:"product_name"` // 产品名称
	CardType    string `json:"card_type"`    // 卡类型
	BrandCode   string `json:"brand_code"`   // 卡品牌
	Currency    string `json:"currency"`     // 币种
	Description string `json:"description"`  // 产品描述
}

// CardProductsResponse 查询可用卡产品列表响应
type CardProductsResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data struct {
		Products []CardProduct `json:"products"` // 产品列表
	} `json:"data"`
}

// CardApplyResultResponse 查询开卡结果响应
type CardApplyResultResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data map[string]interface{} `json:"data"` // 开卡结果数据
}

// CardListRequest 查询卡列表请求
type CardListRequest struct {
	Page         int    `json:"page"`          // 页码，从1开始，默认1
	Limit        int    `json:"limit"`         // 每页记录数，默认20
	ProductCode  string `json:"product_code"`  // 卡产品编码
	BrandCode    string `json:"brand_code"`    // 卡品牌 Enum: "VISA" "MASTER"
	CardHolderID string `json:"card_holder_id"` // 持卡人ID
	CreateStart  string `json:"create_start"`  // 查询创建卡起始时间（包），使用ISO-8601时间格式
	CreateEnd    string `json:"create_end"`    // 查询创建卡截止时间（不含），使用ISO-8601时间格式
	Status       string `json:"status"`        // 卡状态 Enum: "PENDING" "INACTIVE" "ACTIVE" "FREEZING" "FROZEN" "UNFREEZING" "EXPIRED" "CANCELLING" "CANCELLED"
}

// Card 卡片信息
type Card struct {
	CardID       string                 `json:"card_id"`       // 卡片ID
	ProductCode  string                 `json:"product_code"`  // 产品代码
	BrandCode    string                 `json:"brand_code"`    // 卡品牌
	CardHolderID string                 `json:"card_holder_id"` // 持卡人ID
	Status       string                 `json:"status"`        // 卡状态
	CreatedAt    string                 `json:"created_at"`    // 创建时间
	UpdatedAt    string                 `json:"updated_at"`    // 更新时间
	Extra        map[string]interface{} `json:"extra,omitempty"` // 额外信息
}

// CardListResponse 查询卡列表响应
type CardListResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data struct {
		Cards      []Card `json:"cards"`       // 卡片列表
		Page       int    `json:"page"`        // 当前页码
		Limit      int    `json:"limit"`       // 每页数量
		TotalCount int    `json:"total_count"` // 总记录数
		TotalPage  int    `json:"total_page"`  // 总页数
	} `json:"data"`
}

// CardInfo 卡片详细信息
type CardInfo struct {
	CardID              string                 `json:"card_id"`               // 卡片ID
	CardName            string                 `json:"card_name"`             // 卡片名称
	MaskCardNumber      string                 `json:"mask_card_number"`      // 掩码卡号，如 41******1111
	CardCurrency        string                 `json:"card_currency"`         // 卡币种
	AvailableBalance    float64                `json:"available_balance"`     // 可用余额
	BrandCode           string                 `json:"brand_code"`            // 卡品牌 VISA/MASTER
	Status              string                 `json:"status"`                // 卡状态
	CardType            string                 `json:"card_type"`             // 卡类型 PHYSICAL/VIRTUAL
	AccountingType      string                 `json:"accounting_type"`       // 账务类型 SHARE/RECHARGE
	CardRegion          string                 `json:"card_region"`           // 卡地区，如 US
	CardHolderID        string                 `json:"card_holder_id"`        // 持卡人ID
	FirstName           string                 `json:"first_name"`            // 名
	LastName            string                 `json:"last_name"`             // 姓
	Mobile              map[string]interface{} `json:"mobile"`                // 手机号
	Email               string                 `json:"email"`                 // 邮箱
	LimitPerDay         float64                `json:"limit_per_day"`         // 每日限额
	LimitPerMonth       float64                `json:"limit_per_month"`       // 每月限额
	LimitPerTransaction float64                `json:"limit_per_transaction"` // 单笔限额
	BillAddress         map[string]interface{} `json:"bill_address"`          // 账单地址
	SupportTdsTrans     bool                   `json:"support_tds_trans"`     // 是否支持3DS交易
	CreateTime          string                 `json:"create_time"`           // 创建时间
}

// CardInfoResponse 查看卡信息响应
type CardInfoResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data CardInfo `json:"data"` // 卡片详细信息
}

// UpdateCardRequest 修改卡信息请求
type UpdateCardRequest struct {
	CardID               string  `json:"card_id"`                          // 卡ID（path参数）
	CardName             string  `json:"card_name,omitempty"`              // 卡昵称
	LimitPerDay          float64 `json:"limit_per_day,omitempty"`          // 每日限额
	LimitPerMonth        float64 `json:"limit_per_month,omitempty"`        // 每月限额
	LimitPerTransaction  float64 `json:"limit_per_transaction,omitempty"`  // 单笔交易限额
}

// UpdateCardResponse 修改卡信息响应
type UpdateCardResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data CardInfo `json:"data"` // 更新后的卡片信息
}

// DeleteCardResponse 销卡响应
type DeleteCardResponse struct {
	Result struct {
		Result  string `json:"result"`  // S-成功，F-失败，U-未知
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
}

// CardSecureInfo 卡机密信息
type CardSecureInfo struct {
	PAN         string `json:"pan"`          // 卡号（明文）
	ExpireYear  string `json:"expire_year"`  // 到期年份
	ExpireMonth string `json:"expire_month"` // 到期月份
	CVV         string `json:"cvv"`          // CVV安全码
}

// CardSecureInfoResponse 获取卡机密信息响应
type CardSecureInfoResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data CardSecureInfo `json:"data"`
}

// AdjustCardBalanceRequest 卡片调额请求
type AdjustCardBalanceRequest struct {
	CardID    string  `json:"card_id"`    // 卡ID
	Amount    float64 `json:"amount"`     // 修改金额（必须>=0）
	Type      string  `json:"type"`       // 修改类型：INCREASE（增加）、DECREASE（减少）
	RequestID string  `json:"request_id"` // 唯一请求ID
}

// BalanceModifyResult 调额结果
type BalanceModifyResult struct {
	GSalaryRequestID string  `json:"gsalary_request_id"` // GSalary请求ID
	RequestID        string  `json:"request_id"`         // 商户请求ID
	CardID           string  `json:"card_id"`            // 卡ID
	Status           string  `json:"status"`             // 状态：PENDING/SUCCESS/FAIL
	CreateTime       string  `json:"create_time"`        // 创建时间
	FinishTime       string  `json:"finish_time"`        // 完成时间
	Amount           float64 `json:"amount"`             // 调额金额
	Type             string  `json:"type"`               // 调额类型
	PostBalance      float64 `json:"post_balance"`       // 调额后余额
}

// AdjustCardBalanceResponse 卡片调额响应
type AdjustCardBalanceResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data BalanceModifyResult `json:"data"`
}

// GetBalanceModifyResultResponse 查询调额结果响应
type GetBalanceModifyResultResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data BalanceModifyResult `json:"data"`
}

// SetCardFreezeStatusRequest 冻结/解冻卡请求
type SetCardFreezeStatusRequest struct {
	CardID string `json:"card_id"` // 卡ID（path参数）
	Freeze bool   `json:"freeze"`   // 冻结状态：true-冻结，false-解冻
}

// SetCardFreezeStatusResponse 冻结/解冻卡响应
type SetCardFreezeStatusResponse struct {
	Result struct {
		Result  string `json:"result"`  // S-成功，F-失败，U-未知
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
}

// CardTransactionsRequest 卡交易列表请求
type CardTransactionsRequest struct {
	Page           int    `json:"page"`                      // 页码，从1开始
	Limit          int    `json:"limit"`                     // 每页记录条数
	TransactionID  string `json:"transaction_id,omitempty"`  // 卡交易ID
	MchRequestID   string `json:"mch_request_id,omitempty"`  // 商户请求ID
	TimeStart      string `json:"time_start,omitempty"`      // 查询起始时间（含）
	TimeEnd        string `json:"time_end,omitempty"`        // 查询截止时间（不含）
	CardID         string `json:"card_id,omitempty"`         // 卡ID
}

// CardTransaction 卡交易信息
type CardTransaction struct {
	TransactionID         string  `json:"transaction_id"`          // 交易ID
	CardID                string  `json:"card_id"`                 // 卡ID
	TransactionType       string  `json:"transaction_type"`        // 交易类型
	Amount                float64 `json:"amount"`                  // 交易金额
	Currency              string  `json:"currency"`                // 币种
	Status                string  `json:"status"`                  // 交易状态
	StatusDescription     string  `json:"status_description"`      // 状态描述
	TransactionTime       string  `json:"transaction_time"`        // 交易时间
	MerchantName          string  `json:"merchant_name"`           // 商户名称
	MerchantCountry       string  `json:"merchant_country"`        // 商户国家
	MerchantCategoryCode  string  `json:"merchant_category_code"`  // 商户类别代码
}

// CardTransactionsResponse 卡交易列表响应
type CardTransactionsResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data struct {
		Query        CardTransactionsRequest `json:"query"`        // 查询条件
		Transactions []CardTransaction       `json:"transactions"` // 交易列表
		Page         int                     `json:"page"`         // 当前页码
		Limit        int                     `json:"limit"`        // 每页数量
		TotalCount   int                     `json:"total_count"`  // 总记录数
		TotalPage    int                     `json:"total_page"`   // 总页数
	} `json:"data"`
}

// BalanceHistoryRequest 卡余额变更记录请求
type BalanceHistoryRequest struct {
	Page          int    `json:"page"`                     // 页码，从1开始
	Limit         int    `json:"limit"`                    // 每页记录条数
	TransactionID string `json:"transaction_id,omitempty"` // 卡交易ID
	LogID         string `json:"log_id,omitempty"`         // 入账ID
	TimeStart     string `json:"time_start,omitempty"`     // 查询起始时间（含）
	TimeEnd       string `json:"time_end,omitempty"`       // 查询截止时间（不含）
	CardID        string `json:"card_id,omitempty"`        // 卡ID
}

// BalanceHistoryRecord 余额变更记录
type BalanceHistoryRecord struct {
	LogID           string  `json:"log_id"`           // 入账ID
	CardID          string  `json:"card_id"`          // 卡ID
	TransactionID   string  `json:"transaction_id"`   // 关联交易ID
	Amount          float64 `json:"amount"`           // 变更金额
	PostBalance     float64 `json:"post_balance"`     // 变更后余额
	TransactionType string  `json:"transaction_type"` // 交易类型
	CreateTime      string  `json:"create_time"`      // 创建时间
	Description     string  `json:"description"`      // 描述
}

// BalanceHistoryResponse 卡余额变更记录响应
type BalanceHistoryResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data struct {
		Query      BalanceHistoryRequest  `json:"query"`   // 查询条件
		History    []BalanceHistoryRecord `json:"history"` // 余额变更记录
		Page       int                    `json:"page"`    // 当前页码
		Limit      int                    `json:"limit"`   // 每页数量
		TotalCount int                    `json:"total_count"` // 总记录数
		TotalPage  int                    `json:"total_page"`  // 总页数
	} `json:"data"`
}

// UpdateCardContactRequest 修改卡联系信息请求
type UpdateCardContactRequest struct {
	CardID string `json:"card_id"`        // 卡ID（path参数）
	Email  string `json:"email,omitempty"` // 持卡人email
}

// UpdateCardContactResponse 修改卡联系信息响应
type UpdateCardContactResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data CardInfo `json:"data"` // 更新后的卡片信息
}

// MobileNumber 手机号信息
type MobileNumber struct {
	CountryCode string `json:"country_code"` // 国家代码
	Number      string `json:"number"`       // 手机号码
}

// Address 地址信息
type Address struct {
	Country    string `json:"country"`     // 国家
	State      string `json:"state"`       // 州/省
	City       string `json:"city"`        // 城市
	PostalCode string `json:"postal_code"` // 邮编
	Line1      string `json:"line1"`       // 地址行1
	Line2      string `json:"line2"`       // 地址行2
}

// CardHolderRequest 添加持卡人请求
type CardHolderRequest struct {
	FirstName   string       `json:"first_name"`   // 持卡人名字，长度需在1至40个字符之间，仅可包含英文字母和空格
	LastName    string       `json:"last_name"`    // 持卡人姓氏，长度需在1至40个字符之间，仅可包含英文字母和空格
	Birth       string       `json:"birth"`        // 持卡人生日，使用ISO-8601日期格式
	Email       string       `json:"email"`        // 请确保持卡人邮箱唯一性
	Mobile      MobileNumber `json:"mobile"`       // 持卡人手机号
	Region      string       `json:"region"`       // 2字符国家码，参考ISO-3166标准国家码清单
	BillAddress Address      `json:"bill_address"` // 账单地址
}

// CardHolderInfo 持卡人信息
type CardHolderInfo struct {
	CardHolderID string       `json:"card_holder_id"` // 持卡人ID
	FirstName    string       `json:"first_name"`     // 名字
	LastName     string       `json:"last_name"`      // 姓氏
	Birth        string       `json:"birth"`          // 生日
	Email        string       `json:"email"`          // 邮箱
	Mobile       MobileNumber `json:"mobile"`         // 手机号
	Region       string       `json:"region"`         // 国家码
	BillAddress  Address      `json:"bill_address"`   // 账单地址
	CreatedAt    string       `json:"created_at"`     // 创建时间
	UpdatedAt    string       `json:"updated_at"`     // 更新时间
}

// CardHolderResponse 添加持卡人响应
type CardHolderResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data CardHolderInfo `json:"data"` // 持卡人信息
}

// UpdateCardHolderRequest 修改持卡人信息请求
type UpdateCardHolderRequest struct {
	FirstName   string       `json:"first_name"`   // 持卡人名字
	LastName    string       `json:"last_name"`    // 持卡人姓氏
	Birth       string       `json:"birth"`        // 持卡人生日，使用ISO-8601日期格式
	Email       string       `json:"email"`        // 邮箱
	Mobile      MobileNumber `json:"mobile"`       // 手机号
	Region      string       `json:"region"`       // 2字符国家码
	BillAddress Address      `json:"bill_address"` // 账单地址
}

// UpdateCardHolderResponse 修改持卡人信息响应
type UpdateCardHolderResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data CardHolderInfo `json:"data"` // 更新后的持卡人信息
}

// CardHolderListRequest 查询持卡人列表请求
type CardHolderListRequest struct {
	Page      int    `json:"page"`       // 页码，从1开始，默认1
	Limit     int    `json:"limit"`      // 每页记录数，默认20
	TimeStart string `json:"time_start"` // 查询创建卡起始时间（包），使用ISO-8601时间格式
	TimeEnd   string `json:"time_end"`   // 查询创建卡截止时间（不含），使用ISO-8601时间格式
}

// CardHolderListResponse 查询持卡人列表响应
type CardHolderListResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data struct {
		Query        map[string]interface{} `json:"query"`         // 查询条件
		CardHolders  []CardHolderInfo       `json:"card_holders"`  // 持卡人列表
		Page         int                    `json:"page"`          // 当前页码
		Limit        int                    `json:"limit"`         // 每页数量
		TotalCount   int                    `json:"total_count"`   // 总记录数
		TotalPage    int                    `json:"total_page"`    // 总页数
	} `json:"data"`
}

// CardHolderDetailInfo 持卡人详细信息
type CardHolderDetailInfo struct {
	CardHolderID string      `json:"card_holder_id"` // 持卡人ID
	FirstName    string      `json:"first_name"`     // 名字
	LastName     string      `json:"last_name"`      // 姓氏
	Birth        string      `json:"birth"`          // 生日
	Email        string      `json:"email"`          // 邮箱
	Region       string      `json:"region"`         // 国家码
	CreateTime   string      `json:"create_time"`    // 创建时间
	BillAddress  interface{} `json:"bill_address"`   // 账单地址
	Status       string      `json:"status"`         // 状态
}

// CardHolderDetailResponse 查看持卡人信息响应
type CardHolderDetailResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data CardHolderDetailInfo `json:"data"` // 持卡人详细信息
}

// ExchangeRateRequest 查询汇率请求
type ExchangeRateRequest struct {
	BuyCurrency  string `json:"buy_currency"`  // 购入币种，参考ISO-4217币种清单
	SellCurrency string `json:"sell_currency"` // 卖出币种，参考ISO-4217币种清单
}

// ExchangeRateResponse 查询汇率响应
type ExchangeRateResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data struct {
		BuyCurrency  string  `json:"buy_currency"`  // 购入币种
		SellCurrency string  `json:"sell_currency"` // 卖出币种
		Rate         float64 `json:"rate"`          // 汇率
		UpdateTime   string  `json:"update_time"`   // 更新时间
	} `json:"data"`
}

// ExchangeQuoteRequest 请求锁汇报价请求
type ExchangeQuoteRequest struct {
	BuyCurrency  string  `json:"buy_currency"`  // 购入币种，参考ISO-4217币种清单
	SellCurrency string  `json:"sell_currency"` // 卖出币种，参考ISO-4217币种清单
	BuyAmount    float64 `json:"buy_amount"`    // 购入金额。购入金额和卖出金额不可同时为空，如果同时提供将忽略购入金额
	SellAmount   float64 `json:"sell_amount"`   // 卖出金额。购入金额和卖出金额不可同时为空，如果同时提供将忽略购入金额
}

// CurrencyAmount 货币金额
type CurrencyAmount struct {
	Currency string  `json:"currency"` // 币种
	Amount   float64 `json:"amount"`   // 金额
}

// ExchangeQuoteResponse 请求锁汇报价响应
type ExchangeQuoteResponse struct {
	Result struct {
		Result  string `json:"result"`
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"result"`
	Data struct {
		QuoteID    string         `json:"quote_id"`    // 报价ID
		Buy        CurrencyAmount `json:"buy"`         // 购入金额信息
		Sell       CurrencyAmount `json:"sell"`        // 卖出金额信息
		Surcharge  CurrencyAmount `json:"surcharge"`   // 附加费
		TotalCost  CurrencyAmount `json:"total_cost"`  // 总成本
		UpdateTime string         `json:"update_time"` // 更新时间
		ExpireTime string         `json:"expire_time"` // 过期时间
	} `json:"data"`
}

// ExchangeSubmitRequest 提交换汇订单请求
type ExchangeSubmitRequest struct {
	RequestID string `json:"request_id"` // 唯一请求ID
	QuoteID   string `json:"quote_id"`   // 锁汇返回的锁汇ID
}

// ExchangeOrder 换汇订单详情
type ExchangeOrder struct {
	OrderID      string         `json:"order_id"`      // 订单ID
	RequestID    string         `json:"request_id"`    // 请求ID
	CreateTime   string         `json:"create_time"`   // 创建时间
	Status       string         `json:"status"`        // 订单状态：PENDING/SUCCESS/FAILED
	Source       string         `json:"source"`        // 来源：PORTAL/API
	Sell         CurrencyAmount `json:"sell"`          // 卖出币种和金额
	Buy          CurrencyAmount `json:"buy"`           // 买入币种和金额
	Surcharge    CurrencyAmount `json:"surcharge"`     // 手续费
	ExchangeRate float64        `json:"exchange_rate"` // 汇率
}

// ExchangeSubmitResponse 提交换汇订单响应
type ExchangeSubmitResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data ExchangeOrder `json:"data"` // 订单详情
}

// ExchangeOrdersRequest 查询换汇订单列表请求
type ExchangeOrdersRequest struct {
	Page         int    `json:"page"`          // 页码，从1开始计算
	Limit        int    `json:"limit"`         // 每页记录条数
	TimeStart    string `json:"time_start"`    // 查询起始时间（不含），ISO-8601时间格式
	TimeEnd      string `json:"time_end"`      // 查询截止时间（不含），ISO-8601时间格式
	Status       string `json:"status"`        // 换汇订单状态：PENDING/SUCCESS/FAIL
	BuyCurrency  string `json:"buy_currency"`  // 购入币种
	SellCurrency string `json:"sell_currency"` // 卖出币种
}

// ExchangeOrdersResponse 查询换汇订单列表响应
type ExchangeOrdersResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		Query struct {
			TimeStart    string `json:"time_start"`    // 查询起始时间
			TimeEnd      string `json:"time_end"`      // 查询截止时间
			Status       string `json:"status"`        // 订单状态
			BuyCurrency  string `json:"buy_currency"`  // 购入币种
			SellCurrency string `json:"sell_currency"` // 卖出币种
		} `json:"query"` // 查询条件
		Orders     []ExchangeOrder `json:"orders"`      // 订单列表
		Page       int             `json:"page"`        // 当前页码
		Limit      int             `json:"limit"`       // 每页数量
		TotalCount int             `json:"total_count"` // 总记录数
		TotalPage  int             `json:"total_page"`  // 总页数
	} `json:"data"`
}

// ==================== Payment Types ====================

// PaymentConsultRequest 支付咨询请求
type PaymentConsultRequest struct {
	MchAppID                      string   `json:"mch_app_id"`                         // 商户应用ID
	PaymentCurrency               string   `json:"payment_currency"`                   // 支付币种（ISO 4217三位代码）
	PaymentAmount                 float64  `json:"payment_amount"`                     // 支付金额（单位：元）
	SettlementCurrency            string   `json:"settlement_currency"`                // 结算币种（ISO 4217三位代码）
	AllowedPaymentMethodRegions   []string `json:"allowed_payment_method_regions"`     // 允许的支付方式所属国家/地区
	AllowedPaymentMethods         []string `json:"allowed_payment_methods"`            // 允许的支付方式列表
	UserRegion                    string   `json:"user_region,omitempty"`              // 用户所在国家/地区（ISO两位代码）
	EnvTerminalType               string   `json:"env_terminal_type"`                  // 终端类型：WEB/WAP/APP/MINI_APP
	EnvOsType                     string   `json:"env_os_type,omitempty"`              // 操作系统类型：IOS/ANDROID
	EnvClientIP                   string   `json:"env_client_ip,omitempty"`            // 客户端设备IP地址
}

// PaymentOption 支付方式选项
type PaymentOption struct {
	PaymentMethodType string `json:"payment_method_type"` // 支付方式类型
	Currency          string `json:"currency"`            // 币种
	Limit             struct {
		Min float64 `json:"min"` // 最小金额
		Max float64 `json:"max"` // 最大金额
	} `json:"limit"` // 限额
	Country string `json:"country"` // 国家/地区代码
}

// PaymentConsultResponse 支付咨询响应
type PaymentConsultResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		PaymentOptions []PaymentOption `json:"payment_options"` // 支付方式选项列表
	} `json:"data"`
}

// OrderInfo 订单信息
type OrderInfo struct {
	ReferenceOrderID   string  `json:"reference_order_id"`             // 商户订单ID
	OrderDescription   string  `json:"order_description,omitempty"`    // 订单描述
	OrderCurrency      string  `json:"order_currency,omitempty"`       // 订单币种
	OrderAmount        float64 `json:"order_amount,omitempty"`         // 订单金额
	OrderBuyerID       string  `json:"order_buyer_id,omitempty"`       // 买家ID
	OrderBuyerEmail    string  `json:"order_buyer_email,omitempty"`    // 买家邮箱
}

// PaymentSessionRequest 支付会话创建请求（收银台）
type PaymentSessionRequest struct {
	MchAppID                  string    `json:"mch_app_id"`                     // 商户应用ID
	PaymentRequestID          string    `json:"payment_request_id"`             // 商户自定义支付请求ID（唯一）
	PaymentCurrency           string    `json:"payment_currency"`               // 支付币种（ISO 4217三位代码）
	PaymentAmount             float64   `json:"payment_amount"`                 // 支付金额（单位：元）
	PaymentMethodType         string    `json:"payment_method_type"`            // 支付方式类型
	PaymentSessionExpiryTime  string    `json:"payment_session_expiry_time,omitempty"` // 会话过期时间（ISO 8601格式）
	PaymentRedirectURL        string    `json:"payment_redirect_url"`           // 支付完成后用户重定向地址
	Order                     OrderInfo `json:"order"`                          // 订单信息
	SettlementCurrency        string    `json:"settlement_currency"`            // 结算币种（ISO 4217三位代码）
	EnvClientIP               string    `json:"env_client_ip,omitempty"`        // 客户端设备IP地址
	ProductScene              string    `json:"product_scene"`                  // 产品场景：CHECKOUT_PAYMENT/ELEMENT_PAYMENT
}

// PaymentSessionResponse 支付会话创建响应
type PaymentSessionResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		PaymentSessionData       string `json:"payment_session_data"`        // 加密的会话数据
		PaymentSessionExpiryTime string `json:"payment_session_expiry_time"` // 会话过期时间
		PaymentSessionID         string `json:"payment_session_id"`          // 会话ID
		NormalURL                string `json:"normal_url"`                  // 支付URL
	} `json:"data"`
}

// EasySafePaySessionRequest 钱包授权支付会话创建请求（第一次支付）
type EasySafePaySessionRequest struct {
	PaymentSessionRequest                    // 继承支付会话请求
	AuthState             string `json:"auth_state"`           // 商户生成的快捷授权专属ID
	UserLoginID           string `json:"user_login_id,omitempty"` // 用户支付方式登录ID（邮箱/手机号）
}

// EasySafePayRequest 钱包授权支付请求（第二次支付）
type EasySafePayRequest struct {
	MchAppID            string    `json:"mch_app_id"`             // 商户应用ID
	PaymentRequestID    string    `json:"payment_request_id"`     // 商户自定义支付请求ID（唯一）
	PaymentCurrency     string    `json:"payment_currency"`       // 支付币种（ISO 4217三位代码）
	PaymentAmount       float64   `json:"payment_amount"`         // 支付金额（单位：元）
	PaymentMethodID     string    `json:"payment_method_id"`      // 支付方式ID（access_token）
	PaymentMethodType   string    `json:"payment_method_type"`    // 支付方式类型
	PaymentRedirectURL  string    `json:"payment_redirect_url"`   // 支付完成后重定向地址
	Order               OrderInfo `json:"order"`                  // 订单信息
	SettlementCurrency  string    `json:"settlement_currency"`    // 结算币种（ISO 4217三位代码）
	EnvClientIP         string    `json:"env_client_ip,omitempty"` // 客户端设备IP地址
	PaymentExpiryTime   string    `json:"payment_expiry_time,omitempty"` // 支付有效期（ISO 8601格式）
	EnvTerminalType     string    `json:"env_terminal_type"`      // 终端类型：WEB/WAP/APP/MINI_APP
	EnvOsType           string    `json:"env_os_type,omitempty"`  // 操作系统类型：IOS/ANDROID
}

// PaymentResponse 支付响应
type PaymentResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		PaymentRequestID string  `json:"payment_request_id"` // 支付请求ID
		PaymentID        string  `json:"payment_id"`         // 支付ID
		PaymentCurrency  string  `json:"payment_currency"`   // 支付币种
		PaymentAmount    float64 `json:"payment_amount"`     // 支付金额
		NormalURL        string  `json:"normal_url"`         // 支付URL
		SchemeURL        string  `json:"scheme_url"`         // Scheme URL
		ApplinkURL       string  `json:"applink_url"`        // Applink URL
		AppIdentifier    string  `json:"app_identifier"`     // App标识符
	} `json:"data"`
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	MchAppID       string `json:"mch_app_id"`               // 商户应用ID
	RefreshToken   string `json:"refresh_token"`            // 刷新令牌
	MerchantRegion string `json:"merchant_region,omitempty"` // 商户业务国家/地区（ISO两位代码）
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		AccessToken           string `json:"access_token"`             // 新的访问令牌
		AccessTokenExpiryTime string `json:"access_token_expiry_time"` // 访问令牌过期时间
		RefreshToken          string `json:"refresh_token"`            // 新的刷新令牌
		RefreshTokenExpiryTime string `json:"refresh_token_expiry_time"` // 刷新令牌过期时间
		UserLoginID           string `json:"user_login_id"`            // 用户登录ID
	} `json:"data"`
}

// RevokeTokenRequest 取消授权请求
type RevokeTokenRequest struct {
	MchAppID    string `json:"mch_app_id"`    // 商户应用ID
	AccessToken string `json:"access_token"`  // 访问令牌
}

// RevokeTokenResponse 取消授权响应
type RevokeTokenResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败，U-未知
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
}

// CancelPaymentRequest 取消支付请求
type CancelPaymentRequest struct {
	MchAppID         string `json:"mch_app_id"`          // 商户应用ID
	PaymentRequestID string `json:"payment_request_id"`  // 支付请求ID
	PaymentID        string `json:"payment_id"`          // 支付ID
}

// CancelPaymentResponse 取消支付响应
type CancelPaymentResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败，U-未知
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		PaymentRequestID string `json:"payment_request_id"` // 支付请求ID
		PaymentID        string `json:"payment_id"`         // 支付ID
		CancelTime       string `json:"cancel_time"`        // 取消时间
	} `json:"data"`
}

// QueryPaymentRequest 查询支付请求
type QueryPaymentRequest struct {
	MchAppID         string `json:"mch_app_id"`          // 商户应用ID
	PaymentRequestID string `json:"payment_request_id"`  // 支付请求ID（与PaymentID二选一）
	PaymentID        string `json:"payment_id"`          // 支付ID（与PaymentRequestID二选一）
}

// QueryPaymentResponse 查询支付响应
type QueryPaymentResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		PaymentRequestID  string                 `json:"payment_request_id"`  // 支付请求ID
		PaymentID         string                 `json:"payment_id"`          // 支付ID
		PaymentAmount     float64                `json:"payment_amount"`      // 支付金额
		PaymentCurrency   string                 `json:"payment_currency"`    // 支付币种
		PaymentStatus     string                 `json:"payment_status"`      // 支付状态：SUCCESS/FAIL/PROCESSING/CANCELLED
		PaymentMethodType string                 `json:"payment_method_type"` // 支付方式类型
		PaymentCreateTime string                 `json:"payment_create_time"` // 支付创建时间
		PaymentUpdateTime string                 `json:"payment_update_time"` // 支付更新时间
		PaymentResultCode string                 `json:"payment_result_code"` // 支付结果代码
		PaymentResultInfo map[string]interface{} `json:"payment_result_info"` // 支付结果信息
	} `json:"data"`
}

// ============ 对外付款相关类型 ============

// PayeeRequest 新增收款人请求
type PayeeRequest struct {
	SubjectType   string `json:"subject_type"`            // 主体类型：INDIVIDUAL/ENTERPRISE
	AccountType   string `json:"account_type"`            // 账户类型：E_WALLET/BANK_ACCOUNT
	Country       string `json:"country"`                 // 国家/地区代码
	FirstName     string `json:"first_name,omitempty"`    // 收款人名（个人类型必填）
	LastName      string `json:"last_name,omitempty"`     // 收款人姓（个人类型必填）
	AccountHolder string `json:"account_holder,omitempty"`// 账户持有人（企业类型必填）
	Currency      string `json:"currency"`                // 收款币种
}

// Payee 收款人信息
type Payee struct {
	PayeeID       string        `json:"payee_id"`                 // 收款人ID
	SubjectType   string        `json:"subject_type"`             // 主体类型
	AccountType   string        `json:"account_type"`             // 账户类型
	FirstName     string        `json:"first_name,omitempty"`     // 收款人名
	LastName      string        `json:"last_name,omitempty"`      // 收款人姓
	Country       string        `json:"country"`                  // 国家/地区代码
	Currencies    []string      `json:"currencies"`               // 支持的币种列表
	AccountHolder string        `json:"account_holder,omitempty"` // 账户持有人
	Mobile        *MobileNumber `json:"mobile,omitempty"`         // 手机号码
	Address       string        `json:"address,omitempty"`        // 地址
}

// PayeeResponse 收款人响应
type PayeeResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data Payee `json:"data"`
}

// PayeeListRequest 查询收款人列表请求
type PayeeListRequest struct {
	Page     int    `json:"page"`               // 页码，从1开始
	Limit    int    `json:"limit"`              // 每页记录条数
	PayeeID  string `json:"payee_id,omitempty"` // 收款人ID
	Name     string `json:"name,omitempty"`     // 姓名或公司名称
	Country  string `json:"country,omitempty"`  // 国家/地区代码
	Currency string `json:"currency,omitempty"` // 币种代码
	Mobile   string `json:"mobile,omitempty"`   // 联系电话（不带国家区号）
}

// PayeeListResponse 查询收款人列表响应
type PayeeListResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		Query      map[string]interface{} `json:"query"`       // 查询条件
		Payees     []Payee                `json:"payees"`      // 收款人列表
		Page       int                    `json:"page"`        // 当前页码
		Limit      int                    `json:"limit"`       // 每页数量
		TotalCount int                    `json:"total_count"` // 总记录数
		TotalPage  int                    `json:"total_page"`  // 总页数
	} `json:"data"`
}

// PayeeAccount 收款账户信息
type PayeeAccount struct {
	PaymentMethod string   `json:"payment_method"`       // 付款方式
	Currencies    []string `json:"currencies"`           // 支持的币种列表
	AccountID     string   `json:"account_id"`           // 账户ID
	AccountNo     string   `json:"account_no,omitempty"` // 账号
	Status        string   `json:"status"`               // 状态：PENDING/ACTIVE/INACTIVE
	RequireClearingNetwork bool `json:"require_clearing_network,omitempty"` // 是否需要清算网络
	FormFields    []FormField `json:"form_fields,omitempty"` // 表单字段（银行账户）
}

// FormField 表单字段
type FormField struct {
	FieldName  string `json:"field_name"`  // 字段名
	FieldValue string `json:"field_value"` // 字段值
}

// PayeeAccountRequest 新增收款账户请求（电子钱包）
type PayeeAccountRequest struct {
	PaymentMethod string `json:"payment_method"` // 付款方式（如ALIPAY）
	AccountNo     string `json:"account_no"`     // 账号ID
}

// PayeeAccountResponse 收款账户响应
type PayeeAccountResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data PayeeAccount `json:"data"`
}

// PayeeAccountsResponse 查询收款账户列表响应
type PayeeAccountsResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		Accounts []PayeeAccount `json:"accounts"` // 账户列表
	} `json:"data"`
}

// PayeeAccountFormRequest 获取收款账户表单请求（银行账户）
type PayeeAccountFormRequest struct {
	PaymentMethod string `json:"payment_method"`       // 付款方式：BANK_TRANSFER
	Currency      string `json:"currency,omitempty"`   // 收款账户币种
	Language      string `json:"language,omitempty"`   // 表单语言，默认：en
}

// PayeeAccountFormResponse 获取收款账户表单响应
type PayeeAccountFormResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	AccountType   string `json:"account_type"`   // 账户类型
	PaymentMethod string `json:"payment_method"` // 付款方式
	SubjectType   string `json:"subject_type"`   // 主体类型
	Currency      string `json:"currency"`       // 币种
	Country       string `json:"country"`        // 国家/地区
	Fields        []struct {
		FieldName   string `json:"field_name"`   // 字段名
		Required    bool   `json:"required"`     // 是否必填
		Description string `json:"description"`  // 字段描述
	} `json:"fields"` // 表单字段列表
}

// PayeeAccountBankRequest 新增/更新收款账户请求（银行账户）
type PayeeAccountBankRequest struct {
	PaymentMethod string      `json:"payment_method"`       // 付款方式：BANK_TRANSFER
	Currency      string      `json:"currency,omitempty"`   // 账户币种
	Fields        []FormField `json:"fields"`               // 表单字段集合
}

// PaymentMethodsResponse 查询可用付款方式响应
type PaymentMethodsResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		PaymentMethods []string `json:"payment_methods"` // 付款方式列表
	} `json:"data"`
}

// PayoutCurrenciesRequest 查询支持付款国家和币种请求
type PayoutCurrenciesRequest struct {
	PaymentMethod string `json:"payment_method"` // 付款方式
}

// PayoutCurrency 付款币种信息
type PayoutCurrency struct {
	Country    string   `json:"country"`    // 国家/地区代码
	Currencies []string `json:"currencies"` // 支持的币种列表
}

// PayoutCurrenciesResponse 查询支持付款国家和币种响应
type PayoutCurrenciesResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		PaymentMethod string           `json:"payment_method"` // 付款方式
		Countries     []PayoutCurrency `json:"countries"`      // 国家币种列表
	} `json:"data"`
}

// UploadAttachmentRequest 上传附件请求
type UploadAttachmentRequest struct {
	Type     string `json:"type"`     // 附件类型：CERT_FILE
	Filename string `json:"filename"` // 文件名
	Base64   string `json:"base64"`   // Base64编码的文件内容
}

// UploadAttachmentResponse 上传附件响应
type UploadAttachmentResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	FileID string `json:"file_id"` // 文件ID
}

// PayerRequest 新增付款人请求
type PayerRequest struct {
	SubjectType    string   `json:"subject_type"`              // 主体类型：INDIVIDUAL/ENTERPRISE
	FirstName      string   `json:"first_name,omitempty"`      // 付款人名（个人类型必填）
	LastName       string   `json:"last_name,omitempty"`       // 付款人姓（个人类型必填）
	CertType       string   `json:"cert_type"`                 // 证件类型：PASSPORT/DRIVING_LICENSE/ID_CARD/BUSINESS_LICENSE
	CertNumber     string   `json:"cert_number"`               // 证件号码
	CertFiles      []string `json:"cert_files"`                // 证件附件ID列表
	Birthday       string   `json:"birthday,omitempty"`        // 生日（ISO-8601日期格式，个人类型必填）
	Region         string   `json:"region"`                    // 国家/地区编码（ISO-3166 2字符）
	CompanyName    string   `json:"company_name,omitempty"`    // 公司名称（企业类型必填）
	RegisterNumber string   `json:"register_number,omitempty"` // 公司注册号（企业类型必填）
	BusinessScopes []string `json:"business_scopes,omitempty"` // 企业业务类型（企业类型必填）
	Address        Address  `json:"address"`                   // 付款人地址信息
}

// Payer 付款人信息
type Payer struct {
	PayerID        string   `json:"payer_id"`                  // 付款人ID
	SubjectType    string   `json:"subject_type"`              // 主体类型
	FirstName      string   `json:"first_name,omitempty"`      // 付款人名
	LastName       string   `json:"last_name,omitempty"`       // 付款人姓
	CertType       string   `json:"cert_type"`                 // 证件类型
	CertNumber     string   `json:"cert_number"`               // 证件号码
	CertFiles      []string `json:"cert_files"`                // 证件附件ID列表
	Birthday       string   `json:"birthday,omitempty"`        // 生日
	Region         string   `json:"region"`                    // 国家/地区编码
	CompanyName    string   `json:"company_name,omitempty"`    // 公司名称
	RegisterNumber string   `json:"register_number,omitempty"` // 公司注册号
	BusinessScopes []string `json:"business_scopes,omitempty"` // 企业业务类型
	Address        Address  `json:"address"`                   // 地址信息
}

// PayerResponse 付款人响应
type PayerResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data Payer `json:"data"`
}

// PayerListResponse 付款人列表响应
type PayerListResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		Payers []Payer `json:"payers"` // 付款人列表
	} `json:"data"`
}

// ClearingNetworkRequest 查询可用清算网络请求
type ClearingNetworkRequest struct {
	PayeeAccountID  string  `json:"payee_account_id"`  // 收款人账户ID
	PayCurrency     string  `json:"pay_currency"`      // 付款币种（ISO-4217）
	Amount          float64 `json:"amount"`            // 金额
	AmountType      string  `json:"amount_type"`       // 金额类型：PAY_AMOUNT/RECEIVE_AMOUNT
	ReceiveCurrency string  `json:"receive_currency"`  // 收款币种（ISO-4217）
}

// ClearingNetwork 清算网络信息
type ClearingNetwork struct {
	Network               string  `json:"network"`                 // 清算网络：SWIFT/ACH/FPS等
	Fee                   float64 `json:"fee"`                     // 手续费
	EstimatedArrivalTime  string  `json:"estimated_arrival_time"`  // 预计到账时间
}

// ClearingNetworkResponse 查询可用清算网络响应
type ClearingNetworkResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	ClearingNetworks []ClearingNetwork `json:"clearing_networks"` // 清算网络列表
	PayeeAccountID   string            `json:"payee_account_id"`  // 收款人账户ID
	PayCurrency      string            `json:"pay_currency"`      // 付款币种
	ReceiveCurrency  string            `json:"receive_currency"`  // 收款币种
	Amount           float64           `json:"amount"`            // 金额
	AmountType       string            `json:"amount_type"`       // 金额类型
	Country          string            `json:"country"`           // 国家/地区
}

// QuoteRequest 申请锁汇请求
type QuoteRequest struct {
	PayeeAccountID         string  `json:"payee_account_id"`                   // 收款人账户ID
	PayerID                string  `json:"payer_id,omitempty"`                 // 付款人ID
	Purpose                string  `json:"purpose"`                            // 汇款目的
	PayCurrency            string  `json:"pay_currency"`                       // 付款币种（ISO-4217）
	ReceiveCurrency        string  `json:"receive_currency,omitempty"`         // 收款币种（ISO-4217）
	Amount                 float64 `json:"amount"`                             // 金额
	AmountType             string  `json:"amount_type"`                        // 金额类型：PAY_AMOUNT/RECEIVE_AMOUNT
	ClearingNetwork        string  `json:"clearing_network,omitempty"`         // 清算网络
	AbaNumber              string  `json:"aba_number,omitempty"`               // ABA码
	FpsBankID              string  `json:"fps_bank_id,omitempty"`              // FPS码
	IfsCode                string  `json:"ifs_code,omitempty"`                 // IFS码
	IntermediarySwiftCode  string  `json:"intermediary_swift_code,omitempty"`  // 中间行Swift码
	Remark                 string  `json:"remark,omitempty"`                   // 汇款备注
}

// AmountInfo 金额信息
type AmountInfo struct {
	Currency string  `json:"currency"` // 币种
	Amount   float64 `json:"amount"`   // 金额
}

// Quote 锁汇单信息
type Quote struct {
	QuoteID       string     `json:"quote_id"`        // 锁汇ID
	PaymentMethod string     `json:"payment_method"`  // 付款方式
	PayAmount     AmountInfo `json:"pay_amount"`      // 付款金额
	ReceiveAmount AmountInfo `json:"receive_amount"`  // 收款金额
	Surcharge     AmountInfo `json:"surcharge"`       // 手续费
	ExchangeRate  float64    `json:"exchange_rate"`   // 汇率
	ExpireAt      string     `json:"expire_at"`       // 过期时间
}

// QuoteResponse 申请锁汇响应
type QuoteResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data Quote `json:"data"`
}

// OrderRequest 提交付款订单请求
type OrderRequest struct {
	QuoteID       string `json:"quote_id"`        // 锁汇ID
	ClientOrderID string `json:"client_order_id"` // 客户系统唯一订单ID
}

// RemittanceOrder 付款订单信息
type RemittanceOrder struct {
	OrderID        string     `json:"order_id"`         // 订单ID
	OrderSource    string     `json:"order_source"`     // 订单来源：API/PORTAL
	ClientOrderID  string     `json:"client_order_id"`  // 客户订单ID
	CreateTime     string     `json:"create_time"`      // 创建时间
	FinishTime     string     `json:"finish_time"`      // 完成时间
	Status         string     `json:"status"`           // 状态：CREATED/SUCCESS/FAILED等
	PayeeID        string     `json:"payee_id"`         // 收款人ID
	PayeeAccountID string     `json:"payee_account_id"` // 收款账户ID
	PaymentMethod  string     `json:"payment_method"`   // 付款方式
	PayAmount      AmountInfo `json:"pay_amount"`       // 付款金额
	ReceiveAmount  AmountInfo `json:"receive_amount"`   // 收款金额
	Surcharge      AmountInfo `json:"surcharge"`        // 手续费
	ExchangeRate   float64    `json:"exchange_rate"`    // 汇率
	ErrorMessage   string     `json:"error_message"`    // 错误消息
}

// OrderResponse 提交付款订单响应
type OrderResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data RemittanceOrder `json:"data"`
}

// OrderListRequest 查询付款单列表请求
type OrderListRequest struct {
	Page          int    `json:"page"`                     // 页码，从1开始
	Limit         int    `json:"limit"`                    // 每页条数
	PayeeID       string `json:"payee_id,omitempty"`       // 收款人ID
	PayerID       string `json:"payer_id,omitempty"`       // 付款人ID
	TimeStart     string `json:"time_start,omitempty"`     // 查询起始时间（含）
	TimeEnd       string `json:"time_end,omitempty"`       // 查询截止时间（不含）
	OrderID       string `json:"order_id,omitempty"`       // 平台付款单号
	ClientOrderID string `json:"client_order_id,omitempty"`// 客户付款单号
}

// OrderListResponse 查询付款单列表响应
type OrderListResponse struct {
	Result struct {
		Result  string `json:"result"`  // 结果：S-成功，F-失败
		Code    string `json:"code"`    // 结果代码
		Message string `json:"message"` // 结果消息
	} `json:"result"`
	Data struct {
		Query      map[string]interface{} `json:"query"`       // 查询条件
		Orders     []RemittanceOrder      `json:"orders"`      // 订单列表
		Page       int                    `json:"page"`        // 当前页码
		Limit      int                    `json:"limit"`       // 每页数量
		TotalCount int                    `json:"total_count"` // 总记录数
		TotalPage  int                    `json:"total_page"`  // 总页数
	} `json:"data"`
}
