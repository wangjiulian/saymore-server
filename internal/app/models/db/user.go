package db

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"type:varchar(50);unique;not null;comment:用户名"`
	Password  string         `json:"-" gorm:"type:varchar(100);not null;comment:密码"` // json:"-" 表示不在JSON中返回密码
	Email     string         `json:"email" gorm:"type:varchar(100);unique;comment:邮箱"`
	Nickname  string         `json:"nickname" gorm:"type:varchar(50);comment:昵称"`
	Avatar    string         `json:"avatar" gorm:"type:varchar(255);comment:头像"`
	Status    uint8          `json:"status" gorm:"type:tinyint(1);default:1;comment:状态 1:正常 0:禁用"`
	Test      string         `json:"test" gorm:"type:varchar(255);comment:测试字段"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"` // 软删除
}
