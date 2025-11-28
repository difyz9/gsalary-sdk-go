package gsalary

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GSalaryClient 客户端
type GSalaryClient struct {
	config     *GSalaryConfig
	httpClient *http.Client
}

// NewClient 创建新的客户端
func NewClient(config *GSalaryConfig) *GSalaryClient {
	return &GSalaryClient{
		config:     config,
		httpClient: &http.Client{},
	}
}

// GSalaryException 业务异常
type GSalaryException struct {
	BizCode   string `json:"biz_result"`
	ErrorCode string `json:"error_code"`
	Message   string `json:"message"`
}

// Error 实现error接口
func (e *GSalaryException) Error() string {
	return fmt.Sprintf("[%s - %s] %s", e.BizCode, e.ErrorCode, e.Message)
}

// Request 发起请求
func (c *GSalaryClient) Request(request *GSalaryRequest) (map[string]interface{}, error) {
	// 验证请求
	if err := request.Valid(); err != nil {
		return nil, err
	}

	// 生成签名
	authHeader, err := request.SignRequest(c.config)
	if err != nil {
		return nil, fmt.Errorf("failed to sign request: %w", err)
	}

	// 准备请求body
	var reqBody io.Reader
	if request.HasBody() {
		bodyBytes, err := json.Marshal(request.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewReader(bodyBytes)
	}

	// 创建HTTP请求
	fullURL := c.config.ConcatPath(request.PathWithArgs(true))
	httpReq, err := http.NewRequest(request.Method, fullURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("X-Appid", c.config.AppID)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", authHeader.ToHeaderValue())

	// 发送请求
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 处理响应
	if resp.StatusCode == http.StatusOK {
		// 验证响应签名
		respAuthHeader := FromHeaderValue(resp.Header.Get("Authorization"))
		if !respAuthHeader.Valid() {
			return nil, fmt.Errorf("invalid authorization header in response")
		}

		if !request.VerifySignature(c.config, respAuthHeader, string(responseBody)) {
			return nil, fmt.Errorf("signature verification failed")
		}

		// 解析响应
		var result map[string]interface{}
		if err := json.Unmarshal(responseBody, &result); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response: %w", err)
		}

		return result, nil
	}

	// 处理错误响应
	var errResp GSalaryException
	if err := json.Unmarshal(responseBody, &errResp); err != nil {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(responseBody))
	}

	return nil, &errResp
}
