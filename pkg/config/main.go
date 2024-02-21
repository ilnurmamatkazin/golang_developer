package config

import (
	"os"
	"sync"
	"test/pkg/log"

	"gopkg.in/yaml.v2"
)

const configPath = "./config.yaml"

var once sync.Once
var cfg *Config

// Переменная необходима для чтения/записи в контекст
var ContextKeyConfig contextKeyConfig

// Инициализируем конфигурацию, загружаем данные из конфигурационного файла
func New() (*Config, error) {

	var err error

	once.Do(func() {
		cfg = &Config{}

		yamlFile, err := os.ReadFile(configPath)
		if err != nil {
			return
		}

		if err = yaml.Unmarshal(yamlFile, cfg); err != nil {
			return
		}

		log.Info.Println("Конфигурация загружена.")
	})

	return cfg, err

}
