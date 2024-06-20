package app

import (
	"context"
	"errors"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"gopkg.in/natefinch/lumberjack.v2"

	"OzonTest/internal/config"
	"OzonTest/internal/entity"
	"OzonTest/internal/repository/contracts"
	"OzonTest/internal/repository/db"
	"OzonTest/internal/repository/inmemory"
	"OzonTest/internal/service"
	"OzonTest/internal/transport/server"
)

func StartApp() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config file")
	}

	log.Logger, err = setupLog(&cfg.LogConfig)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to setup logger")
	}

	var MediaRepository contracts.MediaRepository

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errGroup, errCtx := errgroup.WithContext(ctx)

	switch cfg.RepositoryType {
	case "postgres":
		commentStorage, err := db.NewStorage[entity.Comment](cfg.DataBaseConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create comment storage")
		}

		log.Info().Msg("Comment storage connected...")

		errGroup.Go(func() error {
			<-errCtx.Done()
			if err = commentStorage.Close(); err != nil {
				return fmt.Errorf("failed to close comment storage %w", err)
			}

			log.Info().Msg("Comment storage disconnected...")
			return nil
		})

		postStorage, err := db.NewStorage[entity.Post](cfg.DataBaseConfig)
		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create post storage")
		}

		log.Info().Msg("Post storage connected...")

		errGroup.Go(func() error {
			<-errCtx.Done()
			if err = postStorage.Close(); err != nil {
				return fmt.Errorf("failed to close comment storage %w", err)
			}

			log.Info().Msg("Post storage disconnected...")
			return nil
		})

		MediaRepository = db.NewDatabaseMedia(commentStorage, postStorage)

	case "in_memory":
		MediaRepository = inmemory.NewMedia()
	default:
		log.Fatal().Msg("Unknown repository type")
	}

	log.Info().Msg("Media repository started...")

	commentService := service.NewCommentService(MediaRepository)
	postService := service.NewPostService(MediaRepository)
	subscriptionService := service.NewSubscriptionService()

	router := server.NewRouter(commentService, postService, subscriptionService)

	controller := server.NewMediaController(router.Mux, cfg.HttpServerConfig)
	controller.Start()
	log.Info().Msg("Media controller started...")

	log.Info().Msg("App started...")

	errGroup.Go(func() error {
		<-errCtx.Done()

		if err = controller.Shutdown(); err != nil {
			return fmt.Errorf("failed to shutdown controller: %w", err)
		}
		log.Info().Msg("Media controller stopped...")
		return nil
	})

	errGroup.Go(func() error {
		<-errCtx.Done()

		timeoutCtx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		select {
		case <-timeoutCtx.Done():
			if errors.Is(timeoutCtx.Err(), context.DeadlineExceeded) {
				return timeoutCtx.Err()
			}

			return nil
		}
	})

	err = errGroup.Wait()

	if errors.Is(err, context.Canceled) || err == nil {
		log.Info().Msg("App stopped...")
	} else if err != nil {
		log.Fatal().Err(err).Msg("Failed to shutdown app")
	}
}

func setupLog(config *config.LogConfig) (zerolog.Logger, error) {
	lumberjackCfg := config.Lumberjack

	lr := &lumberjack.Logger{
		Filename:   config.Path,
		MaxSize:    int(lumberjackCfg.MaxSize),
		MaxAge:     int(lumberjackCfg.MaxAge),
		MaxBackups: int(lumberjackCfg.MaxBackups),
		LocalTime:  lumberjackCfg.LocalTime,
		Compress:   lumberjackCfg.Compress,
	}

	level, err := zerolog.ParseLevel(config.Level)

	if err != nil {
		return zerolog.Logger{}, fmt.Errorf("invalid log level: %s", config.Level)
	}

	return zerolog.New(lr).Level(level).With().Timestamp().Logger(), nil
}
