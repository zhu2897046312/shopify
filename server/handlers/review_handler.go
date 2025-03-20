package handlers

import (
	"net/http"
	"strconv"
	"shopify/models"
	"shopify/service"
	"shopify/handlers/request"
	"shopify/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

// CreateReview 创建评论
// @Summary 创建评论
// @Description 用户根据订单对商品进行评论
// @Tags 评论相关
// @Accept json
// @Produce json
// @Param review body request.ReviewUserRequest true "评论内容"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=models.Review} "创建成功"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /reviews [post]
func CreateReview(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	var req  request.ReviewUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	review := &models.Review{
		UserID:    userID.(uint),
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Content:   req.Content,
		Images:    req.Images,
	}

	svc := c.MustGet("reviewService").(*service.ReviewService)
	if err := svc.CreateReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(review))
}

// GetProductReviews 获取商品评论列表
// @Summary 获取商品评论列表
// @Description 获取指定商品的评论列表，支持分页查询
// @Tags 评论相关
// @Produce json
// @Param id path int true "商品ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.SuccessResponse{data=gin.H{"reviews":[]models.Review,"total":int,"rating_stats":models.RatingStats}} "获取成功"
// @Failure 400 {object} response.ErrorResponse "商品ID无效"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /products/{id}/reviews [get]
func GetProductReviews(c *gin.Context) {
	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid product ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	svc := c.MustGet("reviewService").(*service.ReviewService)
	reviews, total, err := svc.GetProductReviews(uint(productID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	// 获取评分统计
	stats, err := svc.GetProductRatingStats(uint(productID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"reviews":    reviews,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"rating_stats": stats,
	}))
}

// GetUserReviews 获取用户评论列表
// @Summary 获取用户评论列表
// @Description 获取当前用户所写的评论列表，支持分页查询
// @Tags 评论相关
// @Produce json
// @Security BearerAuth
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.SuccessResponse{data=gin.H{"reviews":[]models.Review,"total":int}} "获取成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /user/reviews [get]
func GetUserReviews(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	svc := c.MustGet("reviewService").(*service.ReviewService)
	reviews, total, err := svc.GetUserReviews(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"reviews":   reviews,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// DeleteReview 删除评论
// @Summary 删除评论
// @Description 删除指定的评论
// @Tags 评论相关
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.SuccessResponse "删除成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 400 {object} response.ErrorResponse "评论ID无效"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /reviews/{id} [delete]
func DeleteReview(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid review ID"))
		return
	}

	svc := c.MustGet("reviewService").(*service.ReviewService)
	if err := svc.DeleteReview(uint(reviewID), userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// AdminCreateReview 管理员创建评论
// @Summary 管理员创建评论
// @Description 管理员为用户创建评论
// @Tags 评论相关
// @Accept json
// @Produce json
// @Param review body request.ReviewUserRequest true "评论内容"
// @Security BearerAuth
// @Success 200 {object} response.SuccessResponse{data=models.Review} "创建成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "请求参数错误"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/reviews [post]
func AdminCreateReview(c *gin.Context) {
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	var req request.AdminReviewUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	review := &models.Review{
		UserID:    req.UserID,
		OrderID:   req.OrderID,
		ProductID: req.ProductID,
		Rating:    req.Rating,
		Content:   req.Content,
		Images:    req.Images,
	}

	svc := c.MustGet("reviewService").(*service.ReviewService)
	if err := svc.CreateReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(review))
}

// AdminDeleteReview 管理员删除评论
// @Summary 管理员删除评论
// @Description 管理员删除指定评论
// @Tags 评论相关
// @Security BearerAuth
// @Param id path int true "评论ID"
// @Success 200 {object} response.SuccessResponse "删除成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "评论ID无效"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/reviews/{id} [delete]
func AdminDeleteReview(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	reviewID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid review ID"))
		return
	}

	svc := c.MustGet("reviewService").(*service.ReviewService)
	if err := svc.AdminDeleteReview(uint(reviewID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// AdminReviewsListOfProduct 管理员获取商品评论列表
// @Summary 管理员获取商品评论列表
// @Description 获取指定商品的评论列表，支持分页查询，只有管理员能访问
// @Tags 评论相关
// @Produce json
// @Security BearerAuth
// @Param id path int true "商品ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.SuccessResponse{data=gin.H{"reviews":[]models.Review,"total":int}} "获取成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "商品ID无效"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/products/{id}/reviews [get]
func AdminReviewsListOfProduct(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	productID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid product ID"))
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	svc := c.MustGet("reviewService").(*service.ReviewService)
	reviews, total, err := svc.GetProductReviews(uint(productID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"reviews":   reviews,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}