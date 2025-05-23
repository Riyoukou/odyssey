package cicd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Riyoukou/odyssey/app/model"
	"github.com/Riyoukou/odyssey/app/repository"
	"github.com/Riyoukou/odyssey/app/service"
	"github.com/Riyoukou/odyssey/app/utils"
	"github.com/Riyoukou/odyssey/pkg/logger"
	"github.com/Riyoukou/odyssey/pkg/response"
	"github.com/gin-gonic/gin"
)

func HandleServiceCICDMap(c *gin.Context) {
	var (
		updateData model.DeployItem
		deployMap  map[string][]model.DeployItem
	)
	intID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	result, err := repository.GetServiceByID(intID)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	err = json.Unmarshal(result.DeployMap, &deployMap)
	if err != nil {
		return
	}

	result.DeployMap = service.ServiceCICDMap(deployMap, c.Query("action"), strings.Split(c.Query("clusters"), ","), updateData)

	err = repository.UpdateServiceByNameAndProjectByEnv(result)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	response.Success(c, nil, "OK")

}

func HandleCICDBuildByJenkins(c *gin.Context) {
	id := c.Param("id")
	buildRecordID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		logger.Errorf("解析构建记录 ID 失败: %v", err)
	}
	buildRecord, err := repository.GetBuildRecordByID(buildRecordID)
	if err != nil {
		logger.Errorf("获取构建记录失败: %v", err)
	}

	if buildRecord.Status != "Pending" {
		response.Success(c, http.StatusBadRequest, "ERROR STATUS")
		return
	}

	buildRecord.Status = "Building" // 更新数据库
	err = repository.UpdateBuildRecordsByID(buildRecordID, *buildRecord)
	if err != nil {
		logger.Errorf("更新构建记录失败: %v", err)
	}

	go func() {
		_, err := service.CICDBuildByJenkins(*buildRecord)
		if err != nil {
			logger.Errorf("构建失败: %v", err)
			return
		}
	}()

	response.Success(c, nil, "BUILD BY JENKINS SUCCESS")

}

func HandleCICDFetch(c *gin.Context) {
	var (
		err    error
		result interface{}
	)
	switch c.Param("type") {
	case "cluster":
		result, err = repository.FetchClusters()
	case "project":
		result, err = repository.FetchProjects()
	case "env":
		result, err = repository.FetchEnvsByProject(c.Query("project"))
	case "service":
		result, err = repository.FetchServicesByProjectAndEnv(c.Query("project"), c.Query("env"))
	case "code_library":
		result, err = repository.FetchCodeLibraries()
	case "build_record":
		result, err = repository.FetchBuildRecordsByProjectName(c.Query("project"))
	case "deploy_record":
		result, err = repository.FetchDeployRecordsByProjectName(c.Query("project"))
	case "credential":
		result, err = repository.FetchCredentials()
	case "cicd_tool":
		result, err = repository.FetchCICDTools()
	case "git_project":
		cicdTool, _ := repository.GetCICDToolByName(c.Query("cicd_tool"))
		credential, _ := repository.GetCredentialByName(cicdTool.CredentialName)
		result = utils.GetGitlabProjects(cicdTool.URL+"/api/v4/projects?simple=true&per_page=100", credential.Data)
	case "gitlab_branch":
		codeLibrary, _ := repository.GetCodeLibraryByNameAndProject(c.Query("code_library"), c.Query("project"))
		cicdTool, _ := repository.GetCICDToolByName(codeLibrary.CodeSourceName)
		credential, _ := repository.GetCredentialByName(cicdTool.CredentialName)
		result = utils.GitGetBranches(codeLibrary.URL, credential.Data)
	case "gitlab_tag":
		codeLibrary, _ := repository.GetCodeLibraryByNameAndProject(c.Query("code_library"), c.Query("project"))
		cicdTool, _ := repository.GetCICDToolByName(codeLibrary.CodeSourceName)
		credential, _ := repository.GetCredentialByName(cicdTool.CredentialName)
		result = utils.GitGetTags(codeLibrary.URL, credential.Data)
	case "build_service_record":
		result, err = repository.GetBuildServiceRecordsByBuildRecordName(c.Query("build_record"))
	case "deploy_service_record":
		result, err = repository.GetDeployServiceRecordsByDeployRecordName(c.Query("deploy_record"))
	}
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, result, fmt.Sprintf("%s fetched successfully", c.Param("type")))
}

