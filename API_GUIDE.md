# GSalary SDK API Complete Guide

Complete API reference for all implemented modules in the GSalary Golang SDK.

## 📋 Table of Contents

1. [Core SDK](#core-sdk)
2. [Card Management](#card-management)
3. [Cardholder Management](#cardholder-management)
4. [Wallet Management](#wallet-management)
5. [Exchange Management](#exchange-management)
6. [Payment Acquiring](#payment-acquiring)
7. [Payee Management](#payee-management)
8. [Payer Management](#payer-management)
9. [Remittance Management](#remittance-management)
10. [Webhook Handler](#webhook-handler)

## Core SDK

All modules implemented with RSA-SHA256 signature verification and automatic request signing.

### Configuration

```go
config := gsalary.NewConfig()
config.AppID = "your-app-id"
config.Endpoint = "https://api.gsalary.com"
config.ConfigClientPrivateKeyPEMFile("private_key.pem")
config.ConfigServerPublicKeyPEMFile("server_public_key.pem")
```

## Card Management (api/card.go)

### Implemented APIs

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| ApplyCard | POST /v1/card_applies | Apply for new card | ✅ |
| GetCardApplyResult | GET /v1/card_applies/{request_id} | Get application result | ✅ |
| GetAvailableQuotas | GET /v1/cards/available_quotas | Query card quotas | ✅ |
| GetProducts | GET /v1/card_support/products | Get product list | ✅ |
| GetCardList | GET /v1/cards | Get card list | ✅ |
| GetCardInfo | GET /v1/cards/{card_id} | Get card details | ✅ |

## Cardholder Management (api/cardholder.go)

### Implemented APIs

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| AddCardHolder | POST /v1/card_holders | Add cardholder | ✅ |
| GetCardHolderList | GET /v1/card_holders | Get cardholder list | ✅ |
| GetCardHolderInfo | GET /v1/card_holders/{holder_id} | Get cardholder details | ✅ |
| UpdateCardHolder | PUT /v1/card_holders/{holder_id} | Update cardholder | ✅ |

## Wallet Management (api/wallet.go)

### Implemented APIs

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GetBalance | GET /v1/wallets/balance | Query wallet balance | ✅ |

## Exchange Management (api/exchange.go)

### Implemented APIs

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GetCurrentExchangeRate | GET /v1/exchange/current_exchange_rate | Get current rate | ✅ |
| RequestQuote | POST /v1/exchange/quotes | Request quote (lock rate) | ✅ |
| SubmitExchangeRequest | POST /v1/exchange/submit_request | Submit exchange order | ✅ |
| GetExchangeOrders | GET /v1/exchange/orders | Query order list | ✅ |

## Payment Acquiring (api/payment.go)

### Implemented APIs

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| PaymentConsult | POST /gateway/v1/acquiring/pay_consult | Payment consultation | ✅ |
| CreatePaymentSession | POST /gateway/v1/acquiring/pay_session | Create checkout session | ✅ |
| CreateEasySafePaySession | POST /gateway/v1/acquiring/easy_safe_pay/pay_session | Create wallet auth session | ✅ |
| EasySafePayPay | POST /gateway/v1/acquiring/easy_safe_pay/pay | Wallet auth payment | ✅ |
| CreateCardAutoDebitSession | POST /gateway/v1/acquiring/card_auto_debit/pay_session | Create card auth session | ✅ |
| CardAutoDebitPay | POST /gateway/v1/acquiring/card_auto_debit/pay | Card auth payment | ✅ |
| RefreshAuthToken | POST /gateway/v1/acquiring/auth_refresh_token | Refresh access token | ✅ |
| RevokeAuthToken | POST /gateway/v1/acquiring/auth_revoke_token | Revoke authorization | ✅ |
| CancelPayment | POST /gateway/v1/acquiring/cancel | Cancel payment | ✅ |
| QueryPayment | GET /gateway/v1/acquiring/query | Query payment status | ✅ |

## Payee Management (api/payee.go)

### Implemented APIs

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| AddPayee | POST /v1/payees | Add payee | ✅ |
| GetPayeeList | GET /v1/payees | Get payee list | ✅ |
| UpdatePayee | PUT /v1/payees/{payee_id} | Update payee | ✅ |
| DeletePayee | DELETE /v1/payees/{payee_id} | Delete payee | ✅ |
| AddPayeeAccount | POST /v1/payees/{payee_id}/accounts | Add payee account | ✅ |
| GetPayeeAccounts | GET /v1/payees/{payee_id}/accounts | Get payee accounts | ✅ |
| UpdatePayeeAccount | PUT /v1/payees/{payee_id}/accounts/{account_id} | Update account | ✅ |
| DeletePayeeAccount | DELETE /v1/payees/{payee_id}/accounts/{account_id} | Delete account | ✅ |

## Payer Management (api/payer.go)

### Implemented APIs

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| UploadAttachment | POST /remittance/attachments | Upload attachment | ✅ |
| AddPayer | POST /remittance/payers | Add payer | ✅ |
| GetPayerList | GET /remittance/payers | Get payer list | ✅ |
| GetPayer | GET /remittance/payers/{payer_id} | Get payer details | ✅ |
| UpdatePayer | PUT /remittance/payers/{payer_id} | Update payer | ✅ |
| DeletePayer | DELETE /remittance/payers/{payer_id} | Delete payer | ✅ |

## Remittance Management (api/remittance.go)

### Implemented APIs

| Method | Endpoint | Description | Status |
|--------|----------|-------------|--------|
| GetClearingNetworks | GET /remittance/clearing_networks | Query clearing networks | ✅ |
| CreateQuote | POST /remittance/quotes | Create remittance quote | ✅ |
| SubmitRemittanceOrder | POST /remittance/orders | Submit remittance order | ✅ |
| GetRemittanceOrderList | GET /remittance/orders | Get order list | ✅ |
| GetRemittanceOrder | GET /remittance/orders/{order_id} | Get order details | ✅ |

## Webhook Handler (api/webhook.go)

### Supported Events

| Event Type | Constant | Description |
|------------|----------|-------------|
| Payment Result | EventAcquiringPaymentResult | Payment result notification |
| Auth Token | EventAcquiringAuthToken | Authorization token notification |
| Card Status Update | EventCardStatusUpdate | Card status change notification |
| Card Transaction | EventCardTransaction | Card transaction notification |
| Card Adjust Result | EventCardAdjustResult | Card adjustment result |
| Card Apply Result | EventCardApplyResult | Card application result |
| Exchange Order Result | EventExchangeOrderResult | Exchange order result |
| Remittance Order Result | EventRemittanceOrderResult | Remittance order result |
| Payee Deactivated | EventPayeeDeactivated | Payee deactivation notification |

### Usage Example

```go
webhookHandler := api.NewWebhookHandler(config)

http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
    resp, err := webhookHandler.HandleWebhook(r)
    // Process webhook...
    json.NewEncoder(w).Encode(resp)
})
```

## Test Coverage

All modules have comprehensive unit tests:
- `api/card_test.go` - Card API tests
- `api/cardholder_test.go` - Cardholder API tests
- `api/wallet_test.go` - Wallet API tests
- `api/exchange_test.go` - Exchange API tests
- `api/payment_test.go` - Payment API tests
- `api/webhook_test.go` - Webhook handler tests

Run tests:
```bash
go test -v ./api
```

## Examples

### Card Operations Example
See `examples/card_apply_demo.go` for complete card workflow.

### Webhook Server Example
See `examples/webhook/main.go` for webhook handling server.

## Notes

- All API methods return `(*Response, error)` format
- Signature verification is automatic
- Webhook signature validation included
- Comprehensive error handling with business result codes
- Support for both test and production environments

