package main

import (
	"context"
	"os"
	"os/signal"
	"runtime/debug"
	"sync"
	"syscall"
	"time"

	"test/internal/storage"
	h "test/internal/transport/http"
	"test/pkg/config"
	"test/pkg/log"
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Помещаем объект конфигурации в контекст
	ctx = initConfig(ctx)

	// Инициализируем хранилище
	stor, err := storage.Init(ctx)
	if err != nil {
		log.Error.Fatalln(err)
	}

	// Инициализируем http сервер
	httpServer := h.Init(ctx, stor)

	// Ловим сигнал завершения работы сервиса
	chSignal := make(chan os.Signal, 2)
	signal.Notify(chSignal, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	tmRun := time.NewTimer(0)
	isContinue := true

	var wg sync.WaitGroup

	for isContinue {
		select {
		case <-chSignal:
			log.Info.Println("Начало остановки сервиса...")

			// Шлём во все горутины сигнал завершения работы
			cancel()

			// Завершаем работу хранилища и http сервера
			httpServer.CloseGracefully(ctx)
			stor.Stor.CloseGracefully(ctx)

			isContinue = false

		case <-tmRun.C:
			wg.Add(1)
			go func() {

				// Паника в горутине уронит приложение (если я ни чего не путаю,
				// то она не поймается вышестоящим кодом), поэтому необходим собственный defer
				defer func() {
					if rec := recover(); rec != nil {
						log.Error.Printf("%v / %v", rec, string(debug.Stack()))
					}

					wg.Done()

					// Если по какой-то причине завершается горутина, шлем сигнал завершения приложения
					chSignal <- syscall.SIGINT
				}()

				var err error

				// Запускаем http сервер
				if err = httpServer.Run(); err != nil {
					log.Error.Println(err)
					return
				}

			}()
		}
	}

	wg.Wait()
	log.Info.Println("Сервис выключен")

}

func initConfig(ctx context.Context) context.Context {

	// Инициализируем конфигурацию
	conf, err := config.New()
	if err != nil {
		log.Error.Fatalf("Ошибка чтения конфигурации: %v", err)
	}

	// Записываем конфигурацию в контекст
	ctx = context.WithValue(ctx, config.ContextKeyConfig, conf)

	return ctx

}
