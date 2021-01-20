package srv_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/db/gorm"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/model/condition"
)

var (
	c condition.Clause
)

func TestParseSingleClause(t *testing.T) {
	var list []*entity.DataSource
	err := gorm.Query(context.TODO(), db.DB(), entity.DataSource{}.TableName(), c.Eq("type", "mysql"), &list)
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

func TestParseCombineClause(t *testing.T) {
	var list []*entity.DataSource
	err := gorm.Query(context.TODO(), db.DB(), entity.DataSource{}.TableName(), c.And(c.Eq("description", "mysql"), c.In("type", []interface{}{"mysql"})), &list)
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
