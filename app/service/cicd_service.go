package service

import (
	"encoding/json"
	"fmt"

	"github.com/Riyoukou/odyssey/app/model"
)

func ServiceCICDMap(cicdMap map[string][]model.ServiceCICDForm, action string, clusters []string, updateData model.ServiceCICDForm) []byte {
	// 定义空模板
	empty := model.ServiceCICDForm{}

	switch action {
	case "create":
		for _, name := range clusters {
			cicdMap[name] = []model.ServiceCICDForm{empty}
		}
	case "add":
		for _, name := range clusters {
			if _, exists := cicdMap[name]; !exists {
				cicdMap[name] = []model.ServiceCICDForm{empty}
			}
		}
	case "delete":
		for _, name := range clusters {
			delete(cicdMap, name)
		}
	case "update":
		fmt.Println(updateData)
		for _, name := range clusters {
			if _, exists := cicdMap[name]; exists {
				cicdMap[name][0] = updateData
			}
		}
	case "get":
	}

	// 输出结果
	jsonBytes, err := json.MarshalIndent(cicdMap, "", "  ")

	if err != nil {
		panic(err)
	}

	return jsonBytes
}
