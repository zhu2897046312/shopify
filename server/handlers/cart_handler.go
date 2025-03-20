package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"shopify/service"
	"shopify/pkg/utils/response"
)

// AddCartItem 添加商品到购物车
// @Summary 添加商品到购物车
// @Description 将商品添加到当前用户的购物车
// @Tags 购物车
// @Produce json
// @Security BearerAuth
// @Param product_id body uint true "商品ID"
// @Param quantity body int true "商品数量"
// @Success 200 {object} response.SuccessResponse{data=nil} "添加成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /cart/items [post]
func AddCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	var req struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("cartService").(*service.CartService)
	if err := svc.AddItem(userID.(uint), req.ProductID, req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// ListCartItems 获取购物车列表
// @Summary 获取购物车列表
// @Description 获取当前用户的购物车内所有商品
// @Tags 购物车
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]models.CartItem} "获取成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /cart/items [get]
func ListCartItems(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	svc := c.MustGet("cartService").(*service.CartService)
	items, err := svc.ListItems(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(items))
}

// UpdateCartItem 更新购物车项数量
// @Summary 更新购物车项数量
// @Description 更新购物车中的某一项商品数量，或更新选中状态
// @Tags 购物车
// @Produce json
// @Security BearerAuth
// @Param id path int true "购物车项ID"
// @Param quantity body int true "商品数量"
// @Param selected body bool false "商品是否被选中"
// @Success 200 {object} response.SuccessResponse{data=nil} "更新成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /cart/items/{id} [put]
func UpdateCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid item ID"))
		return
	}

	var req struct {
		Quantity int  `json:"quantity" binding:"required,min=1"`
		Selected *bool `json:"selected,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("cartService").(*service.CartService)

	// 如果提供了 selected 字段，更新选中状态
	if req.Selected != nil {
		if err := svc.UpdateSelected(userID.(uint), uint(itemID), *req.Selected); err != nil {
			c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
			return
		}
	}

	// 更新数量
	if err := svc.UpdateQuantity(userID.(uint), uint(itemID), req.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// RemoveCartItem 删除购物车项
// @Summary 删除购物车项
// @Description 从购物车中删除一个商品
// @Tags 购物车
// @Produce json
// @Security BearerAuth
// @Param id path int true "购物车项ID"
// @Success 200 {object} response.SuccessResponse{data=nil} "删除成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /cart/items/{id} [delete]
func RemoveCartItem(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	itemID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid item ID"))
		return
	}

	svc := c.MustGet("cartService").(*service.CartService)
	if err := svc.RemoveItem(userID.(uint), uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// SelectAllCartItems 全选/取消全选购物车项
// @Summary 全选/取消全选购物车项
// @Description 设置购物车中的所有商品的选中状态
// @Tags 购物车
// @Produce json
// @Security BearerAuth
// @Param selected body bool true "是否全选"
// @Success 200 {object} response.SuccessResponse{data=nil} "更新成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /cart/items/select-all [put]
func SelectAllCartItems(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	var req struct {
		Selected bool `json:"selected" `
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("cartService").(*service.CartService)
	if err := svc.SelectAll(userID.(uint), req.Selected); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// GetSelectedCartItems 获取选中的购物车项
// @Summary 获取选中的购物车项
// @Description 获取当前用户购物车中已选中的所有商品
// @Tags 购物车
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=[]models.CartItem} "获取成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /cart/items/selected [get]
func GetSelectedCartItems(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	svc := c.MustGet("cartService").(*service.CartService)
	items, err := svc.GetSelectedItems(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(items))
}

// SelectCartItems 选择部分购物车商品
// @Summary 选择部分购物车商品
// @Description 设置指定购物车商品的选中状态
// @Tags 购物车
// @Produce json
// @Security BearerAuth
// @Param request body struct{ItemIDs []uint "商品ID列表";Selected bool "是否选中"} true "请求参数"
// @Success 200 {object} response.SuccessResponse{data=nil} "更新成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /cart/select-items [put]
func SelectCartItems(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	var req struct {
		ItemIDs  []uint `json:"item_ids" binding:"required"`
		Selected bool   `json:"selected" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("cartService").(*service.CartService)
	if err := svc.SelectItems(userID.(uint), req.ItemIDs, req.Selected); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}