package srv

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/xuanbo/ohmydata/pkg/cache"
	"github.com/xuanbo/ohmydata/pkg/db"
	"github.com/xuanbo/ohmydata/pkg/entity"
	"github.com/xuanbo/ohmydata/pkg/log"
	"github.com/xuanbo/ohmydata/pkg/model"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// DataSet 数据集服务
type DataSet struct {
	db     *gorm.DB
	tpl    *template.Template
	router *Node
}

// NewDataSet 创建实例
func NewDataSet() *DataSet {
	tpl := template.New("api_template.md")
	// 自定义方法
	tpl.Funcs(template.FuncMap{
		"pl": pl,
		"pt": pt,
	})
	tpl, err := tpl.ParseFiles("./config/api_template.md")
	if err != nil {
		log.Logger().Info("初始化API文档模板错误", zap.Error(err))
	}
	return &DataSet{db: db.DB, tpl: tpl, router: new(Node)}
}

// Create 新增
func (s *DataSet) Create(ctx context.Context, dataSet *entity.DataSet) error {
	dataSet.ID = db.NewID()
	if err := s.validDataSet(ctx, dataSet); err != nil {
		return err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if len(dataSet.RequestParams) > 0 {
			for _, requestParam := range dataSet.RequestParams {
				if err := validRequestParam(requestParam); err != nil {
					return err
				}
				requestParam.ID = db.NewID()
				requestParam.DataSetID = dataSet.ID
			}
			if err := tx.CreateInBatches(dataSet.RequestParams, len(dataSet.RequestParams)).Error; err != nil {
				return err
			}
		}
		if len(dataSet.ResponseParams) > 0 {
			for _, responseParam := range dataSet.ResponseParams {
				if err := validResponseParam(responseParam); err != nil {
					return err
				}
				responseParam.ID = db.NewID()
				responseParam.DataSetID = dataSet.ID
			}
			if err := tx.CreateInBatches(dataSet.ResponseParams, len(dataSet.ResponseParams)).Error; err != nil {
				return err
			}
		}
		if err := tx.Create(dataSet).Error; err != nil {
			return err
		}
		return nil
	})
	// 清除缓存
	s.clearCache(ctx, "all")
	return err
}

// Modify 修改
func (s *DataSet) Modify(ctx context.Context, dataSet *entity.DataSet) error {
	if dataSet.ID == "" {
		return errors.New("更新时主键不能为空")
	}
	if err := s.validDataSet(ctx, dataSet); err != nil {
		return err
	}
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除参数
		if err := tx.Delete(entity.RequestParam{}, "data_set_id = ?", dataSet.ID).Error; err != nil {
			return err
		}
		if err := tx.Delete(entity.ResponseParam{}, "data_set_id = ?", dataSet.ID).Error; err != nil {
			return err
		}
		// 新增参数
		if len(dataSet.RequestParams) > 0 {
			for _, requestParam := range dataSet.RequestParams {
				if err := validRequestParam(requestParam); err != nil {
					return err
				}
				requestParam.ID = db.NewID()
				requestParam.DataSetID = dataSet.ID
			}
			if err := tx.CreateInBatches(dataSet.RequestParams, len(dataSet.RequestParams)).Error; err != nil {
				return err
			}
		}
		if len(dataSet.ResponseParams) > 0 {
			for _, responseParam := range dataSet.ResponseParams {
				if err := validResponseParam(responseParam); err != nil {
					return err
				}
				responseParam.ID = db.NewID()
				responseParam.DataSetID = dataSet.ID
			}
			if err := tx.CreateInBatches(dataSet.ResponseParams, len(dataSet.ResponseParams)).Error; err != nil {
				return err
			}
		}
		if err := tx.Save(dataSet).Error; err != nil {
			return err
		}
		return nil
	})
	// 清除数据缓存
	s.clearCache(ctx, "all")
	s.clearCache(ctx, dataSet.ID)
	return err
}

// Remove 主键删除
func (s *DataSet) Remove(ctx context.Context, id string) error {
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(entity.RequestParam{}, "data_set_id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Delete(entity.ResponseParam{}, "data_set_id = ?", id).Error; err != nil {
			return err
		}
		return tx.Delete(entity.DataSet{}, "id = ?", id).Error
	})
	// 清除数据缓存
	s.clearCache(ctx, "all")
	s.clearCache(ctx, id)
	return err
}

