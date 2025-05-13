package repository

import (
	"encoding/base64"
	"fmt"

	"github.com/Riyoukou/odyssey/app/model"
	"github.com/Riyoukou/odyssey/pkg/logger"
)

// cicd_cluster
func FetchClusters() ([]model.ClusterTable, error) {
	var clusters []model.ClusterTable
	if err := DB.Find(&clusters).Error; err != nil {
		logger.Errorf("Failed to fetch clusters: %v", err)
		return nil, err
	}

	return clusters, nil
}

func CreateCluster(cluster model.ClusterTable) error {
	if err := DB.Where("name = ? AND api_server = ?", cluster.Name, cluster.APIServer).
		First(&model.ClusterTable{}).Error; err == nil {
		logger.Errorf("Cluster already exists: name=%s, api_server=%s", cluster.Name, cluster.APIServer)
		return err
	}

	if err := DB.Create(&cluster).Error; err != nil {
		logger.Errorf("Failed to create cluster: %v", err)
		return err
	}

	return nil
}

func DeleteCluster(clusterID int64) error {
	if err := DB.Delete(&model.ClusterTable{}, clusterID).Error; err != nil {
		logger.Errorf("Failed to delete cluster: %v", err)
		return err
	}

	return nil
}

func GetClusterByName(name string) (model.ClusterTable, error) {
	var cluster model.ClusterTable
	if err := DB.Where("name = ?", name).
		First(&cluster).Error; err != nil {
		return model.ClusterTable{}, err
	}
	return cluster, nil
}

func UpdateCluster(cluster model.ClusterTable) error {
	if err := DB.Model(&model.ClusterTable{}).Where("id = ?", cluster.ID).Updates(cluster).Error; err != nil {
		logger.Errorf("Failed to update cluster: %v", err)
		return err
	}

	return nil
}

// cicd_project
func FetchProjects() ([]model.ProjectTable, error) {
	var projects []model.ProjectTable
	if err := DB.Find(&projects).Error; err != nil {
		logger.Errorf("Failed to fetch projects: %v", err)
		return nil, err
	}

	return projects, nil
}

func CreateProject(project model.ProjectTable) error {
	if err := DB.Where("name = ?", project.Name).
		First(&model.ProjectTable{}).Error; err == nil {
		logger.Errorf("Project already exists: name=%s", project.Name)
		return err
	}

	if err := DB.Create(&project).Error; err != nil {
		logger.Errorf("Failed to create project: %v", err)
		return err
	}

	return nil
}

func DeleteProject(projectID int64) error {
	if err := DB.Delete(&model.ProjectTable{}, projectID).Error; err != nil {
		logger.Errorf("Failed to delete project: %v", err)
		return err
	}

	return nil
}

func GetProjectByName(name string) (model.ProjectTable, error) {
	var project model.ProjectTable
	if err := DB.Where("name = ?", name).
		First(&project).Error; err != nil {
		return model.ProjectTable{}, err
	}

	return project, nil
}

func UpdateProject(project model.ProjectTable) error {
	if err := DB.Model(&model.ProjectTable{}).Where("id = ?", project.ID).Updates(project).Error; err != nil {
		logger.Errorf("Failed to update project: %v", err)
		return err
	}

	return nil
}

// cicd_env
func FetchEnvs() ([]model.EnvTable, error) {
	var envs []model.EnvTable
	if err := DB.Find(&envs).Error; err != nil {
		logger.Errorf("Failed to fetch envs: %v", err)
		return nil, err
	}

	return envs, nil
}

func FetchEnvsByProject(projectName string) ([]model.EnvTable, error) {
	var envs []model.EnvTable
	if err := DB.Find(&envs, "project_name = ?", projectName).Error; err != nil {
		logger.Errorf("Failed to fetch envs: %v", err)
		return nil, err
	}

	return envs, nil
}

func CreateEnv(env model.EnvTable) error {
	if err := DB.Where("name = ? AND project_name = ?", env.Name, env.ProjectName).
		First(&model.EnvTable{}).Error; err == nil {
		logger.Errorf("Env already exists: name=%s project_name=%s", env.Name, env.ProjectName)
		return err
	}

	if err := DB.Create(&env).Error; err != nil {
		logger.Errorf("Failed to create env: %v", err)
		return err
	}

	return nil
}

func DeleteEnv(envID int64) error {
	if err := DB.Delete(&model.EnvTable{}, envID).Error; err != nil {
		logger.Errorf("Failed to delete env: %v", err)
		return err
	}

	return nil
}

