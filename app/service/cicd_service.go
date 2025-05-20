package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Riyoukou/odyssey/app/model"
	"github.com/Riyoukou/odyssey/app/repository"
	"github.com/Riyoukou/odyssey/pkg/logger"
	"github.com/bndr/gojenkins"
)

func ServiceCICDMap(cicdMap map[string][]model.ServiceDeployMap, action string, clusters []string, updateData model.ServiceDeployMap) []byte {
	// 定义空模板
	empty := model.ServiceDeployMap{}

	switch action {
	case "create":
		for _, name := range clusters {
			cicdMap[name] = []model.ServiceDeployMap{empty}
		}
	case "add":
		for _, name := range clusters {
			if _, exists := cicdMap[name]; !exists {
				cicdMap[name] = []model.ServiceDeployMap{empty}
			}
		}
	case "delete":
		for _, name := range clusters {
			delete(cicdMap, name)
		}
	case "update":
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

func CreateDeployRecord(deployRecord model.DeployRecordTable) error {
	deployRecord.Status = "Pending"
	err := repository.CreateDeployRecord(deployRecord)
	if err != nil {
		fmt.Printf("CreateDeployRecord error: %v", err)
		return err
	}

	buildServiceRecords, err := repository.GetBuildServiceRecordsByBuildRecordName(deployRecord.BuildRecordName)
	if err != nil {
		fmt.Printf("GetBuildServiceRecordsByBuildRecordName error: %v", err)
		return err
	}
	for _, cluster := range strings.Split(deployRecord.ClusterNames, ",") {
		for _, buildServiceRecord := range buildServiceRecords {
			deployServiceRecord := model.DeployServiceRecordTable{
				ServiceName:      buildServiceRecord.ServiceName,
				ProjectName:      deployRecord.ProjectName,
				Image:            buildServiceRecord.Image,
				Status:           "Pending",
				ClusterName:      cluster,
				Env:              buildServiceRecord.Env,
				DeployRecordName: deployRecord.Name,
			}
			err = repository.CreateDeployServiceRecord(deployServiceRecord)
			if err != nil {
				fmt.Printf("CreateDeployServiceRecord error: %v", err)
				return err
			}
		}
	}
	return nil
}

func CICDBuildByJenkins(buildRecord model.BuildRecordTable) (string, error) {
	var (
		wg           sync.WaitGroup
		mu           sync.Mutex
		successCount int
	)

	// 获取构建服务记录
	buildServiceRecords, err := repository.GetBuildServiceRecordsByBuildRecordName(buildRecord.Name)
	if err != nil {
		logger.Errorf("获取构建服务记录失败: %v", err)
		return "", err
	}

	for _, buildServiceRecord := range buildServiceRecords {
		time.Sleep(1 * time.Second)
		wg.Add(1)

		go func(record model.BuildServiceRecordTable) {
			defer wg.Done()
			service, err := repository.GetServiceByNameAndProjectByEnv(record.ServiceName, record.ProjectName, record.Env)
			if err != nil {
				logger.Errorf("获取服务失败: %v", err)
			}
			var serviceBuildMap model.ServiceBuildMap
			if err := json.Unmarshal([]byte(service.BuildMap), &serviceBuildMap); err != nil {
				logger.Errorf("解析BuildMap失败: %v", err)
				return
			}
			jenkinsData, err := repository.GetCICDToolByName(serviceBuildMap.CICDTool)
			if err != nil {
				logger.Errorf("获取代码库失败: %v", err)
			}

			jenkinsCredential, err := repository.GetCredentialByName(jenkinsData.CredentialName)
			if err != nil {
				logger.Errorf("获取jenkins密码失败: %v", err)
			}

			var credential map[string]string
			err = json.Unmarshal([]byte(jenkinsCredential.Data), &credential)
			if err != nil {
				panic(err)
			}

			jenkins := gojenkins.CreateJenkins(nil, jenkinsData.URL, credential["username"], credential["password"])
			codelibrary, err := repository.GetCodeLibraryByNameAndProject(service.CodeLibraryName, record.ProjectName)
			if err != nil {
				logger.Errorf("获取代码库失败: %v", err)
			}
			buildID, buildURL, err := BuildByJenkins(context.Background(), jenkins, serviceBuildMap.JobURL, record.Branch, buildServiceRecord.Env, codelibrary.URL)
			if err != nil {
				logger.Errorf("构建服务失败: %v", err)
				return
			}

			mu.Lock()
			record.BuildID = buildID
			record.BuildURL = buildURL
			record.Status = "Building"
			if err := repository.UpdateBuildServiceRecordsByID(record.ID, record); err != nil {
				logger.Errorf("更新构建服务记录失败: %v", err)
			}
			mu.Unlock()

			buildStatus, err := GetJenkinsBuildResult(context.Background(), jenkins, serviceBuildMap.JobURL, buildID)
			if err != nil {
				logger.Errorf("获取构建结果失败: %v", err)
			}
			image, err := GetJenkinsBuildImage(context.Background(), jenkins, serviceBuildMap.JobURL, buildID)
			if err != nil {
				logger.Errorf("获取构建镜像失败: %v", err)
			}

			mu.Lock()
			record.Image = image
			record.Status = buildStatus
			if buildStatus == "SUCCESS" {
				successCount++
			}
			if err := repository.UpdateBuildServiceRecordsByID(record.ID, record); err != nil {
				logger.Errorf("更新构建服务记录失败: %v", err)
			}
			mu.Unlock()
		}(buildServiceRecord)
	}

	wg.Wait()

	// 异步更新总状态
	go func() {
		if successCount == len(buildServiceRecords) {
			buildRecord.Status = "Success"
		} else {
			buildRecord.Status = "Failure"
		}
		if err := repository.UpdateBuildRecordsByID(buildRecord.ID, buildRecord); err != nil {
			logger.Errorf("最终更新构建记录失败: %v", err)
		}
	}()

	return buildRecord.Status, nil
}

func BuildByJenkins(ctx context.Context, jenkins *gojenkins.Jenkins, jobName, branch, env, url string) (int64, string, error) {
	// 触发构建
	queueID, err := jenkins.BuildJob(ctx, jobName, map[string]string{
		"GIT_BRANCH": branch,
		"GIT_URL":    url,
		"ENV_NAME":   env,
	})
	if err != nil {
		return 0, "", err
	}

	// 查询构建 ID（等待队列执行完成）
	var queueItem *gojenkins.Task
	for {
		queueItem, err = jenkins.GetQueueItem(ctx, queueID)
		if err != nil {
			return 0, "", err
		}
		if queueItem.Raw.Executable.Number != 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}
	// 查询构建 ID
	queueItem, _ = jenkins.GetQueueItem(ctx, queueID)
	buildID := queueItem.Raw.Executable.Number
	buildURL := queueItem.Raw.Executable.URL
	return buildID, buildURL, nil
}

func GetJenkinsBuildResult(ctx context.Context, jenkins *gojenkins.Jenkins, jobName string, buildID int64) (string, error) {
	// 持续轮询构建状态
	for {
		// 拉取构建状态并检查结果
		build, err := jenkins.GetBuild(ctx, jobName, buildID)
		if err != nil {
			return "UNKNOW", err
		}
		buildStatus := build.GetResult()
		switch buildStatus {
		case "SUCCESS":
			return buildStatus, nil
		case "FAILURE":
			return buildStatus, nil
		case "ABORTED":
			return buildStatus, nil
		default:
			time.Sleep(5 * time.Second)
		}
	}
}

func GetJenkinsBuildImage(ctx context.Context, jenkins *gojenkins.Jenkins, jobName string, buildID int64) (string, error) {
	build, err := jenkins.GetBuild(ctx, jobName, buildID)
	if err != nil {
		return "UNKNOW", err
	}
	imageTag := build.Raw.Description
	// 类型断言，将 imageTag 转换为 string
	if tag, ok := imageTag.(string); ok {
		//fmt.Println("IMAGE_TAG:", tag)
		return tag, nil
	} else {
		return "", fmt.Errorf("failed to assert imageTag to string")
	}
}
