package user

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Riyoukou/odyssey/app/model"
	"github.com/Riyoukou/odyssey/app/repository"
	"github.com/Riyoukou/odyssey/pkg/logger"
	"github.com/Riyoukou/odyssey/pkg/response"
	"github.com/gin-gonic/gin"
)

func HandleUserFetch(c *gin.Context) {
	var (
		err    error
		result interface{}
	)
	switch c.Param("type") {
	case "user":
		result, err = repository.FetchUsers()
		if err != nil {
			response.Error(c, http.StatusBadRequest, err)
			return
		}
		response.Success(c, result, fmt.Sprintf("%s fetched successfully", c.Param("type")))
	}
}
func HandleUserGet(c *gin.Context) {

}

func HandleUserCreate(c *gin.Context) {

}

func HandleUserDelete(c *gin.Context) {
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
	case "user":
		err = repository.DeleteUser(intID)
		if err != nil {
			response.Error(c, http.StatusBadRequest, err)
			return
		}
		response.Success(c, nil, fmt.Sprintf("%s deleted successfully", c.Param("type")))
	}
}
func HandleUserUpdate(c *gin.Context) {

}

func HandleUserLogin(c *gin.Context) {
	var req model.UserTable
	err := c.ShouldBind(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}

	userInfo, err := repository.UserLogin(req.Name, req.Password)

	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		logger.Error(err)
		return
	}

	response.Success(c, map[string]interface{}{
		"id":    userInfo.ID,
		"name":  userInfo.Name,
		"token": userInfo.Token,
		"role":  userInfo.Role,
	}, "登录成功")
}

func HandleUserRegister(c *gin.Context) {
	var req model.UserTable
	err := c.ShouldBind(&req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		return
	}
	err = repository.UserRegister(req.Name, req.Password, req.Email, req.Phone)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err)
		logger.Error(err)
		return
	}
	response.Success(c, nil, "注册成功")
}
