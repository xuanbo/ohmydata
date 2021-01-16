package model

// Dict 字典
type Dict struct {
	Name  string      `json:"name"`
	Text  string      `json:"text"`
	Value interface{} `json:"value"`
}

// NewDict 创建
func NewDict(name, text string, value interface{}) *Dict {
	return &Dict{name, text, value}
}
