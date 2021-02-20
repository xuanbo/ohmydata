package entity

import (
	"github.com/xuanbo/ohmydata/pkg/api/util"

	"gorm.io/gorm"
)

// User 用户
type User struct {
	Entity
	Name     string `json:"name" gorm:"type:string;size:50"`
	Username string `json:"username" gorm:"type:string;size:50"`
	Password string `json:"password" gorm:"type:string;size:200"`
}

// TableName 表名
func (User) TableName() string {
	return "oh_user"
}

// BeforeCreate 创建前
func (u *User) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	user := ctx.Value(util.UserID)
	if user != nil {
		if userID, ok := user.(string); ok {
			u.CreatedBy = userID
		}
	}
	return nil
}

// BeforeUpdate 更新前
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	user := ctx.Value(util.UserID)
	if user != nil {
		if userID, ok := user.(string); ok {
			u.UpdatedBy = userID
		}
	}
	return nil
}
