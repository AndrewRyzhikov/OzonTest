package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"

	"OzonTest/internal/config"
)

type MediaController struct {
	server *http.Server
	config config.HttpServerConfig
}

func NewMediaController(handler http.Handler, config config.HttpServerConfig) *MediaController {
	server := &http.Server{
		Addr:         ":" + config.Port,
		Handler:      handler,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	}
	return &MediaController{server: server, config: config}
}

func (mc *MediaController) Start() {
	go func() {
		err := mc.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Failed to started media Controller")
		}
	}()
}

func (mc *MediaController) Shutdown() error {
	if err := mc.server.Shutdown(context.Background()); err != nil {
		return fmt.Errorf("failed to stopped Healtcontrollerhecker : %w", err)
	}
	return nil
}
