package srv_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/srv"
)

func TestDataSourceCreate(t *testing.T) {
	dataSource := srv.NewDataSource()
	err := dataSource.Create(context.TODO(), &entity.DataSource{
		Type: "mysql",
		Name: "MySQL",
		URL:  "root:123456@tcp(127.0.0.1:3306)/ohmydata?charset=utf8",
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDataSourceModify(t *testing.T) {
	dataSource := srv.NewDataSource()
	err := dataSource.Modify(context.TODO(), &entity.DataSource{
		ID:          "1346731818388295680",
		Type:        "mysql",
		Description: "123",
		URL:         "root:123456@tcp(127.0.0.1:3306)/ohmydata?charset=utf8",
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDataSourceAll(t *testing.T) {
	dataSource := srv.NewDataSource()
	list, err := dataSource.All(context.TODO())
	if err != nil {
		t.Error(err)
		return
	}
	b, err := json.Marshal(list)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("list: %s", string(b))
}

func TestDataSourceRemove(t *testing.T) {
	dataSource := srv.NewDataSource()
	err := dataSource.Remove(context.TODO(), "1346729420760551424")
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDataSourceTableNames(t *testing.T) {
	dataSource := srv.NewDataSource()
	list, err := dataSource.TableNames("1346729481838006272")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("list: %s", list)
}
