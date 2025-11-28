# GSalary SDK for Golang

GSalary API 的 Golang SDK 实现，提供完整的 RSA 签名验证和 API 调用功能。

## API 文档

- [中文文档](https://api.gsalary.com/doc/index.html?lang=cn)
- [英文文档](https://api.gsalary.com/doc/index.html?lang=en)

## 特性

- ✅ 完整的 RSA 签名和验证
- ✅ 支持 GET、POST、PUT、DELETE 请求
- ✅ 自动处理请求签名和响应验证
- ✅ 类型安全的 API 调用
- ✅ 详细的错误处理

## 安装

```bash
go get github.com/difyz9/gsalary-sdk-go
```

## 快速开始

### 1. 生成密钥对

首先需要生成 RSA 密钥对（如果还没有的话）：

```bash
# 使用提供的 Python 脚本生成
python3 generate_keys.py
```

这将生成 `private_key.pem`（客户端私钥）和 `public_key.pem`（客户端公钥）。

### 2. 配置客户端

```go
package main

import (
    "log"
    gsalary "github.com/difyz9/gsalary-sdk-go"
)

func main() {
    // 创建配置
    config := gsalary.NewConfig()
    config.AppID = "your_app_id"
    config.Endpoint = "https://api-test.gsalary.com"
    
    // 加载客户端私钥
    if err := config.ConfigClientPrivateKeyPEMFile("private_key.pem"); err != nil {
        log.Fatal(err)
    }
    
    // 加载服务端公钥
    if err := config.ConfigServerPublicKeyPEMFile("server_public_key.pem"); err != nil {
        log.Fatal(err)
    }
    
    // 创建客户端
    client := gsalary.NewClient(config)
}
```

### 3. 发起 GET 请求

```go
// 查询卡列表
request := gsalary.NewRequest("GET", "/v1/cards")
request.QueryArgs = map[string]string{
    "create_start": "2024-02-01T00:00:00+00:00",
    "create_end":   "2024-05-01T00:00:00+00:00",
    "page":         "1",
    "limit":        "20",
}

resp, err := client.Request(request)
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp)
```

### 4. 发起 POST 请求

```go
// 创建汇率报价
request := gsalary.NewRequest("POST", "/v1/exchange/quotes")
request.Body = map[string]interface{}{
    "sell_currency": "USD",
    "buy_currency":  "CNY",
    "sell_amount":   0.1,
}

resp, err := client.Request(request)
if err != nil {
    log.Fatal(err)
}

fmt.Println(resp)
```

## 配置方式

### 方式 1: 从文件加载密钥

```go
config := gsalary.NewConfig()
config.ConfigClientPrivateKeyPEMFile("path/to/private_key.pem")
config.ConfigServerPublicKeyPEMFile("path/to/server_public_key.pem")
```

### 方式 2: 直接传入 PEM 字符串

```go
privateKeyPEM := `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQD...
-----END PRIVATE KEY-----`

config := gsalary.NewConfig()
config.ConfigClientPrivateKeyPEM(privateKeyPEM)
```

### 方式 3: 从环境变量加载

```go
config := gsalary.NewConfig()
config.AppID = os.Getenv("GSALARY_APPID")
config.ConfigClientPrivateKeyPEMFile(os.Getenv("GSALARY_CLIENT_PRIVATE_KEY_PEM_FILE"))
config.ConfigServerPublicKeyPEMFile(os.Getenv("GSALARY_SERVER_PUBLIC_KEY_PEM_FILE"))
```

## 错误处理

SDK 提供了详细的错误信息：

```go
resp, err := client.Request(request)
if err != nil {
    // 检查是否是业务异常
    if gsalaryErr, ok := err.(*gsalary.GSalaryException); ok {
        fmt.Printf("业务错误: %s - %s: %s\n", 
            gsalaryErr.BizCode, 
            gsalaryErr.ErrorCode, 
            gsalaryErr.Message)
    } else {
        // 其他错误（网络错误、签名错误等）
        fmt.Printf("错误: %v\n", err)
    }
    return
}
```

## 签名机制

SDK 自动处理以下签名流程：

1. **请求签名**：使用客户端私钥对请求进行 RSA-SHA256 签名
2. **签名内容**：`METHOD PATH\nAPPID\nTIMESTAMP\nBODY_HASH\n`
3. **响应验证**：使用服务端公钥验证响应签名

签名格式：
```
Authorization: algorithm=RSA2,time=1234567890000,signature=base64_encoded_signature
```

## 项目结构

```
.
├── config.go          # 配置管理
├── entities.go        # 鉴权头部信息
├── request.go         # 请求对象和签名逻辑
├── client.go          # HTTP 客户端
├── example/
│   └── main.go        # 使用示例
├── generate_keys.py   # 密钥生成脚本
├── go.mod
└── README.md
```

## 运行示例

```bash
# 1. 生成密钥（如果还没有）
python3 generate_keys.py

# 2. 运行示例代码
cd example
go run main.go
```

## 环境变量配置

可以使用以下环境变量：

```bash
export GSALARY_APPID="your_app_id"
export GSALARY_CLIENT_PRIVATE_KEY_PEM_FILE="./private_key.pem"
export GSALARY_SERVER_PUBLIC_KEY_PEM_FILE="./server_public_key.pem"
```

## 依赖

- Go 1.16+
- 标准库（无第三方依赖）

## 参考

- Python SDK: [gsalary-sdk-python](https://github.com/gsalary-develop/gsalary-sdk-python)

## License

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
