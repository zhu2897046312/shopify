package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"shopify/handlers/request"
	"shopify/models"
	"shopify/service"
	"shopify/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

// CreateProduct 创建产品(管理员)
// @Summary 创建产品
// @Description 由管理员创建新产品
// @Tags 产品管理
// @Produce json
// @Security BearerAuth
// @Param product body request.ProductRequest true "产品信息"
// @Success 200 {object} response.SuccessResponse{data=models.Product} "创建成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "请求参数无效"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/products [post]
func CreateProduct(c *gin.Context) {
    // 检查是否是管理员
    role, exists := c.Get("userRole")
    if !exists || role.(string) != "admin" {
        c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
        return
    }

    var req request.ProductRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
        return
    }

    // 创建产品对象
    product := &models.Product{
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        Category:    req.Category,
        Images:      req.Images,
        Tags:        req.Tags,
        Status:      req.Status,
    }

    svc := c.MustGet("productService").(*service.ProductService)
    if err := svc.CreateProduct(product); err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
        return
    }

    c.JSON(http.StatusOK, response.Success(product))
}

// UpdateProduct 更新产品(管理员)
// @Summary 更新产品信息
// @Description 由管理员更新现有产品信息
// @Tags 产品管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Param product body request.ProductRequest true "产品更新信息"
// @Success 200 {object} response.SuccessResponse{data=models.Product} "更新成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "无效的产品ID"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/products/{id} [put]
func UpdateProduct(c *gin.Context) {
    // 检查是否是管理员
    role, exists := c.Get("userRole")
    if !exists || role.(string) != "admin" {
        c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
        return
    }

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, response.Error(400, "Invalid product ID"))
        return
    }

    // 使用单独的请求结构体，只包含可更新的字段
    var req request.ProductRequest

    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
        return
    }

    // 创建产品对象，只设置需要更新的字段
    product := &models.Product{
        ID:          uint(id),
        Name:        req.Name,
        Description: req.Description,
        Price:       req.Price,
        Stock:       req.Stock,
        Rating:      req.Rating,
        Category:    req.Category,
        Images:      req.Images,
        Tags:        req.Tags,
        Status:      req.Status,
    }

    // 确保数组字段不为 nil
    if product.Images == nil {
        product.Images = make([]string, 0)
    }
    if product.Tags == nil {
        product.Tags = make([]string, 0)
    }

    svc := c.MustGet("productService").(*service.ProductService)
    if err := svc.UpdateProduct(product); err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
        return
    }

    // 获取更新后的完整产品信息
    updatedProduct, err := svc.GetProduct(uint(id))
    if err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
        return
    }

    c.JSON(http.StatusOK, response.Success(updatedProduct))
}

// DeleteProduct 删除产品(管理员)
// @Summary 删除产品
// @Description 由管理员删除产品
// @Tags 产品管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "产品ID"
// @Success 200 {object} response.SuccessResponse{data=nil} "删除成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "无效的产品ID"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/products/{id} [delete]
func DeleteProduct(c *gin.Context) {
    // 检查是否是管理员
    role, exists := c.Get("userRole")
    if !exists || role.(string) != "admin" {
        c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
        return
    }

    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, response.Error(400, "Invalid product ID"))
        return
    }

    svc := c.MustGet("productService").(*service.ProductService)
    if err := svc.DeleteProduct(uint(id)); err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
        return
    }

    c.JSON(http.StatusOK, response.Success(nil))
}

// ListProducts 获取产品列表
// @Summary 获取产品列表
// @Description 获取产品列表并支持分页、筛选
// @Tags 产品管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param category query string false "类别"
// @Param min_price query float false "最小价格"
// @Param max_price query float false "最大价格"
// @Param tags query string false "标签"
// @Param keyword query string false "搜索关键字"
// @Success 200 {object} response.SuccessResponse{data=gin.H{"items":[]models.Product, "total":int, "page":int, "page_size":int}} "获取成功"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /products [get]
func ListProducts(c *gin.Context) {
    page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
    fmt.Println(page, pageSize)
    
    category := c.Query("category")
    minPrice, _ := strconv.ParseFloat(c.Query("min_price"), 64)
    maxPrice, _ := strconv.ParseFloat(c.Query("max_price"), 64)
    tags := c.Query("tags")
    keyword := c.Query("keyword")

    svc := c.MustGet("productService").(*service.ProductService)
    var products []models.Product
    var total int64
    var err error

    // 根据不同的查询参数调用不同的服务方法
    if category != "" {
        products, total, err = svc.ListProductsByCategory(category, page, pageSize)
    } else if minPrice > 0 || maxPrice > 0 {
        products, total, err = svc.ListProductsByPriceRange(minPrice, maxPrice, page, pageSize)
    } else if tags != "" {
        tagList := strings.Split(tags, ",")
        products, total, err = svc.ListProductsByTags(tagList, page, pageSize)
    } else if keyword != "" {
        products, total, err = svc.SearchProducts(keyword, page, pageSize)
    } else {
        products, total, err = svc.ListProducts(page, pageSize)
    }

    if err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
        return
    }

    c.JSON(http.StatusOK, response.Success(gin.H{
        "items": products,
        "total": total,
        "page": page,
        "page_size": pageSize,
    }))
}

// GetProduct 获取产品详情
// @Summary 获取产品详情
// @Description 获取单个产品的详细信息
// @Tags 产品管理
// @Produce json
// @Param id path int true "产品ID"
// @Success 200 {object} response.SuccessResponse{data=models.Product} "获取成功"
// @Failure 400 {object} response.ErrorResponse "无效的产品ID"
// @Failure 404 {object} response.ErrorResponse "产品未找到"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /products/{id} [get]
func GetProduct(c *gin.Context) {
    id, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, response.Error(400, "Invalid product ID"))
        return
    }

    svc := c.MustGet("productService").(*service.ProductService)
    product, err := svc.GetProduct(uint(id))
    if err != nil {
        c.JSON(http.StatusNotFound, response.Error(404, "Product not found"))
        return
    }

    c.JSON(http.StatusOK, response.Success(product))
} 