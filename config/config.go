package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Output   OutputConfig   `mapstructure:"output"`
}

type DatabaseConfig struct {
	Type     string `mapstructure:"type"`     // mysql, oracle, pgsql
	DSN      string `mapstructure:"dsn"`      // 数据库连接字符串模板
	Username string `mapstructure:"username"` // 用户名
	Password string `mapstructure:"password"` // 密码
}

type OutputConfig struct {
	Directory     string `mapstructure:"directory"`
	Format        string `mapstructure:"format"`
	ShowInConsole bool   `mapstructure:"show_in_console"`
	SaveToFile    bool   `mapstructure:"save_to_file"`
}

func LoadConfig(configPath string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(configPath)

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return &config, nil
}
