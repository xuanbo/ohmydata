package gorm_test

import (
	"testing"

	"github.com/xuanbo/ohmydata/pkg/db/gorm"
	"github.com/xuanbo/ohmydata/pkg/model/condition"
)

func TestParseSingleClause(t *testing.T) {
	s, v, err := gorm.ParseClause(condition.Eq("name", "zhangsan"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(condition.Like("name", "zhangsan"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(condition.IsNull("name"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(condition.In("name", nil))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)
}

func TestParseCombineClause(t *testing.T) {
	s, v, err := gorm.ParseClause(condition.And(condition.Eq("name", "zhangsan"), condition.Gt("age", 10)))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(condition.And(condition.Eq("name", "zhangsan"), condition.In("status", []interface{}{1, 2, 3})))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(condition.And(condition.Eq("name", nil), condition.In("status", []interface{}{1, 2, 3})))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(condition.And(condition.Eq("name", "zhangsan"), condition.In("status", []interface{}{})))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(condition.And(condition.Eq("name", nil), condition.In("status", []interface{}{})))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)
}
