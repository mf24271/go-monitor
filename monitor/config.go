package monitor

type config struct {
	// where to store the data
	LogPath string
}

// NewConfig returns a config with default values.
func NewConfig() *config {
	return &config{
		LogPath: "./monitor_log",
	}
}
