package config

type LoggerConfig struct {
	ConsoleEnabled bool   `envconfig:"LOG_CONSOLE_ENABLED" default:"true"`
	ConsoleLevel   string `envconfig:"LOG_CONSOLE_LEVEL" default:"debug"`
	ConsoleJSON    bool   `envconfig:"LOG_CONSOLE_JSON" default:"false"`

	FileEnabled bool   `envconfig:"LOG_FILE_ENABLED" default:"false"`
	FileLevel   string `envconfig:"LOG_FILE_LEVEL" default:"info"`
	FileJSON    bool   `envconfig:"LOG_FILE_JSON" default:"true"`

	// Directory to log to to when filelogging is enabled
	Directory string `envconfig:"LOG_FILE_DIR" default:"./logs"`
	// Filename is the name of the logfile which will be placed inside the directory
	Filename string `envconfig:"LOG_FILE_NAME" default:"collector.log"`
	// MaxSize the max size in MB of the logfile before it's rolled
	MaxSize int `envconfig:"LOG_MAX_SIZE"`
	// MaxBackups the max number of rolled files to keep
	MaxBackups int `envconfig:"LOG_MAX_BACKUPS" default:"7"`
	// MaxAge the max age in days to keep a logfile
	MaxAge int `enconfig:"LOG_MAX_AGE" default:"1"`
}
