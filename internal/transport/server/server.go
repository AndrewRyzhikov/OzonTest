package server

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
)

type MediaController struct {
	server *http.Server
}

func NewMediaController(handler http.Handler, port string) *MediaController {
	server := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}
	return &MediaController{server: server}
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