func GetEnvByNameAndProject(name, projectName string) (model.EnvTable, error) {
	var env model.EnvTable
	if err := DB.Where("name = ? AND project_name = ?", name, projectName).
		First(&env).Error; err != nil {
		return model.EnvTable{}, err
	}

	return env, nil
}

func UpdateEnvByNameAndProject(env model.EnvTable) error {
	if err := DB.Model(&model.EnvTable{}).Where("name = ? AND project_name = ?", env.Name, env.ProjectName).Updates(env).Error; err != nil {
		logger.Errorf("Failed to update env: %v", err)
		return err
	}
	return nil
}

// cicd_service
func FetchServicesByProjectAndEnv(projectName, envName string) ([]model.ServiceTable, error) {
	var services []model.ServiceTable
	if err := DB.Find(&services, "project_name = ? AND env_name = ?", projectName, envName).Error; err != nil {
		logger.Errorf("Failed to fetch services: %v", err)
		return nil, err
	}

	return services, nil
}

func CreateService(service model.ServiceTable) error {
	if err := DB.Where("name = ? AND project_name = ? ", service.Name, service.ProjectName).
		First(&model.ServiceTable{}).Error; err == nil {
		logger.Errorf("Service already exists: name=%s project_name=%s", service.Name, service.ProjectName)
		return err
	}

	if err := DB.Create(&service).Error; err != nil {
		logger.Errorf("Failed to create service: %v", err)
		return err
	}

	return nil
}

func DeleteService(serviceID int64) error {
	if err := DB.Delete(&model.ServiceTable{}, serviceID).Error; err != nil {
		logger.Errorf("Failed to delete service: %v", err)
		return err
	}

	return nil
}

func GetServiceByID(id int64) (model.ServiceTable, error) {
	var service model.ServiceTable
	if err := DB.Where("id = ?", id).
		First(&service).Error; err != nil {
		return model.ServiceTable{}, err
	}

	return service, nil
}

func GetServiceByNameAndProjectByEnv(name, projectName, envName string) (model.ServiceTable, error) {
	var service model.ServiceTable
	if err := DB.Where("name = ? AND project_name = ? AND env_name = ?", name, projectName, envName).
		First(&service).Error; err != nil {
		return model.ServiceTable{}, err
	}

	return service, nil
}

func UpdateServiceByNameAndProjectByEnv(service model.ServiceTable) error {
	if err := DB.Model(&model.ServiceTable{}).Where("name = ? AND project_name = ? AND env_name = ?", service.Name, service.ProjectName, service.EnvName).Updates(service).Error; err != nil {
		logger.Errorf("Failed to update service: %v", err)
		return err
	}

	return nil
}

func UpdateServiceBuildMap(id int64, buildMap []byte) error {
	if err := DB.Model(&model.ServiceTable{}).Where("id = ?", id).Update("build_map", buildMap).Error; err != nil {
		logger.Errorf("Failed to update service build map: %v", err)
		return err
	}

	return nil
}

// code_library
func FetchCodeLibraries() ([]model.CodeLibraryTable, error) {
	var codeLibraries []model.CodeLibraryTable
	if err := DB.Find(&codeLibraries).Error; err != nil {
		logger.Errorf("Failed to fetch code libraries: %v", err)
		return nil, err
	}

	return codeLibraries, nil
}

func CreateCodeLibrary(codeLibrary model.CodeLibraryTable) error {
	if err := DB.Where("name = ? AND project_name = ?", codeLibrary.Name, codeLibrary.ProjectName).
		First(&model.CodeLibraryTable{}).Error; err == nil {
		logger.Errorf("Code library already exists: name=%s project_name=%s", codeLibrary.Name, codeLibrary.ProjectName)
		return err
	}
	if err := DB.Create(&codeLibrary).Error; err != nil {
		logger.Errorf("Failed to create code library: %v", err)
		return err
	}

	return nil
}

func DeleteCodeLibrary(codeLibraryID int64) error {
	if err := DB.Delete(&model.CodeLibraryTable{}, codeLibraryID).Error; err != nil {
		logger.Errorf("Failed to delete code library: %v", err)
		return err
	}

	return nil
}

func GetCodeLibraryByNameAndProject(name, projectName string) (model.CodeLibraryTable, error) {
	var codeLibrary model.CodeLibraryTable
	if err := DB.Where("name = ? AND project_name = ?", name, projectName).
		First(&codeLibrary).Error; err != nil {
		return model.CodeLibraryTable{}, err
	}

	return codeLibrary, nil
}