// Page 分页查询
func (s *DataSet) Page(ctx context.Context, dataSet *entity.DataSet, page *model.Pagination) error {
	where := make([]string, 0, 8)
	values := make([]interface{}, 0, 8)
	if dataSet.Name != "" {
		where = append(where, "name LIKE ?")
		values = append(values, "%"+dataSet.Name+"%")
	}
	if dataSet.Path != "" {
		where = append(where, "path LIKE ?")
		values = append(values, "%"+dataSet.Path+"%")
	}
	var (
		total int64
		list  []*entity.DataSet
	)
	condition := strings.Join(where, " AND ")
	if err := s.db.WithContext(ctx).WithContext(ctx).Model(&entity.DataSet{}).Where(condition, values...).
		Count(&total).Error; err != nil {
		return err
	}
	if total < 1 {
		return nil
	}
	if err := s.db.WithContext(ctx).Model(&entity.DataSet{}).Where(condition, values...).
		Limit(int(page.Size)).Offset(int(page.Offset)).
		Find(&list).Error; err != nil {
		return err
	}
	page.Set(uint64(total), list)
	return nil
}

// All 查询所有
func (s *DataSet) All(ctx context.Context) ([]*entity.DataSet, error) {
	var (
		list []*entity.DataSet
		key  = "ohmydata:dataset:all"
		err  error
	)
	if err = cache.Get(ctx, key, &list); errors.Is(err, redis.Nil) {
		// 查询db
		if err = s.db.WithContext(ctx).Find(&list).Error; err != nil {
			return nil, err
		}
		// 写入缓存
		cache.Set(ctx, key, list, cacheTTL)
	}
	return list, err
}

// ID 主键查询
func (s *DataSet) ID(ctx context.Context, id string) (*entity.DataSet, error) {
	var (
		dataSet entity.DataSet
		key     = "ohmydata:dataset:" + id
		err     error
	)
	if err = cache.Get(ctx, key, &dataSet); errors.Is(err, redis.Nil) {
		// 查询db
		if err = s.db.WithContext(ctx).Where("id = ?", id).Find(&dataSet).Error; err != nil {
			return nil, err
		}
		if dataSet.ID == "" {
			return nil, nil
		}
		// 写入缓存
		cache.Set(ctx, key, &dataSet, cacheTTL)
	}
	return &dataSet, err
}

// Detail 主键详情查询
func (s *DataSet) Detail(ctx context.Context, id string) (*entity.DataSet, error) {
	var (
		dataSet        entity.DataSet
		requestParams  []*entity.RequestParam
		responseParams []*entity.ResponseParam
		key            = "ohmydata:dataset:" + id + ":detail"
		err            error
	)
	if err = cache.Get(ctx, key, &dataSet); errors.Is(err, redis.Nil) {
		// 查询db
		if err = s.db.WithContext(ctx).Where("id = ?", id).Find(&dataSet).Error; err != nil {
			return nil, err
		}
		if dataSet.ID == "" {
			return nil, nil
		}
		if err = s.db.WithContext(ctx).Where("data_set_id = ?", id).Find(&requestParams).Error; err != nil {
			return nil, err
		}
		if err = s.db.WithContext(ctx).Where("data_set_id = ?", id).Find(&responseParams).Error; err != nil {
			return nil, err
		}
		dataSet.RequestParams = requestParams
		dataSet.ResponseParams = responseParams
		// 写入缓存
		cache.Set(ctx, key, &dataSet, cacheTTL)
	}
	return &dataSet, err
}

// ChagePublishStatus 修改发布状态
func (s *DataSet) ChagePublishStatus(ctx context.Context, id string, status bool) error {
	if err := s.db.WithContext(ctx).Model(&entity.DataSet{}).Where("id = ?", id).
		Update("publish_status", status).Error; err != nil {
		return err
	}
	// 清除缓存
	s.clearCache(ctx, "all")
	s.clearCache(ctx, id)
	return nil
}

// RenderAPIDoc 渲染API文档
func (s *DataSet) RenderAPIDoc(ctx context.Context, id string) (string, error) {
	dataSet, err := s.Detail(ctx, id)
	if err != nil {
		return "", err
	}
	if s == nil {
		return "", errors.New("数据集不存在")
	}
	if !dataSet.PublishStatus {
		return "", errors.New("数据集未发布")
	}
	var buff bytes.Buffer
	if err := s.tpl.Execute(&buff, dataSet); err != nil {
		return "", err
	}
	return buff.String(), nil
}

