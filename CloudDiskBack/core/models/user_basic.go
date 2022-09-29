package models

import "time"

type UserBasic struct {
	Id        int
	Identity  string
	Name      string
	Password  string
	Email     string
	Avatar    string
	CreatedAt time.Time  `gorm:"column:created_time"`
	UpdatedAt time.Time  `gorm:"column:updated_time"`
	DeletedAt *time.Time `gorm:"column:deleted_time"`
}

func (UserBasic) TableName() string {
	return "user_basic"
}
