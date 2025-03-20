package service

import (
	"errors"
	"shopify/models"
	"gorm.io/gorm"
)

type ReviewService struct {
	*Service
}

func NewReviewService(base *Service) *ReviewService {
	return &ReviewService{Service: base}
}

// CreateReview 创建评论
func (s *ReviewService) CreateReview(review *models.Review) error {
	// 检查订单是否存在
	order, err := s.repoFactory.GetOrderRepository().GetByID(review.OrderID)
	if err != nil {
		return errors.New("order not found")
	}

	// 检查订单是否属于该用户
	if order.UserID != review.UserID {
		return errors.New("order does not belong to this user")
	}

	// 检查订单状态是否为已完成
	if order.Status != "completed" {
		return errors.New("can only review completed orders")
	}

	// 检查是否已经评论过
	reviewed, err := s.repoFactory.GetReviewRepository().CheckUserReviewed(review.UserID, review.OrderID)
	if err != nil {
		return err
	}
	if reviewed {
		return errors.New("order has already been reviewed")
	}

	// 检查评分范围
	if review.Rating < 1 || review.Rating > 5 {
		return errors.New("rating must be between 1 and 5")
	}

	// 开启事务
	err = s.repoFactory.GetDB().Transaction(func(tx *gorm.DB) error {
		// 创建评论
		if err := s.repoFactory.GetReviewRepository().Create(review); err != nil {
			return err
		}

		// 更新商品评分
		if err := s.repoFactory.GetReviewRepository().UpdateProductRating(review.ProductID); err != nil {
			return err
		}

		return nil
	})

	return err
}

// GetProductReviews 获取商品评论列表
func (s *ReviewService) GetProductReviews(productID uint, page, pageSize int) ([]models.Review, int64, error) {
	return s.repoFactory.GetReviewRepository().ListByProduct(productID, page, pageSize)
}

// GetUserReviews 获取用户评论列表
func (s *ReviewService) GetUserReviews(userID uint, page, pageSize int) ([]models.Review, int64, error) {
	return s.repoFactory.GetReviewRepository().ListByUser(userID, page, pageSize)
}

// DeleteReview 删除评论
func (s *ReviewService) DeleteReview(reviewID, userID uint) error {
	// 检查评论是否存在
	review, err := s.repoFactory.GetReviewRepository().GetByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	// 检查评论是否属于该用户
	if review.UserID != userID {
		return errors.New("review does not belong to this user")
	}

	// 开启事务
	err = s.repoFactory.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除评论
		if err := s.repoFactory.GetReviewRepository().Delete(reviewID); err != nil {
			return err
		}

		// 更新商品评分
		if err := s.repoFactory.GetReviewRepository().UpdateProductRating(review.ProductID); err != nil {
			return err
		}

		return nil
	})

	return err
}

// GetProductRatingStats 获取商品评分统计
func (s *ReviewService) GetProductRatingStats(productID uint) (map[int]int64, error) {
	return s.repoFactory.GetReviewRepository().GetProductRatingStats(productID)
}

// AdminDeleteReview 管理员删除评论
func (s *ReviewService) AdminDeleteReview(reviewID uint) error {
	// 检查评论是否存在
	review, err := s.repoFactory.GetReviewRepository().GetByID(reviewID)
	if err != nil {
		return errors.New("review not found")
	}

	// 开启事务
	err = s.repoFactory.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除评论
		if err := s.repoFactory.GetReviewRepository().Delete(reviewID); err != nil {
			return err
		}

		// 更新商品评分
		if err := s.repoFactory.GetReviewRepository().UpdateProductRating(review.ProductID); err != nil {
			return err
		}

		return nil
	})

	return err
} 