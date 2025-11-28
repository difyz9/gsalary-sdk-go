# G-Salary API 集成项目

## 项目概述

本项目包含 G-Salary API 的 Golang SDK 实现，提供完整的 RSA 签名验证和 API 调用功能。

## 📚 开发文档

- **API 文档**: https://api.gsalary.com/doc/index.html
- **配置管理**: https://b.gsalary.com/#/config/developer

## 🔑 配置说明

⚠️ **安全提示**: 
- 请从 GSalary Portal 获取您的 AppID
- 不要将私钥、AppID 等敏感信息提交到公开仓库
- 使用环境变量或配置文件管理敏感信息

## 📁 项目结构

```
.
├── config.go              # 配置管理（密钥加载、配置）
├── entities.go            # 鉴权头部信息实体
├── request.go             # 请求对象和签名逻辑
├── client.go              # HTTP 客户端实现
├── gsalary_test.go        # 单元测试
├── example/
│   └── main.go            # 使用示例
├── generate_keys.py       # RSA 密钥生成工具
├── README_GO.md           # Golang SDK 详细文档
├── SDK_COMPARISON.md      # Python vs Golang SDK 对比
└── go.mod                 # Go 模块配置
```

## 🚀 快速开始

### 1. 生成密钥对

```bash
python3 generate_keys.py
```

生成后会得到：
- `private_key.pem` - 客户端私钥（保密）
- `public_key.pem` - 客户端公钥（需提交到 G-Salary 平台配置）

### 2. 配置公钥

访问 https://b.gsalary.com/#/config/developer 将生成的公钥配置到平台。

### 3. 使用 SDK

```go
package main

import (
    "log"
    gsalary "github.com/difyz9/gsalary-sdk-go"
)

func main() {
    // 创建配置
    config := gsalary.NewConfig()
    config.AppID = "6bc1xxxxxxxxxae10-9105dcad37ad"
    config.Endpoint = "https://api-test.gsalary.com"
    
    // 加载密钥
    config.ConfigClientPrivateKeyPEMFile("private_key.pem")
    config.ConfigServerPublicKeyPEMFile("server_public_key.pem")
    
    // 创建客户端
    client := gsalary.NewClient(config)
    
    // 发起请求
    request := gsalary.NewRequest("GET", "/v1/cards")
    request.QueryArgs = map[string]string{
        "page":  "1",
        "limit": "20",
    }
    
    resp, err := client.Request(request)
    if err != nil {
        log.Fatal(err)
    }
    
    log.Println(resp)
}
```

## 📖 文档

- **[README_GO.md](README_GO.md)** - Golang SDK 完整使用文档
- **[SDK_COMPARISON.md](SDK_COMPARISON.md)** - Python 和 Golang SDK 对比文档

## ✅ 功能特性

- ✅ RSA-SHA256 签名验证
- ✅ 自动处理请求签名
- ✅ 自动验证响应签名
- ✅ 支持 GET/POST/PUT/DELETE 请求
- ✅ 完整的错误处理
- ✅ 单元测试覆盖
- ✅ 类型安全

## 🧪 运行测试

```bash
go test -v
```

## 📝 参考资料

- **Python SDK**: [gsalary-sdk-python](https://github.com/gsalary-develop/gsalary-sdk-python)
- **API 文档**: https://api.gsalary.com/doc/index.html

## 密钥文件说明

- `private_key_prod.pem` - 生产环境客户端私钥
- `public_key_prod.pem` - 生产环境客户端公钥
- `plate_key_prod.pem` - 平台公钥（用于验证响应）

⚠️ **重要**: 请勿将私钥文件提交到版本控制系统！