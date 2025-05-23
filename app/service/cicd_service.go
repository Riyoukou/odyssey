package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Riyoukou/odyssey/app/model"
	"github.com/Riyoukou/odyssey/app/repository"
	"github.com/Riyoukou/odyssey/app/utils"
	"github.com/Riyoukou/odyssey/pkg/logger"
	"github.com/bndr/gojenkins"
	kruiseclientset "github.com/openkruise/kruise-api/client/clientset/versioned"
	"gopkg.in/yaml.v3"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func ServiceCICDMap(cicdMap map[string][]model.DeployItem, action string, clusters []string, updateData model.DeployItem) []byte {
	// 定义空模板
	empty := model.DeployItem{}

	switch action {
	case "create":
		for _, name := range clusters {
			cicdMap[name] = []model.DeployItem{empty}
		}
	case "add":
		for _, name := range clusters {
			if _, exists := cicdMap[name]; !exists {
				cicdMap[name] = []model.DeployItem{empty}
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

func CICDSyncV2(deployServiceRecords []model.DeployServiceRecordTable) {
	go func() {
		for _, deployServiceRecord := range deployServiceRecords {
			serviceDetail, err := repository.GetServiceByNameAndProjectByEnv(deployServiceRecord.ServiceName, deployServiceRecord.ProjectName, deployServiceRecord.Env)
			if err != nil {
				fmt.Printf("❌ Get cluster failed: %v\n", err)
				return
			}

			var deployConfig model.ServiceDeployMap

			err = json.Unmarshal(serviceDetail.DeployMap, &deployConfig)
			if err != nil {
				panic(err)
			}

			argocdData, err := repository.GetCICDToolByName(deployConfig[deployServiceRecord.ClusterName][0].Release.CicdTool)
			if err != nil {
				logger.Errorf("获取代码库失败: %v", err)
			}

			argocdCredential, err := repository.GetCredentialByName(argocdData.CredentialName)
			if err != nil {
				logger.Errorf("获取jenkins密码失败: %v", err)
			}

			syncURL := fmt.Sprintf("%s/api/v1/applications/%s/sync", argocdData.URL, deployConfig[deployServiceRecord.ClusterName][0].Release.ArgoCDApplication)

			err = utils.ArgoCDSyncV2LiteV2(syncURL, argocdCredential.Data)
			if err != nil {
				fmt.Printf("❌ Sync service failed: %v\n", err)
				return
			}
			fmt.Printf("✅ Service %s sync success.\n", deployConfig[deployServiceRecord.ClusterName][0].Release.ArgoCDApplication)
			time.Sleep(1 * time.Second)
		}
	}()
}

var deployLock sync.Mutex

func StartDeployBykustomization(deployServiceRecords []model.DeployServiceRecordTable) {
	deployLock.Lock()
	defer deployLock.Unlock()
	serviceDetail, err := repository.GetServiceByNameAndProjectByEnv(deployServiceRecords[0].ServiceName, deployServiceRecords[0].ProjectName, deployServiceRecords[0].Env)
	if err != nil {
		fmt.Printf("❌ Get cluster failed: %v\n", err)
		return
	}

	var deployConfig model.ServiceDeployMap

	err = json.Unmarshal([]byte(serviceDetail.DeployMap), &deployConfig)
	if err != nil {
		panic(err)
	}

	gitData, err := repository.GetCodeLibraryByNameAndProject(deployConfig[deployServiceRecords[0].ClusterName][0].Yaml.GitOpsRepo, deployServiceRecords[0].ProjectName)
	if err != nil {
		logger.Errorf("获取代码库失败: %v", err)
	}

	codeSource, err := repository.GetCICDToolByName(gitData.CodeSourceName)
	if err != nil {
		logger.Errorf("获取代码源失败: %v", err)
	}

	gitCredential, err := repository.GetCredentialByName(codeSource.CredentialName)
	if err != nil {
		logger.Errorf("获取jenkins密码失败: %v", err)
	}

	// 克隆仓库
	repoURL := gitData.URL
	branch := "master"

	token := gitCredential.Data // 强烈建议从环境变量读取
	tmpDir := filepath.Join(os.TempDir(), fmt.Sprintf("gitops-%d", time.Now().Unix()))
	fmt.Println("✅ Cloning repo to", tmpDir)

	//结束后删除文件夹
	defer func() {
		os.RemoveAll(tmpDir)
		fmt.Println("✅ tmpDir deleted.")
	}()
	//clone仓库
	repo := utils.GitClone(repoURL, branch, tmpDir, token)
	//更新kustomization.yaml
	for _, deployServiceRecord := range deployServiceRecords {
		parts := strings.Split(deployServiceRecord.Image, ":")
		if len(parts) != 2 {
			fmt.Printf("invalid image format: %s", deployServiceRecord.Image)
		}
		dockerRegistry := parts[0]
		imageTag := parts[1]

		serviceDetail, err := repository.GetServiceByNameAndProjectByEnv(deployServiceRecord.ServiceName, deployServiceRecord.ProjectName, deployServiceRecord.Env)
		if err != nil {
			fmt.Printf("❌ Get cluster failed: %v\n", err)
			return
		}

		var deployConfig model.ServiceDeployMap

		err = json.Unmarshal([]byte(serviceDetail.DeployMap), &deployConfig)
		if err != nil {
			panic(err)
		}

		err = updateKustomizationImageTag(filepath.Join(tmpDir, fmt.Sprintf("%s/kustomization.yaml", deployConfig[deployServiceRecord.ClusterName][0].Yaml.FilePath)), dockerRegistry, imageTag)
		if err != nil {
			fmt.Printf("❌ Update kustomization.yaml failed: %v\n", err)
			return
		}

		time.Sleep(1 * time.Second)
	}

	utils.GitCommit(repo, ".")

	utils.GitPush(repo, tmpDir, repoURL, branch, token)
}

// 更新 kustomization.yaml 中镜像 tag（保留注释）
func updateKustomizationImageTag(path, imageName, newTag string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return err
	}

	if root.Kind != yaml.DocumentNode || len(root.Content) == 0 {
		return fmt.Errorf("invalid YAML")
	}
	doc := root.Content[0]

	var imagesNode *yaml.Node
	for i := 0; i < len(doc.Content); i += 2 {
		key := doc.Content[i]
		if key.Value == "images" {
			imagesNode = doc.Content[i+1]
			break
		}
	}
	if imagesNode == nil {
		return fmt.Errorf("images section not found")
	}

	updated := false
	for _, imageItem := range imagesNode.Content {
		var nameNode, tagNode *yaml.Node
		for i := 0; i < len(imageItem.Content); i += 2 {
			key := imageItem.Content[i]
			val := imageItem.Content[i+1]
			if key.Value == "name" && val.Value == imageName {
				nameNode = val
			}
			if key.Value == "newTag" {
				tagNode = val
			}
		}
		if nameNode != nil {
			if tagNode != nil {
				tagNode.Value = newTag
			} else {
				imageItem.Content = append(imageItem.Content,
					&yaml.Node{Kind: yaml.ScalarNode, Value: "newTag"},
					&yaml.Node{Kind: yaml.ScalarNode, Value: newTag},
				)
			}
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("image %s not found", imageName)
	}

	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()

	encoder := yaml.NewEncoder(outFile)
	encoder.SetIndent(2)
	return encoder.Encode(&root)
}

func GetWorkloadStatus(clusterName, workloadType, namespace, serviceName string) string {
	cluster, err := repository.GetClusterByName(clusterName)
	if err != nil {
		fmt.Printf("❌ Get cluster failed: %v\n", err)
		return "Failed"
	}

	credential, err := repository.GetCredentialByName(cluster.Config)
	if err != nil {
		log.Fatalf("❌ Get credential failed: %v\n", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(string(credential.Data)))
	if err != nil {
		log.Fatalf("Failed to create REST config: %v", err)
	}

	clientset := utils.CreateKubernetesClientset(credential.Data)
	kruiseClient := kruiseclientset.NewForConfigOrDie(config)

	for {
		switch workloadType {
		case "StatefulSet":
			// 官方 StatefulSet
			sts, err := clientset.AppsV1().StatefulSets(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
			if err == nil {
				if sts.Status.ReadyReplicas == *sts.Spec.Replicas {
					return "Running"
				}
			}

			// Kruise StatefulSet
			if errors.IsNotFound(err) {
				ksts, kerr := kruiseClient.AppsV1alpha1().StatefulSets(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
				if kerr == nil {
					if ksts.Status.ReadyReplicas == *ksts.Spec.Replicas {
						return "Running"
					}
				}
			}
		case "DaemonSet":
			ds, err := clientset.AppsV1().DaemonSets(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
			if err == nil {
				if ds.Status.NumberReady == ds.Status.DesiredNumberScheduled {
					return "Running"
				}
			}
			// Kruise DaemonSet
			if errors.IsNotFound(err) {
				kds, kerr := kruiseClient.AppsV1alpha1().DaemonSets(namespace).Get(context.Background(), serviceName, metav1.GetOptions{})
				if kerr == nil {
					if kds.Status.NumberReady == kds.Status.DesiredNumberScheduled {
						return "Running"
					}
				}
			}
		default:
			status := GetKruiseStatus(clusterName, serviceName, namespace)
			return status
		}
		time.Sleep(10 * time.Second)
	}
}

func GetKruiseStatus(clusterName, service, namespace string) string {
	clusterDetail, err := repository.GetClusterByName(clusterName)
	if err != nil {
		logger.Errorf("获取集群失败: %v", err)
		return "Failed"
	}

	credential, err := repository.GetCredentialByName(clusterDetail.Config)
	if err != nil {
		log.Fatalf("❌ Get credential failed: %v\n", err)
	}

	kruiseClientset := utils.CreateKruiseClientset(credential.Data)
	maxRetryCount := 120 // 最大重试次数，可以根据需要调整
	retryCount := 0
	time.Sleep(5 * time.Second)
	for {
		for retryCount < maxRetryCount {
			// 获取 Rollout 状态
			kruiseClient, err := kruiseClientset.RolloutsV1beta1().Rollouts(namespace).Get(context.Background(), service, metav1.GetOptions{})
			if err != nil {
				logger.Errorf("Failed to get rollout status: %v", err)
				return "Failed"
			}

			// 获取步骤信息
			var (
				currentStepIndex int32
				currentStepState string
			)
			maxSteps := int32(len(kruiseClient.Spec.Strategy.Canary.Steps))
			if kruiseClient.Status.CanaryStatus == nil {
				currentStepIndex = maxSteps
				currentStepState = "Completed"
			} else {
				currentStepIndex = kruiseClient.Status.CanaryStatus.CommonStatus.CurrentStepIndex
				currentStepState = string(kruiseClient.Status.CanaryStatus.CommonStatus.CurrentStepState)
			}

			// StepPaused 状态
			if currentStepState == "StepPaused" {
				return fmt.Sprintf("Pending(%d/%d)", currentStepIndex, maxSteps)
			}

			// 完成状态
			if currentStepState == "Completed" && currentStepIndex == maxSteps {
				return fmt.Sprintf("Running(%d/%d)", currentStepIndex, maxSteps)
			}

			retryCount++
			log.Printf("Service %s Waiting for next step, current state: %s, step: %d/%d", service, currentStepState, currentStepIndex, maxSteps)

			// 如果没有完成，可以加入延迟再尝试
			time.Sleep(10 * time.Second) // 设置等待时间，避免频繁轮询
		}
	}
}

func ApproveRolloutKruise(clusterName, namespace, appName string) {

	cluster, err := repository.GetClusterByName(clusterName)
	if err != nil {
		fmt.Printf("❌ Get cluster failed: %v\n", err)
		return
	}

	credential, err := repository.GetCredentialByName(cluster.Config)
	if err != nil {
		log.Fatalf("❌ Get credential failed: %v\n", err)
	}

	config, err := clientcmd.RESTConfigFromKubeConfig([]byte(credential.Data))
	if err != nil {
		logger.Errorf("Failed to create REST config: %v", err)
	}

	dynamicClientset, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("创建 Kubernetes DynamicClient 失败: %v", err)
	}

	rollout := utils.RolloutGetOptions(dynamicClientset, namespace, appName)
	// 断言获取到的 Rollout 状态为 map[string]interface{}
	status, ok := rollout.Object["status"].(map[string]interface{})
	if !ok {
		logger.Error("无法获取 Rollout 状态")
		return
	}

	// 获取 canaryStatus
	canaryStatus, ok := status["canaryStatus"].(map[string]interface{})
	if !ok {
		logger.Error("无法获取 canaryStatus")
		return
	}

	if canaryStatus["currentStepState"].(string) != "StepPaused" {
		fmt.Printf("%s Rollout 未暂停，无法批准！当前状态为: %s\n", appName, canaryStatus["currentStepState"].(string))
		return
	}

	canaryStatus["currentStepState"] = "StepReady"

	// 更新 Rollout 资源
	err = utils.RolloutUpdateOptions(dynamicClientset, namespace, appName, rollout)
	if err != nil {
		logger.Errorf("Rollout 更新失败: %v", err)
		return
	}

	fmt.Printf("%s Rollout 更新成功！\n", appName)
}
