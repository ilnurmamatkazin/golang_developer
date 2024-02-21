package maps

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"test/internal/model"
	"test/pkg/config"
	"test/pkg/log"
)

// Структура, которая реализует физическое хранилище
type Storage struct {
	mutex     sync.RWMutex           // необходим для потокобезопасной работы с мапой
	m         map[string]storageData // сама мапа, где хранятся объекты
	wg        sync.WaitGroup         // используем для синхронизации горутин при завершении работы
	isNotInit bool                   // проверяем, что не начинаем работать с неинициализированным хранилищем
}

type storageData struct {
	Expires time.Time   `json:"expires"`
	Data    interface{} `json:"data"`
}

var (
	once sync.Once
	stor *Storage
)

// Инициализируем хранилище
func Init(ctx context.Context) (*Storage, error) {

	var err error

	once.Do(func() {
		stor = &Storage{isNotInit: false}

		if err = stor.load(ctx); err != nil {
			return
		}

		log.Info.Println("Создан Storage.")

		// Используем sync.WaitGroup, что бы при завершении работы функция CloseGracefully
		// могла подождать пока не завершиться данная горутина.
		stor.wg.Add(1)
		go func() {

			// Паника в горутине уронит приложение (если я ни чего не путаю,
			// то она не поймается вышестоящим кодом), поэтому необходим собственный defer
			defer func() {
				if rec := recover(); rec != nil {
					log.Error.Printf("%v / %v", rec, string(debug.Stack()))
				}

				stor.wg.Done()
			}()

			stor.run(ctx)

		}()

	})

	return stor, err

}

// Функция добавления нового значения в мапу
func (s *Storage) Add(key string, data interface{}, expires time.Time) error {

	if s.isNotInit {
		return errors.New("попытка работать с не инициализированным хранилищем")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.m[key] = storageData{Expires: expires, Data: data}

	// Фиксируем в метрике факт свершения действия
	model.ObjectCountGauge.Set(float64(len(s.m)))

	return nil

}

// При получении значения из мапы, делаем проверку на его время жизни
// и если оно меньше текущего, значит данные просрочены и подлежат удалению.
func (s *Storage) Get(key string) (interface{}, error) {

	if s.isNotInit {
		return nil, errors.New("попытка работать с не инициализированным хранилищем")
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	return s.checkExpires(key), nil

}

// Метод завершения работы и сохранения данных в файл
func (s *Storage) CloseGracefully(ctx context.Context) {

	s.wg.Wait()

	s.mutex.Lock()
	data, err := json.Marshal(s.m)
	s.mutex.Unlock()

	// Скорее всего, просто так в случае ошибки выйти из функции нельзя, нужно поробовать
	// сохранить данные как-то по другому... Возможно построчно в цикле.
	if err != nil {
		log.Error.Println(err)
		return
	}

	cfg := ctx.Value(config.ContextKeyConfig).(*config.Config)

	if err = os.WriteFile(cfg.App.StorageFile, data, 0644); err != nil {
		log.Error.Println(err)
		return
	}

	log.Info.Println("Хранилище выгруженно в файл.")

}

func (s *Storage) IsNotInit() bool {

	return s.isNotInit

}
