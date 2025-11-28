package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	gsalary "github.com/difyz9/gsalary-sdk-go"
)

func main() {
	// 创建配置
	config := gsalary.NewConfig()
	config.AppID = "your_app_id" // 替换为实际的AppID
	config.Endpoint = "https://api-test.gsalary.com"

	// 加载私钥和公钥
	// 方式1: 从文件加载
	if err := config.ConfigClientPrivateKeyPEMFile("private_key.pem"); err != nil {
		log.Fatalf("Failed to load private key: %v", err)
	}
	if err := config.ConfigServerPublicKeyPEMFile("server_public_key.pem"); err != nil {
		log.Fatalf("Failed to load server public key: %v", err)
	}

	// 方式2: 或者直接传入PEM字符串
	// privateKeyPEM := `-----BEGIN PRIVATE KEY-----
	// ...
	// -----END PRIVATE KEY-----`
	// if err := config.ConfigClientPrivateKeyPEM(privateKeyPEM); err != nil {
	//     log.Fatalf("Failed to load private key: %v", err)
	// }

	// 创建客户端
	client := gsalary.NewClient(config)

	// 示例1: GET请求 - 查询卡列表
	fmt.Println("=== 示例1: GET请求 ===")
	getRequest := gsalary.NewRequest("GET", "/v1/cards")
	getRequest.QueryArgs = map[string]string{
		"create_start": "2024-02-01T00:00:00+00:00",
		"create_end":   "2024-05-01T00:00:00+00:00",
		"page":         "1",
		"limit":        "20",
	}

	resp, err := client.Request(getRequest)
	if err != nil {
		log.Printf("GET request failed: %v", err)
	} else {
		printJSON(resp)
	}

	// 示例2: POST请求 - 创建汇率报价
	fmt.Println("\n=== 示例2: POST请求 ===")
	postRequest := gsalary.NewRequest("POST", "/v1/exchange/quotes")
	postRequest.Body = map[string]interface{}{
		"sell_currency": "USD",
		"buy_currency":  "CNY",
		"sell_amount":   0.1,
	}

	resp, err = client.Request(postRequest)
	if err != nil {
		log.Printf("POST request failed: %v", err)
	} else {
		printJSON(resp)
	}

	// 从环境变量加载配置的示例
	fmt.Println("\n=== 从环境变量加载配置 ===")
	configFromEnv := gsalary.NewConfig()
	configFromEnv.AppID = os.Getenv("GSALARY_APPID")

	privateKeyFile := os.Getenv("GSALARY_CLIENT_PRIVATE_KEY_PEM_FILE")
	if privateKeyFile != "" {
		if err := configFromEnv.ConfigClientPrivateKeyPEMFile(privateKeyFile); err != nil {
			log.Printf("Failed to load private key from env: %v", err)
		}
	}

	publicKeyFile := os.Getenv("GSALARY_SERVER_PUBLIC_KEY_PEM_FILE")
	if publicKeyFile != "" {
		if err := configFromEnv.ConfigServerPublicKeyPEMFile(publicKeyFile); err != nil {
			log.Printf("Failed to load public key from env: %v", err)
		}
	}

	if configFromEnv.AppID != "" && configFromEnv.GetClientPrivateKey() != nil && configFromEnv.GetServerPublicKey() != nil {
		clientFromEnv := gsalary.NewClient(configFromEnv)
		fmt.Println("使用环境变量配置创建客户端成功")
		_ = clientFromEnv
	}
}

// printJSON 格式化输出JSON
func printJSON(data interface{}) {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
