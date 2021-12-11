package logging

import (
	"fmt"
	"log"

	"github.com/aphistic/golf"
)

/*
ServiceLogger main type for logging
*/
type serviceLogger struct {
	GelfURL    string
	GelfPort   int
	SystemID   string
	Attrs      map[string]interface{}
	gelfActive bool
	c          *golf.Client
}

var Logger serviceLogger

/*
InitGelf initialise gelf logging
*/
func (s *serviceLogger) InitGelf() {
	s.gelfActive = false
	if s.GelfURL != "" {
		s.c, _ = golf.NewClient()
		s.c.Dial(fmt.Sprintf("udp://%s:%d", s.GelfURL, s.GelfPort))

		l, _ := s.c.NewLogger()

		golf.DefaultLogger(l)
		for key, value := range s.Attrs {
			l.SetAttr(key, value)
		}
		l.SetAttr("system_id", s.SystemID)
		s.gelfActive = true
	}
}

/*
Debug log this message at debug level
*/
func (s *serviceLogger) Debug(msg string) {
	if s.gelfActive {
		golf.Dbg(msg)
	}
	log.Println(msg)
}

/*
Debugf log this message at debug level with formatting
*/
func (s *serviceLogger) Debugf(format string, va ...interface{}) {
	if s.gelfActive {
		golf.Dbgf(format, va...)
	}
	log.Printf(format+"\n", va...)
}

/*
Info log this message at info level
*/
func (s *serviceLogger) Info(msg string) {
	if s.gelfActive {
		golf.Info(msg)
	}
	log.Println(msg)
}

/*
Infof log this message at info level with formatting
*/
func (s *serviceLogger) Infof(format string, va ...interface{}) {
	if s.gelfActive {
		golf.Infof(format, va...)
	}
	log.Printf(format+"\n", va...)
}

/*
Alert log this message at alert level
*/
func (s *serviceLogger) Alert(msg string) {
	if s.gelfActive {
		golf.Alert(msg)
	}
	log.Printf("Alert: %s\n", msg)
}

/*
Alertf log this message at alert level with formatting.
*/
func (s *serviceLogger) Alertf(format string, va ...interface{}) {
	if s.gelfActive {
		golf.Alertf(format, va...)
	}
	log.Printf("Alert: %s\n", fmt.Sprintf(format, va...))
}

// Fatal logs a message at level Fatal on the standard logger.
func (s *serviceLogger) Fatal(msg string) {
	if s.gelfActive {
		golf.Crit(msg)
	}
	log.Fatalf("Fatal: %s\n", msg)
}

// Fatalf logs a message at level Fatal on the standard logger with formatting.
func (s *serviceLogger) Fatalf(format string, va ...interface{}) {
	if s.gelfActive {
		golf.Critf(format, va...)
	}
	log.Fatalf("Fatal: %s\n", fmt.Sprintf(format, va...))
}

// Error logs a message at level Error on the standard logger.
func (s *serviceLogger) Error(msg string) {
	if s.gelfActive {
		golf.Err(msg)
	}
	log.Printf("Fatal: %s\n", msg)
}

// Errorf logs a message at level Error on the standard logger with formatting.
func (s *serviceLogger) Errorf(format string, va ...interface{}) {
	if s.gelfActive {
		golf.Errf(format, va...)
	}
	log.Printf("Fatal: %s\n", fmt.Sprintf(format, va...))
}

/*
Close this logging client
*/
func (s *serviceLogger) Close() {
	if s.gelfActive {
		s.c.Close()
	}
}
