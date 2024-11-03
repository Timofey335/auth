package env

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	logLevelEnv          = "LOG_LEVEL"
	logFileNameEnv       = "LOG_FILENAME"
	logFileMaxSizeEnv    = "LOG_FILE_MAX_SIZE"
	logFileMaxBackupsEnv = "LOG_FILE_MAX_BACKUPS"
	logFileMaxAgeEnv     = "LOG_FILE_MAX_AGE"
)

type loggerConfig struct {
	logLevel          string
	logFilename       string
	logFileMaxSize    int64
	logFileMaxBackups int64
	logFileMaxAge     int64
}

// NewLoggerConfig - конфигурация для логгера
func NewLoggerConfig() (*loggerConfig, error) {
	logLevel := os.Getenv(logLevelEnv)
	if len(logLevel) == 0 {
		return nil, errors.New("log level not found")
	}

	logFilename := os.Getenv(logFileNameEnv)
	if len(logFilename) == 0 {
		return nil, errors.New("log filename not found")
	}

	logFileMaxSizeStr := os.Getenv(logFileMaxSizeEnv)
	if len(logFileMaxSizeStr) == 0 {
		return nil, errors.New("log file size not found")
	}

	logFileMaxSize, err := strconv.ParseInt(logFileMaxSizeStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log file max size")
	}

	logFileMaxBackupsStr := os.Getenv(logFileMaxBackupsEnv)
	if len(logFileMaxBackupsStr) == 0 {
		return nil, errors.New("log file max size not found")
	}

	logFileMaxBackups, err := strconv.ParseInt(logFileMaxBackupsStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log file max backups")
	}

	logFileMaxAgeStr := os.Getenv(logFileMaxAgeEnv)
	if len(logFileMaxAgeStr) == 0 {
		return nil, errors.New("log file max age not found")
	}

	logFileMaxAge, err := strconv.ParseInt(logFileMaxAgeStr, 10, 64)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse log file max age")
	}

	return &loggerConfig{
		logLevel:          logLevel,
		logFilename:       logFilename,
		logFileMaxSize:    logFileMaxSize,
		logFileMaxBackups: logFileMaxBackups,
		logFileMaxAge:     logFileMaxAge,
	}, nil
}

// LogLevel
func (cfg *loggerConfig) LogLevel() string {
	return cfg.logLevel
}

// LogFilename
func (cfg *loggerConfig) LogFilename() string {
	return cfg.logFilename
}

// LogFileMaxSize
func (cfg *loggerConfig) LogFileMaxSize() int {
	return int(cfg.logFileMaxSize)
}

// LogFileMaxBackups
func (cfg *loggerConfig) LogFileMaxBackups() int {
	return int(cfg.logFileMaxBackups)
}

// LogFileMaxAge
func (cfg *loggerConfig) LogFileMaxAge() int {
	return int(cfg.logFileMaxAge)
}
