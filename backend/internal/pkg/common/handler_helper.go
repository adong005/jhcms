package common

import (
	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
)

// PagedRequest 通用分页请求
type PagedRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
}

// ParseIDParam 解析路由中的 UUID 主键
func ParseIDParam(c *gin.Context) (string, error) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		return "", errors.New("invalid id")
	}
	return idStr, nil
}

func normalizePagination(req interface{}, defaultPageSize int) {
	if defaultPageSize <= 0 {
		defaultPageSize = 10
	}

	v := reflect.ValueOf(req)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return
	}

	elem := v.Elem()
	if elem.Kind() != reflect.Struct {
		return
	}

	pageField := elem.FieldByName("Page")
	if pageField.IsValid() && pageField.CanSet() && pageField.Kind() == reflect.Int && pageField.Int() == 0 {
		pageField.SetInt(1)
	}

	pageSizeField := elem.FieldByName("PageSize")
	if pageSizeField.IsValid() && pageSizeField.CanSet() && pageSizeField.Kind() == reflect.Int && pageSizeField.Int() == 0 {
		pageSizeField.SetInt(int64(defaultPageSize))
	}
}

// HandleListRequest 处理通用列表请求（统一参数解析、统一分页、统一响应）
func HandleListRequest(c *gin.Context, req interface{}, defaultPageSize int, fetchFunc func() (interface{}, int64, error), errorMsg string) {
	if err := c.ShouldBindJSON(req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	normalizePagination(req, defaultPageSize)

	items, total, err := fetchFunc()
	if err != nil {
		if errorMsg == "" {
			errorMsg = "获取列表失败"
		}
		response.Error(c, errorMsg)
		return
	}

	response.PageSuccess(c, items, total)
}

// HandleTreeListRequest 处理树形列表请求（与列表页保持统一返回格式）
func HandleTreeListRequest(c *gin.Context, req interface{}, defaultPageSize int, fetchFunc func() (interface{}, int64, error), errorMsg string) {
	HandleListRequest(c, req, defaultPageSize, fetchFunc, errorMsg)
}

// HandleAllRequest 处理 all 场景（统一成功/失败响应）
func HandleAllRequest(c *gin.Context, fetchFunc func() (interface{}, error), errorMsg string) {
	data, err := fetchFunc()
	if err != nil {
		if errorMsg == "" {
			errorMsg = "获取数据失败"
		}
		response.Error(c, errorMsg)
		return
	}

	response.Success(c, data)
}

// HandleDetailRequest 处理详情请求的通用逻辑
func HandleDetailRequest(c *gin.Context, fetchFunc func(id string) (interface{}, error)) {
	id, err := ParseIDParam(c)
	if err != nil {
		response.Error(c, "无效的ID")
		return
	}

	item, err := fetchFunc(id)
	if err != nil {
		response.Error(c, "获取详情失败")
		return
	}

	response.Success(c, item)
}

// HandleCreateRequest 处理创建请求的通用逻辑
func HandleCreateRequest(c *gin.Context, data interface{}, createFunc func(data interface{}) error) {
	if err := c.ShouldBindJSON(data); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	if err := createFunc(data); err != nil {
		response.Error(c, "创建失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", nil)
}

// HandleUpdateRequest 处理更新请求的通用逻辑
func HandleUpdateRequest(c *gin.Context, data interface{}, updateFunc func(id string, data interface{}) error) {
	id, err := ParseIDParam(c)
	if err != nil {
		response.Error(c, "无效的ID")
		return
	}

	if err := c.ShouldBindJSON(data); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	if err := updateFunc(id, data); err != nil {
		response.Error(c, "更新失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// HandleDeleteRequest 处理删除请求的通用逻辑
func HandleDeleteRequest(c *gin.Context, deleteFunc func(id string) error) {
	id, err := ParseIDParam(c)
	if err != nil {
		response.Error(c, "无效的ID")
		return
	}

	if err := deleteFunc(id); err != nil {
		response.Error(c, "删除失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// BuildTreeData 构建树形数据结构
func BuildTreeData(items []map[string]interface{}, idField, parentField string) []map[string]interface{} {
	itemMap := make(map[interface{}]map[string]interface{})
	var roots []map[string]interface{}

	// 创建映射
	for _, item := range items {
		id := item[idField]
		itemMap[id] = item
		item["children"] = []map[string]interface{}{}
	}

	// 构建树
	for _, item := range items {
		parentID := item[parentField]
		if parentID == nil {
			roots = append(roots, item)
		} else {
			if parent, exists := itemMap[parentID]; exists {
				children := parent["children"].([]map[string]interface{})
				parent["children"] = append(children, item)
			}
		}
	}

	return roots
}

// FilterItems 过滤数据
func FilterItems(items []map[string]interface{}, filterFunc func(map[string]interface{}) bool) []map[string]interface{} {
	var filtered []map[string]interface{}
	for _, item := range items {
		if filterFunc(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// PaginateItems 分页数据
func PaginateItems(items []map[string]interface{}, page, pageSize int) ([]map[string]interface{}, int64) {
	total := int64(len(items))
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(items) {
		return []map[string]interface{}{}, total
	}

	if end > len(items) {
		end = len(items)
	}

	return items[start:end], total
}
