package entity

import (
	"time"

	"github.com/xuanbo/ohmydata/pkg/api/util"

	"gorm.io/gorm"
)

// ResponseParam 响应参数
type ResponseParam struct {
	ID          string `json:"id" gorm:"primaryKey;type:string;size:30"`
	DataSetID   string `json:"dataSetId" gorm:"type:string;size:30"`
	Name        string `json:"name" gorm:"type:string;size:50"`
	Description string `json:"description" gorm:"type:string;size:100"`
	// 参数类型
	ParamType ParamType `json:"paramType" gorm:"type:uint;size:2"`
	// 转换方式
	ConvertType  ConvertType `json:"convertType" gorm:"type:uint;size:1"`
	ConvertValue string      `json:"convertValue" gorm:"type:string;size:255"`
	CreatedAt    *time.Time  `json:"createdAt" gorm:"<-:create"`
	UpdatedAt    *time.Time  `json:"updatedAt" gorm:"<-:create;<-:update"`
	CreatedBy    string      `json:"createdBy" gorm:"<-:create"`
	UpdatedBy    string      `json:"updatedBy" gorm:"<-:create;<-:update"`
}

// TableName 表名
func (ResponseParam) TableName() string {
	return "oh_response_param"
}

// BeforeCreate 创建前
func (p *ResponseParam) BeforeCreate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	user := ctx.Value(util.UserID)
	if user != nil {
		if userID, ok := user.(string); ok {
			p.CreatedBy = userID
		}
	}
	return nil
}

// BeforeUpdate 更新前
func (p *ResponseParam) BeforeUpdate(tx *gorm.DB) error {
	ctx := tx.Statement.Context
	user := ctx.Value(util.UserID)
	if user != nil {
		if userID, ok := user.(string); ok {
			p.UpdatedBy = userID
		}
	}
	return nil
}
