package forms

import (
	"blog/model"
	"gorm.io/gorm"
)

type UserForm struct {
	UserName string `json:"name" form:"name" binding:"gte=3,lte=13"`
	PassWord string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
	Avatar   string `json:"avatar" form:"avatar"`
}

func (u *UserForm) GetUser() *model.User {
	return &model.User{
		UserName: u.UserName,
		PassWord: u.PassWord,
		Email:    u.Email,
		Avatar:   u.Avatar,
	}
}

// UserListForm
type UserListForm struct {
	// 页数
	PageNum int `forms:"pagenum" json:"pagenum" binding:"required"`
	// 每页个数
	PageSize int `forms:"pagesize" json:"pagesize" binding:"required"`
}

type UpdateUserForm struct {
	Name     string `json:"name" form:"name" binding:"required"`
	Password string `json:"password" form:"password" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
}

func (u *UpdateUserForm) GetUser(uid uint) *model.User {
	_ = gorm.Model{
		ID: uid,
	}
	return &model.User{
		UserName: u.Name,
		PassWord: u.Password,
		Email:    u.Email,
	}
}
