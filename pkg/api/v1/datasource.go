package v1

import (
	"net/http"

	"github.com/xuanbo/ohmydata/pkg/api/middleware"
	"github.com/xuanbo/ohmydata/pkg/api/util"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/model"
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
func (d *DataSource) Init() error {
	// 同步适配层
	return srv.SyncDataSource(d.srv)
}

// AddRoutes 添加路由
func (d *DataSource) AddRoutes(e *echo.Echo) {
	g := e.Group("/v1")
	{
		// 数据源管理
		g.GET("/data-source/list", d.List)
		g.POST("/data-source", d.Create)
		g.PUT("/data-source", d.Modify)
		g.DELETE("/data-source/:id", d.Remove)

		// 数据库操作
		g.GET("/data-source/:id/tables", d.TableNames)
		g.GET("/data-source/:id/table", d.Table)
		g.GET("/data-source/:id/data", d.QueryTable)
		g.POST("/data-source/:id/query", d.Query)
	}
}

// List 列表查询
func (d *DataSource) List(ctx echo.Context) error {
	c := ctx.(*middleware.Context).Ctx()
	list, err := d.srv.List(c)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(list))
}

// Create 创建
func (d *DataSource) Create(ctx echo.Context) error {
	var s entity.DataSource
	if err := ctx.Bind(&s); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := d.srv.Create(c, &s); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(&s))
}

// Modify 修改
func (d *DataSource) Modify(ctx echo.Context) error {
	var s entity.DataSource
	if err := ctx.Bind(&s); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := d.srv.Modify(c, &s); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(&s))
}

// Remove 删除
func (d *DataSource) Remove(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	if err := d.srv.Remove(c, id); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(id))
}

// TableNames 查询数据源表
func (d *DataSource) TableNames(ctx echo.Context) error {
	id := ctx.Param("id")
	list, err := d.srv.TableNames(id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(list))
}

// Table 查询数据源表结构
func (d *DataSource) Table(ctx echo.Context) error {
	id := ctx.Param("id")
	name := ctx.QueryParam("name")
	if name == "" {
		return ctx.JSON(http.StatusBadRequest, model.Fail("请求参数name必须"))
	}
	table, err := d.srv.Table(id, name)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(table))
}

// QueryTable 查询表数据
func (d *DataSource) QueryTable(ctx echo.Context) error {
	id := ctx.Param("id")
	name := ctx.QueryParam("name")
	if name == "" {
		return ctx.JSON(http.StatusBadRequest, model.Fail("请求参数name必须"))
	}
	pagination, err := util.BindPagination(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.Fail(err.Error()))
	}
	if err := d.srv.QueryTable(id, name, pagination); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(pagination))
}

type queryModel struct {
	SQL string `json:"sql"`
}

// Query 查询数据
func (d *DataSource) Query(ctx echo.Context) error {
	id := ctx.Param("id")
	q := new(queryModel)
	if err := ctx.Bind(q); err != nil || q.SQL == "" {
		return ctx.JSON(http.StatusBadRequest, model.Fail("请求参数sql必须"))
	}
	pagination, err := util.BindPagination(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.Fail(err.Error()))
	}
	if err := d.srv.Query(id, q.SQL, pagination); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(pagination))
}
