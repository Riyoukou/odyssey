package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/Riyoukou/odyssey/pkg/logger"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
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

func GitGetTags(repoURL, token string) []string {
	// 创建 Remote 实例
	remote := git.NewRemote(nil, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})

	// 可选：如果是私有仓库，需要加上 token（或留空）
	var auth transport.AuthMethod = &githttp.BasicAuth{
		Username: "odyssey", // 可以随便写（GitLab 要求不为空）
		Password: token,     // GitLab 的 Personal Access Token
	}
	// 获取远程引用（tags、branches 等）
	refs, err := remote.List(&git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		logger.Errorf("Failed to list remote refs: %v\n", err)
		return nil
	}
	var tags []string
	// 分类标签和分支
	for _, ref := range refs {
		if ref.Name().IsTag() {
			tags = append(tags, ref.Name().Short())
		}
	}

	return tags
}

func GitGetBranches(repoURL, token string) []string {
	// 创建 Remote 实例
	remote := git.NewRemote(nil, &config.RemoteConfig{
		Name: "origin",
		URLs: []string{repoURL},
	})

	// 可选：如果是私有仓库，需要加上 token（或留空）
	var auth transport.AuthMethod = &githttp.BasicAuth{
		Username: "odyssey", // 可以随便写（GitLab 要求不为空）
		Password: token,     // GitLab 的 Personal Access Token
	}
	// 获取远程引用（tags、branches 等）
	refs, err := remote.List(&git.ListOptions{
		Auth: auth,
	})
	if err != nil {
		logger.Errorf("Failed to list remote refs: %v\n", err)
		return nil
	}
	var branches []string
	// 分类标签和分支
	for _, ref := range refs {
		if ref.Name().IsBranch() {
			branches = append(branches, ref.Name().Short())
		}
	}

	return branches
}

func GitIsBehind(repo *git.Repository, branch, token string) bool {
	// 1. Fetch 最新远程引用
	err := repo.Fetch(&git.FetchOptions{
		Auth: &githttp.BasicAuth{
			Username: "odyssey",
			Password: token,
		},
		RemoteName: "origin",
		RefSpecs:   []config.RefSpec{config.RefSpec("+refs/heads/*:refs/remotes/origin/*")},
		Force:      true,
		Progress:   nil,
		Tags:       git.NoTags,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		logger.Errorf("Failed to fetch remote: %v\n", err)
		return false
	}

	// 2. 获取本地 HEAD
	headRef, err := repo.Head()
	if err != nil {
		logger.Errorf("Failed to get HEAD: %v\n", err)
		return false
	}
	localHash := headRef.Hash()

	// 3. 获取远程分支引用
	remoteRefName := plumbing.NewRemoteReferenceName("origin", branch)
	remoteRef, err := repo.Reference(remoteRefName, true)
	if err != nil {
		logger.Errorf("Failed to get remote ref: %v\n", err)
		return false
	}
	remoteHash := remoteRef.Hash()

	// 4. 比较是否落后
	if localHash != remoteHash {
		return true
	}
	return false
}

func GitClone(repoURL, branch, tmpDir, token string) *git.Repository {
	// Clone repo
	repo, err := git.PlainClone(tmpDir, false, &git.CloneOptions{
		URL:           repoURL,
		ReferenceName: plumbing.NewBranchReferenceName(branch),
		SingleBranch:  true,
		Auth: &githttp.BasicAuth{
			Username: "odyssey",
			Password: token,
		},
	})
	if err != nil {
		logger.Errorf("Git clone failed: %v\n", err)
		return nil
	}
	return repo
}

func GitPull(repo *git.Repository, repoURL, branch, token string, force bool) {
	// ✅ 获取工作区并执行 pull
	w, err := repo.Worktree()
	if err != nil {
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth: &githttp.BasicAuth{
			Username: "odyssey",
			Password: token,
		},
		Force: force,
	})
	// 如果已经是最新，会报 "already up-to-date"
	if err == git.NoErrAlreadyUpToDate {
		fmt.Println("✅ Already up-to-date")
	} else if err != nil {
		logger.Errorf("Git pull failed: %v\n", err)
	}
}

func GitCommit(repo *git.Repository, addPath string) {
	// Git add + commit
	worktree, _ := repo.Worktree()
	_, err := worktree.Add(addPath)
	if err != nil {
		logger.Errorf("Failed to add file to git: %v\n", err)
	}
	_, err = worktree.Commit("odyssey commit message", &git.CommitOptions{
		Author: &object.Signature{ // ⬅️ 是 object.Signature，不是 git.Signature
			Name:  "odyssey",
			Email: "odyssey@neolix.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		logger.Errorf("Git commit failed: %v\n", err)
	}
}

func GitPush(repo *git.Repository, repoURL, branch, tmpDir, token string) {
	// Push
	err := repo.Push(&git.PushOptions{
		Auth: &githttp.BasicAuth{
			Username: "odyssey",
			Password: token,
		},
	})
	if err != nil {
		logger.Errorf("Git push failed: %v\n", err)
	}
	fmt.Println("✅ Git push success.")
}
