package entity

import (
	"github.com/xuanbo/ohmydata/pkg/api/util"

	"gorm.io/gorm"
)

// DataSource 数据源
type DataSource struct {
	Entity
	Type         string `json:"type" gorm:"type:string;size:20"`
	Name         string `json:"name" gorm:"type:string;size:50"`
	Description  string `json:"description" gorm:"type:string;size:100"`
	URL          string `json:"url" gorm:"type:string;size:200"`
	Username     string `json:"username" gorm:"type:string;size:50"`
	Password     string `json:"password" gorm:"type:string;size:100"`
	MaxIdleConns int    `json:"maxIdleConns" gorm:"type:uint;size:3"`
	MaxOpenConns int    `json:"maxOpenConns" gorm:"type:uint;size:5"`
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

// BeforeUpdate 更新前
func (s *DataSource) BeforeUpdate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	user := ctx.Value(util.UserID)
	if user != nil {
		if userID, ok := user.(string); ok {
			s.UpdatedBy = userID
		}
	}
	return nil
}
