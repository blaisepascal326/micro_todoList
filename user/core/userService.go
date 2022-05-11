package core

import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"user/model"
	"user/services"
)

func BuildUser(item model.User) *services.UserModel{
	userModel := services.UserModel{
		ID: uint32(item.ID),
		UserName: item.UserName,
		CreateAt: item.CreatedAt.Unix(),
		UpdateAt: item.UpdatedAt.Unix(),
		DeleteAt: item.DeletedAt.Unix(),
	}
	return &userModel
}

func (*UserService) UserLogin (ctx context.Context, req *services.UserRequest, resp *services.UserDetailResponse) error {
	var user model.User
	resp.Code = 200
	if err := model.DB.Where("user_name = ?", req.UserName).First(&user).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			resp.Code = 400
			return nil
		}
		resp.Code = 500
		return nil
	}
	if !user.CheckPassword(req.PassWord) {
		resp.Code = 400
		return nil
	}
	resp.UserDetail = BuildUser(user)
	return nil
}

func (*UserService) UserRegister(ctx context.Context, req *services.UserRequest, resp *services.UserDetailResponse) error {
	if req.PassWord != req.PasswordConfirm {
		err := errors.New("两次密码输入不一致，请重新输入")
		return err
	}
	count := 0
	if err := model.DB.Model(&model.User{}).Where("user_name = ?", req.UserName).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		err := errors.New("用户名已存在")
		return err
	}
	var user model.User
	user.UserName = req.UserName
	if err := user.SetPassword(req.PassWord); err != nil {
		return err
	}
	if err := model.DB.Create(&user).Error; err != nil {
		return err
	}
	resp.UserDetail = BuildUser(user)
	return nil
}
