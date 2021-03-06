package v1

import (
	"fmt"
	"net/http"

	"github.com/xuanbo/ohmydata/pkg/api/middleware"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/model"
	"github.com/xuanbo/ohmydata/pkg/model/condition"
	"github.com/xuanbo/ohmydata/pkg/srv"

	"github.com/labstack/echo/v4"
)

// DataSource 数据源API管理
type DataSource struct {
	srv *srv.DataSource
}

// NewDataSource 创建
func NewDataSource(srv *srv.DataSource) *DataSource {
	return &DataSource{srv}
}

// Init 初始化
func (s *DataSource) Init() error {
	// 同步适配层
	return s.srv.SyncDataSource()
}

// AddRoutes 添加路由
func (s *DataSource) AddRoutes(e *echo.Echo) {
	g := e.Group("/v1")
	{
		// 数据源管理
		g.GET("/data-source/list", s.List)
		g.POST("/data-source/test", s.Test)
		g.POST("/data-source", s.Create)
		g.PUT("/data-source", s.Modify)
		g.DELETE("/data-source/:id", s.Remove)

		// 数据库操作
		g.GET("/data-source/:id/tables", s.TableNames)
		g.GET("/data-source/:id/table", s.Table)
		g.POST("/data-source/:id/data", s.QueryTable)
		g.POST("/data-source/:id/query", s.Query)
	}
}

// List 列表查询
func (s *DataSource) List(ctx echo.Context) error {
	c := ctx.(*middleware.Context).Ctx()
	list, err := s.srv.All(c)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(list))
}

// Test 测试数据源连接
func (s *DataSource) Test(ctx echo.Context) error {
	var dataSource entity.DataSource
	if err := ctx.Bind(&dataSource); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.Test(c, &dataSource); err != nil {
		return fmt.Errorf("数据源连接失败: %w", err)
	}
	return ctx.JSON(http.StatusOK, model.OK("数据源连接成功"))
}

// Create 创建
func (s *DataSource) Create(ctx echo.Context) error {
	var dataSource entity.DataSource
	if err := ctx.Bind(&dataSource); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.Create(c, &dataSource); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(&dataSource))
}

// Modify 修改
func (s *DataSource) Modify(ctx echo.Context) error {
	var dataSource entity.DataSource
	if err := ctx.Bind(&dataSource); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.Modify(c, &dataSource); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(&dataSource))
}

// Remove 删除
func (s *DataSource) Remove(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.Remove(c, id); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(id))
}

// TableNames 查询数据源表
func (s *DataSource) TableNames(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	list, err := s.srv.TableNames(c, id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(list))
}

// Table 查询数据源表结构
func (s *DataSource) Table(ctx echo.Context) error {
	id := ctx.Param("id")
	name := ctx.QueryParam("name")
	if name == "" {
		return ctx.JSON(http.StatusBadRequest, model.Fail("请求参数name必须"))
	}
	c := ctx.(*middleware.Context).Ctx()
	table, err := s.srv.Table(c, id, name)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(table))
}

type dataSourceQuery struct {
	Name    string      `json:"name"`
	Operate string      `json:"operate"`
	Value   interface{} `json:"value"`
}

type dataSourceQueryTableCondition struct {
	Name    string             `json:"name" query:"name"`
	Page    uint64             `json:"page" query:"page"`
	Size    uint64             `json:"size" query:"size"`
	Queries []*dataSourceQuery `json:"queries"`
}

// QueryTable 查询表数据
func (s *DataSource) QueryTable(ctx echo.Context) error {
	id := ctx.Param("id")
	var queryCondition dataSourceQueryTableCondition
	if err := ctx.Bind(&queryCondition); err != nil {
		return err
	}
	if queryCondition.Name == "" {
		return ctx.JSON(http.StatusBadRequest, model.Fail("请求参数name必须"))
	}
	pagination := model.NewPagination(queryCondition.Page, queryCondition.Size)

	if len(queryCondition.Queries) > 0 {
		combineClause := condition.NewCombineClause(condition.CombineAnd)
		for _, query := range queryCondition.Queries {
			switch query.Operate {
			case "EQ":
				combineClause.Add(condition.Eq(query.Name, query.Value))
			case "GT":
				combineClause.Add(condition.Gt(query.Name, query.Value))
			case "GTE":
				combineClause.Add(condition.Gte(query.Name, query.Value))
			case "LT":
				combineClause.Add(condition.Lt(query.Name, query.Value))
			case "LTE":
				combineClause.Add(condition.Lte(query.Name, query.Value))
			case "LIKE":
				combineClause.Add(condition.Like(query.Name, query.Value))
			case "IS_NOT_NULL":
				combineClause.Add(condition.IsNotNull(query.Name))
			case "IS_NULL":
				combineClause.Add(condition.IsNull(query.Name))
			}
		}
		clause := condition.WrapCombineClause(combineClause)
		pagination.Clause = clause
	}

	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.QueryTable(c, id, queryCondition.Name, pagination); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(pagination))
}

type dataSourceQueryCondition struct {
	SQL  string `json:"sql"`
	Page uint64 `json:"page" query:"page"`
	Size uint64 `json:"size" query:"size"`
}

// Query 查询数据
func (s *DataSource) Query(ctx echo.Context) error {
	id := ctx.Param("id")
	var condition dataSourceQueryCondition
	if err := ctx.Bind(&condition); err != nil {
		return err
	}
	if condition.SQL == "" {
		return ctx.JSON(http.StatusBadRequest, model.Fail("请求参数sql必须"))
	}
	pagination := model.NewPagination(condition.Page, condition.Size)
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.Query(c, id, condition.SQL, pagination); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(pagination))
}
