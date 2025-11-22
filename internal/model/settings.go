package model

type Settings struct {
	LogLevel      string `json:"log_level,omitempty"`
	LogFile       string `json:"log_file,omitempty"`
	LogMaxSize    int64  `json:"log_max_size,omitempty"`
	LogMaxBackups int    `json:"log_max_backups,omitempty"`
	LogShowCaller bool   `json:"log_show_caller,omitempty"`

	MaxConcurrentOperations int `json:"max_concurrent_operations,omitempty"`
	OperationTimeoutSeconds int `json:"operation_timeout_seconds,omitempty"`

	FileWatcherBufferSize uint32 `json:"file_watcher_buffer_size,omitempty"`

	EventChannelBufferSize int `json:"event_channel_buffer_size,omitempty"`
	MinStabilityTimeMs     int `json:"min_stability_time_ms,omitempty"`

	ExecutionDedupIntervalSeconds int `json:"execution_dedup_interval_seconds,omitempty"`

	DirectoryStabilityQuietMs        int `json:"directory_stability_quiet_ms,omitempty"`
	DirectoryStabilityTimeoutSeconds int `json:"directory_stability_timeout_seconds,omitempty"`

	RetryAttempts     int `json:"retry_attempts,omitempty"`
	RetryDelaySeconds int `json:"retry_delay_seconds,omitempty"`

	HealthCheckIntervalSeconds int `json:"health_check_interval_seconds,omitempty"`
}
