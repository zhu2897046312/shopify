package service

import (
    "errors"
    "shopify/models"
)

type CartService struct {
    *Service
}

func NewCartService(base *Service) *CartService {
    return &CartService{Service: base}
}

// AddItem 添加商品到购物车
func (s *CartService) AddItem(userID, productID uint, quantity int) error {
    // 检查商品是否存在
    product, err := s.repoFactory.GetProductRepository().GetByID(productID)
    if err != nil {
        return errors.New("product not found")
    }

    // 检查库存
    if product.Stock < quantity {
        return errors.New("insufficient stock")
    }

    // 检查购物车是否已有该商品
    existingItem, err := s.repoFactory.GetCartRepository().GetCartItem(userID, productID)
    if err == nil {
        // 更新数量
        newQuantity := existingItem.Quantity + quantity
        if product.Stock < newQuantity {
            return errors.New("insufficient stock")
        }
        return s.repoFactory.GetCartRepository().UpdateQuantity(existingItem.ID, newQuantity)
    }

    // 创建新的购物车项
    cartItem := &models.CartItem{
        UserID:    userID,
        ProductID: productID,
        Quantity:  quantity,
        Selected:  true,
    }

    return s.repoFactory.GetCartRepository().Create(cartItem)
}

// ListItems 获取购物车列表
func (s *CartService) ListItems(userID uint) ([]models.CartItem, error) {
    return s.repoFactory.GetCartRepository().GetUserCart(userID)
}

// UpdateQuantity 更新购物车项数量
func (s *CartService) UpdateQuantity(userID, itemID uint, quantity int) error {
    // 检查购物车项是否存在且属于该用户
    item, err := s.repoFactory.GetCartRepository().GetCartItemByID(itemID)
    if err != nil || item.UserID != userID {
        return errors.New("cart item not found")
    }

    // 检查库存
    product, err := s.repoFactory.GetProductRepository().GetByID(item.ProductID)
    if err != nil {
        return errors.New("product not found")
    }

    if product.Stock < quantity {
        return errors.New("insufficient stock")
    }

    return s.repoFactory.GetCartRepository().UpdateQuantity(itemID, quantity)
}

// UpdateSelected 更新购物车项选中状态
func (s *CartService) UpdateSelected(userID, itemID uint, selected bool) error {
    // 检查购物车项是否存在且属于该用户
    item, err := s.repoFactory.GetCartRepository().GetCartItemByID(itemID)
    if err != nil || item.UserID != userID {
        return errors.New("cart item not found")
    }

    return s.repoFactory.GetCartRepository().UpdateSelected(itemID, selected)
}

// RemoveItem 删除购物车项
func (s *CartService) RemoveItem(userID, itemID uint) error {
    // 检查购物车项是否存在且属于该用户
    item, err := s.repoFactory.GetCartRepository().GetCartItemByID(itemID)
    if err != nil || item.UserID != userID {
        return errors.New("cart item not found")
    }

    return s.repoFactory.GetCartRepository().Delete(itemID)
}

// ClearCart 清空购物车
func (s *CartService) ClearCart(userID uint) error {
    return s.repoFactory.GetCartRepository().ClearCart(userID)
}

// GetSelectedItems 获取选中的购物车项
func (s *CartService) GetSelectedItems(userID uint) ([]models.CartItem, error) {
    return s.repoFactory.GetCartRepository().GetSelectedItems(userID)
}

// SelectAll 全选/取消全选
func (s *CartService) SelectAll(userID uint, selected bool) error {
    return s.repoFactory.GetCartRepository().BatchUpdateSelected(userID, selected)
}

// SelectItems 选择部分购物车商品
func (s *CartService) SelectItems(userID uint, itemIDs []uint, selected bool) error {
    return s.repoFactory.GetCartRepository().SelectItems(userID, itemIDs, selected)
}