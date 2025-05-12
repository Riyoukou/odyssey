package service

import (
	"encoding/json"
	"fmt"

	"github.com/Riyoukou/odyssey/app/model"
	"github.com/Riyoukou/odyssey/app/repository"
	"github.com/Riyoukou/odyssey/pkg/logger"
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

func CreateBuildRecord(apiBuildRecord model.ApiBuildRecord) error {
	buildRecord := model.BuildRecordTable{
		Name:        apiBuildRecord.Name,
		Env:         apiBuildRecord.Env,
		Tag:         apiBuildRecord.Tag,
		Status:      "Pending",
		ProjectName: apiBuildRecord.ProjectName,
		BuildUser:   apiBuildRecord.BuildUser,
		Description: apiBuildRecord.Describe,
	}

	err := repository.CreateBuildRecord(buildRecord)
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, service := range apiBuildRecord.Services {
		buildServiceRecord := model.BuildServiceRecordTable{
			ServiceName:     service.ServiceName,
			ProjectName:     apiBuildRecord.ProjectName,
			Image:           "Unkonwn",
			BuildURL:        "Unkonwn",
			Status:          "Pending",
			Env:             apiBuildRecord.Env,
			BuildRecordName: apiBuildRecord.Name,
			Branch:          service.Branch,
		}
		err = repository.CreateBuildServiceRecord(buildServiceRecord)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}
