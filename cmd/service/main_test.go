package main

import (
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/willie68/go-micro/internal/utils"
)

func TestMain(t *testing.T) {
	ast := assert.New(t)
	configFile = "../../testdata/service_local_minimal.yaml"

	fakeExit := func(int) {
		// fake function does nothing
	}

	p := utils.PatchOSExit(t, fakeExit)
	defer p.Unpatch()

	go func() {
		main()
	}()

	time.Sleep(10 * time.Second)
	ast.NotNil(c)
	c <- syscall.SIGINT
	time.Sleep(1 * time.Second)

	// Assert that os.Exit gets called
	if !p.Called {
		t.Errorf("Expected os.Exit to be called but it was not called")
		return
	}

	// Also, Assert that os.Exit gets called with the correct code
	expectedCalledWith := 0 // no error

	if p.CalledWith != expectedCalledWith {
		t.Errorf("Expected os.Exit to be called with %d but it was called with %d", expectedCalledWith, p.CalledWith)
		return
	}
}