func HandleCICDGet(c *gin.Context) {
	var (
		err    error
		result interface{}
	)
	switch c.Param("type") {
	case "cluster":
		result, err = repository.GetClusterByName(c.Query("name"))
	case "project":
		result, err = repository.GetProjectByName(c.Query("name"))
	case "env":
		result, err = repository.GetEnvByNameAndProject(c.Query("name"), c.Query("project"))
	case "service":
		result, err = repository.GetServiceByNameAndProjectByEnv(c.Query("name"), c.Query("project"), c.Query("env"))
	case "code_library":
		result, err = repository.GetCodeLibraryByNameAndProject(c.Query("name"), c.Query("project"))
	case "build_service_record":
		result, err = repository.GetBuildServiceRecordsByBuildRecordName(c.Query("build_record"))
	case "deploy_service_record":
		result, err = repository.GetDeployServiceRecordsByDeployRecordName(c.Query("deploy_record"))
	case "credential":
		result, err = repository.GetCredentialByName(c.Query("name"))
	}
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, result, fmt.Sprintf("%s fetched successfully", c.Param("type")))
}

func HandleCICDCreate(c *gin.Context) {
	var err error
	switch c.Param("type") {
	case "cluster":
		var req model.ClusterTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.CreateCluster(req)
	case "project":
		var req model.ProjectTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.CreateProject(req)
	case "env":
		var (
			req model.EnvTable
			//project model.ProjectTable
			//envs    []string
		)
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		if err = repository.CreateEnv(req); err != nil {
			break
		}
		/*project, err = repository.GetProjectByName(req.ProjectName)
		if err != nil {
			break
		}
		if len(project.Env) > 0 && project.Env != nil {
			err = json.Unmarshal(project.Env, &envs)
			if err != nil {
				return
			}
		}
		// 追加新环境
		envs = append(envs, req.Name)
		// 写回数据库字段
		project.Env, err = json.Marshal(envs)
		if err != nil {
			return
		}
		err = repository.UpdateProject(project)
		if err != nil {
			break
		}*/
	case "service":
		var (
			req      model.ServiceTable
			clusters []string
		)
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		_ = json.Unmarshal(req.Clusters, &clusters)
		deployMap := service.ServiceCICDMap(
			map[string][]model.DeployItem{}, "create",
			clusters,
			model.DeployItem{},
		)

		jsonBuildMap, err := json.MarshalIndent(model.ServiceBuildMap{}, "", "  ")
		if err != nil {
			break
		}
		req.BuildMap = jsonBuildMap
		req.DeployMap = deployMap

		err = repository.CreateService(req)
		if err != nil {
			break
		}
	case "code_library":
		var req model.CodeLibraryTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.CreateCodeLibrary(req)
	case "credential":
		var req model.CredentialTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.CreateCredential(req)
	case "cicd_tool":
		var req model.CICDToolTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.CreateCICDTool(req)
	case "build_record":
		var req model.ApiBuildRecord
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = service.CreateBuildRecord(req)
	case "build_service_record":
		var req model.BuildServiceRecordTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.CreateBuildServiceRecord(req)
	case "deploy_record":
		var req model.DeployRecordTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = service.CreateDeployRecord(req)
	}
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, nil, fmt.Sprintf("%s created successfully", c.Param("type")))
}

func HandleCICDUpdate(c *gin.Context) {
	var err error
	switch c.Param("type") {
	case "cluster":
		var req model.ClusterTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.UpdateCluster(req)
	case "project":
		var req model.ProjectTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.UpdateProject(req)
	case "env":
		var req model.EnvTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.UpdateEnvByNameAndProject(req)
	case "service":
		var req model.ServiceTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.UpdateServiceByNameAndProjectByEnv(req)
	case "service_build_map":
		intID, err := strconv.ParseInt(c.Query("id"), 10, 64)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err)
			return
		}

		var req model.ServiceBuildMap
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		jsonStr, err := json.Marshal(req)
		if err != nil {
			break
		}
		err = repository.UpdateServiceBuildMap(intID, jsonStr)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err)
			return
		}
	case "code_library":
		var req model.CodeLibraryTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.UpdateCodeLibraryByNameAndProject(req)

	case "cicd_tool":
		var req model.CICDToolTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.UpdateCICDTool(req)
	}
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, nil, fmt.Sprintf("%s updated successfully", c.Param("type")))
}

func HandleCICDDelete(c *gin.Context) {
	var (
		err   error
		intID int64
	)
	intID, err = strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	switch c.Param("type") {
	case "cluster":
		err = repository.DeleteCluster(intID)
	case "project":
		err = repository.DeleteProject(intID)
	case "service":
		err = repository.DeleteService(intID)
	case "code_library":
		err = repository.DeleteCodeLibrary(intID)
	case "credential":
		err = repository.DeleteCredential(intID)
	case "cicd_tool":
		err = repository.DeleteCICDTool(intID)
	case "build_record":
		err = repository.DeleteBuildRecord(intID)
	case "deploy_record":
		err = repository.DeteleDeployRecord(intID)
	}
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, nil, fmt.Sprintf("%s deleted successfully", c.Param("type")))
}

