package handlers

import (
	"net/http"
	"strconv"
	"time"
	"shopify/handlers/request"
	"shopify/models"
	"shopify/service"
	"shopify/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

// CreateOrder 创建订单
// @Summary 创建订单
// @Description 创建一个新的订单，用户选择商品并提供送货地址。
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param order_items body []request.OrderItemRequest true "订单商品列表"
// @Param address_id body uint true "送货地址ID"
// @Success 200 {object} response.SuccessResponse{data=models.Order} "创建成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	var req request.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("orderService").(*service.OrderService)
	order, err := svc.CreateOrder(userID.(uint), req.OrderItems, req.AddressID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(order))
}

// GetOrder 获取订单详情
// @Summary 获取订单详情
// @Description 获取指定订单的详细信息。
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} response.SuccessResponse{data=models.Order} "获取成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 404 {object} response.ErrorResponse "订单未找到"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /orders/{id} [get]
func GetOrder(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid order ID"))
		return
	}

	svc := c.MustGet("orderService").(*service.OrderService)
	order, err := svc.GetOrder(uint(orderID), userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "Order not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(order))
}

// ListOrders 获取用户订单列表
// @Summary 获取用户订单列表
// @Description 获取当前用户的订单列表，可以分页查询。
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.SuccessResponse{data=object} "获取成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /orders [get]
func ListOrders(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	svc := c.MustGet("orderService").(*service.OrderService)
	orders, total, err := svc.ListUserOrders(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"orders":    orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// AdminListOrders 管理员查看订单列表
// @Summary 管理员查看订单列表
// @Description 通过订单状态查询所有订单，管理员可分页查看。
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param status query string false "订单状态"
// @Success 200 {object} response.SuccessResponse{data=object} "获取成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/orders [get]
func AdminListOrders(c *gin.Context) {
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	status := c.Query("status")

	svc := c.MustGet("orderService").(*service.OrderService)
	orders, total, err := svc.ListOrdersByStatus(status, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"orders":    orders,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// AdminGetOrder 管理员查看订单详情
// @Summary 管理员查看订单详情
// @Description 管理员查看指定订单的详细信息
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} response.SuccessResponse{data=models.Order} "获取成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 404 {object} response.ErrorResponse "订单未找到"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/orders/{id} [get]
func AdminGetOrder(c *gin.Context) {
	// 获取订单ID
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid order ID"))
		return
	}

	// 获取订单服务
	svc := c.MustGet("orderService").(*service.OrderService)

	// 获取订单详情
	order, err := svc.GetOrderByID(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "Order not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(order))
}

// UpdateOrderStatus 更新订单状态
// @Summary 更新订单状态
// @Description 管理员可以更新订单的状态。
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Param status body string true "订单状态"
// @Success 200 {object} response.SuccessResponse{data=nil} "更新成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 404 {object} response.ErrorResponse "订单未找到"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	// role, exists := c.Get("userRole")
	// if !exists || role.(string) != "admin" {
	// 	c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
	// 	return
	// }

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid order ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("orderService").(*service.OrderService)
	if err := svc.UpdateOrderStatus(uint(orderID), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// UpdateLogistics 更新物流信息
// @Summary 更新物流信息
// @Description 管理员更新订单的物流信息。
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Param logistics body models.Logistics true "物流信息"
// @Success 200 {object} response.SuccessResponse{data=models.Logistics} "更新成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 404 {object} response.ErrorResponse "订单未找到"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /orders/{id}/logistics [put]
func UpdateLogistics(c *gin.Context) {
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid order ID"))
		return
	}

	var logistics models.Logistics
	if err := c.ShouldBindJSON(&logistics); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	logistics.OrderID = uint(orderID)

	svc := c.MustGet("orderService").(*service.OrderService)

	_, err = svc.GetOrder(uint(orderID), 0)
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "Order not found"))
		return
	}

	if err := svc.UpdateLogistics(&logistics); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(logistics))
}

// GetLogistics 获取订单物流信息
// @Summary 获取订单物流信息
// @Description 获取指定订单的物流信息。
// @Tags 订单
// @Produce json
// @Param id path int true "订单ID"
// @Success 200 {object} response.SuccessResponse{data=models.Logistics} "获取成功"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 404 {object} response.ErrorResponse "物流信息未找到"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /orders/{id}/logistics [get]
func GetLogistics(c *gin.Context) {
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid order ID"))
		return
	}

	svc := c.MustGet("orderService").(*service.OrderService)
	logistics, err := svc.GetLogistics(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "Logistics not found"))
		return
	}

	c.JSON(http.StatusOK, response.Success(logistics))
}

// AddLogisticsTrace 添加物流跟踪记录
// @Summary 添加物流跟踪记录
// @Description 管理员可以添加物流跟踪记录。
// @Tags 订单
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Param trace body models.LogisticsTrace true "物流跟踪记录"
// @Success 200 {object} response.SuccessResponse{data=models.LogisticsTrace} "添加成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 404 {object} response.ErrorResponse "物流信息未找到"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /orders/{id}/logistics/trace [post]
func AddLogisticsTrace(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	// 获取订单ID
	orderID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid order ID"))
		return
	}

	var req struct {
		Status      string    `json:"status" binding:"required"`
		Location    string    `json:"location" binding:"required"`
		Description string    `json:"description"`
		TraceTime   time.Time `json:"trace_time"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("orderService").(*service.OrderService)

	// 先获取物流信息
	logistics, err := svc.GetLogistics(uint(orderID))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "Logistics not found"))
		return
	}

	// 创建物流跟踪记录
	trace := &models.LogisticsTrace{
		LogisticsID: logistics.ID,  // 使用获取到的物流信息ID
		Location:    req.Location,
		Status:      req.Status,
		Description: req.Description,
		TraceTime:   req.TraceTime,
	}

	if err := svc.AddLogisticsTrace(trace); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(trace))
}