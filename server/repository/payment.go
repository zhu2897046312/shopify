package repository

import (
	"shopify/models"
	"gorm.io/gorm"
	"time"
)

type PaymentRepository struct {
	*BaseRepository
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建支付记录
func (r *PaymentRepository) Create(payment *models.Payment) error {
	return r.db.Create(payment).Error
}

// GetByID 获取支付记录
func (r *PaymentRepository) GetByID(id uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.First(&payment, id).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByOrderID 获取订单的支付记录
func (r *PaymentRepository) GetByOrderID(orderID uint) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// GetByTradeNo 通过交易号获取支付记录
func (r *PaymentRepository) GetByTradeNo(tradeNo string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Where("trade_no = ?", tradeNo).First(&payment).Error
	if err != nil {
		return nil, err
	}
	return &payment, nil
}

// UpdateStatus 更新支付状态
func (r *PaymentRepository) UpdateStatus(id uint, status string, tradeNo string) error {
	updates := map[string]interface{}{
		"status": status,
		"trade_no": tradeNo,
	}
	if status == models.PaymentStatusPaid {
		now := time.Now()
		updates["pay_time"] = &now
	}
	return r.db.Model(&models.Payment{}).Where("id = ?", id).Updates(updates).Error
}

// CreateCallback 创建支付回调记录
func (r *PaymentRepository) CreateCallback(callback *models.PaymentCallback) error {
	return r.db.Create(callback).Error
} 