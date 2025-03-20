package service

import (
	"errors"
	"shopify/models"
	"shopify/pkg/payment"
	"shopify/config"
)

type PaymentService struct {
	*Service
}

func NewPaymentService(base *Service) *PaymentService {
	return &PaymentService{Service: base}
}

// CreatePayment 创建支付
func (s *PaymentService) CreatePayment(orderID uint, method string) (*models.Payment, string, error) {
	// 获取订单信息
	order, err := s.repoFactory.GetOrderRepository().GetByID(orderID)
	if err != nil {
		return nil, "", errors.New("order not found")
	}

	// 检查订单状态
	if order.Status != "pending" {
		return nil, "", errors.New("invalid order status")
	}

	// 创建支付记录
	paymentRecord := &models.Payment{
		OrderID:       orderID,
		PaymentMethod: method,
		Amount:        order.TotalAmount,
		Status:        models.PaymentStatusPending,
	}

	if err := s.repoFactory.GetPaymentRepository().Create(paymentRecord); err != nil {
		return nil, "", err
	}

	// 从配置文件获取支付配置
	paymentConfig := config.GlobalConfig.Payment

	// 根据支付方式选择支付提供者
	var provider payment.PaymentProvider
	switch method {
	case models.PaymentMethodWechat:
		provider = payment.NewWechatPayProvider(payment.WechatPayConfig{
			AppID:     paymentConfig.WechatPay.AppID,
			MchID:     paymentConfig.WechatPay.MchID,
			ApiKey:    paymentConfig.WechatPay.ApiKey,
			NotifyURL: paymentConfig.WechatPay.NotifyURL,
		})
	case models.PaymentMethodAlipay:
		provider = payment.NewAlipayProvider(payment.AlipayConfig{
			AppID:      paymentConfig.Alipay.AppID,
			PrivateKey: paymentConfig.Alipay.PrivateKey,
			PublicKey:  paymentConfig.Alipay.PublicKey,
			NotifyURL:  paymentConfig.Alipay.NotifyURL,
		})
	default:
		return nil, "", errors.New("unsupported payment method")
	}

	// 调用支付接口
	paymentURL, err := provider.CreatePayment(paymentRecord.ID, paymentRecord.Amount, order.OrderNumber)
	if err != nil {
		return nil, "", err
	}

	return paymentRecord, paymentURL, nil
}

// HandleCallback 处理支付回调
func (s *PaymentService) HandleCallback(method string, data map[string]string) error {
	// 从配置文件获取支付配置
	paymentConfig := config.GlobalConfig.Payment

	// 根据支付方式选择支付提供者
	var provider payment.PaymentProvider
	switch method {
	case models.PaymentMethodWechat:
		provider = payment.NewWechatPayProvider(payment.WechatPayConfig{
			AppID:     paymentConfig.WechatPay.AppID,
			MchID:     paymentConfig.WechatPay.MchID,
			ApiKey:    paymentConfig.WechatPay.ApiKey,
			NotifyURL: paymentConfig.WechatPay.NotifyURL,
		})
	case models.PaymentMethodAlipay:
		provider = payment.NewAlipayProvider(payment.AlipayConfig{
			AppID:      paymentConfig.Alipay.AppID,
			PrivateKey: paymentConfig.Alipay.PrivateKey,
			PublicKey:  paymentConfig.Alipay.PublicKey,
			NotifyURL:  paymentConfig.Alipay.NotifyURL,
		})
	default:
		return errors.New("unsupported payment method")
	}

	// 验证回调数据
	paymentID, tradeNo, status, err := provider.VerifyCallback(data)
	if err != nil {
		return err
	}

	// 更新支付状态
	if err := s.repoFactory.GetPaymentRepository().UpdateStatus(paymentID, status, tradeNo); err != nil {
		return err
	}

	// 记录回调信息
	callback := &models.PaymentCallback{
		PaymentID: paymentID,
		TradeNo:   tradeNo,
		Status:    status,
		RawData:   provider.SerializeCallback(data),
	}
	if err := s.repoFactory.GetPaymentRepository().CreateCallback(callback); err != nil {
		return err
	}

	// 如果支付成功，更新订单状态
	if status == models.PaymentStatusPaid {
		payment, err := s.repoFactory.GetPaymentRepository().GetByID(paymentID)
		if err != nil {
			return err
		}
		if err := s.repoFactory.GetOrderRepository().UpdateStatus(payment.OrderID, "paid"); err != nil {
			return err
		}
	}

	return nil
}

// QueryPaymentStatus 查询支付状态
func (s *PaymentService) QueryPaymentStatus(paymentID uint) (string, error) {
	payment, err := s.repoFactory.GetPaymentRepository().GetByID(paymentID)
	if err != nil {
		return "", err
	}
	return payment.Status, nil
}
