package logging

import "log/slog"

type Level string

const (
	LevelDebug Level = "debug"
	LevelInfo  Level = "info"
	LevelWarn  Level = "warn"
	LevelError Level = "error"
)

func (l Level) String() string {
	return string(l)
}

func (l Level) Slog() (slog.Level, error) {
	var level slog.Level
	if err := level.UnmarshalText([]byte(l)); err != nil {
		return 0, err
	}
	return level, nil
}
