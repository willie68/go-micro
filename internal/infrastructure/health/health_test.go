package health

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/samber/do/v2"
	"github.com/stretchr/testify/assert"
)

var (
	hs *Service

	_ Check = &MyCheck{}
)

type MyCheck struct {
	name  string
	fired atomic.Bool
	times atomic.Int32
	ret   bool
	err   error
}

// Check implements Check.
func (m *MyCheck) Check() (bool, error) {
	m.fired.Store(true)
	m.times.Add(1)
	return m.ret, m.err
}

// CheckName implements Check.
func (m *MyCheck) CheckName() string {
	return m.name
}

func InitHealth(inj do.Injector, ast *assert.Assertions) {
	if hs == nil {
		do.ProvideValue(inj, Config{
			Period:     10,
			StartDelay: 1,
		})
		err := Provide(inj)
		ast.Nil(err)
		hs = do.MustInvoke[*Service](inj)
		ast.NotNil(hs)
	}
}

func ShutdownHealth(inj do.Injector) {
	_ = do.Shutdown[*Service](inj)
	hs = nil
}

func TestHealthBase(t *testing.T) {
	ast := assert.New(t)
	inj := do.New()

	InitHealth(inj, ast)

	// check if the service is injected
	hsdi, err := do.Invoke[*Service](inj)
	ast.Nil(err)
	ast.NotNil(hsdi)

	chk := MyCheck{
		name: "myname",
		ret:  true,
		err:  nil,
	}
	hs.Register(&chk)

	time.Sleep(12 * time.Second)

	ast.True(chk.fired.Load())
	ast.Equal(int32(1), chk.times.Load())

	hs.CheckHealthCheckTimer()
	ast.True(hs.Readyz())
	ast.Equal(0, len(hs.Message().Messages))

	ok := hs.Unregister(chk.name)
	ast.True(ok)

	ShutdownHealth(inj)
}

func TestMessage(t *testing.T) {
	ast := assert.New(t)
	inj := do.New()

	n := time.Now()
	InitHealth(inj, ast)
	hs.reg.Lock()
	hs.lastChecked = n
	hs.messages = make([]string, 0)
	hs.reg.Unlock()
	msg := hs.Message()
	ast.Equal(n.String(), msg.LastCheck)
	ast.Equal(0, len(msg.Messages))

	ShutdownHealth(inj)
}

func TestHealthUnhealthy(t *testing.T) {
	ast := assert.New(t)
	inj := do.New()

	InitHealth(inj, ast)

	hsdi, err := do.Invoke[*Service](inj)
	ast.Nil(err)
	ast.NotNil(hsdi)

	chk := MyCheck{
		name: "myname",
		ret:  false,
		err:  errors.New("error"),
	}
	err = Register(inj, &chk)
	ast.Nil(err)

	time.Sleep(12 * time.Second)

	ast.True(chk.fired.Load())
	ast.Equal(int32(1), chk.times.Load())

	msg := hs.Message()
	ast.NotNil(msg.LastCheck)
	ast.Equal(1, len(msg.Messages))
	ast.Equal("myname: error", msg.Messages[0])

	err = Unregister(inj, chk.name)
	ast.Nil(err)

	ShutdownHealth(inj)
}

type fakeServiceName struct{}

func (fakeServiceName) ServiceName() string { return "coverage" }

func TestHealthyzAndHandler(t *testing.T) {
	ast := assert.New(t)
	inj := do.New()
	InitHealth(inj, ast)
	defer ShutdownHealth(inj)

	ast.True(hs.Healthyz())
	_ = hs.LastChecked()
	hs.Register(&MyCheck{name: "dup", ret: true})
	hs.Register(&MyCheck{name: "dup", ret: true})

	do.ProvideValue(inj, fakeServiceName{})
	h := NewHealthHandler(inj)
	path, mux := h.Routes()
	ast.Equal("/", path)

	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/livez", nil))
	ast.Equal(http.StatusOK, rec.Code)

	hs.reg.Lock()
	hs.healthy = false
	hs.readyz = false
	hs.reg.Unlock()
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/livez", nil))
	ast.Equal(http.StatusServiceUnavailable, rec.Code)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodHead, "/livez", nil))
	ast.Equal(http.StatusNoContent, rec.Code)
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodHead, "/readyz", nil))
	ast.Equal(http.StatusNoContent, rec.Code)

	hs.reg.Lock()
	hs.healthy = true
	hs.readyz = true
	hs.reg.Unlock()
	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodHead, "/livez", nil))
	ast.Equal(http.StatusNoContent, rec.Code)

	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/readyz", nil))
	ast.Equal(http.StatusOK, rec.Code)

	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodHead, "/readyz", nil))
	ast.Equal(http.StatusNoContent, rec.Code)

	rec = httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	ast.Contains(rec.Body.String(), "coverage")
}
