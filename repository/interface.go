package repository

import (
	"blog/model"
	"context"
)

// 工厂模式接口
type Repository interface {
	User() UserRepository
	Close() error
	Ping(ctx context.Context) error
	Migrant
}

type UserRepository interface {
	List(pageSize int, pageNum int) (int, []interface{})
	Create(user *model.User) (*model.User, error)
	Delete(id int) error
	Update(user *model.User) (*model.User, error)
	GetAll(userlist []model.User) ([]model.User, error)
	GetUserByID(id int) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
	Migrate() error
}

type Migrant interface {
	Migrate() error
}
