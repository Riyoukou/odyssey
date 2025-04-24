package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Riyoukou/odyssey/app/repository"
	"github.com/Riyoukou/odyssey/app/router"
	"github.com/Riyoukou/odyssey/cmd/cron"
	"github.com/Riyoukou/odyssey/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	if err := initConfig(); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}
	// 初始化日志
	if err := logger.InitLogger(viper.GetString("server.log_file_name")); err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	// 初始化数据库
	repository.Init()

	// 设置运行模式
	gin.SetMode(viper.GetString("server.mode"))

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	setupRouter(r)

	go cron.CronGeneratedRegister()

	// 启动服务器
	address := viper.GetString("server.address")
	if err := r.Run(address); err != nil {
		log.Fatalf("启动服务器失败: %v", err)
	}
}

func initConfig() error {
	e := os.Getenv("GO_ENV")
	if e == "" {
		e = "dev" // 如果没有设置 GO_ENV 默认加载 dev 环境
	}

	viper.SetConfigName("conf.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(fmt.Sprintf("configs/%s", e)) // 根据环境选择文件夹
	return viper.ReadInConfig()
}

func setupRouter(r *gin.Engine) {
	// TODO: 注册路由
	router.SetupRouter(r)
}
