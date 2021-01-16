package util

import (
	"errors"
	"strconv"

	"github.com/xuanbo/ohmydata/pkg/model"

	"github.com/labstack/echo/v4"
)

// BindPagination 提取分页参数
func BindPagination(ctx echo.Context) (*model.Pagination, error) {
	page := ctx.QueryParam("page")
	size := ctx.QueryParam("size")
	if page == "" {
		page = "1"
	}
	if size == "" {
		size = "10"
	}
	pageUint64, err := strconv.ParseUint(page, 10, 64)
	if err != nil {
		return nil, errors.New("请求参数page必须为uint64")
	}
	sizeUint64, err := strconv.ParseUint(size, 10, 64)
	if err != nil {
		return nil, errors.New("请求参数size必须为uint64")
	}
	pagination := model.NewPagination(pageUint64, sizeUint64)
	return pagination, err
}
