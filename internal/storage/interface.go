package storage

import (
	"context"
	"time"

	"test/internal/storage/maps"
)

// Интерфейс хранилища
type Storager interface {
	Add(key string, data interface{}, expires time.Time) error
	Get(key string) (interface{}, error)
	IsNotInit() bool
	CloseGracefully(ctx context.Context)
}

// Объект который реализует интерфейс хранилища
type Storage struct {
	Stor Storager
}

func Init(ctx context.Context) (*Storage, error) {

	// Получаем объект, который реализует физическое хранилище
	stor, err := maps.Init(ctx)

	return &Storage{Stor: stor}, err

}
