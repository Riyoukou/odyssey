package cicd

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Riyoukou/odyssey/app/model"
	"github.com/Riyoukou/odyssey/app/repository"
	"github.com/Riyoukou/odyssey/pkg/response"
	"github.com/gin-gonic/gin"
)

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
	case "code_source":
		result, err = repository.FetchCodeSources()
	case "build_record":
		result, err = repository.FetchBuildRecordsByProjectName(c.Query("project"))
	case "deploy_record":
		result, err = repository.FetchDeployRecordsByProjectName(c.Query("project"))
	case "credential":
		result, err = repository.FetchCredentials()
	case "cicd_tool":
		result, err = repository.FetchCICDTools()
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
	case "code_source":
		result, err = repository.GetCodeSourceByName(c.Query("name"))
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
		var req model.ServiceTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.CreateService(req)
	case "code_library":
		var req model.CodeLibraryTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.CreateCodeLibrary(req)
	case "code_source":
		var req model.CodeSourceTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.CreateCodeSource(req)
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
	case "code_library":
		var req model.CodeLibraryTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.UpdateCodeLibraryByNameAndProject(req)
	case "code_source":
		var req model.CodeSourceTable
		if err = c.ShouldBind(&req); err != nil {
			break
		}
		err = repository.UpdateCodeSourceByName(req)
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
	}
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	response.Success(c, nil, fmt.Sprintf("%s deleted successfully", c.Param("type")))
}
