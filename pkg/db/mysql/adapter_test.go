package mysql_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/db/mysql"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/model"
	"github.com/xuanbo/ohmydata/pkg/model/condition"
)

func init() {
	// 日志
	if err := log.Init(); err != nil {
		panic(err)
	}
	//  注册
	if err := mysql.Register(); err != nil {
		panic(err)
	}
}

func TestTableNames(t *testing.T) {
	adapterFactory, err := db.GetAdapterFactory("mysql")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:          "root:123456@tcp(127.0.0.1:3306)/ohmydata?charset=utf8",
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
	adapterFactory, err := db.GetAdapterFactory("mysql")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:          "root:123456@tcp(127.0.0.1:3306)/ohmydata?charset=utf8",
		MaxIdleConns: 1,
		MaxOpenConns: 8,
	}
	adapter, err := adapterFactory.Create(dataSource)
	if err != nil {
		t.Error(err)
		return
	}
	defer adapter.Close()

	table, err := adapter.Table(context.TODO(), "oh_data_source")
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
	adapterFactory, err := db.GetAdapterFactory("mysql")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		ID:           "test",
		URL:          "root:123456@tcp(127.0.0.1:3306)/ohmydata?charset=utf8",
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
	page.Clause = condition.Eq("name", "MySQL")
	err = adapter.QueryTable(context.TODO(), "oh_data_source", page)
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
	adapterFactory, err := db.GetAdapterFactory("mysql")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:          "root:123456@tcp(127.0.0.1:3306)/ohmydata?charset=utf8",
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
	err = adapter.Query(context.TODO(), "select * from oh_data_source", page)
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
