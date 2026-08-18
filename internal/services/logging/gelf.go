package logging

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	slogcommon "github.com/samber/slog-common"
	sloggraylog "github.com/samber/slog-graylog/v2"
)

const (
	GELFProtocolUDP = "udp"
	GELFProtocolTCP = "tcp"

	// gelfTCPMaxQueue caps how many messages are buffered while the TCP
	// connection to the GELF server is down. Oldest messages are dropped first.
	gelfTCPMaxQueue = 10000
)

var gelfStandardFields = map[string]struct{}{
	"version":       {},
	"host":          {},
	"short_message": {},
	"full_message":  {},
	"timestamp":     {},
	"level":         {},
	"facility":      {},
}

func gelfExtraConverter(addSource bool, replaceAttr func(groups []string, a slog.Attr) slog.Attr, loggerAttr []slog.Attr, groups []string, record *slog.Record) map[string]any {
	extra := sloggraylog.DefaultConverter(addSource, replaceAttr, loggerAttr, groups, record)
	result := make(map[string]any, len(extra))

	for key, value := range extra {
		if key == "" {
			continue
		}
		if _, ok := gelfStandardFields[key]; ok {
			result[key] = value
			continue
		}
		if key[0] == '_' {
			result[key] = value
			continue
		}
		result["_"+key] = value
	}

	return result
}

// gelfTCPHandler is a self-contained slog.Handler that sends GELF messages
// over TCP (plain JSON, null-byte delimited). It has no dependency on an
// external graylog client library; the connection itself is managed by a
// gelfTCPSink that reconnects automatically and queues messages meanwhile.
type gelfTCPHandler struct {
	sink     *gelfTCPSink
	level    slog.Leveler
	hostname string
	attrs    []slog.Attr
	groups   []string
}

func newGelfTCPHandler(addr string, level slog.Leveler) *gelfTCPHandler {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return &gelfTCPHandler{
		sink:     newGelfTCPSink(addr),
		level:    level,
		hostname: hostname,
	}
}

func (h *gelfTCPHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.level.Level()
}

func (h *gelfTCPHandler) Handle(_ context.Context, record slog.Record) error {
	extra := gelfExtraConverter(false, nil, h.attrs, h.groups, &record)

	msg := make(map[string]any, len(extra)+5)
	for k, v := range extra {
		msg[k] = v
	}
	msg["version"] = "1.1"
	msg["host"] = h.hostname
	msg["short_message"] = short(&record)
	msg["full_message"] = strings.TrimSpace(record.Message)
	msg["timestamp"] = float64(record.Time.UnixNano()) / 1e9
	msg["level"] = sloggraylog.LogLevels[record.Level]

	payload, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	h.sink.send(payload)
	return nil
}

func (h *gelfTCPHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	nh := *h
	nh.attrs = slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...)
	return &nh
}

func (h *gelfTCPHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	nh := *h
	nh.groups = append(append([]string{}, h.groups...), name)
	return &nh
}

// gelfTCPSink owns the TCP connection to the GELF server. On write failure
// it closes the connection, queues the message (dropping the oldest once
// gelfTCPMaxQueue is reached) and starts a single reconnect loop with
// exponential backoff. Once reconnected, queued messages are flushed in order.
type gelfTCPSink struct {
	addr string

	mu           sync.Mutex
	conn         net.Conn
	queue        [][]byte
	reconnecting bool
}

func newGelfTCPSink(addr string) *gelfTCPSink {
	s := &gelfTCPSink{addr: addr}
	if !s.connect() {
		s.reconnecting = true
		go s.reconnectLoop()
	}
	return s
}

func (s *gelfTCPSink) connect() bool {
	conn, err := net.Dial("tcp", s.addr)
	if err != nil {
		log.Printf("warn: gelf tcp dial %s: %v", s.addr, err)
		return false
	}
	s.mu.Lock()
	s.conn = conn
	s.mu.Unlock()
	log.Printf("gelf: tcp connected to %s", s.addr)
	return true
}

// send writes payload (a GELF JSON message, without the trailing null byte)
// to the current connection. On failure it queues the message for later
// delivery and (re-)starts the reconnect loop if not already running.
func (s *gelfTCPSink) send(payload []byte) {
	s.mu.Lock()
	conn := s.conn
	s.mu.Unlock()

	if conn != nil {
		if _, err := conn.Write(append(payload, 0)); err == nil {
			return
		}
		s.mu.Lock()
		_ = conn.Close()
		s.conn = nil
		s.mu.Unlock()
	}

	s.mu.Lock()
	if len(s.queue) >= gelfTCPMaxQueue {
		s.queue = s.queue[1:] // drop oldest
	}
	s.queue = append(s.queue, payload)
	needReconnect := !s.reconnecting
	s.reconnecting = true
	s.mu.Unlock()

	if needReconnect {
		go s.reconnectLoop()
	}
}

// reconnectLoop retries connecting with exponential backoff until it
// succeeds and the whole queue has been flushed.
func (s *gelfTCPSink) reconnectLoop() {
	backoff := time.Second
	const maxBackoff = 30 * time.Second
	for {
		if s.connect() && s.flushQueue() {
			s.mu.Lock()
			s.reconnecting = false
			s.mu.Unlock()
			return
		}
		time.Sleep(backoff)
		if backoff < maxBackoff {
			backoff *= 2
		}
	}
}

// flushQueue drains queued messages over the current connection in order.
// Returns false if the connection failed partway through; the remaining
// (unsent) messages stay queued for the next reconnect attempt.
func (s *gelfTCPSink) flushQueue() bool {
	for {
		s.mu.Lock()
		if len(s.queue) == 0 {
			s.mu.Unlock()
			return true
		}
		payload := s.queue[0]
		conn := s.conn
		s.mu.Unlock()

		if conn == nil {
			return false
		}
		if _, err := conn.Write(append(payload, 0)); err != nil {
			log.Printf("warn: gelf tcp flush: %v", err)
			s.mu.Lock()
			_ = conn.Close()
			s.conn = nil
			s.mu.Unlock()
			return false
		}

		s.mu.Lock()
		s.queue = s.queue[1:]
		s.mu.Unlock()
	}
}
