package elastic_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/db/elastic"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/model"
)

func init() {
	// 日志
	if err := log.Init(); err != nil {
		panic(err)
	}
	//  注册
	if err := elastic.Register(); err != nil {
		panic(err)
	}
}

func TestPing(t *testing.T) {
	adapterFactory, err := db.GetAdapterFactory("elastic")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:      "http://192.168.101.32:9200",
		Username: "elastic",
		Password: "123456",
	}
	adapter, err := adapterFactory.Create(dataSource)
	if err != nil {
		t.Error(err)
		return
	}
	defer adapter.Close()

	if err := adapter.Ping(context.TODO()); err != nil {
		t.Error(err)
		return
	}
	t.Logf("ping: ok")
}

func TestTableNames(t *testing.T) {
	adapterFactory, err := db.GetAdapterFactory("elastic")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:      "http://192.168.101.32:9200",
		Username: "elastic",
		Password: "123456",
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
	adapterFactory, err := db.GetAdapterFactory("elastic")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:      "http://192.168.101.32:9200",
		Username: "elastic",
		Password: "123456",
	}
	adapter, err := adapterFactory.Create(dataSource)
	if err != nil {
		t.Error(err)
		return
	}
	defer adapter.Close()

	table, err := adapter.Table(context.TODO(), "records")
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
	adapterFactory, err := db.GetAdapterFactory("elastic")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:      "http://192.168.101.32:9200",
		Username: "elastic",
		Password: "123456",
	}
	adapter, err := adapterFactory.Create(dataSource)
	if err != nil {
		t.Error(err)
		return
	}
	defer adapter.Close()

	page := model.NewPagination(1, 10)
	err = adapter.QueryTable(context.TODO(), "records", page)
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
	adapterFactory, err := db.GetAdapterFactory("elastic")
	if err != nil {
		t.Error(err)
		return
	}
	dataSource := &entity.DataSource{
		URL:      "http://192.168.101.32:9200",
		Username: "elastic",
		Password: "123456",
	}
	adapter, err := adapterFactory.Create(dataSource)
	if err != nil {
		t.Error(err)
		return
	}
	defer adapter.Close()

	page := model.NewPagination(1, 10)
	err = adapter.Query(context.TODO(), "select * from records", page)
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
