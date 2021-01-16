package entity

import (
	"time"

	"github.com/xuanbo/ohmydata/pkg/api/util"

	"gorm.io/gorm"
)

// DataSource 数据源
type DataSource struct {
	ID           string     `json:"id" gorm:"primaryKey;type:string;size:30"`
	Type         string     `json:"type" gorm:"type:string;size:20"`
	Name         string     `json:"name" gorm:"type:string;size:50"`
	Description  string     `json:"description" gorm:"type:string;size:100"`
	URL          string     `json:"url" gorm:"type:string;size:200"`
	Username     string     `json:"username" gorm:"type:string;size:50"`
	Password     string     `json:"password" gorm:"type:string;size:100"`
	MaxIdleConns int        `json:"maxIdleConns" gorm:"type:uint;size:3"`
	MaxOpenConns int        `json:"maxOpenConns" gorm:"type:uint;size:5"`
	CreatedAt    *time.Time `json:"createdAt" gorm:"<-:create"`
	UpdatedAt    *time.Time `json:"updatedAt" gorm:"<-:create;<-:update"`
	CreatedBy    string     `json:"createdBy" gorm:"<-:create"`
	UpdatedBy    string     `json:"updatedBy" gorm:"<-:create;<-:update"`
}

// TableName 表名
func (DataSource) TableName() string {
	return "oh_data_source"
}

// BeforeCreate 创建前
func (s *DataSource) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	user := ctx.Value(util.UserID)
	if user != nil {
		if userID, ok := user.(string); ok {
			s.CreatedBy = userID
		}
	}
	return nil
}

// BeforeSave 更新前
func (s *DataSource) BeforeSave(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	user := ctx.Value(util.UserID)
	if user != nil {
		if userID, ok := user.(string); ok {
			s.UpdatedBy = userID
		}
	}
	return nil
}
