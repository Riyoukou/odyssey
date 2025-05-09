package model

import (
	"time"

	"gorm.io/datatypes"
)

// ServiceCICDForm

type ServiceCICDForm struct {
	Yaml    YamlSection    `json:"yaml"`
	Build   BuildSection   `json:"build"`
	Release ReleaseSection `json:"release"`
}

type YamlSection struct {
	IsGitOps   bool   `json:"isGitOps"`
	GitOpsRepo string `json:"gitopsrepo"`
	GitOpsType string `json:"gitopsType"`
	FilePath   string `json:"filePath"`
	Content    string `json:"content"`
}

type BuildSection struct {
	Type     string        `json:"type"`
	CICDTool string        `json:"cicd_tool"`
	JobURL   string        `json:"job_url"`
	JobParam []interface{} `json:"job_param"` // 可根据实际结构替换为具体类型
}

type ReleaseSection struct {
	DeployType        string `json:"deployType"`
	Workload          string `json:"workload"`
	Type              string `json:"type"`
	CICDTool          string `json:"cicd_tool"`
	ArgoCDApplication string `json:"argocd_application"`
}

// repository_model
type ClusterTable struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	APIServer   string    `json:"api_server"`
	Config      string    `json:"config"`
	Region      string    `json:"region"`
	Version     string    `json:"version"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ClusterTable) TableName() string {
	return "clusters"
}

type ProjectTable struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (ProjectTable) TableName() string {
	return "projects"
}

type EnvTable struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	ProjectName string    `json:"project_name"`
	Type        string    `json:"type"`
	Namespace   string    `json:"namespace"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (EnvTable) TableName() string {
	return "envs"
}

type ServiceTable struct {
	ID              int64          `json:"id"`
	Name            string         `json:"name"`
	ProjectName     string         `json:"project_name"`
	CodeLibraryName string         `json:"code_library_name"`
	Clusters        datatypes.JSON `json:"clusters"`
	CICDMap         datatypes.JSON `json:"cicd_map"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
}

func (ServiceTable) TableName() string {
	return "services"
}

type CodeLibraryTable struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	ProjectName    string    `json:"project_name"`
	Type           string    `json:"type"`
	URL            string    `json:"url"`
	CodeSourceName string    `json:"code_source_name"`
	ProjectID      int64     `json:"project_id"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (CodeLibraryTable) TableName() string {
	return "code_libraries"
}

type CodeSourceTable struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	URL          string    `json:"url"`
	Type         string    `json:"type"`
	PrivateToken string    `json:"private_token"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (CodeSourceTable) TableName() string {
	return "code_sources"
}

type BuildRecordTable struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Env         string    `json:"env"`
	Tag         string    `json:"tag"`
	Status      string    `json:"status"`
	ProjectName string    `json:"project_name"`
	BuildUser   string    `json:"build_user"`
	CreatedAt   time.Time `json:"created_at"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (BuildRecordTable) TableName() string {
	return "build_records"
}

type ApiBuildRecord struct {
	Describe    string                      `json:"describe"`
	Env         string                      `json:"env"`
	ProjectName string                      `json:"project_name"`
	Services    []ApiCICDBuildRecordService `json:"services"`
	Tag         string                      `json:"tag"`
	BuildUser   string                      `json:"build_user"`
	Name        string                      `json:"name"`
}

type ApiCICDBuildRecordService struct {
	ServiceName string `json:"service_name"`
	Branch      string `json:"branch"`
}

type BuildServiceRecordTable struct {
	ID              int64     `json:"id"`
	ServiceName     string    `json:"service_name"`
	ProjectName     string    `json:"project_name"`
	Env             string    `json:"env"`
	Image           string    `json:"image"`
	BuildRecordName string    `json:"build_record_name"`
	BuildURL        string    `json:"build_url"`
	Status          string    `json:"status"`
	Branch          string    `json:"branch"`
	BuildID         int64     `json:"build_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (BuildServiceRecordTable) TableName() string {
	return "build_service_records"
}

type DeployRecordTable struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Env             string    `json:"env"`
	ProjectName     string    `json:"project_name"`
	DeployUser      string    `json:"deploy_user"`
	BuildRecordName string    `json:"build_record_name"`
	Status          string    `json:"status"`
	Tag             string    `json:"tag"`
	ClusterNames    string    `json:"cluster_names"`
	CreatedAt       time.Time `json:"created_at"`
	Description     string    `json:"description"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (DeployRecordTable) TableName() string {
	return "deploy_records"
}

type DeployServiceRecordTable struct {
	ID               int64     `json:"id"`
	ServiceName      string    `json:"service_name"`
	ProjectName      string    `json:"project_name"`
	Env              string    `json:"env"`
	DeployRecordName string    `json:"deploy_record_name"`
	ClusterName      string    `json:"cluster_name"`
	Status           string    `json:"status"`
	Image            string    `json:"image"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (DeployServiceRecordTable) TableName() string {
	return "deploy_service_records"
}

type CredentialTable struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Data        string `json:"data"`
}

func (CredentialTable) TableName() string {
	return "credentials"
}

type CICDToolTable struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	Type           string    `json:"type"`
	URL            string    `json:"url"`
	CredentialType string    `json:"credential_type"`
	CredentialName string    `json:"credential_name"`
	Description    string    `json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (CICDToolTable) TableName() string {
	return "cicd_tools"
}
