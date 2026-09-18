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

// victoriaLogsSink owns the HTTP client to VictoriaLogs. send() only enqueues
// (dropping the oldest once victoriaLogsMaxQueue is reached) so slog.Handle
// never blocks on the network. A background loop delivers queued messages and
// retries with exponential backoff while VictoriaLogs is down.
type victoriaLogsSink struct {
	url    string
	client *http.Client

	mu    sync.Mutex
	cond  *sync.Cond
	queue [][]byte
}

func newVictoriaLogsSink(url string) *victoriaLogsSink {
	s := &victoriaLogsSink{
		url: url,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
	s.cond = sync.NewCond(&s.mu)
	go s.sendLoop()
	return s
}

// send queues payload for background delivery. It must not perform HTTP I/O:
// slog-multi Fanout runs handlers sequentially, so a blocked VictoriaLogs
// send would stall stdout/file/gelf logging as well.
func (s *victoriaLogsSink) send(payload []byte) {
	s.mu.Lock()
	if len(s.queue) >= victoriaLogsMaxQueue {
		s.queue = s.queue[1:] // drop oldest
	}
	s.queue = append(s.queue, payload)
	s.cond.Signal()
	s.mu.Unlock()
}

func (s *victoriaLogsSink) sendLoop() {
	backoff := time.Second
	const maxBackoff = 30 * time.Second
	connected := false

	for {
		s.mu.Lock()
		for len(s.queue) == 0 {
			s.cond.Wait()
		}
		payload := s.queue[0]
		s.mu.Unlock()

		if err := s.post(payload); err != nil {
			log.Printf("warn: victorialogs http send %s: %v", s.url, err)
			connected = false
			time.Sleep(backoff)
			if backoff < maxBackoff {
				backoff *= 2
			}
			continue
		}

		if !connected {
			log.Printf("victorialogs: http connected to %s", s.url)
			connected = true
		}
		backoff = time.Second

		s.mu.Lock()
		s.queue = s.queue[1:]
		s.mu.Unlock()
	}
}

func (s *victoriaLogsSink) post(payload []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/stream+json")

	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	_, _ = io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}
