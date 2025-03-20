package payment

import (
    "github.com/shopspring/decimal"
)

// PaymentProvider 支付提供者接口
type PaymentProvider interface {
    // CreatePayment 创建支付订单，返回支付URL
    CreatePayment(paymentID uint, amount decimal.Decimal, orderNo string) (string, error)
    
    // VerifyCallback 验证支付回调，返回支付ID、交易号和支付状态
    VerifyCallback(data map[string]string) (paymentID uint, tradeNo string, status string, err error)
    
    // SerializeCallback 序列化回调数据
    SerializeCallback(data map[string]string) string
} 