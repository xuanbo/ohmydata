package gorm

import (
	"fmt"
	"strings"

	"github.com/xuanbo/ohmydata/pkg/model/condition"
)

var columnOptionFunc ColumnOptionFunc

// ColumnOption column选项
type ColumnOption struct {
	columnPrefix string
	columnSuffix string
}

// ColumnOptionFunc column选项方法
type ColumnOptionFunc func(op *ColumnOption)

// WithColumnPrefix 配置字段前缀
func (ColumnOptionFunc) WithColumnPrefix(columnPrefix string) ColumnOptionFunc {
	return func(op *ColumnOption) {
		op.columnPrefix = columnPrefix
	}
}

// WithColumnSuffix 配置字段后缀
func (ColumnOptionFunc) WithColumnSuffix(columnSuffix string) ColumnOptionFunc {
	return func(op *ColumnOption) {
		op.columnSuffix = columnSuffix
	}
}

// ParseClause 解析短语
func ParseClause(clause *condition.Clause, opts ...ColumnOptionFunc) (string, interface{}, error) {
	if clause == nil || clause.IsEmpty() {
		return "", nil, nil
	}
	if clause.SingleClause != nil {
		s, v, err := ParseSingleClause(clause)
		if err != nil {
			return "", nil, err
		}
		return s, v, nil
	}
	if clause.CombineClause != nil {
		s, v, err := ParseCombineClause(clause)
		if err != nil {
			return "", nil, err
		}
		return s, v, nil
	}
	return "", nil, nil
}

// ParseSingleClause 解析短语
func ParseSingleClause(clause *condition.Clause, opts ...ColumnOptionFunc) (string, interface{}, error) {
	if clause == nil || clause.IsEmpty() || clause.SingleClause == nil {
		return "", nil, nil
	}
	op := &ColumnOption{}
	for _, opt := range opts {
		opt(op)
	}
	column := op.columnPrefix + clause.Name + op.columnSuffix

	switch clause.Op {
	case condition.OpEq:
		return column + " = ?", clause.Value, nil
	case condition.OpNotEq:
		return column + " <> ?", clause.Value, nil
	case condition.OpGt:
		return column + " > ?", clause.Value, nil
	case condition.OpGte:
		return column + " >= ?", clause.Value, nil
	case condition.OpLt:
		return column + " < ?", clause.Value, nil
	case condition.OpLte:
		return column + " <= ?", clause.Value, nil
	case condition.OpLike:
		return column + " LIKE ?", "%" + fmt.Sprintf("%v", clause.Value) + "%", nil
	case condition.OpNotLike:
		return column + " NOT LIKE ?", "%" + fmt.Sprintf("%v", clause.Value) + "%", nil
	case condition.OpIn:
		return column + " IN ?", clause.Value, nil
	case condition.OpNotIn:
		return column + " NOT IN ?", clause.Value, nil
	case condition.OpIsNull:
		return column + " IS NULL", nil, nil
	case condition.OpIsNotNull:
		return column + " IS NOT NULL", nil, nil
	default:
		return "", nil, fmt.Errorf("unsupported op: %v", condition.OpEq)
	}
}

// ParseCombineClause 解析多短语
func ParseCombineClause(clause *condition.Clause, opts ...ColumnOptionFunc) (string, []interface{}, error) {
	if clause == nil || clause.IsEmpty() || clause.CombineClause == nil {
		return "", nil, nil
	}
	sl := make([]string, 0, 8)
	vl := make([]interface{}, 0, 8)
	for _, c := range clause.Clauses {
		s, v, err := ParseClause(c)
		if err != nil {
			return "", nil, err
		}
		if s == "" || v == nil {
			continue
		}
		sl = append(sl, "("+s+")")
		vl = append(vl, v)
	}
	if len(vl) == 0 {
		return "", nil, nil
	}
	if clause.Combine == condition.CombineAnd {
		return strings.Join(sl, " AND "), vl, nil
	}
	return strings.Join(sl, " OR "), vl, nil
}
