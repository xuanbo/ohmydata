package postgres_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/db/postgres"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/model"
)

func init() {
	// 日志
	if err := log.Init(); err != nil {
		panic(err)
	}
	// 注册
	if err := postgres.Register(); err != nil {
		panic(err)
	}
}

func TestTableNames(t *testing.T) {
	adapterFactory, err := db.GetAdapterFactory("postgres")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:          "host=localhost user=postgres password=123456 dbname=ohmydata port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConns: 1,
		MaxOpenConns: 8,
	}
	adapter, err := adapterFactory.Create(dataSource)
	if err != nil {
		t.Error(err)
		return
	}
	defer adapter.Close()

	tableNames, err := adapter.TableNames(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("tableNames: %s", tableNames)
}

func TestTable(t *testing.T) {
	adapterFactory, err := db.GetAdapterFactory("postgres")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:          "host=localhost user=postgres password=123456 dbname=ohmydata port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConns: 1,
		MaxOpenConns: 8,
	}
	adapter, err := adapterFactory.Create(dataSource)
	if err != nil {
		t.Error(err)
		return
	}
	defer adapter.Close()

	table, err := adapter.Table(context.TODO(), "oh_data_set")
	if err != nil {
		t.Error(err)
		return
	}
	b, err := json.Marshal(table)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("table: %s", string(b))
}

func TestQueryTable(t *testing.T) {
	adapterFactory, err := db.GetAdapterFactory("postgres")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:          "host=localhost user=postgres password=123456 dbname=ohmydata port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConns: 1,
		MaxOpenConns: 8,
	}
	adapter, err := adapterFactory.Create(dataSource)
	if err != nil {
		t.Error(err)
		return
	}
	defer adapter.Close()

	page := model.NewPagination(1, 10)
	err = adapter.QueryTable(context.TODO(), "oh_data_set", page)
	if err != nil {
		t.Error(err)
		return
	}

	b, err := json.Marshal(page)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("page: %s", string(b))
}

func TestQuery(t *testing.T) {
	adapterFactory, err := db.GetAdapterFactory("postgres")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:          "host=localhost user=postgres password=123456 dbname=ohmydata port=5432 sslmode=disable TimeZone=Asia/Shanghai",
		MaxIdleConns: 1,
		MaxOpenConns: 8,
	}
	adapter, err := adapterFactory.Create(dataSource)
	if err != nil {
		t.Error(err)
		return
	}
	defer adapter.Close()

	page := model.NewPagination(1, 10)
	err = adapter.Query(context.TODO(), "select * from oh_data_set", page)
	if err != nil {
		t.Error(err)
		return
	}

	b, err := json.Marshal(page)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("page: %s", string(b))
}
