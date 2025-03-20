package models

import (
    "time"
    "github.com/shopspring/decimal"
    "gorm.io/gorm"
)

// 支付方式常量
const (
    PaymentMethodWechat  = "wechat"
    PaymentMethodAlipay  = "alipay"
)

// 支付状态常量
const (
    PaymentStatusPending   = "pending"    // 待支付
    PaymentStatusPaid      = "paid"       // 已支付
    PaymentStatusFailed    = "failed"     // 支付失败
    PaymentStatusRefunded  = "refunded"   // 已退款
)

// Payment 支付记录
type Payment struct {
    ID            uint            `gorm:"primarykey;autoIncrement" json:"id"`
    OrderID       uint            `gorm:"not null;index" json:"order_id"`           // 关联的订单ID
    Order         Order           `gorm:"foreignKey:OrderID" json:"-"`              // 关联的订单
    PaymentMethod string          `gorm:"type:varchar(20);not null" json:"payment_method"` // 支付方式
    Amount        decimal.Decimal `gorm:"type:decimal(10,2);not null" json:"amount"`       // 支付金额
    TradeNo       string          `gorm:"type:varchar(100)" json:"trade_no"`              // 第三方支付交易号
    Status        string          `gorm:"type:varchar(20);not null" json:"status"`        // 支付状态
    PayTime       *time.Time      `json:"pay_time"`                                       // 支付时间
    CreatedAt     time.Time       `json:"created_at"`
    UpdatedAt     time.Time       `json:"updated_at"`
    DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

// PaymentCallback 支付回调记录
type PaymentCallback struct {
    ID        uint           `gorm:"primarykey;autoIncrement" json:"id"`
    PaymentID uint           `gorm:"not null;index" json:"payment_id"`          // 关联的支付记录ID
    Payment   Payment        `gorm:"foreignKey:PaymentID" json:"-"`            // 关联的支付记录
    TradeNo   string         `gorm:"type:varchar(100)" json:"trade_no"`        // 第三方支付交易号
    Status    string         `gorm:"type:varchar(20);not null" json:"status"`  // 回调状态
    RawData   string         `gorm:"type:text" json:"raw_data"`               // 原始回调数据
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
} 