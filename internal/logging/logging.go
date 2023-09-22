package logging

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/aphistic/golf"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LoggingConfig configuration for the gelf logging
type LoggingConfig struct {
	Level    string `yaml:"level"`
	Filename string `yaml:"filename"`

	Gelfurl  string `yaml:"gelf-url"`
	Gelfport int    `yaml:"gelf-port"`
}

// constants for logging levels
const (
	Debug string = "DEBUG"
	Info  string = "INFO"
	Alert string = "ALERT"
	Error string = "ERROR"
	Fatal string = "FATAL"
)

// Levels defining a list of levels
var Levels = []string{Debug, Info, Alert, Error, Fatal}

// ServiceLogger main type for logging
type Logger struct {
	Level      string
	LevelInt   int
	GelfURL    string
	GelfPort   int
	SystemID   string
	Attrs      map[string]any
	Filename   string
	name       string
	gelfActive bool
	c          *golf.Client
	l          *golf.Logger
}

// Root to use for all logging
var Root Logger

func init() {
	Root.SetLevel(Debug)
	Root.name = "Root"
}

func Init(cfg LoggingConfig) {
	Root.SetLevel(cfg.Level)
	Root.GelfURL = cfg.Gelfurl
	Root.GelfPort = cfg.Gelfport
	Root.Init()
}

func New() *Logger {
	var lo *golf.Logger
	if Root.gelfActive {
		lo = Root.l.Clone()
	}
	l := Logger{
		gelfActive: Root.gelfActive,
		c:          Root.c,
		l:          lo,
	}
	return &l
}

// Init initialise logging
func (s *Logger) Init() {
	s.gelfActive = false
	s.dail()
	s.output()
}

func (s *Logger) output() {
	var w io.Writer
	if s.Filename == "" {
		w = os.Stdout
	} else {
		w = io.MultiWriter(&lumberjack.Logger{
			Filename:   s.Filename,
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28,    // days
			Compress:   false, // disabled by default
		}, os.Stdout)
	}
	log.SetOutput(w)
}

func (s *Logger) dail() {
	if s.GelfURL != "" {
		s.c, _ = golf.NewClient()
		err := s.c.Dial(fmt.Sprintf("udp://%s:%d", s.GelfURL, s.GelfPort))
		if err != nil {
			s.Errorf("can't connect to gelf: %v", err)
			return
		}

		l, _ := s.c.NewLogger()

		golf.DefaultLogger(l)
		for key, value := range s.Attrs {
			l.SetAttr(key, value)
		}
		l.SetAttr("system_id", s.SystemID)
		s.l = l
		s.gelfActive = true
	}
}

func (s *Logger) WithLevel(level string) *Logger {
	s.SetLevel(level)
	return s
}

// SetLevel setting the level of this logger
func (s *Logger) SetLevel(level string) {
	switch strings.ToUpper(level) {
	case Debug:
		s.LevelInt = 0
	case Info:
		s.LevelInt = 1
	case Alert:
		s.LevelInt = 2
	case Error:
		s.LevelInt = 3
	case Fatal:
		s.LevelInt = 4
	}
}

func (s *Logger) SetName(name string) {
	s.name = name
}

func (s *Logger) WithName(name string) *Logger {
	s.SetName(name)
	return s
}

// Debug log this message at debug level
func (s *Logger) Debug(m string) {
	if s.LevelInt <= 0 {
		msg := s.format(m)
		if s.gelfActive {
			_ = s.l.Dbg(msg)
		}
		log.Printf("Debug: %s\n", msg)
	}
}

// Debugf log this message at debug level with formatting
func (s *Logger) Debugf(format string, va ...any) {
	if s.LevelInt <= 0 {
		msg := s.format(format, va...)
		if s.gelfActive {
			_ = s.l.Dbg(msg)
		}
		log.Printf("Debug: %s\n", msg)
	}
}

// Info log this message at info level
func (s *Logger) Info(m string) {
	if s.LevelInt <= 1 {
		msg := s.format(m)
		if s.gelfActive {
			_ = s.l.Info(msg)
		}
		log.Printf("Info: %s\n", msg)
	}
}

// Infof log this message at info level with formatting
func (s *Logger) Infof(format string, va ...any) {
	if s.LevelInt <= 1 {
		msg := s.format(format, va...)
		if s.gelfActive {
			_ = s.l.Info(msg)
		}
		log.Printf("Info: %s\n", msg)
	}
}

// Alert log this message at alert level
func (s *Logger) Alert(m string) {
	if s.LevelInt <= 2 {
		msg := s.format(m)
		if s.gelfActive {
			_ = s.l.Alert(msg)
		}
		log.Printf("Alert: %s\n", msg)
	}
}

// Alertf log this message at alert level with formatting.
func (s *Logger) Alertf(format string, va ...any) {
	if s.LevelInt <= 2 {
		msg := s.format(format, va...)
		if s.gelfActive {
			_ = s.l.Alert(msg)
		}
		log.Printf("Alert: %s\n", msg)
	}
}

// Fatal logs a message at level Fatal on the standard logger.
func (s *Logger) Fatal(m string) {
	if s.LevelInt <= 4 {
		msg := s.format(m)
		if s.gelfActive {
			_ = s.l.Crit(msg)
		}
		log.Printf("Fatal: %s\n", msg)
	}
}

// Fatalf logs a message at level Fatal on the standard logger with formatting.
func (s *Logger) Fatalf(format string, va ...any) {
	if s.LevelInt <= 4 {
		msg := s.format(format, va...)
		if s.gelfActive {
			_ = s.l.Crit(msg)
		}
		log.Printf("Fatal: %s\n", msg)
	}
}

// Error logs a message at level Error on the standard logger.
func (s *Logger) Error(m string) {
	if s.LevelInt <= 3 {
		msg := s.format(m)
		if s.gelfActive {
			_ = s.l.Err(msg)
		}
		log.Printf("Error: %s\n", msg)
	}
}

// Errorf logs a message at level Error on the standard logger with formatting.
func (s *Logger) Errorf(format string, va ...any) {
	if s.LevelInt <= 3 {
		msg := s.format(format, va...)
		if s.gelfActive {
			_ = s.l.Err(msg)
		}
		log.Printf("Error: %s\n", msg)
	}
}

// IsDebug this logger is set to debug level
func (s *Logger) IsDebug() bool {
	return s.LevelInt <= 0
}

// IsInfo this logger is set to debug or info level
func (s *Logger) IsInfo() bool {
	return s.LevelInt <= 1
}

// IsAlert this logger is set to debug or info level
func (s *Logger) IsAlert() bool {
	return s.LevelInt <= 2
}

// IsError this logger is set to debug or info level
func (s *Logger) IsError() bool {
	return s.LevelInt <= 3
}

// IsFatal this logger is set to debug or info level
func (s *Logger) IsFatal() bool {
	return s.LevelInt <= 4
}

// format the central format method
func (s *Logger) format(format string, v ...any) string {
	return fmt.Sprintf("%s: %s", s.name, fmt.Sprintf(format, v...))
}

// Close this logging client
func (s *Logger) Close() {
	if s.gelfActive && s.name == Root.name {
		_ = s.c.Close()
	}
}
