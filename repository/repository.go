package repository

import (
	"context"
	"gorm.io/gorm"
)

func NewRepository(db *gorm.DB) Repository {
	r := &repository{
		user: NewUserRepository(db),
		db:   db,
	}
	r.migrants = getMigrants(r.user)
	return r
}

func getMigrants(obj ...interface{}) []Migranter {
	var migrants []Migranter
	for _, obj := range migrants {
		if m, ok := obj.(Migranter); ok {
			migrants = append(migrants, m)
		}
	}
	return migrants
}

type Migranter interface {
	Migrant() error
}

// 创建数据库的表结构
func (r *repository) Migrant() error {
	for _, m := range r.migrants {
		if err := m.Migrant(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Migrate() error {
	return r.Migrant()
}

type repository struct {
	user     UserRepository
	db       *gorm.DB
	migrants []Migranter
}

func (r *repository) User() UserRepository {
	return r.user
}

func (r *repository) Close() error {
	db, _ := r.db.DB()
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
	}
	return nil
}

func (r *repository) Ping(ctx context.Context) error {
	db, _ := r.db.DB()
	if db != nil {
		if err := db.PingContext(ctx); err != nil {
			return err
		}
	}
	return nil
}
