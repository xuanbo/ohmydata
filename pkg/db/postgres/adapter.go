package postgres

import (
	"fmt"
	"time"

	"github.com/xuanbo/ohmydata/pkg/db"
	orm "github.com/xuanbo/ohmydata/pkg/db/gorm"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/model"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// adapter PostgreSQL实现
type adapter struct {
	db *gorm.DB
}

func (a *adapter) Ping() error {
	db, err := a.db.DB()
	if err != nil {
		return err
	}
	return db.Ping()
}

func (a *adapter) Close() error {
	db, err := a.db.DB()
	if err != nil {
		return err
	}
	return db.Close()
}

func (a *adapter) TableNames() ([]string, error) {
	var tableNames []string
	if err := a.db.Raw("select tablename from pg_tables where schemaname='public'").Scan(&tableNames).Error; err != nil {
		return nil, err
	}
	return tableNames, nil
}

func (a *adapter) Table(name string) (*db.Table, error) {
	querySQL := fmt.Sprintf("SELECT * FROM %s LIMIT 1", name)
	rows, err := a.db.Raw(querySQL).Rows()
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

func (a *adapter) QueryTable(tableName string, page *model.Pagination) error {
	var (
		total uint64
		data  []map[string]interface{}
	)
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)
	if err := a.db.Raw(countSQL).Scan(&total).Error; err != nil {
		return err
	}

	if total == 0 {
		return nil
	}

	pageSQL := fmt.Sprintf("SELECT * FROM %s LIMIT %d OFFSET %d", tableName, page.Size, page.Offset)
	if err := a.db.Raw(pageSQL).Scan(&data).Error; err != nil {
		return err
	}

	page.Set(total, data)
	return nil
}

func (a *adapter) Query(exp string, page *model.Pagination) error {
	var (
		total uint64
		data  []map[string]interface{}
	)
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM (%s) TMP_COUNT", exp)
	if err := a.db.Raw(countSQL).Scan(&total).Error; err != nil {
		return err
	}

	if total == 0 {
		return nil
	}

	pageSQL := fmt.Sprintf("SELECT * FROM (%s) TMP_PAGE LIMIT %d OFFSET %d", exp, page.Size, page.Offset)
	if err := a.db.Raw(pageSQL).Scan(&data).Error; err != nil {
		return err
	}

	page.Set(total, data)
	return nil
}

// adapter MySQL实现
type adapterFactory struct {
}

func (a *adapterFactory) Create(dataSource *entity.DataSource) (db.Adapter, error) {
	gormDB, err := gorm.Open(postgres.Open(dataSource.URL), &gorm.Config{
		Logger: orm.NewZapLogger(log.Logger(), 200*time.Millisecond),
	})
	if err != nil {
		return nil, err
	}
	// 设置日志级别
	gormDB.Logger.LogMode(logger.Warn)
	// 设置数据库连接池
	db, err := gormDB.DB()
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(dataSource.MaxIdleConns)
	db.SetMaxOpenConns(dataSource.MaxOpenConns)
	return &adapter{gormDB}, nil
}

// Register 注册
func Register() error {
	log.Logger().Info("注册驱动适配", zap.String("name", "postgres"), zap.String("text", "PostgreSQL"))
	return db.RegisterAdapterFactory("postgres", "PostgreSQL", &adapterFactory{})
}
