package config

import (
    "fmt"
    "github.com/spf13/viper"
)

type Config struct {
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Email    EmailConfig    `mapstructure:"email"`
    Payment  PaymentConfig  `mapstructure:"payment"`
}

type ServerConfig struct {
    Port int
    Mode string
}

type DatabaseConfig struct {
    Driver   string
    Host     string
    Port     int
    Username string
    Password string
    DBName   string
    Charset  string
}

// DSN 返回数据库连接字符串
func (c *DatabaseConfig) DSN() string {
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
        c.Username,
        c.Password,
        c.Host,
        c.Port,
        c.DBName,
        c.Charset,
    )
}

type RedisConfig struct {
    Host     string
    Port     int
    Password string
    DB       int
}

type EmailConfig struct {
    Host     string
    Port     int
    Username string
    Password string
    From     string
}

type PaymentConfig struct {
    WechatPay WechatPayConfig `mapstructure:"wechat"`
    Alipay    AlipayConfig    `mapstructure:"alipay"`
}

type WechatPayConfig struct {
    AppID     string `mapstructure:"app_id"`
    MchID     string `mapstructure:"mch_id"`
    ApiKey    string `mapstructure:"api_key"`
    NotifyURL string `mapstructure:"notify_url"`
}

type AlipayConfig struct {
    AppID      string `mapstructure:"app_id"`
    PrivateKey string `mapstructure:"private_key"`
    PublicKey  string `mapstructure:"public_key"`
    NotifyURL  string `mapstructure:"notify_url"`
}

var GlobalConfig Config

func Init() error {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    
    // 只设置相对于 go.mod 的配置路径
    viper.AddConfigPath("configs")
    
    if err := viper.ReadInConfig(); err != nil {
        return fmt.Errorf("failed to read config file: %v\nSearched paths: %v", err, viper.ConfigFileUsed())
    }

    if err := viper.Unmarshal(&GlobalConfig); err != nil {
        return fmt.Errorf("failed to unmarshal config: %v", err)
    }

    // 打印服务器配置
    fmt.Printf("\n=== Server Configuration ===\n")
    fmt.Printf("Port: %d\n", GlobalConfig.Server.Port)
    fmt.Printf("Mode: %s\n", GlobalConfig.Server.Mode)

    // 打印数据库配置
    fmt.Printf("\n=== Database Configuration ===\n")
    fmt.Printf("Driver: %s\n", GlobalConfig.Database.Driver)
    fmt.Printf("Host: %s\n", GlobalConfig.Database.Host)
    fmt.Printf("Port: %d\n", GlobalConfig.Database.Port)
    fmt.Printf("Username: %s\n", GlobalConfig.Database.Username)
    fmt.Printf("Database: %s\n", GlobalConfig.Database.DBName)
    fmt.Printf("Charset: %s\n", GlobalConfig.Database.Charset)
    fmt.Printf("DSN: %s\n", GlobalConfig.Database.DSN())

    // 打印Redis配置
    fmt.Printf("\n=== Redis Configuration ===\n")
    fmt.Printf("Host: %s\n", GlobalConfig.Redis.Host)
    fmt.Printf("Port: %d\n", GlobalConfig.Redis.Port)
    fmt.Printf("DB: %d\n", GlobalConfig.Redis.DB)

    // 打印邮件配置
    fmt.Printf("\n=== Email Configuration ===\n")
    fmt.Printf("Host: %s\n", GlobalConfig.Email.Host)
    fmt.Printf("Port: %d\n", GlobalConfig.Email.Port)
    fmt.Printf("Username: %s\n", GlobalConfig.Email.Username)
    fmt.Printf("From: %s\n", GlobalConfig.Email.From)

    // 打印支付配置
    fmt.Printf("\n=== Payment Configuration ===\n")
    fmt.Printf("=== Wechat Pay ===\n")
    fmt.Printf("App ID: %s\n", GlobalConfig.Payment.WechatPay.AppID)
    fmt.Printf("Mch ID: %s\n", GlobalConfig.Payment.WechatPay.MchID)
    fmt.Printf("Notify URL: %s\n", GlobalConfig.Payment.WechatPay.NotifyURL)

    fmt.Printf("\n=== Alipay ===\n")
    fmt.Printf("App ID: %s\n", GlobalConfig.Payment.Alipay.AppID)
    fmt.Printf("Notify URL: %s\n", GlobalConfig.Payment.Alipay.NotifyURL)

    fmt.Printf("\n=== Configuration End ===\n\n")


    return nil
} 