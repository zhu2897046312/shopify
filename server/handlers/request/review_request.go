package request


type AdminReviewUserRequest struct {
	UserID      uint     `json:"user_id" binding:"required"`
	OrderID     uint     `json:"order_id" binding:"required"`
	ProductID   uint     `json:"product_id" binding:"required"`
	Rating      int      `json:"rating" binding:"required,min=1,max=5"`
	Content     string   `json:"content" binding:"required"`
	Images      []string `json:"images"`
}

type ReviewUserRequest struct {
	OrderID     uint     `json:"order_id" binding:"required"`
	ProductID   uint     `json:"product_id" binding:"required"`
	Rating      int      `json:"rating" binding:"required,min=1,max=5"`
	Content     string   `json:"content" binding:"required"`
	Images      []string `json:"images"`
}
