package gorm

import (
	"context"
	"errors"
	"fmt"

	"github.com/xuanbo/ohmydata/pkg/model/condition"

	"gorm.io/gorm"
)

var (
	// ErrClauseNotSupported not supported
	ErrClauseNotSupported = errors.New("clause not supported")
)

// Engine 引擎
type Engine struct {
	db *gorm.DB
}

// SelectOption select选项
type SelectOption struct {
	clause       *condition.Clause
	page         uint64
	size         uint64
	ctx          context.Context
	tablePrefix  string
	tableSuffix  string
	columnPrefix string
	columnSuffix string
}

// SelectOptionFunc select选项方法
type SelectOptionFunc func(op *SelectOption)

// WithClause 配置短语
func (SelectOptionFunc) WithClause(clause *condition.Clause) SelectOptionFunc {
	return func(op *SelectOption) {
		op.clause = clause
	}
}

// WithContext 配置上下文
func (SelectOptionFunc) WithContext(ctx context.Context) SelectOptionFunc {
	return func(op *SelectOption) {
		op.ctx = ctx
	}
}

// WithPageSize 配置分页
func (SelectOptionFunc) WithPageSize(page, size uint64) SelectOptionFunc {
	return func(op *SelectOption) {
		op.page = page
		op.size = size
	}
}

// WithTablePrefix 配置表前缀
func (SelectOptionFunc) WithTablePrefix(tablePrefix string) SelectOptionFunc {
	return func(op *SelectOption) {
		op.tablePrefix = tablePrefix
	}
}

// WithTableSuffix 配置表后缀
func (SelectOptionFunc) WithTableSuffix(tableSuffix string) SelectOptionFunc {
	return func(op *SelectOption) {
		op.tableSuffix = tableSuffix
	}
}

// WithColumnPrefix 配置字段前缀
func (SelectOptionFunc) WithColumnPrefix(columnPrefix string) SelectOptionFunc {
	return func(op *SelectOption) {
		op.columnPrefix = columnPrefix
	}
}

// WithColumnSuffix 配置字段后缀
func (SelectOptionFunc) WithColumnSuffix(columnSuffix string) SelectOptionFunc {
	return func(op *SelectOption) {
		op.columnSuffix = columnSuffix
	}
}

// New 创建
func New(db *gorm.DB) *Engine {
	return &Engine{db: db}
}

// Count count查询
func (e *Engine) Count(table string, opts ...SelectOptionFunc) (uint, error) {
	var (
		total uint
		op    = &SelectOption{}
		db    = e.db
	)
	for _, opt := range opts {
		opt(op)
	}
	if op.ctx != nil {
		db = db.WithContext(op.ctx)
	}
	table = op.tablePrefix + table + op.tableSuffix

	if op.clause == nil || op.clause.IsEmpty() {
		sql := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		err := db.Raw(sql).Scan(&total).Error
		return total, err
	}
	if op.clause.SingleClause != nil {
		s, v, err := ParseSingleClause(op.clause, columnOptionFunc.WithColumnPrefix(op.columnPrefix), columnOptionFunc.WithColumnSuffix(op.columnSuffix))
		if err != nil {
			return 0, err
		}
		sql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, s)
		if err := db.Raw(sql, v).Scan(&total).Error; err != nil {
			return 0, err
		}
		return total, nil
	}
	if op.clause.CombineClause != nil {
		s, v, err := ParseCombineClause(op.clause, columnOptionFunc.WithColumnPrefix(op.columnPrefix), columnOptionFunc.WithColumnSuffix(op.columnSuffix))
		if err != nil {
			return 0, err
		}
		sql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, s)
		if err := db.Raw(sql, v...).Scan(&total).Error; err != nil {
			return 0, err
		}
		return total, nil
	}
	return 0, ErrClauseNotSupported
}

