package handler

import (
	"fmt"

	"adcms-backend/internal/pkg/ids"
	"adcms-backend/internal/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BaseRepository 基础仓库接口
type BaseRepository interface {
	FindByID(id string) (interface{}, error)
	List(page, pageSize int, filters map[string]interface{}) ([]interface{}, int64, error)
	Create(data interface{}) error
	Update(id string, data interface{}) error
	Delete(id string) error
}

// BaseHandler 基础处理器
type BaseHandler struct {
	repo BaseRepository
	db   *gorm.DB
}

// NewBaseHandler 创建基础处理器
func NewBaseHandler(repo BaseRepository, db *gorm.DB) *BaseHandler {
	return &BaseHandler{
		repo: repo,
		db:   db,
	}
}

// GetList 获取列表（通用方法）
func (h *BaseHandler) GetList(c *gin.Context, convertFunc func(interface{}) map[string]interface{}) {
	var req struct {
		Page     int                    `json:"page"`
		PageSize int                    `json:"pageSize"`
		Filters  map[string]interface{} `json:"filters"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	items, total, err := h.repo.List(req.Page, req.PageSize, req.Filters)
	if err != nil {
		response.Error(c, "获取列表失败")
		return
	}

	result := make([]map[string]interface{}, 0)
	for _, item := range items {
		if convertFunc != nil {
			result = append(result, convertFunc(item))
		}
	}

	response.PageSuccess(c, result, total)
}

// GetDetail 获取详情（通用方法）
func (h *BaseHandler) GetDetail(c *gin.Context, convertFunc func(interface{}) map[string]interface{}) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		response.Error(c, "无效的ID")
		return
	}

	item, err := h.repo.FindByID(idStr)
	if err != nil {
		response.Error(c, "获取详情失败")
		return
	}

	var result map[string]interface{}
	if convertFunc != nil {
		result = convertFunc(item)
	}

	response.Success(c, result)
}

// Create 创建（通用方法）
func (h *BaseHandler) Create(c *gin.Context, data interface{}) {
	if err := c.ShouldBindJSON(data); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	err := h.repo.Create(data)
	if err != nil {
		response.Error(c, "创建失败")
		return
	}

	response.SuccessWithMessage(c, "创建成功", nil)
}

// Update 更新（通用方法）
func (h *BaseHandler) Update(c *gin.Context, data interface{}) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		response.Error(c, "无效的ID")
		return
	}

	if err := c.ShouldBindJSON(data); err != nil {
		response.Error(c, "请求参数错误")
		return
	}

	err := h.repo.Update(idStr, data)
	if err != nil {
		response.Error(c, "更新失败")
		return
	}

	response.SuccessWithMessage(c, "更新成功", nil)
}

// Delete 删除（通用方法）
func (h *BaseHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	if !ids.Valid(idStr) {
		response.Error(c, "无效的ID")
		return
	}

	err := h.repo.Delete(idStr)
	if err != nil {
		response.Error(c, "删除失败")
		return
	}

	response.SuccessWithMessage(c, "删除成功", nil)
}

// BuildTree 构建树形结构（通用方法）；id / parentId 字段值为字符串 UUID。
func BuildTree(items []interface{}, idField, parentField string) []map[string]interface{} {
	itemMap := make(map[string]map[string]interface{})
	var rootItems []map[string]interface{}

	for _, item := range items {
		itemData := item.(map[string]interface{})
		idKey, ok := treeMapKey(itemData[idField])
		if !ok {
			continue
		}
		itemMap[idKey] = itemData
	}

	for _, item := range items {
		itemData := item.(map[string]interface{})
		parentID := itemData[parentField]

		if parentID == nil {
			rootItems = append(rootItems, itemData)
			continue
		}
		parentKey, ok := treeMapKey(parentID)
		if !ok {
			continue
		}
		if parent, exists := itemMap[parentKey]; exists {
			if parent["children"] == nil {
				parent["children"] = []map[string]interface{}{}
			}
			children := parent["children"].([]map[string]interface{})
			parent["children"] = append(children, itemData)
		}
	}

	return rootItems
}

func treeMapKey(v interface{}) (string, bool) {
	switch x := v.(type) {
	case string:
		if x == "" {
			return "", false
		}
		return x, true
	case nil:
		return "", false
	default:
		s := fmt.Sprint(x)
		if s == "" || s == "<nil>" {
			return "", false
		}
		return s, true
	}
}

// FilterList 过滤列表（通用方法）
func FilterList(items []interface{}, filterFunc func(interface{}) bool) []interface{} {
	var filtered []interface{}
	for _, item := range items {
		if filterFunc(item) {
			filtered = append(filtered, item)
		}
	}
	return filtered
}

// PaginateList 分页列表（通用方法）
func PaginateList(items []interface{}, page, pageSize int) ([]interface{}, int64) {
	total := int64(len(items))
	start := (page - 1) * pageSize
	end := start + pageSize

	if start >= len(items) {
		return []interface{}{}, total
	}

	if end > len(items) {
		end = len(items)
	}

	return items[start:end], total
}
