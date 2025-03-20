package repository

import (
	"shopify/models"
	"gorm.io/gorm"
)

type ReviewRepository struct {
	*BaseRepository
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建评论
func (r *ReviewRepository) Create(review *models.Review) error {
	// 确保 Images 字段不为 nil
	if review.Images == nil {
		review.Images = make([]string, 0)
	}

	// 使用事务来处理创建操作
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 创建评论记录
		if err := tx.Create(review).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetByID 获取评论详情
func (r *ReviewRepository) GetByID(id uint) (*models.Review, error) {
	var review models.Review
	err := r.db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id, nickname, avatar") // 只选择需要的用户字段
	}).First(&review, id).Error
	if err != nil {
		return nil, err
	}
	return &review, nil
}

// ListByProduct 获取商品的评论列表(支持分页)
func (r *ReviewRepository) ListByProduct(productID uint, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64
	
	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&models.Review{}).Where("product_id = ?", productID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Where("product_id = ?", productID).
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nickname, avatar")
		}).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, total, err
}

// ListByUser 获取用户的评论列表(支持分页)
func (r *ReviewRepository) ListByUser(userID uint, page, pageSize int) ([]models.Review, int64, error) {
	var reviews []models.Review
	var total int64
	
	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&models.Review{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := r.db.Where("user_id = ?", userID).
		Preload("Product").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&reviews).Error

	return reviews, total, err
}

// Delete 删除评论
func (r *ReviewRepository) Delete(id uint) error {
	return r.db.Delete(&models.Review{}, id).Error
}

// UpdateProductRating 更新商品评分
func (r *ReviewRepository) UpdateProductRating(productID uint) error {
	// 计算商品的平均评分
	var avgRating float64
	err := r.db.Model(&models.Review{}).
		Select("COALESCE(AVG(rating), 0)").
		Where("product_id = ?", productID).
		Scan(&avgRating).Error
	if err != nil {
		return err
	}

	// 更新商品评分
	return r.db.Model(&models.Product{}).
		Where("id = ?", productID).
		Update("rating", avgRating).Error
}

// CheckUserReviewed 检查用户是否已评论过该订单
func (r *ReviewRepository) CheckUserReviewed(userID, orderID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Review{}).
		Where("user_id = ? AND order_id = ?", userID, orderID).
		Count(&count).Error
	return count > 0, err
}

// GetProductRatingStats 获取商品评分统计
func (r *ReviewRepository) GetProductRatingStats(productID uint) (map[int]int64, error) {
	var stats []struct {
		Rating int
		Count  int64
	}
	
	err := r.db.Model(&models.Review{}).
		Select("rating, count(*) as count").
		Where("product_id = ?", productID).
		Group("rating").
		Find(&stats).Error

	if err != nil {
		return nil, err
	}

	// 转换为map格式
	result := make(map[int]int64)
	for _, stat := range stats {
		result[stat.Rating] = stat.Count
	}
	return result, nil
} 