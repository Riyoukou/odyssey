package repository

import (
	"errors"
	"time"

	"github.com/Riyoukou/odyssey/app/model"
	"github.com/Riyoukou/odyssey/app/utils"
	"github.com/Riyoukou/odyssey/pkg/logger"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func UserLogin(name, password string) (*model.UserTable, error) {
	var user model.UserTable
	if err := DB.Where("name = ?", name).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名错误")
		}
		return nil, errors.New("数据库查询错误")
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}

	// 生成 Token 并存入数据库
	user.Token = utils.GenerateAuthToken(user.Name)
	user.LastLogin = time.Now().Format("2006-01-02 15:04:05")
	if err := DB.Save(&user).Error; err != nil {
		return nil, errors.New("数据库保存 Token 失败")
	}
	return &user, nil
}

func UserRegister(name, password, email, phone string) error {
	var user model.UserTable
	if err := DB.Where("name = ?", name).First(&user).Error; err == nil {
		return errors.New("用户名已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("加密密码失败")
	}

	// 保存用户信息到数据库
	user = model.UserTable{
		Name:      name,
		Password:  string(hashedPassword),
		Email:     email,
		Phone:     phone,
		Role:      "user",
		LastLogin: time.Now().Format("2006-01-02 15:04:05"),
		Token:     utils.GenerateAuthToken(name),
		Type:      "local",
	}
	if err := DB.Create(&user).Error; err != nil {
		return errors.New("数据库保存用户信息失败")
	}
	return nil
}

func FetchUsers() ([]model.UserTable, error) {
	var user []model.UserTable
	if err := DB.Find(&user).Error; err != nil {
		return nil, errors.New("数据库查询用户信息失败")
	}
	return user, nil
}

func DeleteUser(userID int64) error {
	if err := DB.Delete(&model.UserTable{}, userID).Error; err != nil {
		logger.Errorf("Failed to delete user: %v", err)
		return err
	}

	return nil
}

func UpdateUser(user model.UserTable) error {
	if err := DB.Save(&user).Error; err != nil {
		return errors.New("数据库更新用户信息失败")
	}
	return nil
}

func UpdateUserPassword(user model.UserUpdatePassword) error {
	var userInfo model.UserTable

	if err := DB.Where("id = ?", user.ID).First(&userInfo).Error; err != nil {
		return errors.New("用户不存在")
	}

	// 校验旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(user.NewPassword)); err == nil {
		return errors.New("新密码不能与旧密码相同")
	}

	// 校验旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(userInfo.Password), []byte(user.OldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 更新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("加密新密码失败")
	}
	userInfo.Password = string(hashedPassword)

	if err := DB.Save(&userInfo).Error; err != nil {
		return errors.New("数据库更新密码失败")
	}
	return nil

}
