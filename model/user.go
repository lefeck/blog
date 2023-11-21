package model

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	//json绑定，前后端交互
	UserName string `json:"username" gorm:"type:varchar(100);not null" `
	PassWord string `json:"password" gorm:"type:varchar(100);not null" `
	Email    string `json:"email" gorm:"type:varchar(256);"`
	Avatar   string `json:"avatar" gorm:"type:varchar(256);"` // images
	//Role     int    `json:"role" gorm:"type:int"`
}

type BaseModel struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u *User) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}
