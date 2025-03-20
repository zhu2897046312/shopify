package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"shopify/models"
	"shopify/service"
	"shopify/pkg/utils/response"
)

// CreatePayment 创建支付
// @Summary 创建支付
// @Description 为指定订单创建支付，支持微信和支付宝支付方式
// @Tags 支付
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body struct{OrderID uint "订单ID";Method string "支付方式(wechat/alipay)"} true "支付请求参数"
// @Success 200 {object} response.SuccessResponse{data=struct{payment_id uint,pay_url string}} "创建成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 400 {object} response.ErrorResponse "无效的请求参数"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /payments [post]
func CreatePayment(c *gin.Context) {
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	var req struct {
		OrderID uint   `json:"order_id" binding:"required"`
		Method  string `json:"method" binding:"required,oneof=wechat alipay"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("paymentService").(*service.PaymentService)
	payment, paymentURL, err := svc.CreatePayment(req.OrderID, req.Method)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"payment_id": payment.ID,
		"pay_url":    paymentURL,
	}))
}

// HandleWechatCallback 处理微信支付回调
// @Summary 处理微信支付回调
// @Description 处理来自微信支付平台的支付结果通知
// @Tags 支付
// @Accept xml
// @Produce xml
// @Success 200 {object} struct{return_code string,return_msg string} "处理成功"
// @Failure 400 {object} struct{return_code string,return_msg string} "无效的请求参数"
// @Router /payments/wechat/callback [post]
func HandleWechatCallback(c *gin.Context) {
	data := make(map[string]string)
	if err := c.ShouldBindXML(&data); err != nil {
		c.XML(http.StatusBadRequest, gin.H{"return_code": "FAIL", "return_msg": "Invalid request"})
		return
	}

	svc := c.MustGet("paymentService").(*service.PaymentService)
	if err := svc.HandleCallback(models.PaymentMethodWechat, data); err != nil {
		c.XML(http.StatusInternalServerError, gin.H{"return_code": "FAIL", "return_msg": err.Error()})
		return
	}

	c.XML(http.StatusOK, gin.H{"return_code": "SUCCESS", "return_msg": "OK"})
}

// HandleAlipayCallback 处理支付宝回调
// @Summary 处理支付宝回调
// @Description 处理来自支付宝平台的支付结果通知
// @Tags 支付
// @Accept form
// @Produce json
// @Success 200 {string} string "success"
// @Failure 400 {string} string "fail"
// @Router /payments/alipay/callback [post]
func HandleAlipayCallback(c *gin.Context) {
	data := c.Request.Form
	params := make(map[string]string)
	for k, v := range data {
		params[k] = v[0]
	}

	svc := c.MustGet("paymentService").(*service.PaymentService)
	if err := svc.HandleCallback(models.PaymentMethodAlipay, params); err != nil {
		c.String(http.StatusInternalServerError, "fail")
		return
	}

	c.String(http.StatusOK, "success")
}

// QueryPaymentStatus 查询支付状态
// @Summary 查询支付状态
// @Description 查询指定支付的当前状态
// @Tags 支付
// @Produce json
// @Security BearerAuth
// @Param id path int true "支付ID"
// @Success 200 {object} response.SuccessResponse{data=models.Payment} "查询成功"
// @Failure 401 {object} response.ErrorResponse "未授权"
// @Failure 404 {object} response.ErrorResponse "支付记录未找到"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /payments/{id}/status [get]
func QueryPaymentStatus(c *gin.Context) {
	paymentID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid payment ID"))
		return
	}

	svc := c.MustGet("paymentService").(*service.PaymentService)
	status, err := svc.QueryPaymentStatus(uint(paymentID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"status": status,
	}))
}