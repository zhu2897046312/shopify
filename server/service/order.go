package service

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
	"shopify/models"
	"shopify/repository"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type OrderService struct {
	*Service
}

func NewOrderService(base *Service) *OrderService {
	return &OrderService{Service: base}
}

// CreateOrder 创建订单
func (s *OrderService) CreateOrder(userID uint, items []models.OrderItem, addressID uint) (*models.Order, error) {
	var result *models.Order

	err := s.repoFactory.GetDB().Transaction(func(tx *gorm.DB) error {
		txRepoFactory := repository.NewRepositoryFactory(tx)
		
		// 验证地址是否存在且属于该用户
		address, err := txRepoFactory.GetUserRepository().GetAddressByID(addressID)
		if err != nil || address.UserID != userID {
			return errors.New("invalid address")
		}

		totalAmount := decimal.NewFromFloat(0)
		
		// 验证商品并计算总金额
		for i := range items {
			product, err := txRepoFactory.GetProductRepository().GetByID(items[i].ProductID)
			if err != nil {
				return err
			}

			if product.Stock < items[i].Quantity {
				return fmt.Errorf("insufficient stock for product: %s", product.Name)
			}

			// 更新库存和销量
			if err := txRepoFactory.GetProductRepository().UpdateStock(product.ID, -items[i].Quantity); err != nil {
				return err
			}
			if err := txRepoFactory.GetProductRepository().UpdateSales(product.ID, items[i].Quantity); err != nil {
				return err
			}

			items[i].Price = product.Price
			itemTotal := product.Price.Mul(decimal.NewFromInt(int64(items[i].Quantity)))
			totalAmount = totalAmount.Add(itemTotal)
		}

		// 创建订单
		order := &models.Order{
			UserID:        userID,
			OrderNumber:   generateOrderNumber(),
			TotalAmount:   totalAmount,
			Status:        "pending",
			AddressID:     addressID,
			PaymentStatus: "unpaid",
		}

		if err := txRepoFactory.GetOrderRepository().Create(order); err != nil {
			return err
		}

		// 创建订单项
		for _, item := range items {
			item.OrderID = order.ID
			if err := txRepoFactory.GetOrderRepository().CreateOrderItem(&item); err != nil {
				return err
			}
		}

		// 创建物流信息
		logistics := &models.Logistics{
			OrderID:     order.ID,
			Status:      "pending",  // 初始状态为待处���
			ShippingFee: decimal.NewFromFloat(0),  // 初始运费为0
		}
		if err := txRepoFactory.GetOrderRepository().CreateLogistics(logistics); err != nil {
			return err
		}

		// 清空购物车
		if err := txRepoFactory.GetCartRepository().ClearCart(userID); err != nil {
			return err
		}

		result = order
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetOrder 获取订单详情
func (s *OrderService) GetOrder(id uint, userID uint) (*models.Order, error) {
	order, err := s.repoFactory.GetOrderRepository().GetByID(id)
	if err != nil {
		return nil, err
	}
	
	// 如果userID为0，表示管理员查询，不检查所属权
	if userID != 0 && order.UserID != userID {
		return nil, errors.New("permission denied")
	}
	
	return order, nil
}

func (s *OrderService) GetOrderByID(orderID uint) (*models.Order, error) {
	return s.repoFactory.GetOrderRepository().GetByID(orderID)
}

// ListUserOrders 获取用户订单列表
func (s *OrderService) ListUserOrders(userID uint, page, pageSize int) ([]models.Order, int64, error) {
	return s.repoFactory.GetOrderRepository().ListByUserID(userID, page, pageSize)
}

// UpdateOrderStatus 更新订单状态
func (s *OrderService) UpdateOrderStatus(orderID uint, status string) error {
	// 验证状态值
	validStatuses := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"completed": true,
		"cancelled": true,
	}
	if !validStatuses[status] {
		return errors.New("invalid status value")
	}

	return s.repoFactory.GetOrderRepository().UpdateStatus(orderID, status)
}

// UpdatePaymentStatus 更新支付状态
func (s *OrderService) UpdatePaymentStatus(orderID uint, status string) error {
	var paymentTime *time.Time
	if status == "paid" {
		now := time.Now()
		paymentTime = &now
	}
	return s.repoFactory.GetOrderRepository().UpdatePaymentStatus(orderID, status, paymentTime)
}

// CreateLogistics 创建物流信息
func (s *OrderService) CreateLogistics(logistics *models.Logistics) error {
	// 开启事务
	err := s.repoFactory.GetDB().Transaction(func(tx *gorm.DB) error {
		txRepoFactory := repository.NewRepositoryFactory(tx)

		// 创建物流信息
		if err := txRepoFactory.GetOrderRepository().CreateLogistics(logistics); err != nil {
			return err
		}

		// 创建初始物流跟踪记录
		trace := &models.LogisticsTrace{
			LogisticsID: logistics.ID,
			Location:    logistics.Carrier,
			Status:     logistics.Status,
			Description: "物流信息已创建",
			TraceTime:   time.Now(),
		}
		if err := txRepoFactory.GetOrderRepository().AddLogisticsTrace(trace); err != nil {
			return err
		}

		return nil
	})

	return err
}

// UpdateLogistics 更新物流信息
func (s *OrderService) UpdateLogistics(logistics *models.Logistics) error {
	// 开启事务
	err := s.repoFactory.GetDB().Transaction(func(tx *gorm.DB) error {
		txRepoFactory := repository.NewRepositoryFactory(tx)

		// 先创建或更新物流信息
		var existingLogistics *models.Logistics
		existingLogistics, err := txRepoFactory.GetOrderRepository().GetLogistics(logistics.OrderID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果不存在，则创建新的物流信息
				if err := txRepoFactory.GetOrderRepository().CreateLogistics(logistics); err != nil {
					return err
				}
				existingLogistics = logistics
			} else {
				return err
			}
		} else {
			// 如果存在，则更新物流信息
			existingLogistics.TrackingNo = logistics.TrackingNo
			existingLogistics.Carrier = logistics.Carrier
			existingLogistics.Status = logistics.Status
			existingLogistics.ShippingFee = logistics.ShippingFee
			existingLogistics.ShippedTime = logistics.ShippedTime
			existingLogistics.DeliveredTime = logistics.DeliveredTime
			
			if err := txRepoFactory.GetOrderRepository().UpdateLogistics(existingLogistics); err != nil {
				return err
			}
		}

		// 如果状态发生变化，添加物流跟踪记录
		if existingLogistics != nil {
			trace := &models.LogisticsTrace{
				LogisticsID: existingLogistics.ID,  // 使用已存在或新创建的物流信息ID
				Location:    logistics.Carrier,
				Status:     logistics.Status,
				Description: getStatusDescription(logistics.Status),
				TraceTime:   time.Now(),
			}
			if err := txRepoFactory.GetOrderRepository().AddLogisticsTrace(trace); err != nil {
				return err
			}
		}

		return nil
	})

	return err
}

// getStatusDescription 根据状态获取描述
func getStatusDescription(status string) string {
	descriptions := map[string]string{
		"pending":    "物流信息已创建",
		"processing": "包裹正在处理中",
		"shipping":   "包裹运输中",
		"delivered":  "包裹已送达",
		"returned":   "包裹已退回",
	}
	
	if desc, ok := descriptions[status]; ok {
		return desc
	}
	return "状态更新"
}

// GetLogistics 获取订单物流信息
func (s *OrderService) GetLogistics(orderID uint) (*models.Logistics, error) {
	return s.repoFactory.GetOrderRepository().GetLogistics(orderID)
}

// AddLogisticsTrace 添加物流跟踪记录
func (s *OrderService) AddLogisticsTrace(trace *models.LogisticsTrace) error {
	return s.repoFactory.GetOrderRepository().AddLogisticsTrace(trace)
}

// ListOrdersByStatus 管理员按状态查询订单
func (s *OrderService) ListOrdersByStatus(status string, page, pageSize int) ([]models.Order, int64, error) {
	return s.repoFactory.GetOrderRepository().ListOrdersByStatus(status, page, pageSize)
}

// generateOrderNumber 生成订单号
func generateOrderNumber() string {
	timestamp := time.Now().Format("20060102150405")
	random := rand.Intn(1000)
	return fmt.Sprintf("ORD%s%03d", timestamp, random)
} 