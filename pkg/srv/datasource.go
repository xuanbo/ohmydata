package srv

import (
	"context"
	"errors"
	"time"

	"github.com/xuanbo/ohmydata/pkg/cache"
	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/model"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DataSource 数据源服务
type DataSource struct {
	db *gorm.DB
}

// NewDataSource 创建实例
func NewDataSource() *DataSource {
	return &DataSource{db: db.DB}
}

// Create 新增
func (s *DataSource) Create(ctx context.Context, dataSource *entity.DataSource) error {
	if dataSource.Type == "" || dataSource.Name == "" || dataSource.URL == "" {
		return errors.New("字段[type、name、url]不能为空")
	}
	dataSource.ID = db.NewID()
	if dataSource.MaxIdleConns < 1 {
		dataSource.MaxIdleConns = 1
	}
	if dataSource.MaxOpenConns < 1 {
		dataSource.MaxOpenConns = 8
	}

	if err := s.db.WithContext(ctx).Create(dataSource).Error; err != nil {
		return err
	}

	// 删除缓存
	s.clearCache(ctx, "all")

	// 适配层新增
	return putAdapter(dataSource)
}

// Modify 修改
func (s *DataSource) Modify(ctx context.Context, dataSource *entity.DataSource) error {
	if dataSource.ID == "" {
		return errors.New("更新时主键不能为空")
	}
	if dataSource.Type == "" || dataSource.Name == "" || dataSource.URL == "" {
		return errors.New("字段[type、name、url]不能为空")
	}
	if dataSource.MaxIdleConns < 1 {
		dataSource.MaxIdleConns = 1
	}
	if dataSource.MaxOpenConns < 1 {
		dataSource.MaxOpenConns = 8
	}

	if err := s.db.WithContext(ctx).Save(dataSource).Error; err != nil {
		return err
	}

	// 删除缓存
	s.clearCache(ctx, "all")
	s.clearCache(ctx, dataSource.ID)

	// 适配层更新
	return putAdapter(dataSource)
}

// All 查询所有
func (s *DataSource) All(ctx context.Context) ([]*entity.DataSource, error) {
	var (
		list []*entity.DataSource
		key  = "ohmydata:datasource:all"
		err  error
	)
	if err = cache.Get(ctx, key, &list); errors.Is(err, redis.Nil) {
		// 查询db
		if err = s.db.WithContext(ctx).Find(&list).Error; err != nil {
			return nil, err
		}
		// 写入缓存
		cache.Set(ctx, key, list, cacheTTL)
	}
	return list, err
}

// Remove 删除
func (s *DataSource) Remove(ctx context.Context, id string) error {
	if err := s.db.WithContext(ctx).Delete(&entity.DataSource{}, id).Error; err != nil {
		return err
	}

	// 删除缓存
	s.clearCache(ctx, "all")
	s.clearCache(ctx, id)

	// 适配层删除
	return db.DelAdapter(id)
}

// Test 测试数据源连接
func (s *DataSource) Test(ctx context.Context, dataSource *entity.DataSource) error {
	if dataSource.Type == "" || dataSource.Name == "" || dataSource.URL == "" {
		return errors.New("字段[type、name、url]不能为空")
	}

	factory, err := db.GetAdapterFactory(dataSource.Type)
	if err != nil {
		return err
	}
	adapter, err := factory.Create(dataSource)
	if err != nil {
		return err
	}
	defer adapter.Close()

	if err := adapter.Ping(ctx); err != nil {
		log.Logger().Warn("驱动适配数据源无法连接", zap.String("type", dataSource.Type), zap.Error(err))
		return err
	}
	log.Logger().Info("驱动适配数据源连接正常", zap.String("type", dataSource.Type))
	return nil
}

// TableNames 查询表
func (s *DataSource) TableNames(ctx context.Context, id string) ([]string, error) {
	adapter, err := db.GetAdapter(id)
	if err != nil {
		return nil, err
	}
	return adapter.TableNames(ctx)
}

// Table 查询表结构
func (s *DataSource) Table(ctx context.Context, id, name string) (*db.Table, error) {
	adapter, err := db.GetAdapter(id)
	if err != nil {
		return nil, err
	}
	return adapter.Table(ctx, name)
}

// QueryTable 查询表数据
func (s *DataSource) QueryTable(ctx context.Context, id, tableName string, page *model.Pagination) error {
	adapter, err := db.GetAdapter(id)
	if err != nil {
		return err
	}
	return adapter.QueryTable(ctx, tableName, page)
}

// Query 查询数据
func (s *DataSource) Query(ctx context.Context, id, exp string, page *model.Pagination) error {
	adapter, err := db.GetAdapter(id)
	if err != nil {
		return err
	}
	return adapter.Query(ctx, exp, page)
}

func (s *DataSource) clearCache(ctx context.Context, id string) {
	log.Logger().Debug("清除数据源缓存", zap.String("id", id))
	// 数据源缓存
	cache.DelMatch(ctx, "ohmydata:datasource:"+id+"*")
}

func putAdapter(dataSource *entity.DataSource) error {
	// 适配层新增
	factory, err := db.GetAdapterFactory(dataSource.Type)
	if err != nil {
		return err
	}

	adapter, err := factory.Create(dataSource)
	if err != nil {
		log.Logger().Warn("驱动适配数据源无法连接", zap.String("id", dataSource.ID), zap.String("type", dataSource.Type), zap.Error(err))
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := adapter.Ping(ctx); err == nil {
		log.Logger().Info("驱动适配数据源连接正常", zap.String("id", dataSource.ID), zap.String("type", dataSource.Type))
	} else {
		log.Logger().Warn("驱动适配数据源无法连接", zap.String("id", dataSource.ID), zap.String("type", dataSource.Type), zap.Error(err))
	}

	if err := db.PutAdapter(dataSource.ID, adapter); err != nil {
		log.Logger().Warn("驱动适配数据源更新错误", zap.String("id", dataSource.ID), zap.String("type", dataSource.Type), zap.Error(err))
		return err
	}
	return nil
}

// SyncDataSource 同步适配层
func SyncDataSource(dataSource *DataSource) error {
	log.Logger().Info("初始化驱动适配数据源")

	list, err := dataSource.All(context.TODO())
	if err != nil {
		return err
	}
	// 异步加载数据源驱动
	go func() {
		for _, e := range list {
			putAdapter(e)
		}
	}()
	return nil
}
