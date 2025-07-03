package ttlcache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSimpleFastFullCircle(t *testing.T) {
	ast := assert.New(t)

	c := New(WithTTL[string, string](1 * time.Second))
	ast.NotNil(c)
	defer c.Close()

	ast.Equal(1*time.Second, c.ttl)

	c.Add("test", "sicher")
	ast.True(c.Has("test"))
	v, ok := c.Get("test")
	ast.True(ok)
	ast.Equal("sicher", *v)

	time.Sleep(2 * time.Second)

	ast.False(c.Has("test"))
	v, ok = c.Get("test")
	ast.False(ok)
	ast.Nil(v)

	time.Sleep(1 * time.Second)
	ast.Equal(0, len(c.items))
}

type kventry struct {
	key   string
	value string
}

func TestPurge(t *testing.T) {
	testdata := make([]kventry, 0)
	for x := range 1000 {
		testdata = append(testdata, kventry{fmt.Sprintf("test%d", x), fmt.Sprintf("value%d", x)})
	}
	ast := assert.New(t)

	c := New(WithTTL[string, string](10 * time.Second))
	ast.NotNil(c)
	defer c.Close()

	for _, kv := range testdata {
		c.Add(kv.key, kv.value)
		ast.True(c.Has(kv.key))
		v, ok := c.Get(kv.key)
		ast.True(ok)
		ast.Equal(kv.value, *v)
	}
	ast.Equal(1000, c.Count())

	c.Purge()
	ast.Equal(0, c.Count())
	ast.Equal(0, len(c.items))
}

func TestDeletion(t *testing.T) {
	ast := assert.New(t)

	c := New(WithTTL[string, string](10 * time.Second))
	ast.NotNil(c)
	defer c.Close()

	c.Add("test", "sicher")
	ast.True(c.Has("test"))
	v, ok := c.Get("test")
	ast.True(ok)
	ast.Equal("sicher", *v)

	c.Delete("test")
	ast.False(c.Has("test"))
	v, ok = c.Get("test")
	ast.False(ok)
	ast.Nil(v)

	ast.Equal(0, len(c.items))
}

func TestNoTTL(t *testing.T) {
	ast := assert.New(t)

	c := New(WithNoTTL[string, string]())
	ast.NotNil(c)
	defer c.Close()

	c.Add("test", "sicher")
	ast.True(c.Has("test"))
	v, ok := c.Get("test")
	ast.True(ok)
	ast.Equal("sicher", *v)

	e, ok := c.items["test"]
	ast.True(ok)
	ast.False(c.isEvicted(e))

	c.deleteEvicted("test")

	ast.Equal(1, c.Count())
	ast.Equal(1, len(c.items))
}

func TestAutoDelete(t *testing.T) {
	ast := assert.New(t)

	c := New[string, string](WithTTL[string, string](2*time.Second), WithAutoDeletion[string, string](5*time.Second))
	ast.NotNil(c)

	c.Add("test", "sicher")
	ast.True(c.Has("test"))
	v, ok := c.Get("test")
	ast.True(ok)
	ast.Equal("sicher", *v)

	time.Sleep(11 * time.Second)

	ast.Equal(0, len(c.items))

	c.Stop()

	c.Add("test", "sicher")
	ast.True(c.Has("test"))

	time.Sleep(11 * time.Second)

	ast.Equal(1, len(c.items))
}

func TestVariableTTL(t *testing.T) {
	ast := assert.New(t)

	c := New(WithTTL[string, string](0))
	ast.NotNil(c)
	defer c.Close()

	c.Add("static", "sicher")
	c.AddWithTTL("onesecond", "sicher", 1*time.Second)
	c.AddWithTTL("tenseconds", "sicher", 10*time.Second)

	ast.True(c.Has("static"))
	v, ok := c.Get("static")
	ast.True(ok)
	ast.Equal("sicher", *v)

	ast.True(c.Has("onesecond"))
	v, ok = c.Get("onesecond")
	ast.True(ok)
	ast.Equal("sicher", *v)

	ast.True(c.Has("tenseconds"))
	v, ok = c.Get("tenseconds")
	ast.True(ok)
	ast.Equal("sicher", *v)

	time.Sleep(2 * time.Second)
	ast.True(c.Has("static"))
	ast.False(c.Has("onesecond"))
	ast.True(c.Has("tenseconds"))

	time.Sleep(10 * time.Second)
	ast.True(c.Has("static"))
	ast.False(c.Has("onesecond"))
	ast.False(c.Has("tenseconds"))
}
