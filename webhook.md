# Webhooks 回调通知文档
## 一、通用说明
### 1.1 推送规则
- **请求方式**：所有 Webhook 通知均通过 `POST` 方法推送
- **端口限制**：仅支持 80（HTTP）或 443（HTTPS）端口
- **回调地址**：需在商户 Portal 提前配置
- **响应要求**：对接系统处理完成后需返回 `200 OK` 状态码，否则视为推送失败
- **重试机制**：推送失败后将按以下间隔重试，累计失败 10 次后停止自动重试：
  ```
  5s → 10s → 15s → 30s → 1min → 5min → 10min × 3 → 停止
  ```
- **手动重试**：失败的通知可在商户 Portal-设置-开发者-开发者通知中手动触发重试
- **签名验证**：Webhook 推送遵循与 API 相同的签名规则，**强烈建议**接收后先验签再处理业务逻辑
- **查询入口**：所有推送记录可在商户 Portal 对应页面查询

### 1.2 通用请求头
| 参数名 | 是否必填 | 类型 | 说明 |
|--------|----------|------|------|
| x-appid | 是 | string | 从 Portal 获取的 appid（参考 Signature Guide） |
| authorization | 是 | string | 请求签名（参考 Signature Guide） |

### 1.3 通用请求体结构
```json
{
  "business_type": "string",  // 事件类型（核心字段，用于区分回调数据格式）
  "event_time": "string",     // 事件发生时间（ISO-8601 格式，如 2025-05-20T13:08:02Z）
  "business_id": "string",    // 事件实体唯一 ID（如交易 ID、订单 ID）
  "data": {}                  // 业务数据（结构随 business_type 变化）
}
```

## 二、发行卡事件推送 Webhook
### 2.1 事件类型枚举
| 事件类型 | 描述 | 备注 |
|----------|------|------|
| CARD_TRANSACTION | 卡交易通知 | 包含授权、消费、退款等卡交易事件 |
| CARD_STATUS_UPDATE | 卡状态更新 | 如激活、冻结、解冻、销卡等状态变更 |
| THREE_DS_VERIFICATION_CODE | 3DS 交易验证码 | 3DS 验证所需的验证码通知 |
| CARD_TOKEN_OTP_CODE | 绑卡验证码 | Apple Pay、Google Pay 等绑卡场景的验证码 |
| CARD_VERIFICATION_CODE | 卡交易验证码 | Paypal、Facebook 等场景的交易验证码 |
| CARD_ACTIVATION_CODE | 卡激活码 | 实体卡激活所需激活码（v1.19 新增） |
| CARD_PIN_RETRY_LIMIT | 卡密重试限制 | 密码输入错误次数达到上限（v1.19 新增） |

### 2.2 示例：卡交易通知（CARD_TRANSACTION）
```json
{
  "business_type": "CARD_TRANSACTION",
  "event_time": "2025-05-20T13:08:02Z",
  "business_id": "111111",
  "data": {
    "transaction_id": "2025052013080222465300228999",
    "card_id": "2025040509335657270600479999",
    "mask_card_number": "536025******4433",
    "transaction_time": "2025-05-20T13:08:02Z",
    "confirm_time": "2025-05-20T13:08:02Z",
    "transaction_amount": {
      "currency": "USD",
      "amount": 100.00
    },
    "accounting_amount": {
      "currency": "USD",
      "amount": 100.00
    },
    "surcharge": {
      "currency": "USD",
      "amount": 2.00
    },
    "biz_type": "AUTH",  // 交易类型：AUTH(授权)、PURCHASE(消费)、REFUND(退款)等
    "status": "AUTHORIZED",  // 交易状态：AUTHORIZED(已授权)、COMPLETED(已完成)、FAILED(失败)等
    "status_description": "succeed",
    "merchant_name": "FACEBK",  // 商户名称
    "merchant_region": "US"     // 商户所在地区
  }
}
```

