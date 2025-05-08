package main

import (
	"encoding/json"
	"fmt"
)

type Task struct {
	Service map[string]string `json:"service"`
	Stages  []Stage           `json:"stages"`
}

type Stage struct {
	Stage int    `json:"stage"`
	Steps []Step `json:"steps"`
}

type Step struct {
	Name     string            `json:"name"`
	Type     string            `json:"type"`
	JobURL   string            `json:"job_url"`
	JobParam []string          `json:"job_param"`
	Script   string            `json:"script"`
	Endpoint string            `json:"endpoint"`
	Method   string            `json:"method"`
	Headers  map[string]string `json:"headers"`
}

func main() {
	data := `{
	  "service": {
	    "GIT_BRANCH": "main"
	  },
	  "stages": [
	    {
	      "stage": 1,
	      "steps": [
	        {"name": "step1", "type": "jenkins", "job_url": "", "job_param": [], "script": "", "endpoint": "", "method": "GET", "headers": {}},
	        {"name": "step2", "type": "jenkins", "job_url": "", "job_param": [], "script": "", "endpoint": "", "method": "GET", "headers": {}}
	      ]
	    },
	    {
	      "stage": 2,
	      "steps": [
	        {"name": "step3", "type": "jenkins", "job_url": "", "job_param": [], "script": "", "endpoint": "", "method": "GET", "headers": {}}
	      ]
	    }
	  ]
	}`

	var task Task
	if err := json.Unmarshal([]byte(data), &task); err != nil {
		panic(err)
	}

	// 指定要提取的 stage 编号
	stageNum := 1
	steps := getStepsForStage(task.Stages, stageNum)
	if steps == nil {
		fmt.Printf("Stage %d not found\n", stageNum)
		return
	}

	fmt.Printf("Stage %d has %d step(s):\n", stageNum, len(steps))
	for i, step := range steps {
		fmt.Printf("  Step %d: %s\n", i+1, step.Name)
	}

	if len(steps) > 1 {
		fmt.Println("This stage is parallel.")
	} else {
		fmt.Println("This stage is sequential.")
	}
}

// 提取指定 stage 的 steps
func getStepsForStage(stages []Stage, stageNum int) []Step {
	for _, s := range stages {
		if s.Stage == stageNum {
			return s.Steps
		}
	}
	return nil
}
