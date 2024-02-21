package maps

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"test/internal/model"
	"test/pkg/config"
	"test/pkg/log"
	"time"
)

// Так как изначально, ничего не сказано про размер хранилища,
// то для уменьшения аллокаций памяти, в самый первый раз инициализируем
// мапу определенной емкости. Размер емкости взят степенью двойки.
func (s *Storage) load(ctx context.Context) error {

	cfg := ctx.Value(config.ContextKeyConfig).(*config.Config)

	file, err := os.Open(cfg.App.StorageFile)
	if err != nil {
		// Если файла не существует, то инициализируем мапу
		if errors.Is(err, os.ErrNotExist) {
			s.mutex.Lock()
			defer s.mutex.Unlock()

			s.m = make(map[string]storageData, 1024)
			log.Info.Println("Мапа проинициализированна.")

			return nil
		}

		return err
	}
	defer file.Close()

	// Если файл существует, то загружаем его в мапу. Так как запись в файл
	// была не построчная, то и чтение тоже не построчное.
	s.mutex.Lock()
	defer s.mutex.Unlock()

	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&s.m); err != nil {
		return err
	}

	log.Info.Println("Файл хранилища загружен.")

	return nil

}

// Данная функция фоново анализирует содержимое мары на предмет устаревших данных
func (s *Storage) run(ctx context.Context) {

	// Раз в сутки пробегаемся по всей мапе и удаляем просроченные объекты.
	// Это необходимо, что бы в мапе не плодились редко запрашиваемые просроченные объекты.
	// В идеале это сделать в полночь через cron, а не через тикер.
	tiCheckExpires := time.NewTicker(time.Hour * 24)
	defer tiCheckExpires.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-tiCheckExpires.C:
			// Нужно иметь в виду, что если мапа будет очень большой, то пока не закончится цикл мы:
			// 1) не сможем выйти по контексту, когда он придет
			// 2) на долго захватим мьютекс и поставим в очередь запросы на добавление и получение данных
			// Если вызывать мьютекс внутри цикла и какое-то время спать (250 мс),
			// то выше описанные проблемы уйдут, но возникнет другая - гонки
			func() {
				s.mutex.Lock()
				defer s.mutex.Unlock()

				for key := range s.m {
					s.checkExpires(key)
				}
			}()

		}
	}

}

// Проверка времени жизни объекта, с последующим его получением.
// Метод не потокобезопасен, поэтому в точках его вызова необходимо
// использовать мьютекс.
func (s *Storage) checkExpires(key string) interface{} {

	data, ok := s.m[key]
	if !ok {
		return nil
	}

	// Проверяем только у объетов, которые имеют не нулевое значение даты
	if (!data.Expires.IsZero()) && data.Expires.UTC().Before(time.Now().UTC()) {
		delete(s.m, key)

		// Фиксируем в метриках факт выполнения действия
		model.ObjectDelCounter.Inc()
		model.ObjectCountGauge.Set(float64(len(s.m)))

		return nil
	}

	return data.Data

}