func HandleStartDeploy(c *gin.Context) {
	var (
		wg  sync.WaitGroup
		mu  sync.Mutex
		req []model.DeployServiceRecordTable
	)

	// 绑定请求数据
	err := c.ShouldBind(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	hasPending := false
	for _, item := range req {
		if item.Status == "Pending" {
			hasPending = true
			break
		}

	}

	if hasPending {
		// 启动并发处理
		for _, deployServiceRecord := range req {
			//防止api超时降速
			time.Sleep(1 * time.Second)
			if deployServiceRecord.Status != "Pending" {
				continue
			}
			wg.Add(1)
			go func() {
				// 同步组减一
				defer wg.Done()
				mu.Lock()
				deployServiceRecord.Status = "Deploying"
				err = repository.UpdateDeployServiceRecordsByID(deployServiceRecord.ID, deployServiceRecord)
				if err != nil {
					logger.Errorf("更新发布服务记录失败: %v", err)
				}
				mu.Unlock()
			}()
		}
		response.Success(c, nil, "START DEPLOY SUCCESS")
		// 执行你的方法
		service.StartDeployBykustomization(req)
		service.CICDSyncV2(req)
	}

	// 启动并发处理
	for _, deployServiceRecord := range req {
		//防止api超时降速
		time.Sleep(1 * time.Second)
		if deployServiceRecord.Status != "Deploying" {
			continue
		}
		wg.Add(1)
		go func() {
			// 同步组减一
			defer wg.Done()
			mu.Lock()

			serviceDetail, err := repository.GetServiceByNameAndProjectByEnv(deployServiceRecord.ServiceName, deployServiceRecord.ProjectName, deployServiceRecord.Env)
			if err != nil {
				logger.Errorf("获取发布服务详情失败: %v", err)
			}

			envDetail, err := repository.GetEnvByNameAndProject(deployServiceRecord.Env, deployServiceRecord.ProjectName)
			if err != nil {
				logger.Errorf("获取环境详情失败: %v", err)
			}
			var deployConfig model.ServiceDeployMap

			err = json.Unmarshal(serviceDetail.DeployMap, &deployConfig)
			if err != nil {
				panic(err)
			}

			status := service.GetWorkloadStatus(deployServiceRecord.ClusterName, deployConfig[deployServiceRecord.ClusterName][0].Release.Workload, envDetail.Namespace, deployServiceRecord.ServiceName)

			deployServiceRecord.Status = status

			mu.Lock()
			err = repository.UpdateDeployServiceRecordsByID(deployServiceRecord.ID, deployServiceRecord)
			if err != nil {
				logger.Errorf("更新发布服务记录失败: %v", err)
			}
			mu.Unlock()
		}()
	}

}

func HandleApproveDeploy(c *gin.Context) {
	var (
		wg  sync.WaitGroup
		mu  sync.Mutex
		req []model.DeployServiceRecordTable
	)

	// 绑定请求数据
	err := c.ShouldBind(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	response.Success(c, nil, "APPROVE DEPLOY SUCCESS")

	// 启动并发处理
	for _, deployServiceRecord := range req {
		if !strings.Contains(deployServiceRecord.Status, "Pending") {
			continue
		}
		//防止api超时降速
		time.Sleep(1 * time.Second)
		wg.Add(1)
		go func() {
			// 同步组减一
			defer wg.Done()
			mu.Lock()
			deployServiceRecord.Status = "Deploying"
			err = repository.UpdateDeployServiceRecordsByID(deployServiceRecord.ID, deployServiceRecord)
			if err != nil {
				logger.Errorf("更新发布服务记录失败: %v", err)
			}
			mu.Unlock()

			serviceDetail, err := repository.GetServiceByNameAndProjectByEnv(deployServiceRecord.ServiceName, deployServiceRecord.ProjectName, deployServiceRecord.Env)
			if err != nil {
				logger.Errorf("获取发布服务详情失败: %v", err)
			}

			var deployConfig model.ServiceDeployMap

			err = json.Unmarshal(serviceDetail.DeployMap, &deployConfig)
			if err != nil {
				panic(err)
			}

			envDetail, err := repository.GetEnvByNameAndProject(deployServiceRecord.Env, deployServiceRecord.ProjectName)
			if err != nil {
				logger.Errorf("获取环境详情失败: %v", err)
			}

			if deployConfig[deployServiceRecord.ClusterName][0].Release.Workload != "StatefulSet" && deployConfig[deployServiceRecord.ClusterName][0].Release.Workload != "DaemonSet" {
				service.ApproveRolloutKruise(deployServiceRecord.ClusterName, envDetail.Namespace, deployServiceRecord.ServiceName)
			}

			status := service.GetWorkloadStatus(deployServiceRecord.ClusterName, deployConfig[deployServiceRecord.ClusterName][0].Release.Workload, envDetail.Namespace, deployServiceRecord.ServiceName)

			deployServiceRecord.Status = status

			mu.Lock()
			err = repository.UpdateDeployServiceRecordsByID(deployServiceRecord.ID, deployServiceRecord)
			if err != nil {
				logger.Errorf("更新发布服务记录失败: %v", err)
			}
			mu.Unlock()
		}()
	}
}
