package model

import "github.com/xuanbo/ohmydata/pkg/model/condition"

// Pagination 分页
type Pagination struct {
	Page   uint64            `json:"page"`
	Size   uint64            `json:"size"`
	Offset uint64            `json:"-"`
	Total  uint64            `json:"total"`
	Clause *condition.Clause `json:"clause"`
	Data   interface{}       `json:"data"`
}

// NewPagination 创建分页对象
func NewPagination(page uint64, size uint64) *Pagination {
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 10
	}
	offset := (page - 1) * size
	return &Pagination{Page: page, Size: size, Offset: offset}
}

// Set 设置数据
func (p *Pagination) Set(total uint64, data interface{}) {
	p.Total = total
	p.Data = data
}
