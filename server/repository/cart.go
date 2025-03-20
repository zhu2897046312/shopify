package repository

import (
    "shopify/models"
    "gorm.io/gorm"
)

type CartRepository struct {
    *BaseRepository
}

func NewCartRepository(db *gorm.DB) *CartRepository {
    return &CartRepository{
        BaseRepository: NewBaseRepository(db),
    }
}

// Create 创建购物车项
func (r *CartRepository) Create(item *models.CartItem) error {
    return r.db.Create(item).Error
}

// GetUserCart 获取用户的购物车列表
func (r *CartRepository) GetUserCart(userID uint) ([]models.CartItem, error) {
    var items []models.CartItem
    err := r.db.Where("user_id = ?", userID).
        Preload("Product").
        Find(&items).Error
    return items, err
}

// GetCartItem 获取特定的购物车项
func (r *CartRepository) GetCartItem(userID, productID uint) (*models.CartItem, error) {
    var item models.CartItem
    err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).
        First(&item).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}

// GetCartItemByID 通过ID获取购物车项
func (r *CartRepository) GetCartItemByID(itemID uint) (*models.CartItem, error) {
    var item models.CartItem
    err := r.db.First(&item, itemID).Error
    if err != nil {
        return nil, err
    }
    return &item, nil
}

// UpdateQuantity 更新购物车项数量
func (r *CartRepository) UpdateQuantity(itemID uint, quantity int) error {
    return r.db.Model(&models.CartItem{}).
        Where("id = ?", itemID).
        Update("quantity", quantity).Error
}

// UpdateSelected 更新购物车项选中状态
func (r *CartRepository) UpdateSelected(itemID uint, selected bool) error {
    return r.db.Model(&models.CartItem{}).
        Where("id = ?", itemID).
        Update("selected", selected).Error
}

// Delete 删除购物车项
func (r *CartRepository) Delete(itemID uint) error {
    return r.db.Delete(&models.CartItem{}, itemID).Error
}

// ClearCart 清空用户购物车
func (r *CartRepository) ClearCart(userID uint) error {
    return r.db.Where("user_id = ?", userID).
        Delete(&models.CartItem{}).Error
}

// GetSelectedItems 获取用户选中的购物车项
func (r *CartRepository) GetSelectedItems(userID uint) ([]models.CartItem, error) {
    var items []models.CartItem
    err := r.db.Where("user_id = ? AND selected = ?", userID, true).
        Preload("Product").
        Find(&items).Error
    return items, err
}

// BatchUpdateSelected 批量更新购物车项选中状态
func (r *CartRepository) BatchUpdateSelected(userID uint, selected bool) error {
    return r.db.Model(&models.CartItem{}).
        Where("user_id = ?", userID).
        Select("selected").
        Update("selected", selected).Error
}

// SelectItems 选择部分购物车商品 false 为零值 不会触发update更新操作，需要使用select进行强制更新
func (r *CartRepository) SelectItems(userID uint, itemIDs []uint, selected bool) error {
    return r.db.Model(&models.CartItem{}).
        Where("user_id = ? AND id IN ?", userID, itemIDs).
        Select("selected").
        Update("selected", selected).Error
}
