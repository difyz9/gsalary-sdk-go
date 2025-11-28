# API 使用指南

## 项目结构

```
g-salary-api/
├── config.go              # 配置和密钥管理
├── entities.go            # 授权头实体
├── request.go             # 请求签名和验证
├── client.go              # HTTP客户端
├── gsalary_test.go        # 核心功能单元测试
├── api/                   # API接口封装
│   ├── client.go          # 统一API客户端
│   ├── card.go            # 卡片API
│   ├── types.go           # 请求/响应数据结构
│   └── card_test.go       # 卡片API测试
└── examples/              # 示例代码
    └── card_apply_demo.go # 卡片申请示例
```

## 快速开始

### 1. 初始化配置

```go
import (
    gsalary "github.com/difyz9/gsalary-sdk-go"
    "github.com/difyz9/gsalary-sdk-go/api"
)

// 创建配置
config := gsalary.NewConfig()
config.AppID = "your-app-id"
config.Endpoint = "https://api-test.gsalary.com" // 测试环境
// config.Endpoint = "https://api.gsalary.com"   // 生产环境

// 加载密钥
err := config.ConfigClientPrivateKeyPEMFile("private_key_prod.pem")
err = config.ConfigServerPublicKeyPEMFile("plate_key_prod.pem")
```

### 2. 创建客户端

```go
// 创建API客户端（推荐，提供高层封装）
client := api.NewClient(config)

// 或者创建底层客户端（需要手动构造请求）
rawClient := gsalary.NewClient(config)
```

### 3. 调用API

#### 申请卡片

```go
req := &api.CardApplyRequest{
    RequestID:           "unique_request_id",
    ProductCode:         "VIRTUAL_CARD_USD",
    Currency:            "USD",
    CardHolderID:        "holder_001",
    LimitPerDay:         1000.00,
    LimitPerMonth:       5000.00,
    LimitPerTransaction: 500.00,
    InitBalance:         100.00,
}

resp, err := client.Card.ApplyCard(req)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("卡片申请成功，状态: %s\n", resp.Data.Status)
```

## 运行测试

### 运行核心SDK测试

```bash
go test -v
```

### 运行API测试（需要真实密钥）

```bash
# 使用环境变量指定密钥文件
export GSALARY_PRIVATE_KEY_FILE=/path/to/private_key_prod.pem
export GSALARY_PUBLIC_KEY_FILE=/path/to/plate_key_prod.pem

go test -v ./api
```

### 运行示例

```bash
cd examples
go run card_apply_demo.go
```

## 如何添加新的API

### 1. 在 `api/types.go` 添加请求/响应结构

```go
type NewAPIRequest struct {
    Field1 string  `json:"field1"`
    Field2 int     `json:"field2"`
}

type NewAPIResponse struct {
    Result ResultInfo `json:"result"`
    Data   struct {
        Field string `json:"field"`
    } `json:"data"`
}
```

### 2. 创建新的API文件（如 `api/newapi.go`）

```go
package api

import (
    "fmt"
    gsalary "github.com/difyz9/gsalary-sdk-go"
)

type NewAPI struct {
    client *gsalary.GSalaryClient
}

func NewNewAPI(client *gsalary.GSalaryClient) *NewAPI {
    return &NewAPI{client: client}
}

func (api *NewAPI) CallNewAPI(req *NewAPIRequest) (*NewAPIResponse, error) {
    request := gsalary.NewRequest("POST", "/v1/new_api_path")
    request.Body = map[string]interface{}{
        "field1": req.Field1,
        "field2": req.Field2,
    }
    
    resp, err := api.client.Request(request)
    if err != nil {
        return nil, fmt.Errorf("api call failed: %w", err)
    }
    
    // 解析响应...
    return &response, nil
}
```

### 3. 在 `api/client.go` 添加新API

```go
type Client struct {
    client *gsalary.GSalaryClient
    Card   *CardAPI
    NewAPI *NewAPI  // 添加这行
}

func NewClient(config *gsalary.GSalaryConfig) *Client {
    gsalaryClient := gsalary.NewClient(config)
    
    return &Client{
        client: gsalaryClient,
        Card:   NewCardAPI(gsalaryClient),
        NewAPI: NewNewAPI(gsalaryClient),  // 添加这行
    }
}
```

## 注意事项

1. **RequestID唯一性**: 每次请求都应该使用唯一的 `request_id`
2. **密钥安全**: 不要将私钥提交到版本控制系统
3. **错误处理**: 所有API调用都应该处理错误
4. **环境切换**: 测试和生产环境使用不同的 Endpoint
5. **签名验证**: SDK会自动处理请求签名和响应验证

## 常见问题

### Q: 如何验证签名是否正确？

运行核心测试：
```bash
go test -v -run TestSignature
```

### Q: 如何调试API请求？

可以在 `client.go` 的 `Request` 方法中添加日志输出查看请求详情。

### Q: 遇到 "Appid invalid" 错误？

检查：
1. AppID 是否正确
2. 密钥文件是否正确加载
3. 使用的是测试还是生产环境的配置

### Q: 如何使用自己的HTTP客户端？

底层 `GSalaryClient` 使用标准的 `http.Client`，可以通过设置 `http.DefaultClient` 或在代码中自定义。