// ParseExpression 解析表达式
func (s *DataSet) ParseExpression(expression string) ([]*entity.RequestParam, error) {
	if expression == "" {
		return nil, nil
	}
	// 解析模板
	tpl, err := template.New(db.NewID()).Parse(expression)
	if err != nil {
		return nil, err
	}
	// 语法树JSON化
	b, err := json.Marshal(tpl.Tree)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	// 语法树中提取变量
	var vars []string
	vars = parseFromMapVariables(m, vars)
	vars = removeDuplicate(vars)
	if len(vars) == 0 {
		return nil, nil
	}
	variables := make([]*entity.RequestParam, len(vars))
	for i, e := range vars {
		variables[i] = &entity.RequestParam{
			Name:          e,
			ParamLocation: entity.ParamBody,
			ParamType:     entity.String,
			Required:      true,
		}
	}
	return variables, nil
}

// ServeAPI 提供API服务
func (s *DataSet) ServeAPI(ctx context.Context, path string, query, body map[string]interface{}) (interface{}, error) {
	// 合并参数
	params := query
	for k, v := range body {
		params[k] = v
	}
	// 路径匹配
	node, nameParams, err := s.router.Match(path)
	if err != nil {
		return nil, err
	}
	// 绑定的数据集ID
	id := node.Handle.(string)
	if id == "" {
		return nil, echo.NewHTTPError(http.StatusNotFound, "API不存在，请检查访问路径")
	}
	for k, v := range nameParams {
		params[k] = v
	}

	// 查询数据集
	dataSet, err := s.Detail(ctx, id)
	if err != nil {
		return nil, err
	}
	if dataSet == nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, "API不存在，请检查访问路径")
	}

	// 分页参数处理
	var page, size uint64
	if page, size, err = parsePagination(params, dataSet); err != nil {
		return nil, err
	}
	pagination := model.NewPagination(page, size)

	if dataSet.EnableCache {
		return doSelectFromCache(ctx, dataSet, pagination, params)
	}

	return doSelect(dataSet, pagination, params)
}

// APIRoutes 当前API路由
func (s *DataSet) APIRoutes() *Node {
	return s.router
}

func (s *DataSet) clearCache(ctx context.Context, id string) {
	log.Logger().Debug("清除数据集缓存", zap.String("id", id))
	// 数据集缓存
	cache.DelMatch(ctx, "ohmydata:dataset:"+id+"*")
	// 数据集API结果缓存
	cache.DelMatch(ctx, "ohmydata:datasetcache:"+id+":*")
}

func (s *DataSet) validDataSet(ctx context.Context, dataSet *entity.DataSet) error {
	if dataSet.Name == "" {
		return errors.New("数据集名称不能为空")
	}
	dataSet.Path = strings.TrimPrefix(dataSet.Path, "/")
	if dataSet.Path == "" {
		return errors.New("数据集自定义请求路径不能为空")
	}
	if dataSet.SourceID == "" {
		return errors.New("数据集数据源不能为空")
	}
	if len(dataSet.ResponseParams) == 0 {
		return errors.New("响应参数不能为空")
	}
	var total int64
	if err := s.db.WithContext(ctx).Model(dataSet).Where("name = ? AND id <> ?", dataSet.Name, dataSet.ID).Count(&total).Error; err != nil {
		return err
	}
	if total > 0 {
		return errors.New("数据集名称已存在")
	}
	if err := s.db.WithContext(ctx).Model(dataSet).Where("path = ? AND id <> ?", dataSet.Path, dataSet.ID).Count(&total).Error; err != nil {
		return err
	}
	if total > 0 {
		return errors.New("数据集自定义请求路径已存在")
	}
	return nil
}

func validRequestParam(requestParam *entity.RequestParam) error {
	if requestParam.Name == "" {
		return errors.New("请求参数名称不能为空")
	}
	return nil
}

func validResponseParam(responseParam *entity.ResponseParam) error {
	if responseParam.Name == "" {
		return errors.New("响应参数名称不能为空")
	}
	if responseParam.ConvertType == entity.ConvertRename && responseParam.ConvertValue == "" {
		return errors.New("响应参数重命名时，转换值不能为空，需要是字段别名")
	}
	return nil
}

func parseFromMapVariables(m map[string]interface{}, variables []string) []string {
	for k, v := range m {
		if k == "Ident" {
			if list, ok := v.([]interface{}); ok {
				for _, e := range list {
					variables = append(variables, fmt.Sprintf("%s", e))
				}
			}
			if a, ok := v.(string); ok {
				variables = append(variables, a)
			}
		}
		if m, ok := v.(map[string]interface{}); ok {
			variables = parseFromMapVariables(m, variables)
		}
		if s, ok := v.([]interface{}); ok {
			variables = parseFromSliceVariables(s, variables)
		}
	}
	return variables
}

