package utils

import (
	"fmt"
)

// ArgoCDSyncV2Lite 用于触发 ArgoCD 同步请求
func ArgoCDSyncV2LiteV2(syncURL, token string) error {
	// 创建请求数据
	headers := map[string]string{
		"Authorization": "Bearer " + token,
		"Content-Type":  "application/json",
	}
	reqData := NewRequestData("POST", syncURL, headers, nil)

	// 调用 HttpRequest 执行请求
	_, err := HttpRequest(reqData)
	if err != nil {
		return fmt.Errorf("sync request failed: %v", err)
	}

	//fmt.Printf("服务 %s 在 %s 命名空间成功同步！\n", service, namespace)
	return nil
}
