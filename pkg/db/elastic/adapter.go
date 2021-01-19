package elastic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/model"

	elasticsearch "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go.uber.org/zap"
)

// ErrNil 未初始化
var ErrNil = errors.New("elastic: es client nil")

// adapter MySQL实现
type adapter struct {
	es *elasticsearch.Client
}

func (a *adapter) Ping(ctx context.Context) error {
	if a.es == nil {
		return ErrNil
	}
	var ping esapi.Ping
	_, err := a.es.Ping(ping.WithContext(ctx))
	if err != nil {
		return err
	}
	return nil
}

func (a *adapter) Close() error {
	if a.es == nil {
		return ErrNil
	}
	return nil
}

func (a *adapter) TableNames(ctx context.Context) ([]string, error) {
	if a.es == nil {
		return nil, ErrNil
	}
	// 执行SQL查询
	var sqlQuery esapi.SQLQuery
	resp, err := a.es.SQL.Query(strings.NewReader(`{"query": "show tables"}`), sqlQuery.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析body
	v, err := a.parseBody(resp)
	if err != nil {
		return nil, err
	}
	// 提取表名
	rows := v["rows"].([]interface{})
	tableNames := make([]string, 0, 8)
	for _, row := range rows {
		e := row.([]interface{})
		if "BASE TABLE" == e[1].(string) {
			tableNames = append(tableNames, e[0].(string))
		}
	}
	return tableNames, nil
}

func (a *adapter) Table(ctx context.Context, name string) (*db.Table, error) {
	if a.es == nil {
		return nil, ErrNil
	}
	// 执行SQL查询
	log.Logger().Debug("查询表结构", zap.String("table", name))
	body := fmt.Sprintf(`{"query": "desc %s"}`, "\\\""+name+"\\\"")
	var sqlQuery esapi.SQLQuery
	resp, err := a.es.SQL.Query(strings.NewReader(body), sqlQuery.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析body
	v, err := a.parseBody(resp)
	if err != nil {
		return nil, err
	}
	// 提取表结构
	rows := v["rows"].([]interface{})
	table := &db.Table{
		Name:    name,
		Columns: make([]*db.Column, 0, 10),
	}
	for _, row := range rows {
		e := row.([]interface{})
		table.Columns = append(table.Columns, &db.Column{
			Name: e[0].(string),
			Type: e[2].(string),
		})
	}
	return table, nil
}

func (a *adapter) QueryTable(ctx context.Context, tableName string, page *model.Pagination) error {
	if a.es == nil {
		return ErrNil
	}
	if page.Size < 1 {
		page.Size = 10
	}
	// 执行SQL查询
	log.Logger().Debug("查询表数据", zap.String("table", tableName))
	body := fmt.Sprintf(`{"query": "SELECT * FROM %s LIMIT %d"}`, "\\\""+tableName+"\\\"", page.Size)
	list, err := a.doQuery(ctx, body)
	if err != nil {
		return err
	}
	page.Set(uint64(len(list)), list)
	return nil
}

func (a *adapter) Query(ctx context.Context, exp string, page *model.Pagination) error {
	if a.es == nil {
		return ErrNil
	}
	if page.Size < 1 {
		page.Size = 10
	}
	// 执行SQL查询
	log.Logger().Debug("查询SQL", zap.String("sql", exp))
	body := fmt.Sprintf(`{"query": "SELECT * FROM (%s) TMP_PAGE LIMIT %d"}`, exp, page.Size)
	list, err := a.doQuery(ctx, body)
	if err != nil {
		return err
	}
	page.Set(uint64(len(list)), list)
	return nil
}

func (a *adapter) parseBody(resp *esapi.Response) (map[string]interface{}, error) {
	var v map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %w", err)
	}
	if resp.IsError() {
		reason := v["error"].(map[string]interface{})["reason"]
		return nil, errors.New(reason.(string))
	}
	return v, nil
}

func (a *adapter) doQuery(ctx context.Context, body string) ([]map[string]interface{}, error) {
	// 执行
	var sqlQuery esapi.SQLQuery
	resp, err := a.es.SQL.Query(strings.NewReader(body), sqlQuery.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// 解析body
	v, err := a.parseBody(resp)
	if err != nil {
		return nil, err
	}
	// 提取数据
	var (
		keys = make([]string, 0, 8)
		list = make([]map[string]interface{}, 0, 8)
	)
	columns := v["columns"].([]interface{})
	for _, column := range columns {
		e := column.(map[string]interface{})
		keys = append(keys, e["name"].(string))
	}
	rows := v["rows"].([]interface{})
	for _, row := range rows {
		e := row.([]interface{})
		m := make(map[string]interface{}, len(keys))
		for i, key := range keys {
			m[key] = e[i]
		}
		list = append(list, m)
	}
	return list, nil
}

// adapter MySQL实现
type adapterFactory struct {
}

func (a *adapterFactory) Create(dataSource *entity.DataSource) (db.Adapter, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{dataSource.URL},
		Username:  dataSource.Username,
		Password:  dataSource.Password,
	})
	if err != nil {
		return &adapter{}, err
	}
	return &adapter{es}, nil
}

// Register 注册
func Register() error {
	log.Logger().Info("注册驱动适配", zap.String("name", "elastic"), zap.String("text", "ElasticSearch"))
	return db.RegisterAdapterFactory("elastic", "ElasticSearch", &adapterFactory{})
}
