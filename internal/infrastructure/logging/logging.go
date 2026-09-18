// logging setup using slog with optional graylog forwarding
package logging

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"strings"
	"sync"

	"github.com/Graylog2/go-gelf/gelf"
	sloggraylog "github.com/samber/slog-graylog/v2"
	slogmulti "github.com/samber/slog-multi"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config for logging.
type Config struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`

	GelfURL      string `yaml:"gelf-url"`
	GelfPort     int    `yaml:"gelf-port"`
	GelfProtocol string `yaml:"gelf-protocol"` // "udp" (default) or "tcp"

	VictoriaLogsURL string `yaml:"victoria-logs-url"`
}

var (
	once sync.Once
	// Root is the global logger.
	Root *slog.Logger
)

func init() {
	Root = slog.Default()
}

// Init sets up global structured logging.
func Init(cfg Config, service string) {
	once.Do(func() {
		hostname, err := os.Hostname()
		if err != nil {
			hostname = ""
		}
		if service == "" {
			service = "go-micro"
		}
		lvl := slog.LevelInfo
		if cfg.Level != "" {
			_ = lvl.UnmarshalText([]byte(cfg.Level))
		}
		hnds := []slog.Handler{
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}),
		}
		if cfg.Filename != "" {
			hnds = append(hnds, slog.NewTextHandler(&lumberjack.Logger{
				Filename:   cfg.Filename,
				MaxSize:    1,
				MaxBackups: 3,
				MaxAge:     28,
				Compress:   true,
			}, &slog.HandlerOptions{Level: lvl}))
		}
		if cfg.GelfURL != "" {
			addr := fmt.Sprintf("%s:%d", cfg.GelfURL, cfg.GelfPort)
			if strings.EqualFold(cfg.GelfProtocol, GELFProtocolTCP) {
				hnds = append(hnds, newGelfTCPHandler(addr, lvl))
			} else {
				gelfWriter, err := gelf.NewWriter(addr)
				if err != nil {
					log.Printf("warn: gelf writer: %v", err)
				} else {
					gelfWriter.CompressionType = gelf.CompressGzip

					hnds = append(hnds, sloggraylog.Option{
						Level:     lvl,
						Writer:    gelfWriter,
						Converter: gelfExtraConverter,
					}.NewGraylogHandler())
				}
			}
		}
		if cfg.VictoriaLogsURL != "" {
			handler := newVictoriaLogsHandler(
				cfg.VictoriaLogsURL, // VictoriaLogs base URL
				map[string]any{"service.name": service, "service.instance": hostname},
			)
			hnds = append(hnds, handler)
		}
		Root = slog.New(slogmulti.Fanout(hnds...)).With("service", service)
		slog.SetDefault(Root)
	})
}

// New returns a logger with a name attribute.
func New(name string) *slog.Logger {
	return Root.With(slog.String("loggername", name))
}

func short(record *slog.Record) string {
	msg := strings.TrimSpace(record.Message)
	if i := strings.IndexRune(msg, '\n'); i > 0 {
		return msg[:i]
	}
	return msg
}
