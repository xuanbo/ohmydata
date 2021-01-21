package condition_test

import (
	"encoding/json"
	"testing"

	"github.com/xuanbo/ohmydata/pkg/model/condition"
)

func TestClause(t *testing.T) {
	clause := condition.NewCombineClause(condition.CombineAnd)
	clause.Add(condition.Eq("name", nil))
	clause.Add(condition.In("type", []interface{}{1}))

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
