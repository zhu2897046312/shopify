package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"shopify/models"
	"shopify/service"
	"shopify/pkg/utils/response"
)

// CreateAdvertisement 创建广告
// @Summary 创建广告
// @Description 由管理员创建广告
// @Tags 广告管理
// @Produce json
// @Security BearerAuth
// @Param advertisement body models.Advertisement true "广告信息"
// @Success 200 {object} response.SuccessResponse{data=models.Advertisement} "创建成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/advertisements [post]
func CreateAdvertisement(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	var ad models.Advertisement
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("advertisementService").(*service.AdvertisementService)
	if err := svc.CreateAd(&ad); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(ad))
}

// GetAdvertisement 获取广告详情
// @Summary 获取广告详情
// @Description 获取单个广告的详细信息
// @Tags 广告管理
// @Produce json
// @Param id path int true "广告ID"
// @Success 200 {object} response.SuccessResponse{data=models.Advertisement} "获取成功"
// @Failure 400 {object} response.ErrorResponse "无效的广告ID"
// @Failure 404 {object} response.ErrorResponse "广告未找到"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /advertisements/{id} [get]
func GetAdvertisement(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid advertisement ID"))
		return
	}

	svc := c.MustGet("advertisementService").(*service.AdvertisementService)
	ad, err := svc.GetAd(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "Advertisement not found"))
		return
	}
	
	c.JSON(http.StatusOK, response.Success(ad))
}

// UpdateAdvertisement 更新广告信息
// @Summary 更新广告信息
// @Description 由管理员更新广告信息
// @Tags 广告管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "广告ID"
// @Param advertisement body models.Advertisement true "广告更新信息"
// @Success 200 {object} response.SuccessResponse{data=models.Advertisement} "更新成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "无效的广告ID"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/advertisements/{id} [put]
func UpdateAdvertisement(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid advertisement ID"))
		return
	}

	var ad models.Advertisement
	if err := c.ShouldBindJSON(&ad); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}
	ad.ID = uint(id)

	svc := c.MustGet("advertisementService").(*service.AdvertisementService)
	if err := svc.UpdateAd(&ad); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(ad))
}

// DeleteAdvertisement 删除广告
// @Summary 删除广告
// @Description 由管理员删除广告
// @Tags 广告管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "广告ID"
// @Success 200 {object} response.SuccessResponse{data=nil} "删除成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "无效的广告ID"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/advertisements/{id} [delete]
func DeleteAdvertisement(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid advertisement ID"))
		return
	}

	svc := c.MustGet("advertisementService").(*service.AdvertisementService)
	if err := svc.DeleteAd(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}


// ListAdvertisements 获取广告列表
// @Summary 获取广告列表
// @Description 获取广告列表并支持分页
// @Tags 广告管理
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Success 200 {object} response.SuccessResponse{data=gin.H{"items":[]models.Advertisement, "total":int, "page":int, "page_size":int}} "获取成功"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /advertisements [get]
func ListAdvertisements(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	svc := c.MustGet("advertisementService").(*service.AdvertisementService)
	ads, total, err := svc.ListAllAds(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"items":     ads,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}))
}

// GetActiveAdvertisements 获取有效广告
// @Summary 获取有效广告
// @Description 获取当前有效的广告，可以按位置过滤
// @Tags 广告管理
// @Produce json
// @Param position query string false "广告位置"
// @Success 200 {object} response.SuccessResponse{data=[]models.Advertisement} "获取成功"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /advertisements/active [get]
func GetActiveAdvertisements(c *gin.Context) {
	position := c.Query("position")
	svc := c.MustGet("advertisementService").(*service.AdvertisementService)

	var ads []models.Advertisement
	var err error

	if position != "" {
		ads, err = svc.GetActiveAds(position)
	} else {
		ads, err = svc.ListActiveAds()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(ads))
}

// UpdateAdvertisementStatus 更新广告状态
// @Summary 更新广告状态
// @Description 由管理员更新广告状态（例如：启用/禁用）
// @Tags 广告管理
// @Produce json
// @Security BearerAuth
// @Param id path int true "广告ID"
// @Param status body struct{ Status string `json:"status" binding:"required"` } true "广告状态"
// @Success 200 {object} response.SuccessResponse{data=nil} "更新成功"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 400 {object} response.ErrorResponse "无效的广告ID"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /admin/advertisements/{id}/status [put]
func UpdateAdvertisementStatus(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid advertisement ID"))
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("advertisementService").(*service.AdvertisementService)
	if err := svc.UpdateAdStatus(uint(id), req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
} 