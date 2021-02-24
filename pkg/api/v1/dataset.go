package v1

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/xuanbo/ohmydata/pkg/api/middleware"
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
func (s *DataSet) Init() error {
	// 同步数据集
	srv.SyncDataSet(s.srv)
	return nil
}

// AddRoutes 添加路由
func (s *DataSet) AddRoutes(e *echo.Echo) {
	g := e.Group("/v1")
	{
		// 数据集管理
		g.POST("/data-set", s.Create)
		g.PUT("/data-set", s.Modify)
		g.GET("/data-set/:id", s.ID)
		g.GET("/data-set/:id/detail", s.Detail)
		g.PUT("/data-set/:id/publish-status/:publishStatus", s.ChangePublishStatus)
		g.GET("/data-set/:id/doc", s.RenderAPIDoc)
		g.POST("/data-set/exp", s.ParseExpression)
		g.POST("/data-set/preview", s.PreviewData)
		g.DELETE("/data-set/:id", s.Remove)
		g.POST("/data-set/page", s.Page)
		g.GET("/data-set/routes", s.APIRoutes)
	}

	// 数据集API
	e.GET("/api/*", s.ServeAPI)
	e.POST("/api/*", s.ServeAPI)
}

// Create 创建
func (s *DataSet) Create(ctx echo.Context) error {
	var dataSet entity.DataSet
	if err := ctx.Bind(&dataSet); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.Create(c, &dataSet); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(dataSet))
}

// Modify 更新
func (s *DataSet) Modify(ctx echo.Context) error {
	var dataSet entity.DataSet
	if err := ctx.Bind(&dataSet); err != nil {
		return err
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.Modify(c, &dataSet); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(dataSet))
}

// Remove 删除
func (s *DataSet) Remove(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.Remove(c, id); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(id))
}

// ID 主键查询
func (s *DataSet) ID(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	dataSet, err := s.srv.ID(c, id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(dataSet))
}

// Detail 主键详情查询
func (s *DataSet) Detail(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	dataSet, err := s.srv.Detail(c, id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(dataSet))
}

// ChangePublishStatus 修改发布状态
func (s *DataSet) ChangePublishStatus(ctx echo.Context) error {
	id := ctx.Param("id")
	publishStatus := ctx.Param("publishStatus")
	status, err := strconv.ParseBool(publishStatus)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "publishStatus必须为bool")
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.ChangePublishStatus(c, id, status); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(id))
}

type dataSetCondition struct {
	Name string `json:"name" query:"name"`
	Path string `json:"path" query:"path"`
	Page uint64 `json:"page" query:"page"`
	Size uint64 `json:"size" query:"size"`
}

// Page 分页查询
func (s *DataSet) Page(ctx echo.Context) error {
	var condition dataSetCondition
	if err := ctx.Bind(&condition); err != nil {
		return err
	}
	pagination := model.NewPagination(condition.Page, condition.Size)
	dataSet := entity.DataSet{
		Name: condition.Name,
		Path: condition.Path,
	}
	c := ctx.(*middleware.Context).Ctx()
	if err := s.srv.Page(c, &dataSet, pagination); err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(pagination))
}

// RenderAPIDoc 渲染API文档
func (s *DataSet) RenderAPIDoc(ctx echo.Context) error {
	id := ctx.Param("id")
	c := ctx.(*middleware.Context).Ctx()
	doc, err := s.srv.RenderAPIDoc(c, id)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(doc))
}

type dataSetExpressionCondition struct {
	Expression string `json:"expression"`
}

// ParseExpression 解析表达式
func (s *DataSet) ParseExpression(ctx echo.Context) error {
	var condition dataSetExpressionCondition
	if err := ctx.Bind(&condition); err != nil {
		return err
	}
	list, err := s.srv.ParseExpression(condition.Expression)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(list))
}

type dataSetPreviewCondition struct {
	DataSet *entity.DataSet        `json:"dataSet"`
	Params  map[string]interface{} `json:"params"`
}

// PreviewData 预览表达式内容
func (s *DataSet) PreviewData(ctx echo.Context) error {
	var condition dataSetPreviewCondition
	if err := ctx.Bind(&condition); err != nil {
		return err
	}
	if condition.DataSet == nil {
		return echo.NewHTTPError(http.StatusBadRequest, "dataSet必须存在")
	}
	if condition.Params == nil {
		condition.Params = make(map[string]interface{})
	}
	c := ctx.(*middleware.Context).Ctx()
	v, err := s.srv.PreviewData(c, condition.DataSet, condition.Params)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(v))
}

// ServeAPI 提供API服务
func (s *DataSet) ServeAPI(ctx echo.Context) error {
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
	// 合并参数
	params := query
	for k, v := range body {
		params[k] = v
	}
	// 去除前缀
	path := ctx.Request().URL.Path
	path = strings.TrimPrefix(path, "/api/")
	c := ctx.(*middleware.Context).Ctx()
	pagination, err := s.srv.ServeAPI(c, path, params)
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, model.OK(pagination))
}

// APIRoutes 当前API路由
func (s *DataSet) APIRoutes(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.OK(s.srv.APIRoutes()))
}