func parseFromSliceVariables(s []interface{}, variables []string) []string {
	for _, e := range s {
		if m, ok := e.(map[string]interface{}); ok {
			variables = parseFromMapVariables(m, variables)
		}
	}
	return variables
}

func removeDuplicate(list []string) []string {
	result := []string{}
	// 存放不重复主键
	tempMap := map[string]byte{}
	for _, e := range list {
		// 排除空字符串以及，内置方法名
		if e == "" || equalAny(e, funcNames) {
			continue
		}
		size := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != size {
			// 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

func equalAny(v string, list []string) bool {
	for _, e := range list {
		if v == e {
			return true
		}
	}
	return false
}

func parsePagination(param map[string]interface{}, dataSet *entity.DataSet) (uint64, uint64, error) {
	if !dataSet.EnablePage {
		param["page"] = 0
		param["size"] = dataSet.BatchLimit
		return 0, uint64(dataSet.BatchLimit), nil
	}
	var (
		page, size uint64
		err        error
	)
	if v, ok := param["page"]; ok {
		s := fmt.Sprintf("%v", v)
		if page, err = strconv.ParseUint(s, 10, 64); err != nil {
			return 0, 0, fmt.Errorf("分页参数page必须是一个正整数: %v", v)
		}
	}
	if v, ok := param["size"]; ok {
		s := fmt.Sprintf("%v", v)
		if size, err = strconv.ParseUint(s, 10, 64); err != nil {
			return 0, 0, fmt.Errorf("分页参数size必须是一个正整数: %v", v)
		}
	}
	return page, size, nil
}

func doSelectFromCache(ctx context.Context, dataSet *entity.DataSet, pagination *model.Pagination, params map[string]interface{}) (interface{}, error) {
	// 从缓存中查询
	var (
		v   interface{}
		key string
		err error
	)
	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	key = "ohmydata:datasetcache:" + dataSet.ID + ":" + hex.EncodeToString(b)
	log.Logger().Debug("从缓存中查询结果", zap.String("id", dataSet.ID), zap.String("key", key))
	if err = cache.Get(ctx, key, &v); errors.Is(err, redis.Nil) {
		// 缓存未命中，查询db
		log.Logger().Debug("缓存未命中", zap.String("id", dataSet.ID), zap.String("key", key))
		v, err = doSelect(dataSet, pagination, params)
		if err != nil {
			return nil, err
		}
		// 写入缓存
		cache.Set(ctx, key, &v, time.Duration(dataSet.ExpireSeconds)*time.Second)
		log.Logger().Debug("数据写入缓存", zap.String("id", dataSet.ID), zap.String("key", key))
	}
	return v, err
}

func doSelect(dataSet *entity.DataSet, pagination *model.Pagination, params map[string]interface{}) (interface{}, error) {
	adapter, err := db.GetAdapter(dataSet.SourceID)
	if err != nil {
		return nil, err
	}

	// 渲染表达式
	log.Logger().Info("表达式模板", zap.String("expression", dataSet.Expression))

	var buff bytes.Buffer
	tpl, err := template.New(dataSet.ID).Parse(dataSet.Expression)
	if err != nil {
		return nil, err
	}
	if err := tpl.Execute(&buff, params); err != nil {
		return nil, err
	}
	exp := buff.String()
	log.Logger().Info("表达式", zap.String("expression", exp))

	// 查询
	if err := adapter.Query(exp, pagination); err != nil {
		return nil, err
	}

	// 结果处理
	return convertResponseParams(pagination, dataSet)
}

func convertResponseParams(pagination *model.Pagination, dataSet *entity.DataSet) (interface{}, error) {
	if pagination.Data == nil {
		return nil, nil
	}
	if list, ok := pagination.Data.([]map[string]interface{}); ok {
		for _, row := range list {
			// 字段过滤
			for k := range row {
				var ok bool
				for _, p := range dataSet.ResponseParams {
					if k == p.Name {
						ok = true
						break
					}
				}
				if !ok {
					delete(row, k)
				}
			}
			// 处理字段转换
			for _, p := range dataSet.ResponseParams {
				// 重命名
				switch p.ConvertType {
				case entity.ConvertNone:
					// nothing
				case entity.ConvertRename:
					if v, ok := row[p.Name]; ok {
						row[p.ConvertValue] = v
						delete(row, p.Name)
					}
				}
			}
		}
	}
	return pagination, nil
}

// template 包中定义的函数
var funcNames = []string{
	"and",
	"call",
	"html",
	"index",
	"slice",
	"js",
	"len",
	"not",
	"or",
	"print",
	"printf",
	"println",
	"urlquery",
	"eq",
	"ge",
	"gt",
	"le",
	"lt",
	"ne",
	// 自定义
	"pl",
	"pt",
}

// pl 自定义template转换方法，将参数类型转中文
func pl(paramType entity.ParamLocation) string {
	switch paramType {
	case entity.ParamPath:
		return "Path"
	case entity.ParamQuery:
		return "Query"
	case entity.ParamBody:
		return "Body"
	}
	return ""
}

// pl 自定义template转换方法，将参数类型转中文
func pt(paramType entity.ParamType) string {
	switch paramType {
	case entity.Boolean:
		return "Boolean"
	case entity.Int:
		return "Int"
	case entity.Long:
		return "Long"
	case entity.Float:
		return "Float"
	case entity.Double:
		return "Double"
	case entity.DateTime:
		return "DateTime"
	case entity.String:
		return "String"
	case entity.Object:
		return "Object"
	case entity.Array:
		return "Array"
	}
	return ""
}

// Node 节点
type Node struct {
	sync.RWMutex
	// 唯一
	ID string `json:"id"`
	// 命名参数
	Name string `json:"name"`
	// 全路径
	Path string `json:"path"`
	// Handle 处理
	Handle interface{}
	// 子节点
	Children []*Node `json:"children"`
}

// Add 添加路由
func (n *Node) Add(path string, handle interface{}) error {
	if path == "" {
		return errors.New("path不能为空")
	}
	path = strings.TrimPrefix(path, "/")
	names := strings.Split(path, "/")
	var (
		node = n
		err  error
	)
	for _, name := range names {
		node, err = node.findChild(name)
		if err != nil {
			return err
		}
	}
	if node.Path == "" {
		node.Path = path
		node.Handle = handle
		return nil
	}
	return errors.New("路径重复注册")
}

// Remove 移除路由
func (n *Node) Remove(path string) error {
	if path == "" {
		return errors.New("path不能为空")
	}
	path = strings.TrimPrefix(path, "/")
	names := strings.Split(path, "/")
	var (
		node = n
		err  error
	)
	for _, name := range names {
		node, err = node.findChild(name)
		if err != nil {
			return err
		}
	}
	if len(node.Children) == 0 {
		node = nil
	} else {
		node.Path = ""
		node.Handle = nil
	}
	return nil
}

// Match 匹配
func (n *Node) Match(path string) (*Node, map[string]string, error) {
	if path == "" {
		return nil, nil, errors.New("path不能为空")
	}
	path = strings.TrimPrefix(path, "/")
	names := strings.Split(path, "/")
	var (
		node         = n
		namingParams = make(map[string]string)
	)
	for _, name := range names {
		if name == "" {
			return nil, nil, errors.New("路径中不能包含空字符串，即/some//path")
		}
		if len(node.Children) == 0 {
			return nil, nil, nil
		}
		// 寻找该层级匹配的路径
		var matchNode *Node
		for _, child := range node.Children {
			if child.ID == name {
				matchNode = child
			}
			if child.ID == "*" {
				matchNode = child
				namingParams[child.Name] = name
			}
		}
		node = matchNode
		if node == nil {
			return nil, nil, nil
		}
	}
	return node, namingParams, nil
}

func (n *Node) findChild(name string) (*Node, error) {
	if name == "" {
		return nil, errors.New("路径中不能包含空字符串，即/some//path")
	}
	// 命名参数处理
	var id string
	if strings.HasPrefix(name, ":") {
		id = "*"
		name = strings.TrimPrefix(name, ":")
	} else {
		id = name
	}

	if n.Children == nil {
		n.Children = make([]*Node, 0, 16)
	}
	for _, node := range n.Children {
		if node.ID == id {
			return node, nil
		}
	}
	node := &Node{
		ID:   id,
		Name: name,
	}
	n.Children = append(n.Children, node)
	return node, nil
}

// SyncDataSet 同步数据集
func SyncDataSet(dataSet *DataSet) {
	go func() {
		for {
			log.Logger().Debug("同步数据集")

			router := new(Node)

			list, err := dataSet.All(context.TODO())
			if err != nil {
				log.Logger().Warn("查询数据集错误", zap.Error(err))
			} else {
				for _, e := range list {
					if !e.PublishStatus {
						continue
					}
					if err := router.Add(e.Path, e.ID); err != nil {
						log.Logger().Warn("加载数据集错误", zap.String("path", e.Path), zap.Error(err))
					}
				}
			}
			dataSet.router = router

			time.Sleep(30 * time.Second)
		}
	}()
}
