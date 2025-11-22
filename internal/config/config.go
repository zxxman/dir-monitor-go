package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"dir-monitor-go/internal/model"
)

const (
	DefaultMaxConcurrentOperations          = 5
	DefaultOperationTimeoutSeconds          = 300
	DefaultEventChannelBufferSize           = 100
	DefaultMinStabilityTimeMs               = 500
	DefaultDirectoryStabilityQuietMs        = 2000
	DefaultExecutionDedupIntervalSeconds    = 5
	DefaultDirectoryStabilityTimeoutSeconds = 30
	DefaultRetryAttempts                    = 3
	DefaultRetryDelaySeconds                = 5
	DefaultHealthCheckIntervalSeconds       = 60
	DefaultLogMaxBackups                    = 5
)

type Config struct {
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Monitors []Monitor         `json:"monitors"`
	Settings model.Settings    `json:"settings"`
	LogFile  string            `json:"log_file,omitempty"`
	LogLevel string            `json:"log_level,omitempty"`
}

type Monitor struct {
	ID              string   `json:"id"`
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description,omitempty"`
	Directory       string   `json:"directory"`
	Command         string   `json:"command"`
	FilePatterns    []string `json:"file_patterns"`
	Timeout         int      `json:"timeout"`
	Schedule        string   `json:"schedule,omitempty"`
	Enabled         bool     `json:"enabled,omitempty"`
	DebounceSeconds int      `json:"debounce_seconds,omitempty"`
}

func (c *Config) Validate() error {
	if len(c.Monitors) == 0 {
		return errors.New("at least one monitor must be configured")
	}

	monitorIDs := make(map[string]bool)
	for _, monitor := range c.Monitors {
		if monitor.Directory == "" {
			return errors.New("monitor directory cannot be empty")
		}
		if monitor.Command == "" {
			return errors.New("monitor command cannot be empty")
		}
		if len(monitor.FilePatterns) == 0 {
			return errors.New("monitor must have at least one file pattern: " + monitor.Directory)
		}
		if monitor.Timeout <= 0 {
			return errors.New("monitor timeout must be greater than 0: " + monitor.Directory)
		}

		if monitor.ID != "" {
			if monitorIDs[monitor.ID] {
				return fmt.Errorf("duplicate monitor ID: %s", monitor.ID)
			}
			monitorIDs[monitor.ID] = true
		}
	}

	for _, monitor := range c.Monitors {
		if monitor.Schedule != "" {
			if err := validateCronExpression(monitor.Schedule); err != nil {
				return fmt.Errorf("invalid cron expression %s: %v", monitor.Schedule, err)
			}
		}
	}

	return nil
}

func validateCronExpression(cron string) error {
	if cron == "" {
		return nil
	}

	parts := strings.Fields(cron)
	if len(parts) != 5 {
		return errors.New("cron expression must contain 5 fields")
	}

	return nil
}

func LoadConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file does not exist: %s", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse configuration file: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %v", err)
	}

	applyDefaults(&cfg)

	return &cfg, nil
}

func applyDefaults(cfg *Config) {
	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	if cfg.Settings.MaxConcurrentOperations <= 0 {
		cfg.Settings.MaxConcurrentOperations = DefaultMaxConcurrentOperations
	}

	if cfg.Settings.OperationTimeoutSeconds <= 0 {
		cfg.Settings.OperationTimeoutSeconds = DefaultOperationTimeoutSeconds
	}

	if cfg.Settings.EventChannelBufferSize <= 0 {
		cfg.Settings.EventChannelBufferSize = DefaultEventChannelBufferSize
	}

	if cfg.Settings.MinStabilityTimeMs <= 0 {
		cfg.Settings.MinStabilityTimeMs = DefaultMinStabilityTimeMs
	}

	if cfg.Settings.DirectoryStabilityQuietMs <= 0 {
		cfg.Settings.DirectoryStabilityQuietMs = DefaultDirectoryStabilityQuietMs
	}

	if cfg.Settings.ExecutionDedupIntervalSeconds <= 0 {
		cfg.Settings.ExecutionDedupIntervalSeconds = DefaultExecutionDedupIntervalSeconds
	}

	if cfg.Settings.DirectoryStabilityTimeoutSeconds <= 0 {
		cfg.Settings.DirectoryStabilityTimeoutSeconds = DefaultDirectoryStabilityTimeoutSeconds
	}

	if cfg.Settings.RetryAttempts <= 0 {
		cfg.Settings.RetryAttempts = DefaultRetryAttempts
	}

	if cfg.Settings.RetryDelaySeconds <= 0 {
		cfg.Settings.RetryDelaySeconds = DefaultRetryDelaySeconds
	}

	if cfg.Settings.HealthCheckIntervalSeconds <= 0 {
		cfg.Settings.HealthCheckIntervalSeconds = DefaultHealthCheckIntervalSeconds
	}

	if cfg.Settings.LogMaxBackups <= 0 {
		cfg.Settings.LogMaxBackups = DefaultLogMaxBackups
	}
}
