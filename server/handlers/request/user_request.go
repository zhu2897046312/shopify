package request

// RegisterUserRequest 用户注册请求
type RegisterUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Nickname string `json:"nickname"`
	Role     string `json:"role"`
	Avatar   string `json:"avatar"`
}

// LoginUserRequest 用户登录请求
type LoginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// ResetPasswordRequest 重置密码请求
type ResetPasswordRequest struct {
	Email       string `json:"email"`
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

// AddressRequest 用户地址请求
type AddressRequest struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Province  string `json:"province"`
	City      string `json:"city"`
	District  string `json:"district"`
	Street    string `json:"street"`
	PostCode  string `json:"post_code"`
}


