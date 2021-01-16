package srv_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/model"
	"github.com/xuanbo/ohmydata/pkg/srv"
)

func TestDataSetCreate(t *testing.T) {
	dataSet := srv.NewDataSet()
	err := dataSet.Create(context.TODO(), &entity.DataSet{
		SourceID:   "1347065257465483264",
		Name:       "测试",
		Path:       "/test/page",
		Expression: "select * from oh_data_source",
		EnablePage: true,
		ResponseParams: []*entity.ResponseParam{
			{
				Name:      "name",
				ParamType: entity.String,
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDataSetModify(t *testing.T) {
	dataSet := srv.NewDataSet()
	err := dataSet.Modify(context.TODO(), &entity.DataSet{
		ID:         "1347447490357497856",
		SourceID:   "1347065257465483264",
		Name:       "测试",
		Path:       "/test/page",
		Expression: "select * from oh_data_source",
		EnablePage: true,
		ResponseParams: []*entity.ResponseParam{
			{
				Name:      "name",
				ParamType: entity.String,
			},
			{
				Name:      "type",
				ParamType: entity.String,
			},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}
}

func TestDataSetRemove(t *testing.T) {
	dataSet := srv.NewDataSet()
	if err := dataSet.Remove(context.TODO(), "1347065257465483264"); err != nil {
		t.Error(err)
		return
	}
}

func TestDataSetPage(t *testing.T) {
	dataSet := srv.NewDataSet()
	pagination := model.NewPagination(1, 10)
	if err := dataSet.Page(context.TODO(), &entity.DataSet{Name: "测试"}, pagination); err != nil {
		t.Error(err)
		return
	}
	b, err := json.Marshal(pagination)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("page: %s", string(b))
}

func TestDataSetID(t *testing.T) {
	dataSet := srv.NewDataSet()
	d, err := dataSet.ID(context.TODO(), "1347469523384537088")
	if err != nil {
		t.Error(err)
		return
	}
	b, err := json.Marshal(d)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("dataSet: %s", string(b))
}

func TestDataSetParseExpression(t *testing.T) {
	dataSet := srv.NewDataSet()
	list, err := dataSet.ParseExpression(`
		select * from oh_response_param where name = {{.name}}

		{{if eq .age ""}}
		    x
		{{else}}
			y
		{{end}}

		{{range .ids}}
			{{.}}
		{{ end }}
	`)
	if err != nil {
		t.Error(err)
		return
	}
	b, err := json.Marshal(list)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("variables: %s", string(b))
}

func TestRouter(t *testing.T) {
	router := new(srv.Node)
	router.Add("/user")
	router.Add("/user/:id")
	router.Add("/user/:name")
	router.Add("/user/:id/detail")
	router.Add("/user/page")
	router.Add("/some/path")
	path, params, err := router.Match("/user/1/detail")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("path: %s\n", path)
	t.Logf("params: %s\n", params)
}
