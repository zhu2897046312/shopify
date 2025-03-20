package handlers

import (
	"log"
	"net/http"
	"strconv"
	"shopify/handlers/request"
	"shopify/models"
	"shopify/service"
	"shopify/pkg/utils/response"

	"github.com/gin-gonic/gin"
)

// RegisterUser 用户注册
// @Summary 用户注册
// @Description 用户通过提供电子邮件、密码、昵称和头像进行注册。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body request.RegisterUserRequest true "注册信息"
// @Success 200 {object} response.SuccessResponse{data=models.User} "注册成功"
// @Failure 400 {object} response.ErrorResponse "请求参数无效"
// @Router /api/v1/users/register [post]
func RegisterUser(c *gin.Context) {
	var req request.RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	user := models.User{
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
	}
	// 清除敏感信息
	req.Password = ""

	svc := c.MustGet("userService").(*service.UserService)
	if err := svc.Register(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	// 清除敏感信息
	user.Password = ""
	c.JSON(http.StatusOK, response.Success(user))
}

// LoginUser 用户登录
// @Summary 用户登录
// @Description 用户通过电子邮件和密码进行登录。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body request.LoginUserRequest true "登录信息"
// @Success 200 {object} response.SuccessResponse{data=models.LoginResponse} "登录成功"
// @Failure 400 {object} response.ErrorResponse "请求参数无效"
// @Failure 401 {object} response.ErrorResponse "未经授权"
// @Router /api/v1/users/login [post]
func LoginUser(c *gin.Context) {
	var req request.LoginUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("userService").(*service.UserService)
	loginResp, err := svc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(401, err.Error()))
		return
	}

	// 清除敏感信息
	req.Password = ""
	loginResp.User.Password = ""
	c.JSON(http.StatusOK, response.Success(loginResp))
}

// GetUserProfile 获取用户信息
// @Summary 获取用户个人资料
// @Description 获取当前用户的个人资料。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=models.User} "成功获取用户资料"
// @Failure 401 {object} response.ErrorResponse "未经授权"
// @Failure 404 {object} response.ErrorResponse "未找到用户"
// @Router /api/v1/users/profile [get]
func GetUserProfile(c *gin.Context) {
	// 从上下文中获取用户ID（由 AuthMiddleware 设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	svc := c.MustGet("userService").(*service.UserService)
	user, err := svc.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusNotFound, response.Error(404, "User not found"))
		return
	}

	// 清除敏感信息
	user.Password = ""
	c.JSON(http.StatusOK, response.Success(user))
}

// UpdateUserProfile 更新用户信息
// @Summary 更新用户个人资料
// @Description 更新用户的昵称和头像。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param user body request.UpdateUserRequest true "更新信息"
// @Success 200 {object} response.SuccessResponse{data=models.User} "成功更新用户资料"
// @Failure 400 {object} response.ErrorResponse "请求参数无效"
// @Failure 401 {object} response.ErrorResponse "未经授权"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /api/v1/users/profile [put]
func UpdateUserProfile(c *gin.Context) {
	// 从上下文中获取用户ID（由 AuthMiddleware 设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	var updateData request.UpdateUserRequest

	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	// 创建用户对象，只包含要更新的字段
	user := &models.User{
		ID:       userID.(uint),
		Nickname: updateData.Nickname,
		Avatar:   updateData.Avatar,
	}

	svc := c.MustGet("userService").(*service.UserService)
	if err := svc.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	// 获取更新后的用户信息
	updatedUser, err := svc.GetUserByID(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	// 清除敏感信息
	updatedUser.Password = ""
	c.JSON(http.StatusOK, response.Success(updatedUser))
}

// AddUserAddress 添加用户地址
// @Summary 添加用户地址
// @Description 添加一个新地址到用户的地址列表。
// @Tags 地址管理
// @Accept json
// @Produce json
// @Param address body request.AddressRequest true "地址信息"
// @Success 200 {object} response.SuccessResponse{data=models.Address} "成功添加地址"
// @Failure 400 {object} response.ErrorResponse "请求参数无效"
// @Failure 401 {object} response.ErrorResponse "未经授权"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /api/v1/users/addresses [post]
func AddUserAddress(c *gin.Context) {
	// 从上下文中获取用户ID（由 AuthMiddleware 设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	var req request.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	address := models.Address{
		Name:        req.Name,
		Phone:       req.Phone,
		Province:    req.Province,
		City:        req.City,
		District:    req.District,
		Street:      req.Street,
		PostCode:    req.PostCode,
	}

	// 确保只能为自己添加地址
	address.UserID = userID.(uint)
	

	svc := c.MustGet("userService").(*service.UserService)
	if err := svc.AddAddress(&address); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(address))
}



