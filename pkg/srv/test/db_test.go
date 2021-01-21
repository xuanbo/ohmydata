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
	engine           = gorm.New(db.DB())
	selectOptionFunc gorm.SelectOptionFunc
)

func TestParseSingleClause(t *testing.T) {
	var list []*entity.DataSource
	err := engine.Query(
		entity.DataSource{}.TableName(),
		&list,
		selectOptionFunc.WithClause(condition.Eq("type", "mysql")),
		selectOptionFunc.WithContext(context.TODO()),
	)
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

	err = engine.Query(
		entity.DataSource{}.TableName(),
		&list,
		selectOptionFunc.WithClause(condition.Eq("type", nil)),
		selectOptionFunc.WithContext(context.TODO()),
	)
	if err != nil {
		t.Error(err)
		return
	}
	b, err = json.Marshal(list)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("list: %s", string(b))

	err = engine.Query(
		entity.DataSource{}.TableName(),
		&list,
		selectOptionFunc.WithClause(condition.In("type", nil)),
		selectOptionFunc.WithContext(context.TODO()),
	)
	if err != nil {
		t.Error(err)
		return
	}
	b, err = json.Marshal(list)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("list: %s", string(b))
}

func TestParseCombineClause(t *testing.T) {
	var list []*entity.DataSource
	err := engine.Query(
		entity.DataSource{}.TableName(),
		&list,
		selectOptionFunc.WithClause(condition.And(condition.Like("name", "MySQL"), condition.In("type", []interface{}{"mysql"}))),
		selectOptionFunc.WithContext(context.TODO()),
	)
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

	err = engine.Query(
		entity.DataSource{}.TableName(),
		&list,
		selectOptionFunc.WithClause(condition.And(condition.Like("name", "MySQL"), nil)),
		selectOptionFunc.WithContext(context.TODO()),
	)
	if err != nil {
		t.Error(err)
		return
	}
	b, err = json.Marshal(list)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("list: %s", string(b))
}
