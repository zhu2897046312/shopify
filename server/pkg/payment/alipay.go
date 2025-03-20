package payment

import (
	"encoding/json"
	"github.com/shopspring/decimal"
)

type AlipayConfig struct {
	AppID      string
	PrivateKey string
	PublicKey  string
	NotifyURL  string
}

type AlipayProvider struct {
	config AlipayConfig
}

func NewAlipayProvider(config AlipayConfig) *AlipayProvider {
	return &AlipayProvider{
		config: config,
	}
}

func (p *AlipayProvider) CreatePayment(paymentID uint, amount decimal.Decimal, orderNo string) (string, error) {
	// TODO: 实现支付宝支付接口调用
	// 1. 生成支付参数
	// 2. 调用支付宝统一下单接口
	// 3. 返回支付URL或表单
	return "https://alipay.url", nil
}

func (p *AlipayProvider) VerifyCallback(data map[string]string) (uint, string, string, error) {
	// TODO: 实现支付宝回调验证
	// 1. 验证签名
	// 2. 解析支付结果
	// 3. 返回支付ID、交易号和状态
	return 1, "alipay_trade_no", "paid", nil
}

func (p *AlipayProvider) SerializeCallback(data map[string]string) string {
	bytes, _ := json.Marshal(data)
	return string(bytes)
} 