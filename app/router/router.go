package router

import (
	"github.com/Riyoukou/odyssey/app/controller/cicd"
	"github.com/Riyoukou/odyssey/app/controller/user"
	"github.com/Riyoukou/odyssey/app/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRouter 设置路由
func SetupRouter(r *gin.Engine) {

	// 使用请求日志中间件
	r.Use(
		middleware.LogRequest(),
		middleware.Cors(),
	)

	// 健康检查
	r.GET("/man", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "what can i say! mamba out!",
		})
	})

	_user := r.Group("/user")
	{
		_user.POST("login", user.HandleUserLogin)
		_user.POST("register", user.HandleUserRegister)
		_user.GET("fetch/:type", user.HandleUserFetch)
		_user.DELETE("delete/:type/:id", user.HandleUserDelete)
	}

	_cicd := r.Group("/cicd")
	{
		_cicd.GET("fetch/:type", cicd.HandleCICDFetch)
		_cicd.GET("get/:type", cicd.HandleCICDGet)
		_cicd.POST("create/:type", cicd.HandleCICDCreate)
		_cicd.POST("update/:type", cicd.HandleCICDUpdate)
		_cicd.DELETE("delete/:type/:id", cicd.HandleCICDDelete)
	}
}
