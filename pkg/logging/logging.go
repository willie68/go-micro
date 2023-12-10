package logging

import "log/slog"

func New() *slog.Logger {
	l := slog.New()
	return l
}