func UpdateCodeLibraryByNameAndProject(codeLibrary model.CodeLibraryTable) error {
	if err := DB.Model(&model.CodeLibraryTable{}).Where("name = ? AND project_name = ?", codeLibrary.Name, codeLibrary.ProjectName).Updates(codeLibrary).Error; err != nil {
		logger.Errorf("Failed to update code library: %v", err)
		return err
	}

	return nil
}

// cicd_build_record
func FetchBuildRecordsByProjectName(projectName string) ([]model.BuildRecordTable, error) {
	var records []model.BuildRecordTable
	if err := DB.Where("project_name = ?", projectName).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func CreateBuildRecord(buildRecord model.BuildRecordTable) error {
	if err := DB.Create(&buildRecord).Error; err != nil {
		return fmt.Errorf("failed to create build record: %w", err)
	}
	return nil
}

func GetBuildRecordByID(id int64) (*model.BuildRecordTable, error) {
	var record model.BuildRecordTable
	if err := DB.Where("id = ?", id).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func UpdateBuildRecordsByID(id int64, buildRecord model.BuildRecordTable) error {
	if err := DB.Model(&model.BuildRecordTable{}).Where("id = ?", id).Updates(buildRecord).Error; err != nil {
		return fmt.Errorf("failed to update build record: %w", err)
	}

	return nil
}

func GetBuildRecordByName(name string) (*model.BuildRecordTable, error) {
	var record model.BuildRecordTable
	if err := DB.Where("name = ?", name).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func DeleteBuildRecord(id int64) error {
	if err := DB.Delete(&model.BuildRecordTable{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete build record: %w", err)
	}

	return nil
}

func UpdateBuildRecord(buildRecord model.BuildRecordTable) error {
	if err := DB.Model(&model.BuildRecordTable{}).Where("id = ?", buildRecord.ID).Updates(buildRecord).Error; err != nil {
		return fmt.Errorf("failed to update build record: %w", err)
	}

	return nil
}

// cicd_build_service_record
func GetBuildServiceRecordsByBuildRecordName(buildRecordName string) ([]model.BuildServiceRecordTable, error) {
	var records []model.BuildServiceRecordTable
	if err := DB.Where("build_record_name = ?", buildRecordName).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func CreateBuildServiceRecord(buildServiceRecord model.BuildServiceRecordTable) error {
	if err := DB.Create(&buildServiceRecord).Error; err != nil {
		return fmt.Errorf("failed to create build service record: %w", err)
	}
	return nil
}

func UpdateBuildServiceRecord(buildServiceRecord model.BuildServiceRecordTable) error {
	if err := DB.Model(&model.BuildServiceRecordTable{}).Where("id = ?", buildServiceRecord.ID).Updates(buildServiceRecord).Error; err != nil {
		return fmt.Errorf("failed to update build service record: %w", err)
	}
	return nil
}

func UpdateBuildServiceRecordsByID(id int64, buildServiceRecord model.BuildServiceRecordTable) error {
	if err := DB.Model(&model.BuildServiceRecordTable{}).Where("id = ?", id).Updates(buildServiceRecord).Error; err != nil {
		return fmt.Errorf("failed to update build service record: %w", err)
	}

	return nil
}

// cicd_deploy_record
func FetchDeployRecordsByProjectName(projectName string) ([]model.DeployRecordTable, error) {
	var records []model.DeployRecordTable
	if err := DB.Where("project_name = ?", projectName).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func CreateDeployRecord(deployRecord model.DeployRecordTable) error {
	if err := DB.Create(&deployRecord).Error; err != nil {
		return fmt.Errorf("failed to create deploy record: %w", err)
	}
	return nil
}

func GetDeployRecordByID(id int64) (*model.DeployRecordTable, error) {
	var record model.DeployRecordTable
	if err := DB.Where("id = ?", id).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func GetDeployRecordByName(name string) (*model.DeployRecordTable, error) {
	var record model.DeployRecordTable
	if err := DB.Where("name = ?", name).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

func GetDeployRecordByBuildRecordName(name string) (*model.DeployRecordTable, error) {
	var record model.DeployRecordTable
	if err := DB.Where("build_record_name = ?", name).First(&record).Error; err != nil {
		return nil, err
	}

	return &record, nil
}

// cicd_deploy_service_record
func CreateDeployServiceRecord(deployServiceRecord model.DeployServiceRecordTable) error {
	if err := DB.Create(&deployServiceRecord).Error; err != nil {
		return fmt.Errorf("failed to create deploy service record: %w", err)
	}
	return nil
}

func GetDeployServiceRecordsByDeployRecordName(deployRecordName string) ([]model.DeployServiceRecordTable, error) {
	var records []model.DeployServiceRecordTable
	if err := DB.Where("deploy_record_name = ?", deployRecordName).Find(&records).Error; err != nil {
		return nil, err
	}

	return records, nil
}

func UpdateDeployServiceRecordsByID(id int64, deployServiceRecord model.DeployServiceRecordTable) error {
	if err := DB.Model(&model.DeployServiceRecordTable{}).Where("id = ?", id).Updates(deployServiceRecord).Error; err != nil {
		return fmt.Errorf("failed to update deploy service record: %w", err)
	}

	return nil
}

// cicd_crendential
func FetchCredentials() ([]model.CredentialTable, error) {
	var credentials []model.CredentialTable
	if err := DB.Find(&credentials).Error; err != nil {
		return nil, err
	}

	return credentials, nil
}

func CreateCredential(credential model.CredentialTable) error {
	if err := DB.Where("name = ?", credential.Name).
		First(&model.CredentialTable{}).Error; err == nil {
		logger.Errorf("Credential already exists: name=%s", credential.Name)
		return err
	}

	credential.Data = base64.StdEncoding.EncodeToString([]byte(credential.Data))

	if err := DB.Create(&credential).Error; err != nil {
		logger.Errorf("Failed to create credential: %v", err)
		return err
	}

	return nil
}

func GetCredentialByName(name string) (model.CredentialTable, error) {
	var credential model.CredentialTable
	if err := DB.Where("name = ?", name).
		First(&credential).Error; err != nil {
		return model.CredentialTable{}, err
	}

	// Base64 解码数据
	dataBytes, err := base64.StdEncoding.DecodeString(credential.Data)
	if err != nil {
		return model.CredentialTable{}, fmt.Errorf("invalid credential data encoding: %w", err)
	}
	credential.Data = string(dataBytes)
	return credential, nil
}

/*
	// 模拟字符串
	jsonStr := `{"username":"111","password":"111"}`

	// 用 map[string]string 接收
	var data map[string]string
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		panic(err)
	}

	// 使用 data["username"]
	fmt.Println("用户名是:", data["username"])
*/

func UpdateCredentialData(credential model.CredentialTable) error {
	credential.Data = base64.StdEncoding.EncodeToString([]byte(credential.Data))
	if err := DB.Model(&model.CredentialTable{}).Where("id = ?", credential.ID).Update("data", credential.Data).Error; err != nil {
		logger.Errorf("Failed to update credential data: %v", err)
		return err
	}
	return nil
}

func DeleteCredential(credentialID int64) error {
	if err := DB.Delete(&model.CredentialTable{}, credentialID).Error; err != nil {
		logger.Errorf("Failed to delete credential: %v", err)
		return err
	}

	return nil
}

// cicd_tool
func FetchCICDTools() ([]model.CICDToolTable, error) {
	var tools []model.CICDToolTable
	if err := DB.Find(&tools).Error; err != nil {
		return nil, err
	}
	return tools, nil
}

func CreateCICDTool(tool model.CICDToolTable) error {
	if err := DB.Where("name = ?", tool.Name).
		First(&model.CICDToolTable{}).Error; err == nil {
		logger.Errorf("CICD tool already exists: name=%s", tool.Name)
		return err
	}

	if err := DB.Create(&tool).Error; err != nil {
		logger.Errorf("Failed to create CICD tool: %v", err)
		return err
	}

	return nil
}

func GetCICDToolByName(name string) (model.CICDToolTable, error) {
	var tool model.CICDToolTable
	if err := DB.Where("name = ?", name).
		First(&tool).Error; err != nil {
		return model.CICDToolTable{}, err
	}
	return tool, nil
}

func UpdateCICDTool(tool model.CICDToolTable) error {
	if err := DB.Model(&model.CICDToolTable{}).Where("id = ?", tool.ID).Updates(tool).Error; err != nil {
		logger.Errorf("Failed to update CICD tool: %v", err)
		return err
	}
	return nil
}

func DeleteCICDTool(toolID int64) error {
	if err := DB.Delete(&model.CICDToolTable{}, toolID).Error; err != nil {
		logger.Errorf("Failed to delete CICD tool: %v", err)
		return err
	}
	return nil
}
