package gorm

import (
	"fmt"
	"strings"

	"github.com/xuanbo/ohmydata/pkg/model/condition"
)

// ParseClause 解析短语
func ParseClause(clause *condition.Clause) (string, interface{}, error) {
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
func ParseSingleClause(clause *condition.Clause) (string, interface{}, error) {
	if clause == nil || clause.IsEmpty() || clause.SingleClause == nil {
		return "", nil, nil
	}
	switch clause.Op {
	case condition.OpEq:
		return "`" + clause.Name + "`" + " = ?", clause.Value, nil
	case condition.OpNotEq:
		return "`" + clause.Name + "`" + " <> ?", clause.Value, nil
	case condition.OpGt:
		return "`" + clause.Name + "`" + " > ?", clause.Value, nil
	case condition.OpGte:
		return "`" + clause.Name + "`" + " >= ?", clause.Value, nil
	case condition.OpLt:
		return "`" + clause.Name + "`" + " < ?", clause.Value, nil
	case condition.OpLte:
		return "`" + clause.Name + "`" + " <= ?", clause.Value, nil
	case condition.OpLike:
		return "`" + clause.Name + "`" + " LIKE ?", "%" + fmt.Sprintf("%v", clause.Value) + "%", nil
	case condition.OpNotLike:
		return "`" + clause.Name + "`" + " NOT LIKE ?", "%" + fmt.Sprintf("%v", clause.Value) + "%", nil
	case condition.OpIn:
		return "`" + clause.Name + "`" + " IN ?", clause.Value, nil
	case condition.OpNotIn:
		return "`" + clause.Name + "`" + " NOT IN ?", clause.Value, nil
	case condition.OpIsNull:
		return "`" + clause.Name + "`" + " IS NULL", nil, nil
	case condition.OpIsNotNull:
		return "`" + clause.Name + "`" + " IS NOT NULL", nil, nil
	default:
		return "", nil, fmt.Errorf("unsupported op: %v", condition.OpEq)
	}
}

// ParseCombineClause 解析多短语
func ParseCombineClause(clause *condition.Clause) (string, []interface{}, error) {
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