## 三、支付事件推送 Webhook
### 3.1 事件类型枚举
| 事件类型 | 描述 | 备注 |
|----------|------|------|
| PAYEE_ACCOUNT_ACTIVE | 收款人账户启用 | 收款人账户状态变更为 ACTIVE |
| REMITTANCE_FAIL | 付款订单失败 | 对外付款订单执行失败 |
| REMITTANCE_COMPLETE | 付款订单完成 | 对外付款订单执行成功 |
| REMITTANCE_REVERSE | 付款订单退票 | 付款完成后被银行退回（仅中国收款账户可能发生，最晚 3 天内推送） |
| PAYEE_DEACTIVATED | 收款人失效 | 收款人状态变更为失效 |

### 3.2 示例：收款人账户启用（PAYEE_ACCOUNT_ACTIVE）
```json
{
  "business_type": "PAYEE_ACCOUNT_ACTIVE",
  "event_time": "2025-05-20T13:08:02Z",
  "business_id": "111111",
  "data": {
    "payee_id": "2025060601553798427600368999",  // 收款人 ID
    "account_id": "2025060601553951627600368999",  // 收款账户 ID
    "payment_method": "PAYPAL_USD"  // 支付方式
  }
}
```

## 四、收单事件推送 Webhook
### 4.1 事件类型枚举
| 事件类型 | 描述 | 备注 |
|----------|------|------|
| ACQUIRING_PAYMENT_SUCCEED | 收单付款成功 | 收单支付订单执行成功 |
| ACQUIRING_PAYMENT_FAILED | 收单付款失败 | 收单支付订单执行失败 |
| ACQUIRING_REFUND_SUCCEED | 收单退款成功 | 收单退款订单执行成功 |
| ACQUIRING_REFUND_FAILED | 收单退款失败 | 收单退款订单执行失败 |
| ACQUIRING_CAPTURE_SUCCEED | 收单请款成功 | 预授权请款执行成功 |
| ACQUIRING_CAPTURE_FAILED | 收单请款失败 | 预授权请款执行失败 |
| ACQUIRING_AUTH_TOKEN | 支付授权结果通知 | 钱包/卡授权支付的授权结果（如 access_token 生成） |

### 4.2 示例：收单付款成功（ACQUIRING_PAYMENT_SUCCEED）
```json
{
  "business_type": "ACQUIRING_PAYMENT_SUCCEED",
  "event_time": "2025-05-20T13:08:02Z",
  "business_id": "111111",
  "data": {
    "payment_method": "ALIPAY_CN",  // 支付方式
    "payment_status": "SUCCESS",    // 支付状态
    "payment_result_message": "Success",  // 结果描述
    "payment_request_id": "PAY_202506271047000000",  // 商户支付请求 ID
    "payment_id": "P2025062703481275748000226141",    // 平台支付订单 ID
    "payment_amount": {
      "currency": "USD",
      "amount": 100.00
    },
    "surcharge": {  // 手续费
      "currency": "USD",
      "amount": 3.00
    },
    "gross_settlement_amount": {  // 结算金额（含手续费）
      "currency": "USD",
      "amount": 97.00
    },
    "settlement_quote": {  // 结算汇率
      "exchange_rate": 7.2,
      "quote_time": "2025-05-20T13:07:02Z"
    },
    "payment_create_time": "2025-05-20T12:06:02Z",  // 支付创建时间
    "payment_time": "2025-05-20T12:08:02Z"          // 支付完成时间
  }
}
```

## 五、回调处理建议
1. **幂等性处理**：基于 `business_id` 做幂等性控制，避免重复处理同一事件
2. **签名验证**：严格按照签名规则验证 `authorization` 头，防止伪造回调
3. **异步处理**：回调接收后立即返回 200，业务逻辑异步处理（如放入消息队列）
4. **日志记录**：完整记录回调原始数据、处理状态、错误信息，便于问题排查
5. **超时控制**：设置合理的业务处理超时时间，避免阻塞回调响应