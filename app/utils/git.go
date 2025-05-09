package utils

import (
	"encoding/json"
	"fmt"
	"log"
)

type GitlabProjectApi struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	WebURL string `json:"web_url"`
}

func GetGitlabProjects(baseURL, token string) []GitlabProjectApi {
	headers := map[string]string{
		"PRIVATE-TOKEN": token, // 请确保这里的 Token 是有效的
	}

	var allProjects []GitlabProjectApi
	page := 1
	for {
		// 构造分页的 URL
		GitLabTagsURL := fmt.Sprintf("%s&page=%d", baseURL, page)

		// 发送请求
		data := make(map[string]interface{})
		_, body, err := SendRequest("GET", GitLabTagsURL, headers, data)
		if err != nil {
			log.Printf("Error fetching GitLab projects: %v", err)
			break
		}

		// 解析响应体
		var projects []GitlabProjectApi
		err = json.Unmarshal(body, &projects)
		if err != nil {
			log.Println("Error unmarshalling response body:", err)
			break
		}

		// 如果没有返回更多的项目，停止分页
		if len(projects) == 0 {
			break
		}

		// 将当前页的数据添加到 allProjects 中
		allProjects = append(allProjects, projects...)

		// 如果当前页的项目数小于 per_page，表示所有项目已加载完毕
		if len(projects) < 100 {
			break
		}

		// 否则，继续请求下一页
		page++
	}

	// 返回所有获取的项目
	return allProjects
}
