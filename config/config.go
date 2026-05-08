package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Database DatabaseConfig `mapstructure:"database" yaml:"database"`
	Output   OutputConfig   `mapstructure:"output" yaml:"output"`
}

type DatabaseConfig struct {
	Type     string `mapstructure:"type" yaml:"type"`         // mysql, oracle, pgsql
	DSN      string `mapstructure:"dsn" yaml:"dsn"`           // 数据库连接字符串模板
	Username string `mapstructure:"username" yaml:"username"` // 用户名
	Password string `mapstructure:"password" yaml:"password"` // 密码
}

type OutputConfig struct {
	Directory     string `mapstructure:"directory" yaml:"directory"`
	Format        string `mapstructure:"format" yaml:"format"`
	ShowInConsole bool   `mapstructure:"show_in_console" yaml:"show_in_console"`
	SaveToFile    bool   `mapstructure:"save_to_file" yaml:"save_to_file"`
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

func SaveConfig(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	return nil
}
