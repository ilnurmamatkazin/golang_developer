package config

import "time"

type contextKeyConfig struct{}

// Структура конфигурации
type Config struct {
	App struct {
		StorageFile string `yaml:"storage-file"`
	} `yaml:"app"`
	HTTP struct {
		Host              string        `yaml:"host"`
		ReadTimeout       time.Duration `yaml:"read-timeout"`
		WriteTimeout      time.Duration `yaml:"write-timeout"`
		ReadHeaderTimeout time.Duration `yaml:"read-header-timeout"`
	} `yaml:"http"`
}