func UpdateUserAddress(c *gin.Context) {
	// 从上下文中获取用户ID（由 AuthMiddleware 设置）
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	
	addressId, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, response.Error(400, "Invalid address ID"))
        return
    }

	var req request.AddressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	address := models.Address{
		ID:          uint(addressId),
		Name:        req.Name,
		Phone:       req.Phone,
		Province:    req.Province,
		City:        req.City,
		District:    req.District,
		Street:      req.Street,
		PostCode:    req.PostCode,
	}
	
	log.Print("IsDefault :  " , address.IsDefault)

	svc := c.MustGet("userService").(*service.UserService)
    if err := svc.UpdateAddress(&address); err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
        return
    }

    c.JSON(http.StatusOK, response.Success(nil))
}

func SetDefaultAddresses(c *gin.Context) {
	// 从上下文中获取用户ID（由 AuthMiddleware 设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	// 获取用户服务
	svc := c.MustGet("userService").(*service.UserService)


	addressId, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, response.Error(400, "Invalid address ID"))
        return
    }

	// 1. 取消所有地址的默认标记
	if err := svc.UnsetDefaultAddresses(userID.(uint)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	// 2. 设置新地址为默认地址
	if err := svc.SetDefaultAddress(uint(addressId)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, response.Success(nil))
}

// DeleteUserAddress 删除用户地址
// @Summary 删除用户地址
// @Description 根据地址ID删除用户地址。
// @Tags 地址管理
// @Accept json
// @Produce json
// @Param id path int true "地址ID"
// @Success 200 {object} response.SuccessResponse "成功删除地址"
// @Failure 400 {object} response.ErrorResponse "无效的地址ID"
// @Failure 500 {object} response.ErrorResponse "服务器内部错误"
// @Router /api/v1/users/addresses/{id} [delete]
func DeleteUserAddress(c *gin.Context) {
	// 从上下文中获取用户ID（由 AuthMiddleware 设置）
	_, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	addressId, err := strconv.ParseUint(c.Param("id"), 10, 32)
    if err != nil {
        c.JSON(http.StatusBadRequest, response.Error(400, "Invalid address ID"))
        return
    }

	svc := c.MustGet("userService").(*service.UserService)
    if err := svc.DeleteAddress(uint(addressId)); err != nil {
        c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
        return
    }

    c.JSON(http.StatusOK, response.Success(nil))
}

// ListUserAddresses 获取用户地址列表
// @Summary 获取用户地址列表
// @Description 获取当前用户所有的地址。
// @Tags 地址管理
// @Accept json
// @Produce json
// @Success 200 {object} response.SuccessResponse{data=[]models.Address} "成功获取地址列表"
// @Failure 401 {object} response.ErrorResponse "未经授权"
// @Router /api/v1/users/addresses [get]
func ListUserAddresses(c *gin.Context) {
	// 从上下文中获取用户ID（由 AuthMiddleware 设置）
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, response.Error(401, "Unauthorized"))
		return
	}

	svc := c.MustGet("userService").(*service.UserService)
	addresses, err := svc.ListAddresses(userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(addresses))
}

