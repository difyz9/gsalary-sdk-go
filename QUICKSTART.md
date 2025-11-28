# G-Salary API Golang SDK - 快速开始指南

## 第一步：生成密钥对

运行密钥生成脚本：

```bash
python3 generate_keys.py
```

**输出文件：**
- `private_key.pem` - 客户端私钥（⚠️ 请妥善保管，不要泄露）
- `public_key.pem` - 客户端公钥（需要提交到平台）

## 第二步：配置公钥到平台

1. 访问：https://b.gsalary.com/#/config/developer
2. 使用邮箱登录：1298741189@qq.com
3. 将 `public_key.pem` 的内容复制粘贴到"客户端公钥"配置中
4. 保存配置

## 第三步：获取服务端公钥

从 G-Salary 平台获取服务端公钥，保存为 `server_public_key.pem`

## 第四步：编写代码

创建 `main.go`：

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    
    gsalary "github.com/difyz9/gsalary-sdk-go"
)

func main() {
    // 1. 创建配置
    config := gsalary.NewConfig()
    config.AppID = "6bc14a48-e6xxxxxxx0-9105dcad37ad"
    config.Endpoint = "https://api-test.gsalary.com"
    
    // 2. 加载密钥
    if err := config.ConfigClientPrivateKeyPEMFile("private_key.pem"); err != nil {
        log.Fatal("加载私钥失败:", err)
    }
    if err := config.ConfigServerPublicKeyPEMFile("server_public_key.pem"); err != nil {
        log.Fatal("加载公钥失败:", err)
    }
    
    // 3. 创建客户端
    client := gsalary.NewClient(config)
    
    // 4. 发起 GET 请求
    fmt.Println("=== 查询卡列表 ===")
    request := gsalary.NewRequest("GET", "/v1/cards")
    request.QueryArgs = map[string]string{
        "page":  "1",
        "limit": "20",
    }
    
    resp, err := client.Request(request)
    if err != nil {
        log.Printf("请求失败: %v\n", err)
    } else {
        printJSON(resp)
    }
    
    // 5. 发起 POST 请求
    fmt.Println("\n=== 创建汇率报价 ===")
    postRequest := gsalary.NewRequest("POST", "/v1/exchange/quotes")
    postRequest.Body = map[string]interface{}{
        "sell_currency": "USD",
        "buy_currency":  "CNY",
        "sell_amount":   100.0,
    }
    
    resp, err = client.Request(postRequest)
    if err != nil {
        log.Printf("请求失败: %v\n", err)
    } else {
        printJSON(resp)
    }
}

func printJSON(data interface{}) {
    jsonBytes, _ := json.MarshalIndent(data, "", "  ")
    fmt.Println(string(jsonBytes))
}
```

## 第五步：运行程序

```bash
go run main.go
```

## 常见问题

### Q: 签名验证失败怎么办？

**A:** 检查以下几点：
1. 确认客户端公钥已正确配置到平台
2. 确认服务端公钥是从平台获取的最新版本
3. 确认 AppID 正确
4. 确认使用的是正确的环境（测试/生产）

### Q: 如何切换到生产环境？

**A:** 修改 Endpoint：
```go
config.Endpoint = "https://api.gsalary.com"  // 生产环境
```

### Q: 如何处理错误？

**A:** SDK 提供了详细的错误信息：
```go
resp, err := client.Request(request)
if err != nil {
    // 检查是否是业务错误
    if gsalaryErr, ok := err.(*gsalary.GSalaryException); ok {
        fmt.Printf("业务错误: [%s - %s] %s\n", 
            gsalaryErr.BizCode, 
            gsalaryErr.ErrorCode, 
            gsalaryErr.Message)
    } else {
        // 系统错误
        fmt.Printf("系统错误: %v\n", err)
    }
    return
}
```

### Q: 如何查看请求详情？

**A:** 可以在发送请求前打印请求信息：
```go
fmt.Printf("请求方法: %s\n", request.Method)
fmt.Printf("请求路径: %s\n", request.PathWithArgs(false))
fmt.Printf("请求Body: %+v\n", request.Body)
```

## 支持的 API 方法

| HTTP 方法 | 说明 |
|----------|------|
| GET | 查询数据 |
| POST | 创建数据 |
| PUT | 更新数据 |
| DELETE | 删除数据 |

## 下一步

- 📖 阅读 [完整文档](README_GO.md)
- 🔍 查看 [API 文档](https://api.gsalary.com/doc/index.html)
- 🆚 查看 [Python vs Golang 对比](SDK_COMPARISON.md)
- 💻 查看更多 [示例代码](example/main.go)

## 测试 SDK

运行单元测试：
```bash
go test -v
```

运行示例代码：
```bash
cd example
go run main.go
```

## 获取帮助

如遇到问题，请检查：
1. API 文档：https://api.gsalary.com/doc/index.html
2. 开发者控制台：https://b.gsalary.com/#/config/developer
3. Python SDK 参考：https://github.com/gsalary-develop/gsalary-sdk-python
