package email

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/smtp"

	"shopify/config"
)


type EmailSender struct {
	config config.EmailConfig
}

func NewEmailSender() *EmailSender {
	// 从配置文件获取支付配置
	emailConfig := config.GlobalConfig.Email
	return &EmailSender{
		config: emailConfig,
	}
}

// // loadEmailConfig 从config.yaml加载邮箱配置
// func loadEmailConfig() EmailConfig {
// 	viper.SetConfigName("config") // 配置文件名 (不带扩展名)
// 	viper.SetConfigType("yaml")   // 配置文件类型
// 	viper.AddConfigPath(".")      // 配置文件路径

// 	if err := viper.ReadInConfig(); err != nil {
// 		panic(fmt.Errorf("Fatal error config file: %w \n", err))
// 	}

// 	return EmailConfig{
// 		Host:     viper.GetString("email.host"),
// 		Port:     viper.GetInt("email.port"),
// 		Username: viper.GetString("email.username"),
// 		Password: viper.GetString("email.password"),
// 		From:     viper.GetString("email.from"),
// 	}
// }

// GenerateVerificationCode 生成6位数字验证码
func GenerateVerificationCode() string {
	numbers := make([]byte, 6)
	for i := range numbers {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			numbers[i] = byte(i + 1)
			continue
		}
		numbers[i] = byte(n.Int64())
	}

	code := ""
	for _, num := range numbers {
		code += fmt.Sprintf("%d", num)
	}
	return code
}

// SendVerificationCode 发送验证码邮件
func (s *EmailSender) SendVerificationCode(to, code string) error {
	subject := "Email Verification Code"
	body := fmt.Sprintf(`
Dear User,

Your verification code is: %s

This code will expire in 15 minutes.

If you did not request this code, please ignore this email.

Best regards,
Your Application Team
`, code)

	message := fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"Content-Type: text/plain; charset=UTF-8\r\n"+
		"\r\n"+
		"%s", to, subject, body)

	auth := smtp.PlainAuth("", s.config.Username, s.config.Password, s.config.Host)
	addr := fmt.Sprintf("%s:%d", s.config.Host, s.config.Port)

	return smtp.SendMail(addr, auth, s.config.From, []string{to}, []byte(message))
}