// Query 查询
func (e *Engine) Query(table string, dest interface{}, opts ...SelectOptionFunc) error {
	var (
		op = &SelectOption{}
		db = e.db
	)
	for _, opt := range opts {
		opt(op)
	}
	if op.ctx != nil {
		db = db.WithContext(op.ctx)
	}
	table = op.tablePrefix + table + op.tableSuffix

	if op.clause == nil || op.clause.IsEmpty() {
		sql := fmt.Sprintf("SELECT * FROM %s", table)
		return db.Raw(sql).Find(&dest).Error
	}
	if op.clause.SingleClause != nil {
		s, v, err := ParseSingleClause(op.clause, columnOptionFunc.WithColumnPrefix(op.columnPrefix), columnOptionFunc.WithColumnSuffix(op.columnSuffix))
		if err != nil {
			return err
		}
		sql := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, s)
		return db.Raw(sql, v).Find(dest).Error
	}
	if op.clause.CombineClause != nil {
		s, v, err := ParseCombineClause(op.clause, columnOptionFunc.WithColumnPrefix(op.columnPrefix), columnOptionFunc.WithColumnSuffix(op.columnSuffix))
		if err != nil {
			return err
		}
		sql := fmt.Sprintf("SELECT * FROM %s WHERE %s", table, s)
		return db.Raw(sql, v...).Find(dest).Error
	}
	return ErrClauseNotSupported
}

// Page 查询
func (e *Engine) Page(table string, dest interface{}, opts ...SelectOptionFunc) (uint64, error) {
	var (
		op    = &SelectOption{}
		db    = e.db
		total uint64
	)
	for _, opt := range opts {
		opt(op)
	}
	if op.ctx != nil {
		db = db.WithContext(op.ctx)
	}
	if op.page == 0 {
		op.page = 1
	}
	if op.size == 0 {
		op.size = 10
	}
	offset := (op.page - 1) * op.size
	table = op.tablePrefix + table + op.tableSuffix

	if op.clause == nil || op.clause.IsEmpty() {
		sql := fmt.Sprintf("SELECT COUNT(*) FROM %s", table)
		if err := db.Raw(sql).Scan(&total).Error; err != nil {
			return 0, err
		}
		if total == 0 {
			return 0, nil
		}
		if offset == 0 {
			sql := fmt.Sprintf("SELECT * FROM %s LIMIT %d", table, op.size)
			return 0, db.Raw(sql).Find(&dest).Error
		}
		sql = fmt.Sprintf("SELECT * FROM %s OFFSET %d LIMIT %d", table, offset, op.size)
		return 0, db.Raw(sql).Find(&dest).Error
	}
	if op.clause.SingleClause != nil {
		s, v, err := ParseSingleClause(op.clause, columnOptionFunc.WithColumnPrefix(op.columnPrefix), columnOptionFunc.WithColumnSuffix(op.columnSuffix))
		if err != nil {
			return 0, err
		}
		sql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, s)
		if err := db.Raw(sql, v).Scan(&total).Error; err != nil {
			return 0, err
		}
		if total == 0 {
			return 0, nil
		}
		if offset == 0 {
			sql := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT %d", table, s, op.size)
			return total, db.Raw(sql, v).Find(dest).Error
		}
		sql = fmt.Sprintf("SELECT * FROM %s WHERE %s OFFSET %d LIMIT %d", table, s, offset, op.size)
		return total, db.Raw(sql, v).Find(dest).Error
	}
	if op.clause.CombineClause != nil {
		s, v, err := ParseCombineClause(op.clause, columnOptionFunc.WithColumnPrefix(op.columnPrefix), columnOptionFunc.WithColumnSuffix(op.columnSuffix))
		if err != nil {
			return 0, err
		}
		sql := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, s)
		if err := db.Raw(sql, v...).Scan(&total).Error; err != nil {
			return 0, err
		}
		if total == 0 {
			return 0, nil
		}
		if offset == 0 {
			sql := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT %d", table, s, op.size)
			return total, db.Raw(sql, v...).Find(dest).Error
		}
		sql = fmt.Sprintf("SELECT * FROM %s WHERE %s OFFSET %d LIMIT %d", table, s, offset, op.size)
		return total, db.Raw(sql, v...).Find(dest).Error
	}
	return 0, ErrClauseNotSupported
}
