package logging

import (
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInitAndNew(t *testing.T) {
	Init(Config{
		Level:           "debug",
		VictoriaLogsURL: "http://127.0.0.1:1",
		GelfURL:         "127.0.0.1",
		GelfPort:        1,
		GelfProtocol:    GELFProtocolTCP,
	}, "")
	log := New("unit")
	assert.NotNil(t, log)
	log.Info("hello")
}

func TestShort(t *testing.T) {
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "first line\nsecond", 0)
	assert.Equal(t, "first line", short(&r))
	r = slog.NewRecord(time.Now(), slog.LevelInfo, "  only  ", 0)
	assert.Equal(t, "only", short(&r))
}

func TestGelfExtraConverter(t *testing.T) {
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "msg", 0)
	r.Add("host", "h")
	r.Add("custom", "v")
	r.Add("_already", "x")
	extra := gelfExtraConverter(false, nil, nil, nil, &r)
	assert.Equal(t, "h", extra["host"])
	assert.Equal(t, "v", extra["_custom"])
	assert.Equal(t, "x", extra["_already"])
}

func TestVictoriaLogsHandler(t *testing.T) {
	h := newVictoriaLogsHandler("http://127.0.0.1:1", map[string]any{"service.name": "t"})
	assert.True(t, h.Enabled(nil, slog.LevelInfo))
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "hello", 0)
	assert.NoError(t, h.Handle(nil, r))
	h2 := h.WithAttrs([]slog.Attr{slog.String("k", "v")})
	assert.NotNil(t, h2)
	h3 := h.WithGroup("g")
	assert.NotNil(t, h3)
	assert.Equal(t, h, h.WithGroup(""))
}

func TestGelfTCPHandler(t *testing.T) {
	h := newGelfTCPHandler("127.0.0.1:1", slog.LevelInfo)
	assert.True(t, h.Enabled(nil, slog.LevelInfo))
	assert.False(t, h.Enabled(nil, slog.LevelDebug))
	r := slog.NewRecord(time.Now(), slog.LevelInfo, "gelf", 0)
	assert.NoError(t, h.Handle(nil, r))
	assert.Equal(t, h, h.WithAttrs(nil))
	assert.NotNil(t, h.WithAttrs([]slog.Attr{slog.String("a", "b")}))
	assert.Equal(t, h, h.WithGroup(""))
	assert.NotNil(t, h.WithGroup("g"))
}
