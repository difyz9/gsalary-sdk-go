package gsalary

import (
	"fmt"
	"net/url"
	"strings"
)

// AuthorizeHeaderInfo 鉴权头部信息
type AuthorizeHeaderInfo struct {
	Algorithm string
	Timestamp string
	Signature string
}

// NewAuthorizeHeaderInfo 创建新的鉴权头部信息
func NewAuthorizeHeaderInfo(algorithm, timestamp, signature string) *AuthorizeHeaderInfo {
	return &AuthorizeHeaderInfo{
		Algorithm: algorithm,
		Timestamp: timestamp,
		Signature: signature,
	}
}

// Valid 检查鉴权头部信息是否有效
func (a *AuthorizeHeaderInfo) Valid() bool {
	return a.Algorithm != "" && a.Timestamp != "" && a.Signature != ""
}

// ToHeaderValue 转换为HTTP头部值
func (a *AuthorizeHeaderInfo) ToHeaderValue() string {
	return fmt.Sprintf("algorithm=%s,time=%s,signature=%s",
		a.Algorithm,
		a.Timestamp,
		url.QueryEscape(a.Signature))
}

// FromHeaderValue 从HTTP头部值解析鉴权信息
func FromHeaderValue(headerValue string) *AuthorizeHeaderInfo {
	if headerValue == "" {
		return &AuthorizeHeaderInfo{}
	}

	parts := strings.Split(headerValue, ",")
	if len(parts) != 3 {
		return &AuthorizeHeaderInfo{}
	}

	var algorithm, timestamp, signature string
	for _, part := range parts {
		keyValue := strings.SplitN(part, "=", 2)
		if len(keyValue) != 2 {
			continue
		}
		key := keyValue[0]
		value := keyValue[1]

		switch key {
		case "algorithm":
			algorithm = value
		case "time":
			timestamp = value
		case "signature":
			sig, err := url.QueryUnescape(value)
			if err == nil {
				signature = sig
			} else {
				signature = value
			}
		}
	}

	return &AuthorizeHeaderInfo{
		Algorithm: algorithm,
		Timestamp: timestamp,
		Signature: signature,
	}
}
