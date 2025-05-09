package service

import (
	"encoding/json"
	"fmt"

	"github.com/Riyoukou/odyssey/app/model"
)

func ServiceCICDMap(action string, clusters []string) {
	// 初始已有集群数据
	clusterMap := map[string][]model.ServiceCICDForm{
		"ctyun-huabei2-ccse01": {{
			Yaml:    model.YamlSection{IsGitOps: true},
			Build:   model.BuildSection{CICDTool: "existing-tool"},
			Release: model.ReleaseSection{},
		}},
	}

	// 定义空模板
	empty := model.ServiceCICDForm{
		Yaml:    model.YamlSection{IsGitOps: true},
		Build:   model.BuildSection{JobParam: []interface{}{}},
		Release: model.ReleaseSection{},
	}

	if action == "add" {
		for _, name := range clusters {
			if _, exists := clusterMap[name]; !exists {
				clusterMap[name] = []model.ServiceCICDForm{empty}
			}
		}
	} else if action == "delete" {
		for _, name := range clusters {
			delete(clusterMap, name)
		}
	} else if action == "update" {
		for _, name := range clusters {
			if _, exists := clusterMap[name]; exists {
				clusterMap[name][0] = model.ServiceCICDForm{
					Yaml:    model.YamlSection{IsGitOps: true, GitOpsRepo: "https://git.example.com/repo"},
					Build:   model.BuildSection{CICDTool: "new-tool", JobURL: "http://ci.example.com/job"},
					Release: model.ReleaseSection{DeployType: "argo"},
				}
			}
		}
	}

	// 输出结果
	jsonBytes, err := json.MarshalIndent(clusterMap, "", "  ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(jsonBytes))
}
