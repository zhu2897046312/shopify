package repository

import (
	"shopify/models"
	"gorm.io/gorm"
	"time"
)

type OrderRepository struct {
	*BaseRepository
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

// Create 创建订单
func (r *OrderRepository) Create(order *models.Order) error {
	return r.db.Create(order).Error
}

// GetByID 获取订单详情
func (r *OrderRepository) GetByID(id uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Preload("OrderItems").
		Preload("OrderItems.Product").
		Preload("Address").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nickname, email") // 只选择需要的用户字段
		}).
		First(&order, id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *OrderRepository) GetUserOrders(userID uint) ([]models.Order, error) {
    var orders []models.Order
    
    // 使用Preload预加载关联数据
    err := r.db.Preload("Items").  // 预加载订单项
                Preload("Items.Product").  // 预加载订单项中的商品信息
                Preload("Address").  // 预加载地址信息
                Where("user_id = ?", userID).
                Order("created_at DESC").
                Find(&orders).Error
    
    if err != nil {
        return nil, err
    }
    
    return orders, nil
}

// ListByUserID 获取用户订单列表
func (r *OrderRepository) ListByUserID(userID uint, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	
	offset := (page - 1) * pageSize

	// 获取总数
	if err := r.db.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 使用正确的关联名称进行预加载
	err := r.db.Where("user_id = ?", userID).
		Preload("OrderItems").  // 修改这里
		Preload("OrderItems.Product").  // 修改这里
		Preload("Address").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, err
}

// Update 更新订单信息
func (r *OrderRepository) Update(order *models.Order) error {
    // 使用 Model 和 Where 来指定更新的记录
    return r.db.Model(&models.Order{}).
        Where("id = ?", order.ID).
        Updates(map[string]interface{}{
            "status":         order.Status,
            "payment_status": order.PaymentStatus,
            "payment_time":   order.PaymentTime,
            "updated_at":     time.Now(),
        }).Error
}

// UpdateStatus 更新订单状态
func (r *OrderRepository) UpdateStatus(orderID uint, status string) error {
    return r.db.Model(&models.Order{}).
        Where("id = ?", orderID).
        Updates(map[string]interface{}{
            "status":     status,
            "updated_at": time.Now(),
        }).Error
}

// UpdatePaymentStatus 更新支付状态
func (r *OrderRepository) UpdatePaymentStatus(orderID uint, status string, paymentTime *time.Time) error {
    updates := map[string]interface{}{
        "payment_status": status,
        "payment_time":   paymentTime,
        "updated_at":     time.Now(),
    }
    return r.db.Model(&models.Order{}).
        Where("id = ?", orderID).
        Updates(updates).Error
}

// CreateLogistics 创建物流信息
func (r *OrderRepository) CreateLogistics(logistics *models.Logistics) error {
	return r.db.Create(logistics).Error
}

// UpdateLogistics 更新物流信息
func (r *OrderRepository) UpdateLogistics(logistics *models.Logistics) error {
	return r.db.Model(&models.Logistics{}).
        Where("id = ?", logistics.ID).
        Updates(map[string]interface{}{
            "tracking_no":     logistics.TrackingNo,
            "carrier":         logistics.Carrier,
            "status":          logistics.Status,
            "shipping_fee":    logistics.ShippingFee,
            "shipped_time":    logistics.ShippedTime,
            "delivered_time":  logistics.DeliveredTime,
            "updated_at":      time.Now(),
        }).Error
}

// GetLogistics 获取订单物流信息
func (r *OrderRepository) GetLogistics(orderID uint) (*models.Logistics, error) {
	var logistics models.Logistics
	err := r.db.Where("order_id = ?", orderID).
		Preload("Traces", func(db *gorm.DB) *gorm.DB {
			return db.Order("trace_time DESC")
		}).
		First(&logistics).Error
	if err != nil {
		return nil, err
	}
	return &logistics, nil
}

// AddLogisticsTrace 添加物流跟踪记录
func (r *OrderRepository) AddLogisticsTrace(trace *models.LogisticsTrace) error {
	return r.db.Create(trace).Error
}

// ListOrdersByStatus 管理员按状态查询订单
func (r *OrderRepository) ListOrdersByStatus(status string, page, pageSize int) ([]models.Order, int64, error) {
	var orders []models.Order
	var total int64
	
	offset := (page - 1) * pageSize

	query := r.db.Model(&models.Order{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	err := query.Preload("OrderItems.Product").  // 修改这里，使用 OrderItems 而不是 Items
		Preload("Address").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, nickname, email")
		}).
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&orders).Error

	return orders, total, err
}

// CreateOrderItem 创建订单项
func (r *OrderRepository) CreateOrderItem(item *models.OrderItem) error {
	return r.db.Create(item).Error
}

// GetOrderItem 获取订单项详情
func (r *OrderRepository) GetOrderItem(id uint) (*models.OrderItem, error) {
	var item models.OrderItem
	err := r.db.Preload("Product").First(&item, id).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ListOrderItems 获取订单的所有订单项
func (r *OrderRepository) ListOrderItems(orderID uint) ([]models.OrderItem, error) {
	var items []models.OrderItem
	err := r.db.Where("order_id = ?", orderID).
		Preload("Product").
		Find(&items).Error
	return items, err
} 