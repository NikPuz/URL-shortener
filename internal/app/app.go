package app

import (
	"context"
	"github.com/go-chi/chi/v5"
	"github.com/nanmu42/gzip"
	"github.com/seefan/gossdb/v2"
	"github.com/seefan/gossdb/v2/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net/http"
	"strconv"
	"url-shortner/internal/api/handler"
	"url-shortner/internal/api/middleware"
	"url-shortner/internal/app/config"
	"url-shortner/internal/entity"
	"url-shortner/internal/repository"
	"url-shortner/internal/service"
)

func Run(ctx context.Context, cfg *config.Config) {
	closer := newCloser()
	logger := newLogger()
	router := chi.NewRouter()
	startDataBase(logger, cfg)

	server := newServer(cfg, router)
	closer.Add(server.Shutdown)

	// Repository
	urlRepository := repository.NewURLRepository()

	// Service
	urlService := service.NewURLService(urlRepository)

	// API
	middleware := middleware.NewMiddleware(logger)
	handler.RegisterURLHandlers(router, urlService, logger, middleware)

	go func() {
		logger.Info("starting the server", zap.String("Address", server.Addr))
		logger.DPanic("ListenAndServe", zap.Any("Error", server.ListenAndServe()))
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := closer.Close(shutdownCtx); err != nil {
		logger.Error("Close err", zap.Error(err))
	}
}

func newServer(cfg *config.Config, router http.Handler) *http.Server {
	return &http.Server{
		Handler:        gzip.DefaultHandler().WrapHandler(router),
		Addr:           ":" + strconv.Itoa(cfg.Server.Port),
		WriteTimeout:   cfg.Server.WriteTimeout,
		ReadTimeout:    cfg.Server.ReadTimeout,
		IdleTimeout:    cfg.Server.IdleTimeout,
		MaxHeaderBytes: 1 << 20,
	}
}

func startDataBase(logger *zap.Logger, cfg *config.Config) {

	ssdbConfig := conf.Config{
		Host:             cfg.DataBase.Address,
		Port:             cfg.DataBase.Port,
		ReadWriteTimeout: cfg.DataBase.ReadWriteTimeout,
		MaxPoolSize:      cfg.DataBase.PoolSize,
		ConnectTimeout:   cfg.DataBase.ConnectTimeout,
	}

	err := gossdb.Start(&ssdbConfig)

	logger.Info("start establishing a connection to the database", zap.Any("Config", ssdbConfig))

	if err != nil {
		panic(err.Error())
	}

	logger.Info("successful connection to the database")
}

func newLogger() *zap.Logger {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}

	logger.Info("Zap Logger", zap.String("Level", logger.Level().String()))

	return logger
}

func newCloser() *entity.Closer {
	return &entity.Closer{}
}
