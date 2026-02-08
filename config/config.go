package config

import (
	"os"
	"path/filepath"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config holds all application configuration
type Config struct {
	YouTube   YouTubeConfig   `yaml:"youtube"`
	Player    PlayerConfig    `yaml:"player"`
	Pomodoro  PomodoroConfig  `yaml:"pomodoro"`
	Storage   StorageConfig   `yaml:"storage"`
	UI        UIConfig        `yaml:"ui"`
}

// YouTubeConfig holds YouTube API settings
type YouTubeConfig struct {
	APIKey     string `yaml:"api_key"`
	MaxResults int    `yaml:"max_results"`
	RegionCode string `yaml:"region_code"`
}

// PlayerConfig holds mpv player settings
type PlayerConfig struct {
	SocketPath    string `yaml:"socket_path"`
	DefaultVolume int    `yaml:"default_volume"`
}

// PomodoroConfig holds pomodoro timer settings
type PomodoroConfig struct {
	DefaultProfile    string `yaml:"default_profile"`
	NotificationSound bool   `yaml:"notification_sound"`
}

// StorageConfig holds persistence settings
type StorageConfig struct {
	QueueFile string `yaml:"queue_file"`
}

// UIConfig holds UI settings
type UIConfig struct {
	ShowThumbnails  bool `yaml:"show_thumbnails"`
	ThumbnailWidth  int  `yaml:"thumbnail_width"`
	ThumbnailHeight int  `yaml:"thumbnail_height"`
}

// Default returns configuration with default values
func Default() *Config {
	home, _ := os.UserHomeDir()
	return &Config{
		YouTube: YouTubeConfig{
			APIKey:     "",
			MaxResults: 25,
			RegionCode: "US",
		},
		Player: PlayerConfig{
			SocketPath:    "/tmp/youtui-mpv.sock",
			DefaultVolume: 80,
		},
		Pomodoro: PomodoroConfig{
			DefaultProfile:    "classic",
			NotificationSound: true,
		},
		Storage: StorageConfig{
			QueueFile: filepath.Join(home, ".config", "youtui", "queue.json"),
		},
		UI: UIConfig{
			ShowThumbnails:  true,
			ThumbnailWidth:  12,
			ThumbnailHeight: 6,
		},
	}
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	cfg := Default()

	// Try to load from config file
	configPaths := []string{
		"config.yaml",
		filepath.Join(os.Getenv("HOME"), ".config", "youtui", "config.yaml"),
	}

	for _, path := range configPaths {
		if data, err := os.ReadFile(path); err == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, err
			}
			break
		}
	}

	// Override with environment variables
	if apiKey := os.Getenv("YOUTUBE_API_KEY"); apiKey != "" {
		cfg.YouTube.APIKey = apiKey
	}

	if maxResults := os.Getenv("YOUTUI_MAX_RESULTS"); maxResults != "" {
		if n, err := strconv.Atoi(maxResults); err == nil {
			cfg.YouTube.MaxResults = n
		}
	}

	if region := os.Getenv("YOUTUI_REGION"); region != "" {
		cfg.YouTube.RegionCode = region
	}

	if socket := os.Getenv("YOUTUI_MPV_SOCKET"); socket != "" {
		cfg.Player.SocketPath = socket
	}

	if vol := os.Getenv("YOUTUI_DEFAULT_VOLUME"); vol != "" {
		if n, err := strconv.Atoi(vol); err == nil {
			cfg.Player.DefaultVolume = n
		}
	}

	return cfg, nil
}
