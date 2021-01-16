package v1

import (
	"net/http"

	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/model"

	"github.com/labstack/echo/v4"
)

// Dict 字典
type Dict struct {
}

// NewDict 创建
func NewDict() *Dict {
	return &Dict{}
}

// Init 初始化
func (d *Dict) Init() error {
	return nil
}

// AddRoutes 添加路由
func (d *Dict) AddRoutes(e *echo.Echo) {
	g := e.Group("/v1")
	{
		g.GET("/dict/adapter-types", d.AdapterTypes)
		g.GET("/dict/param-locations", d.ParamLocations)
		g.GET("/dict/param-types", d.ParamTypes)
		g.GET("/dict/convert-types", d.ConvertTypes)
	}
}

// AdapterTypes 适配层类型
func (d *Dict) AdapterTypes(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.OK(db.GetAdapterTypeNames()))
}

var paramLocations = []*model.Dict{
	{
		Name:  "paramPath",
		Text:  "path",
		Value: entity.ParamPath,
	},
	{
		Name:  "paramQuery",
		Text:  "query",
		Value: entity.ParamQuery,
	},
	{
		Name:  "paramBody",
		Text:  "body",
		Value: entity.ParamBody,
	},
}

// ParamLocations 参数位置
func (d *Dict) ParamLocations(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.OK(paramLocations))
}

var paramTypes = []*model.Dict{
	{
		Name:  "boolean",
		Text:  "boolean",
		Value: entity.Boolean,
	},
	{
		Name:  "int",
		Text:  "int",
		Value: entity.Int,
	},
	{
		Name:  "long",
		Text:  "long",
		Value: entity.Long,
	},
	{
		Name:  "float",
		Text:  "float",
		Value: entity.Float,
	},
	{
		Name:  "double",
		Text:  "double",
		Value: entity.Double,
	},
	{
		Name:  "dateTime",
		Text:  "dateTime",
		Value: entity.DateTime,
	},
	{
		Name:  "string",
		Text:  "string",
		Value: entity.String,
	},
	{
		Name:  "object",
		Text:  "object",
		Value: entity.Object,
	},
	{
		Name:  "array",
		Text:  "array",
		Value: entity.Array,
	},
}

// ParamTypes 参数类型
func (d *Dict) ParamTypes(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.OK(paramTypes))
}

var convertTypes = []*model.Dict{
	{
		Name:  "none",
		Text:  "无",
		Value: entity.ConvertNone,
	},
	{
		Name:  "rename",
		Text:  "重命名",
		Value: entity.ConvertRename,
	},
}

// ConvertTypes 转换
func (d *Dict) ConvertTypes(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, model.OK(convertTypes))
}
