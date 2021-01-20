package gorm

import (
	"errors"
	"fmt"
	"strings"

	"github.com/xuanbo/ohmydata/pkg/model/condition"
)

var (
	// ErrClauseNotSupported not supported
	ErrClauseNotSupported = errors.New("clause not supported")
)

// ParseClause 解析短语
func ParseClause(clause *condition.Clause) (string, interface{}, error) {
	if clause == nil {
		return "", nil, condition.ErrClauseNil
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
	return "", nil, ErrClauseNotSupported
}

// ParseSingleClause 解析短语
func ParseSingleClause(clause *condition.Clause) (string, interface{}, error) {
	if clause.IsError() {
		return "", nil, clause.Error()
	}
	switch clause.Op {
	case condition.Eq:
		return "`" + clause.Name + "`" + " = ?", clause.Value, nil
	case condition.NotEq:
		return "`" + clause.Name + "`" + " <> ?", clause.Value, nil
	case condition.Gt:
		return "`" + clause.Name + "`" + " > ?", clause.Value, nil
	case condition.Gte:
		return "`" + clause.Name + "`" + " >= ?", clause.Value, nil
	case condition.Lt:
		return "`" + clause.Name + "`" + " < ?", clause.Value, nil
	case condition.Lte:
		return "`" + clause.Name + "`" + " <= ?", clause.Value, nil
	case condition.Like:
		return "`" + clause.Name + "`" + " LIKE ?", "%" + fmt.Sprintf("%s", clause.Value) + "%", nil
	case condition.NotLike:
		return "`" + clause.Name + "`" + " NOT LIKE ?", "%" + fmt.Sprintf("%s", clause.Value) + "%", nil
	case condition.In:
		return "`" + clause.Name + "`" + " IN ?", clause.Value, nil
	case condition.NotIn:
		return "`" + clause.Name + "`" + " NOT IN ?", clause.Value, nil
	case condition.IsNull:
		return "`" + clause.Name + "`" + " IS NULL", nil, nil
	case condition.IsNotNull:
		return "`" + clause.Name + "`" + " IS NOT NULL", nil, nil
	default:
		return "", nil, fmt.Errorf("unsupported op: %v", condition.Eq)
	}
}

// ParseCombineClause 解析多短语
func ParseCombineClause(clause *condition.Clause) (string, []interface{}, error) {
	if clause.IsError() {
		return "", nil, clause.Error()
	}
	sl := make([]string, 0, 8)
	vl := make([]interface{}, 0, 8)
	for _, c := range clause.Clauses {
		s, v, err := ParseClause(c)
		if err != nil {
			return "", nil, err
		}
		sl = append(sl, s)
		if v == nil {
			continue
		}
		vl = append(vl, v)
	}
	if len(vl) == 0 {
		return "", nil, condition.ErrClauseNil
	}
	return strings.Join(sl, " AND "), vl, nil
}
