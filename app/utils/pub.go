package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
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
