package models

import "time"

type UserRepository struct {
	Id                 int
	Identity           string
	UserIdentity       string
	ParentId           int
	RepositoryIdentity string `gorm:"default:NULL"`
	Ext                string `gorm:"default:NULL"`
	Name               string
	CreatedAt          time.Time  `gorm:"column:created_time"`
	UpdatedAt          time.Time  `gorm:"column:updated_time"`
	DeletedAt          *time.Time `gorm:"column:deleted_time"`
}

func (table UserRepository) TableName() string {
	return "user_repository"
}
