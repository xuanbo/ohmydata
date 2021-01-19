package db

import (
	"context"
	"fmt"
	"sync"

	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/model"
)

var (
	adapters         Adapters
	adapterFactories AdapterFactories
)

// Adapter 数据库适配层
type Adapter interface {
	// Ping 数据库是否连通
	Ping(ctx context.Context) error
	// Close 关闭
	Close() error
	// TableNames 表名
	TableNames(ctx context.Context) ([]string, error)
	// Table 表结构
	Table(ctx context.Context, name string) (*Table, error)
	// QueryTable 查询表
	QueryTable(ctx context.Context, tableName string, page *model.Pagination) error
	// Query 数据库查询
	Query(ctx context.Context, exp string, page *model.Pagination) error
}

// Table 表
type Table struct {
	Name    string    `json:"name"`
	Columns []*Column `json:"columns"`
}

// Column 字段
type Column struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Length   int64  `json:"length"`
	Scale    int64  `json:"scale"`
	Nullable bool   `json:"nullable"`
}

// AdapterFactory 数据库适配层工厂
type AdapterFactory interface {
	// Accept 创建适配层
	Create(dataSource *entity.DataSource) (Adapter, error)
}

// AdapterFactories 数据库适配层工厂管理
type AdapterFactories struct {
	sync.RWMutex
	typeNames map[string]string
	store     map[string]AdapterFactory
}

// Get 获取
func (a *AdapterFactories) Get(adapterType string) (AdapterFactory, error) {
	var (
		adapterFactory AdapterFactory
		exist          bool
	)
	a.RLock()
	adapterFactory, exist = a.store[adapterType]
	a.RUnlock()
	if !exist {
		return nil, fmt.Errorf("数据库适配层不支持: %s", adapterType)
	}
	return adapterFactory, nil
}

// Put 添加
func (a *AdapterFactories) Put(adapterType, adapterName string, adapterFactory AdapterFactory) error {
	a.Lock()
	a.store[adapterType] = adapterFactory
	a.typeNames[adapterType] = adapterName
	a.Unlock()
	return nil
}

// GetAdapterTypeNames 适配类型名称
func (a *AdapterFactories) GetAdapterTypeNames() []*model.Dict {
	list := make([]*model.Dict, 0, 10)
	a.RLock()
	for k, v := range a.typeNames {
		list = append(list, model.NewDict(k, v, k))
	}
	a.RUnlock()
	return list
}

// Adapters 数据库适配层管理
type Adapters interface {
	// Get 获取
	Get(id string) (Adapter, error)
	// Put 添加
	Put(id string, adapter Adapter) error
	// Del 删除
	Del(id string) error
}

// MemoryAdapters 数据库适配层管理，基于内存管理
type MemoryAdapters struct {
	sync.RWMutex
	store map[string]Adapter
}

// Get 获取
func (m *MemoryAdapters) Get(id string) (Adapter, error) {
	var (
		adapter Adapter
		exist   bool
	)
	m.RLock()
	adapter, exist = m.store[id]
	m.RUnlock()
	if !exist {
		return nil, fmt.Errorf("数据库适配层不存在或无法连接: %s", id)
	}
	return adapter, nil
}

// Put 添加
func (m *MemoryAdapters) Put(id string, adapter Adapter) error {
	m.Lock()
	if a, ok := m.store[id]; ok {
		a.Close()
	}
	m.store[id] = adapter
	m.Unlock()
	return nil
}

// Del 删除
func (m *MemoryAdapters) Del(id string) error {
	adapter, err := m.Get(id)
	if err != nil {
		return err
	}
	return adapter.Close()
}

// Init 初始化
func init() {
	adapters = &MemoryAdapters{
		store: make(map[string]Adapter),
	}
	adapterFactories = AdapterFactories{
		typeNames: make(map[string]string),
		store:     make(map[string]AdapterFactory),
	}
}

// GetAdapter 获取
func GetAdapter(id string) (Adapter, error) {
	return adapters.Get(id)
}

// PutAdapter 添加
func PutAdapter(id string, adapter Adapter) error {
	return adapters.Put(id, adapter)
}

// DelAdapter 删除
func DelAdapter(id string) error {
	return adapters.Del(id)
}

// GetAdapterFactory 获取
func GetAdapterFactory(adapterType string) (AdapterFactory, error) {
	return adapterFactories.Get(adapterType)
}

// RegisterAdapterFactory 注册
func RegisterAdapterFactory(adapterType, adapterName string, adapterFactory AdapterFactory) error {
	return adapterFactories.Put(adapterType, adapterName, adapterFactory)
}

// GetAdapterTypeNames 适配类型名称
func GetAdapterTypeNames() []*model.Dict {
	return adapterFactories.GetAdapterTypeNames()
}
