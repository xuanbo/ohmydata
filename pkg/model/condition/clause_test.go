package condition_test

import (
	"encoding/json"
	"testing"

	"github.com/xuanbo/ohmydata/pkg/model/condition"
)

var (
	c condition.Clause
)

func TestClause(t *testing.T) {
	clause := c.And(c.Eq("name", "zhangsan"), c.In("type", []interface{}{1, 2, 3}))
	if clause.IsError() {
		t.Error(clause.Error())
		return
	}
	b, err := json.Marshal(clause)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("json: %s", string(b))

	var c condition.Clause
	if err := json.Unmarshal(b, &c); err != nil {
		t.Error(err)
		return
	}
	b, err = json.Marshal(c)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("json: %s", string(b))
}
