package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/xuanbo/ohmydata/pkg/db"
	orm "github.com/xuanbo/ohmydata/pkg/db/gorm"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/model"

	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var selectOptionFunc orm.SelectOptionFunc

// adapter MySQL实现
type adapter struct {
	engine *orm.Engine
	db     *gorm.DB
}

func (a *adapter) Ping(ctx context.Context) error {
	db, err := a.db.DB()
	if err != nil {
		return err
	}
	return db.PingContext(ctx)
}

func (a *adapter) Close() error {
	db, err := a.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (a *adapter) TableNames(ctx context.Context) ([]string, error) {
	var tableNames []string
	if err := a.db.WithContext(ctx).Raw("SHOW TABLES").Scan(&tableNames).Error; err != nil {
		return nil, err
	}
	return tableNames, nil
}

func (a *adapter) Table(ctx context.Context, name string) (*db.Table, error) {
	querySQL := fmt.Sprintf("SELECT * FROM %s LIMIT 1", name)
	rows, err := a.db.WithContext(ctx).Raw(querySQL).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	types, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}

	table := &db.Table{
		Name:    name,
		Columns: make([]*db.Column, len(types)),
	}
	for i, columnType := range types {
		column := &db.Column{
			Name: columnType.Name(),
			Type: columnType.DatabaseTypeName(),
		}
		table.Columns[i] = column

		var (
			length   int64
			scale    int64
			nullable bool
			ok       bool
		)
		length, ok = columnType.Length()
		if ok {
			table.Columns[i].Length = length
		}
		length, scale, ok = columnType.DecimalSize()
		if ok {
			table.Columns[i].Length = length
			table.Columns[i].Scale = scale
		}
		nullable, ok = columnType.Nullable()
		if ok {
			table.Columns[i].Nullable = nullable
		}
	}
	return table, nil
}

func (a *adapter) QueryTable(ctx context.Context, tableName string, page *model.Pagination) error {
	var (
		total uint64
		data  []map[string]interface{}
		err   error
	)
	if total, err = a.engine.Page(
		tableName,
		&data,
		selectOptionFunc.WithContext(ctx),
		selectOptionFunc.WithClause(page.Clause),
		selectOptionFunc.WithPageSize(page.Page, page.Size),
		selectOptionFunc.WithTablePrefix("`"),
		selectOptionFunc.WithTableSuffix("`"),
		selectOptionFunc.WithColumnPrefix("`"),
		selectOptionFunc.WithColumnSuffix("`"),
	); err != nil {
		return err
	}
	page.Set(total, data)
	return nil
}

func (a *adapter) Query(ctx context.Context, exp string, page *model.Pagination) error {
	var (
		total uint64
		data  []map[string]interface{}
	)

	// 未分页限制查询
	if page.Page == 0 {
		pageSQL := fmt.Sprintf("SELECT * FROM (%s) TMP_PAGE LIMIT %d", exp, page.Offset)
		if err := a.db.WithContext(ctx).Raw(pageSQL).Scan(&data).Error; err != nil {
			return err
		}
		total = uint64(len(data))

		page.Set(total, data)
		return nil
	}

	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) TMP_COUNT", exp)
	if err := a.db.WithContext(ctx).Raw(countSQL).Scan(&total).Error; err != nil {
		return err
	}

	if total == 0 {
		return nil
	}

	pageSQL := fmt.Sprintf("SELECT * FROM (%s) TMP_PAGE LIMIT %d, %d", exp, page.Offset, page.Size)
	if err := a.db.WithContext(ctx).Raw(pageSQL).Scan(&data).Error; err != nil {
		return err
	}

	page.Set(total, data)
	return nil
}

// adapter MySQL实现
type adapterFactory struct {
}

func (a *adapterFactory) Create(dataSource *entity.DataSource) (db.Adapter, error) {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dataSource.URL,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{
		DisableAutomaticPing: true,
		Logger:               orm.NewZapLogger(log.Logger(), 200*time.Millisecond, fmt.Sprintf("驱动 [%s] ", dataSource.ID)),
	})
	if err != nil {
		return &adapter{engine: orm.New(gormDB), db: gormDB}, err
	}
	// 设置日志级别
	gormDB.Logger.LogMode(logger.Info)
	// 设置数据库连接池
	db, err := gormDB.DB()
	if err != nil {
		return &adapter{engine: orm.New(gormDB), db: gormDB}, err
	}
	db.SetMaxIdleConns(dataSource.MaxIdleConns)
	db.SetMaxOpenConns(dataSource.MaxOpenConns)
	return &adapter{engine: orm.New(gormDB), db: gormDB}, nil
}

// Register 注册
func Register() error {
	log.Logger().Info("注册驱动适配", zap.String("name", "mysql"), zap.String("text", "MySQL"))
	return db.RegisterAdapterFactory("mysql", "MySQL", &adapterFactory{})
}
