package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"shopify/pkg/utils/response"
	"gorm.io/gorm"
)

// HealthCheck 健康检查
// @Summary 系统健康检查
// @Description 检查系统各组件（包括数据库连接）的运行状态
// @Tags 系统
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=struct{status string,time time.Time,db string}} "系统正常"
// @Router /health [get]
func HealthCheck(c *gin.Context) {
	// 获取数据库状态
	db := c.MustGet("db").(*gorm.DB)
	sqlDB, err := db.DB()
	dbStatus := "connected"
	if err != nil || sqlDB.Ping() != nil {
		dbStatus = "disconnected"
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"status": "ok",
		"time":   time.Now(),
		"db":     dbStatus,
	}))
}