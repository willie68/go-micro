// Package victoriaslog implements a trivially simple slog Handler to upload
// logs to VictoriaMetrics. It is useful when you just need to upload logs
// without pulling the entire opentelemetry package into your system.
package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"maps"
	"net/http"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

const (
	// victoriaLogsMaxQueue caps how many messages are buffered while the HTTP
	// connection to VictoriaLogs is down. Oldest messages are dropped first.
	victoriaLogsMaxQueue = 10000
)

type victoriaLogsHandler struct {
	sink         *victoriaLogsSink
	preset       map[string]any
	currentGroup []string
}

// New takes the base URL of the VictoriaLogs service, the service name, version, and instance
// and returns a new Handler that sends logs to VictoriaLogs.
// The log records are sent in JSON format to the /insert/jsonline endpoint.
// The service name, version, and instance are added as stream fields.
func newVictoriaLogsHandler(baseURL string, preset map[string]any) *victoriaLogsHandler {
	streamFields := strings.Join(slices.Collect(maps.Keys(preset)), ",")
	url := fmt.Sprintf("%s/insert/jsonline?_stream_fields=%s", baseURL, streamFields)

	return &victoriaLogsHandler{
		sink:   newVictoriaLogsSink(url),
		preset: preset,
	}
}

func (h *victoriaLogsHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Implement the logic to check if the given log level is enabled
	return true
}

func findCurrentGroup(m map[string]any, currentGroup []string) map[string]any {
	for _, key := range currentGroup {
		if m[key] == nil {
			x := make(map[string]any)
			m[key] = x
			m = x
		} else {
			v, ok := m[key].(map[string]any)
			if !ok {
				x := make(map[string]any)
				m[key] = x
				m = x
			} else {
				m = v
			}
		}
	}
	return m
}

func (h *victoriaLogsHandler) Handle(ctx context.Context, record slog.Record) error {
	m := maps.Clone(h.preset)

	m["level"] = record.Level.String()
	m["_msg"] = record.Message
	m["_time"] = record.Time.UnixNano()

	fs := runtime.CallersFrames([]uintptr{record.PC})
	f, _ := fs.Next()
	m["source"] = &slog.Source{
		Function: f.Function,
		File:     f.File,
		Line:     f.Line,
	}

	g := findCurrentGroup(m, h.currentGroup)

	record.Attrs(func(a slog.Attr) bool {
		g[a.Key] = a.Value.Any()
		return true
	})

	b, err := json.Marshal(m)
	if err != nil {
		return err
	}

	h.sink.send(b)
	return nil
}

func (h *victoriaLogsHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	m := maps.Clone(h.preset)
	g := findCurrentGroup(m, h.currentGroup)
	for _, a := range attrs {
		g[a.Key] = a.Value.Any()
	}

	return &victoriaLogsHandler{
		sink:         h.sink,
		preset:       m,
		currentGroup: h.currentGroup,
	}
}

func (h *victoriaLogsHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return &victoriaLogsHandler{
		sink:         h.sink,
		preset:       h.preset,
		currentGroup: append(h.currentGroup[:len(h.currentGroup):len(h.currentGroup)], name),
	}
}

// victoriaLogsSink owns the HTTP connection/client to VictoriaLogs. On write failure
// it queues the message (dropping the oldest once victoriaLogsMaxQueue is reached)
// and starts a single reconnect loop with exponential backoff. Once reconnected,
// queued messages are flushed in order.
type victoriaLogsSink struct {
	url string

	mu           sync.Mutex
	queue        [][]byte
	reconnecting bool
}

func newVictoriaLogsSink(url string) *victoriaLogsSink {
	s := &victoriaLogsSink{url: url}
	// Try initial connection, but don't block on startup
	if !s.testConnection() {
		s.mu.Lock()
		s.reconnecting = true
		s.mu.Unlock()
		go s.reconnectLoop()
	} else {
		log.Printf("victorialogs: http connected to %s", url)
	}
	return s
}

// testConnection checks if VictoriaLogs is reachable
func (s *victoriaLogsSink) testConnection() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewReader([]byte("{}")))
	if err != nil {
		log.Printf("warn: victorialogs test connection create request: %v", err)
		return false
	}
	req.Header.Set("Content-Type", "application/stream+json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("warn: victorialogs http test %s: %v", s.url, err)
		return false
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		log.Printf("warn: victorialogs http test status %d", resp.StatusCode)
		return false
	}
	return true
}

// send writes payload (a JSON-formatted log entry) to VictoriaLogs.
// On failure it queues the message for later delivery and (re-)starts
// the reconnect loop if not already running.
func (s *victoriaLogsSink) send(payload []byte) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewReader(payload))
	if err != nil {
		log.Printf("warn: victorialogs http send create request: %v", err)
		goto queue
	}
	req.Header.Set("Content-Type", "application/stream+json")

	{
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("warn: victorialogs http send: %v", err)
			goto queue
		}
		defer resp.Body.Close()
		_, _ = io.ReadAll(resp.Body)

		if resp.StatusCode < 400 {
			return // success
		}

		log.Printf("warn: victorialogs http send status %d", resp.StatusCode)
	}

queue:
	s.mu.Lock()
	if len(s.queue) >= victoriaLogsMaxQueue {
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
func (s *victoriaLogsSink) reconnectLoop() {
	backoff := time.Second
	const maxBackoff = 30 * time.Second
	for {
		if s.testConnection() && s.flushQueue() {
			s.mu.Lock()
			s.reconnecting = false
			s.mu.Unlock()
			log.Printf("victorialogs: reconnected and queue flushed")
			return
		}
		time.Sleep(backoff)
		if backoff < maxBackoff {
			backoff *= 2
		}
	}
}

// flushQueue drains queued messages over the HTTP connection in order.
// Returns false if the connection failed partway through; the remaining
// (unsent) messages stay queued for the next reconnect attempt.
func (s *victoriaLogsSink) flushQueue() bool {
	for {
		s.mu.Lock()
		if len(s.queue) == 0 {
			s.mu.Unlock()
			return true
		}
		payload := s.queue[0]
		s.mu.Unlock()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewReader(payload))
		cancel()
		if err != nil {
			log.Printf("warn: victorialogs flush create request: %v", err)
			return false
		}
		req.Header.Set("Content-Type", "application/stream+json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("warn: victorialogs flush: %v", err)
			return false
		}
		defer resp.Body.Close()
		_, _ = io.ReadAll(resp.Body)

		if resp.StatusCode >= 400 {
			log.Printf("warn: victorialogs flush status %d", resp.StatusCode)
			return false
		}

		s.mu.Lock()
		s.queue = s.queue[1:]
		s.mu.Unlock()
	}
}
