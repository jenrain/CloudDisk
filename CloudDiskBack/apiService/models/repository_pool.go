package models

import "time"

type RepositoryPool struct {
	Id        int
	Identity  string
	Hash      string
	Name      string
	Ext       string
	Size      int64
	Node      string
	CreatedAt time.Time  `gorm:"column:created_time"`
	UpdatedAt time.Time  `gorm:"column:updated_time"`
	DeletedAt *time.Time `gorm:"column:deleted_time"`
}

func (table RepositoryPool) TableName() string {
	return "repository_pool"
}
