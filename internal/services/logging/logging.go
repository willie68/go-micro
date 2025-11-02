package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config configuration for the gelf logging
type Config struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`

	Gelfurl  string `yaml:"gelf-url"`
	Gelfport int    `yaml:"gelf-port"`
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
			gelfWriter, err := gelf.NewWriter(fmt.Sprintf("%s:%d", cfg.Gelfurl, cfg.Gelfport))
			if err != nil {
				log.Fatalf("gelf.NewWriter: %s", err)
			}
			gelfWriter.CompressionType = gelf.CompressNone // for debugging only

			hnds = append(hnds, sloggraylog.Option{Level: slog.LevelDebug, Writer: gelfWriter}.NewGraylogHandler())
		}
		Root = slog.New(slogmulti.Fanout(hnds...))
		slog.SetDefault(Root)
	})
}

func New(name string) *slog.Logger {
	return Root.With(slog.String("name", name))
}
