package app

import (
	"context"
	"fmt"

	"github.com/mganuesquest/youtui/config"
	"github.com/mganuesquest/youtui/internal/adapter/mosaic"
	"github.com/mganuesquest/youtui/internal/adapter/mpv"
	"github.com/mganuesquest/youtui/internal/adapter/storage"
	"github.com/mganuesquest/youtui/internal/adapter/youtube"
	"github.com/mganuesquest/youtui/internal/adapter/ytdlp"
	"github.com/mganuesquest/youtui/internal/service"
)

// App orchestrates the application components.
// It wires together services, adapters, and the UI.
type App struct {
	Search   *service.SearchService
	Queue    *service.QueueService
	Player   *service.PlayerService
	Pomodoro *service.PomodoroService
}

// New creates a fully-wired App from the given configuration.
func New(ctx context.Context, cfg *config.Config) (*App, error) {
	// --- adapters ---

	ytClient, err := youtube.New(ctx, cfg.YouTube.APIKey, cfg.YouTube.RegionCode)
	if err != nil {
		return nil, fmt.Errorf("app: youtube adapter: %w", err)
	}

	thumbRenderer := mosaic.New()
	extractor := ytdlp.New()
	store := storage.New(cfg.Storage.QueueFile)

	mpvPlayer, err := mpv.New(cfg.Player.SocketPath)
	if err != nil {
		return nil, fmt.Errorf("app: mpv adapter: %w", err)
	}

	// --- services ---

	searchSvc := service.NewSearchService(
		ytClient,
		thumbRenderer,
		cfg.UI.ThumbnailWidth,
		cfg.UI.ThumbnailHeight,
	)

	queueSvc := service.NewQueueService(store)

	playerSvc := service.NewPlayerService(
		mpvPlayer,
		extractor,
		queueSvc,
		cfg.Player.DefaultVolume,
	)

	pomodoroSvc := service.NewPomodoroService(
		playerSvc,
		cfg.Pomodoro.DefaultProfile,
	)

	return &App{
		Search:   searchSvc,
		Queue:    queueSvc,
		Player:   playerSvc,
		Pomodoro: pomodoroSvc,
	}, nil
}

// Close shuts down all services that hold resources.
func (a *App) Close() error {
	if a.Player != nil {
		return a.Player.Close()
	}
	return nil
}
