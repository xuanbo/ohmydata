package gorm_test

import (
	"testing"

	"github.com/xuanbo/ohmydata/pkg/db/gorm"
	"github.com/xuanbo/ohmydata/pkg/model/condition"
)

var c condition.Clause

func TestParseSingleClause(t *testing.T) {
	s, v, err := gorm.ParseClause(c.Eq("name", "zhangsan"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(c.Like("name", "zhangsan"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(c.IsNull("name"))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(c.In("name", nil))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)
}

func TestParseCombineClause(t *testing.T) {
	s, v, err := gorm.ParseClause(c.And(c.Eq("name", "zhangsan"), c.Gt("age", 10)))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)

	s, v, err = gorm.ParseClause(c.And(c.Eq("name", "zhangsan"), c.In("status", []interface{}{1, 2, 3})))
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("s: %s, v: %v", s, v)
}
