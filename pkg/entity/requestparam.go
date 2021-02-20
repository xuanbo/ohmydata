package entity

import (
	"github.com/xuanbo/ohmydata/pkg/api/util"

	"gorm.io/gorm"
)

// RequestParam 请求参数
type RequestParam struct {
	Entity
	DataSetID   string `json:"dataSetId" gorm:"type:string;size:30"`
	Name        string `json:"name" gorm:"type:string;size:50"`
	Description string `json:"description" gorm:"type:string;size:100"`
	// 请求参数位置
	ParamLocation ParamLocation `json:"paramLocation" gorm:"type:uint;size:1"`
	// 参数类型
	ParamType ParamType `json:"paramType" gorm:"type:uint;size:2"`
	// 是否必须
	Required bool `json:"required" gorm:"type:bool"`
	// 默认值
	DefaultValue string `json:"defaultValue" gorm:"type:string;size:100"`
}

// TableName 表名
func (RequestParam) TableName() string {
	return "oh_request_param"
}

// BeforeCreate 创建前
func (p *RequestParam) BeforeCreate(tx *gorm.DB) (err error) {
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
func (p *RequestParam) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	user := ctx.Value(util.UserID)
	if user != nil {
		if userID, ok := user.(string); ok {
			p.UpdatedBy = userID
		}
	}
	return nil
}
