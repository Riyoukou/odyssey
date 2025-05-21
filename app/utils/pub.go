package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Hash 密码
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// 生成 AuthToken
func GenerateAuthToken(username string) string {
	data := username + ":" + GenerateSecret(8)
	return base64.StdEncoding.EncodeToString([]byte(data))
}

func GenerateSecret(length int) string {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length] // 截取指定长度
}

// SendRequest 发送 HTTP 请求
func SendRequest(method, url string, headers map[string]string, data map[string]interface{}) (int, []byte, error) {
	// 将请求数据转换为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(context.Background(), method, url, bytes.NewReader(jsonData))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 创建 HTTP 客户端并发送请求
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应数据
	respBody, err := io.ReadAll(resp.Body) // 推荐使用 io.ReadAll 替代 ioutil.ReadAll
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read response: %w", err)
	}

	return resp.StatusCode, respBody, nil
}

// RequestData 用于存储请求数据
type RequestData struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    interface{}
}

// NewRequestData 创建一个新的请求数据
func NewRequestData(method, url string, headers map[string]string, body interface{}) *RequestData {
	return &RequestData{
		Method:  method,
		URL:     url,
		Headers: headers,
		Body:    body,
	}
}

// HttpRequest 执行 HTTP 请求
func HttpRequest(reqData *RequestData) (string, error) {
	// 创建一个 HTTP 客户端
	client := &http.Client{
		Timeout: 10 * time.Second, // 设置超时
	}

	// 准备请求的 body 数据
	var reqBody []byte
	var err error
	if reqData.Body != nil {
		reqBody, err = json.Marshal(reqData.Body) // 将 body 转换成 JSON
		if err != nil {
			return "", errors.New("failed to marshal body")
		}
	}

	// 创建请求对象
	req, err := http.NewRequest(reqData.Method, reqData.URL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	// 设置请求头
	for key, value := range reqData.Headers {
		req.Header.Set(key, value)
	}

	// 执行请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 使用 io.ReadAll 替代 ioutil.ReadAll
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 返回响应内容
	return string(body), nil
}
