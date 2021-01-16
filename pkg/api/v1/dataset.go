package v1

import (
	"net/http"
	"strings"

	"github.com/xuanbo/ohmydata/pkg/api/middleware"
	"github.com/xuanbo/ohmydata/pkg/api/util"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/model"
	"github.com/xuanbo/ohmydata/pkg/srv"

	"github.com/labstack/echo/v4"
)

// DataSet 数据集API管理
type DataSet struct {
	srv *srv.DataSet
}

// NewDataSet 创建
func NewDataSet(srv *srv.DataSet) *DataSet {
	return &DataSet{srv}
}

// Init 初始化
func (d *DataSet) Init() error {
	// 同步数据集
	srv.SyncDataSet(d.srv)
	return nil
}

// AddRoutes 添加路由
func (d *DataSet) AddRoutes(e *echo.Echo) {
	g := e.Group("/v1")
	{
		// 数据集管理
		g.POST("/data-set", d.Create)
		g.PUT("/data-set", d.Modify)
		g.GET("/data-set/:id", d.ID)
		g.GET("/data-set/:id/detail", d.Detail)
		g.GET("/data-set/:id/doc", d.RenderAPIDoc)
		g.POST("/data-set/exp", d.ParseExpression)
		g.DELETE("/data-set/:id", d.Remove)
		g.POST("/data-set/page", d.Page)
		g.GET("/data-set/routes", d.APIRoutes)
	}

	// 数据集API
	e.GET("/api/*", d.ServeAPI)
	e.POST("/api/*", d.ServeAPI)
}

// Create 创建
func (d *DataSet) Create(ctx echo.Context) error {
	var s entity.DataSet
	if err := ctx.Bind(&s); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := d.srv.Create(c, &s); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(s))
}

// Modify 更新
func (d *DataSet) Modify(ctx echo.Context) error {
	var s entity.DataSet
	if err := ctx.Bind(&s); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := d.srv.Modify(c, &s); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(s))
}

// Remove 删除
func (d *DataSet) Remove(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	if err := d.srv.Remove(c, id); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(id))
}

// ID 主键查询
func (d *DataSet) ID(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	s, err := d.srv.ID(c, id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(s))
}

// Detail 主键详情查询
func (d *DataSet) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	s, err := d.srv.Detail(c, id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(s))
}

// Page 分页查询
func (d *DataSet) Page(ctx echo.Context) error {
	var condition entity.DataSet
	if err := ctx.Bind(&condition); err != nil {
		return err
	}
	pagination, err := util.BindPagination(ctx)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.Fail(err.Error()))
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := d.srv.Page(c, &condition, pagination); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(pagination))
}

// RenderAPIDoc 渲染API文档
func (d *DataSet) RenderAPIDoc(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	doc, err := d.srv.RenderAPIDoc(c, id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(doc))
}

type exp struct {
	Expression string `json:"expression"`
}

// ParseExpression 解析表达式
func (d *DataSet) ParseExpression(ctx echo.Context) error {
	var exp exp
	if err := ctx.Bind(&exp); err != nil {
		return err
	}
	list, err := d.srv.ParseExpression(exp.Expression)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(list))
}

// ServeAPI 提供API服务
func (d *DataSet) ServeAPI(ctx echo.Context) error {
	var (
		query = make(map[string]interface{})
		body  = make(map[string]interface{})
	)
	// 参数绑定
	queryParams := ctx.QueryParams()
	for k, v := range queryParams {
		query[k] = v[0]
	}
	if err := ctx.Bind(&body); err != nil {
		return err
	}
	// 去除前缀
	path := ctx.Request().URL.Path
	path = strings.TrimPrefix(path, "/api/")
	c := ctx.(*middleware.Context).Ctx()
	pagination, err := d.srv.ServeAPI(c, path, query, body)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(pagination))
}

// APIRoutes 当前API路由
func (d *DataSet) APIRoutes(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.OK(d.srv.APIRoutes()))
}
