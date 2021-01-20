package gorm

import (
	"context"
	"fmt"

	"github.com/xuanbo/ohmydata/pkg/model/condition"

	"gorm.io/gorm"
)

// Query 查询
func Query(ctx context.Context, db *gorm.DB, table string, clause *condition.Clause, dest interface{}) error {
	if clause == nil {
		return condition.ErrClauseNil
	}
	if clause.SingleClause != nil {
		s, v, err := ParseSingleClause(clause)
		if err != nil {
			return err
		}
		sql := fmt.Sprintf("SELECT * FROM `%s` WHERE %s", table, s)
		return db.WithContext(ctx).Raw(sql, v).Find(dest).Error
	}
	if clause.CombineClause != nil {
		s, v, err := ParseCombineClause(clause)
		if err != nil {
			return err
		}
		sql := fmt.Sprintf("SELECT * FROM `%s` WHERE %s", table, s)
		return db.WithContext(ctx).Raw(sql, v...).Find(dest).Error
	}
	return ErrClauseNotSupported
}
