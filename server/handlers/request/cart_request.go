package request

// 选择部分购物车商品
type CartItemsRequest struct {
	ItemIDs  []uint `json:"item_ids" binding:"required"`
	Selected bool   `json:"selected" binding:"required"`
}