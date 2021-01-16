package entity

import (
	"time"

	"github.com/xuanbo/ohmydata/pkg/api/util"

	"gorm.io/gorm"
)

// User 用户
type User struct {
	ID        string     `json:"id" gorm:"primaryKey;type:string;size:30"`
	Name      string     `json:"name" gorm:"type:string;size:50"`
	Username  string     `json:"username" gorm:"type:string;size:50"`
	Password  string     `json:"password" gorm:"type:string;size:200"`
	CreatedAt *time.Time `json:"createdAt" gorm:"<-:create"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"<-:create;<-:update"`
	CreatedBy string     `json:"createdBy" gorm:"<-:create"`
	UpdatedBy string     `json:"updatedBy" gorm:"<-:create;<-:update"`
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
