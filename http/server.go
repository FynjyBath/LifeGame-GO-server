package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fynjybath/LifeGame-GO-server/http/handler"
	"github.com/fynjybath/LifeGame-GO-server/internal/service"
	"go.uber.org/zap"
)

func new(ctx context.Context,
	logger *zap.Logger,
	lifeService service.LifeService,
) (http.Handler, error) {
	muxHandler, err := handler.New(ctx, lifeService)
	if err != nil {
		return nil, fmt.Errorf("handler initialization error: %w", err)
	}
	// middleware для обработчиков
	muxHandler = handler.Decorate(muxHandler, loggingMiddleware(logger))

	return muxHandler, nil
}

func Run(
	ctx context.Context,
	logger *zap.Logger,
	height, width int,
) (func(context.Context) error, error) {
	// сервис с игрой
	lifeService, err := service.New(height, width)
	if err != nil {
		return nil, err
	}

	muxHandler, err := new(ctx, logger, *lifeService)
	if err != nil {
		return nil, err
	}

	srv := &http.Server{Addr: ":8081", Handler: muxHandler}

	go func() {
		// Запускаем сервер
		if err := srv.ListenAndServe(); err != nil {
			logger.Error("ListenAndServe",
				zap.String("err", err.Error()))
		}
	}()
	// вернем функцию для завершения работы сервера
	return srv.Shutdown, nil
}

// middleware для логированя запросов
func loggingMiddleware(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Пропуск запроса к следующему обработчику
			next.ServeHTTP(w, r)

			// Завершение логирования после выполнения запроса
			duration := time.Since(start)
			logger.Info("HTTP request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}
