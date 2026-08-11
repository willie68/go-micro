package logging

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
	slogmulti "github.com/samber/slog-multi"
	"github.com/slukits/graylog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config configuration for the gelf logging
type Config struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`

	Gelfurl      string `yaml:"gelf-url"`
	Gelfport     int    `yaml:"gelf-port"`
	Gelfprotocol string `yaml:"gelf-protocol"` // "udp" (default) or "tcp"
}

var (
	once sync.Once
	Root *slog.Logger
)

func init() {
	Root = slog.Default()
}

func Init(cfg Config) {
	once.Do(func() {
		lvl := slog.LevelDebug
		if cfg.Level != "" {
			lvl.UnmarshalText([]byte(cfg.Level))
		}
		hnds := make([]slog.Handler, 0)
		hnds = append(hnds, slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))
		if cfg.Filename != "" {
			hnds = append(hnds, slog.NewTextHandler(&lumberjack.Logger{
				Filename:   cfg.Filename,
				MaxSize:    1, // megabytes
				MaxBackups: 3,
				MaxAge:     28,   //days
				Compress:   true, // disabled by default
			}, nil))
		}
		if cfg.Gelfurl != "" {
			if strings.EqualFold(cfg.Gelfprotocol, "tcp") {
				gLogger, err := graylog.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Gelfurl, cfg.Gelfport))
				if err != nil {
					log.Fatalf("graylog.Dial: %s", err)
				}
				hnds = append(hnds, &levelHandler{level: slog.LevelDebug, Handler: gLogger.Logger.Handler()})
			} else {
				gelfWriter, err := gelf.NewWriter(fmt.Sprintf("%s:%d", cfg.Gelfurl, cfg.Gelfport))
				if err != nil {
					log.Fatalf("gelf.NewWriter: %s", err)
				}
				gelfWriter.CompressionType = gelf.CompressNone // for debugging only

				hnds = append(hnds, sloggraylog.Option{Level: slog.LevelDebug, Writer: gelfWriter}.NewGraylogHandler())
			}
		}
		Root = slog.New(slogmulti.Fanout(hnds...))
		slog.SetDefault(Root)
	})
}

func New(name string) *slog.Logger {
	return Root.With(slog.String("name", name))
}

// levelHandler enforces a minimum level on a slog.Handler.
// Needed because the TCP GELF handler (slukits/graylog) always accepts everything from Debug up.
type levelHandler struct {
	level slog.Leveler
	slog.Handler
}

func (h *levelHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *levelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &levelHandler{level: h.level, Handler: h.Handler.WithAttrs(attrs)}
}

func (h *levelHandler) WithGroup(name string) slog.Handler {
	return &levelHandler{level: h.level, Handler: h.Handler.WithGroup(name)}
}
