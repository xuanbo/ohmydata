package entity

import (
	"time"

	"github.com/xuanbo/ohmydata/pkg/api/util"

	"gorm.io/gorm"
)

// DataSet 数据集
type DataSet struct {
	ID          string `json:"id" gorm:"primaryKey;type:string;size:30"`
	SourceID    string `json:"sourceID" gorm:"column:;type:string;size:30"`
	Name        string `json:"name" gorm:"type:string;size:50"`
	Description string `json:"description" gorm:"type:string;size:100"`
	// 自定义请求路径
	Path string `json:"path" gorm:"type:string;size:100"`
	// 查询模板
	Expression string `json:"expression" gorm:"type:string;size:1000"`
	// 发布状态
	PublishStatus bool `json:"status" gorm:"type:bool"`
	// 分页
	EnablePage bool `json:"page" gorm:"type:bool"`
	BatchLimit uint `json:"batchLimit" gorm:"type:uint;size:10"`
	// 缓存
	EnableCache   bool       `json:"cache" gorm:"type:bool"`
	ExpireSeconds uint       `json:"expireSeconds" gorm:"type:uint;size:10"`
	CreatedAt     *time.Time `json:"createdAt" gorm:"<-:create"`
	UpdatedAt     *time.Time `json:"updatedAt" gorm:"<-:create;<-:update"`
	CreatedBy     string     `json:"createdBy" gorm:"<-:create"`
	UpdatedBy     string     `json:"updatedBy" gorm:"<-:create;<-:update"`

	// 参数
	RequestParams  []*RequestParam  `json:"requestParams" gorm:"-"`
	ResponseParams []*ResponseParam `json:"responseParam" gorm:"-"`
}

// TableName 表名
func (DataSet) TableName() string {
	return "oh_data_set"
}

// BeforeCreate 创建前
func (s *DataSet) BeforeCreate(tx *gorm.DB) (err error) {
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
func (s *DataSet) BeforeUpdate(tx *gorm.DB) (err error) {
	ctx := tx.Statement.Context
	user := ctx.Value(util.UserID)
	if user != nil {
		if userID, ok := user.(string); ok {
			s.UpdatedBy = userID
		}
	}
	return nil
}