// AdminListUsers 管理员获取用户列表
// @Summary 获取用户列表（管理员）
// @Description 该接口用于获取用户列表，管理员可以根据分页参数查询用户。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param page query int false "当前页，默认为1"
// @Param page_size query int false "每页显示的数量，默认为10"
// @Success 200 {object} response.SuccessResponse{data=gin.H{users=[]models.User,total=int,page=int,page_size=int}} "成功获取用户列表"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /api/v1/admin/users [get]
func AdminListUsers(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	// 获取分页参数
	page := c.DefaultQuery("page", "1")
	pageSize := c.DefaultQuery("page_size", "10")
	
	pageNum, _ := strconv.Atoi(page)
	pageSizeNum, _ := strconv.Atoi(pageSize)

	svc := c.MustGet("userService").(*service.UserService)
	users, total, err := svc.ListUsers(pageNum, pageSizeNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	// 清除敏感信息
	for i := range users {
		users[i].Password = ""
		users[i].VerifyToken = ""
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"users": users,
		"total": total,
		"page": pageNum,
		"page_size": pageSizeNum,
	}))
}

// AdminUpdateUser 管理员更新用户信息
// @Summary 更新用户信息（管理员）
// @Description 该接口允许管理员更新用户的基本信息。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body models.User true "用户信息"
// @Success 200 {object} response.SuccessResponse{data=models.User} "成功更新用户信息"
// @Failure 400 {object} response.ErrorResponse "请求参数无效"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /api/v1/admin/users/{id} [put]
func AdminUpdateUser(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid user ID"))
		return
	}

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	user.ID = uint(userID)

	svc := c.MustGet("userService").(*service.UserService)
	if err := svc.UpdateUserByAdmin(&user); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, response.Success(user))
}

// AdminDeleteUser 管理员删除用户
// @Summary 删除用户（管理员）
// @Description 该接口允许管理员删除指定的用户。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} response.SuccessResponse "成功删除用户"
// @Failure 400 {object} response.ErrorResponse "无效的用户ID"
// @Failure 403 {object} response.ErrorResponse "权限不足"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /api/v1/admin/users/{id} [delete]
func AdminDeleteUser(c *gin.Context) {
	// 检查是否是管理员
	role, exists := c.Get("userRole")
	if !exists || role.(string) != "admin" {
		c.JSON(http.StatusForbidden, response.Error(403, "Permission denied"))
		return
	}

	userID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid user ID"))
		return
	}

	svc := c.MustGet("userService").(*service.UserService)
	if err := svc.DeleteUser(uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(nil))
}

// SendResetPasswordCode 发送重置密码的验证码
// @Summary 发送重置密码验证码
// @Description 该接口允许用户通过电子邮件发送验证码来重置密码。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param email body struct { Email string `json:"email" binding:"required,email"` } true "电子邮件"
// @Success 200 {object} response.SuccessResponse{data=gin.H{message=string}} "验证码已发送"
// @Failure 400 {object} response.ErrorResponse "无效的电子邮件地址"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /api/v1/users/reset-password/code [post]
func SendResetPasswordCode(c *gin.Context) {
	var req struct {
		Email string `json:"email" binding:"required,email"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid email format"))
		return
	}

	svc := c.MustGet("userService").(*service.UserService)
	if err := svc.SendResetPasswordCode(req.Email); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"message": "Verification code has been sent to your email",
	}))
}

// ResetPasswordByCode 通过验证码重置密码
// @Summary 通过验证码重置密码
// @Description 该接口允许用户通过验证码来重置密码。
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param email body struct { Email string `json:"email" binding:"required,email"` } true "电子邮件"
// @Param code body struct { Code string `json:"code" binding:"required"` } true "验证码"
// @Param new_password body struct { NewPassword string `json:"new_password" binding:"required,min=6"` } true "新密码"
// @Success 200 {object} response.SuccessResponse{data=gin.H{message=string}} "密码重置成功"
// @Failure 400 {object} response.ErrorResponse "请求参数无效"
// @Failure 500 {object} response.ErrorResponse "服务器错误"
// @Router /api/v1/users/reset-password [post]
func ResetPasswordByCode(c *gin.Context) {
	var req struct {
		Email       string `json:"email" binding:"required,email"`
		Code        string `json:"code" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "Invalid request parameters"))
		return
	}

	svc := c.MustGet("userService").(*service.UserService)
	if err := svc.ResetPasswordByCode(req.Email, req.Code, req.NewPassword); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"message": "Password has been reset successfully",
	}))
}

