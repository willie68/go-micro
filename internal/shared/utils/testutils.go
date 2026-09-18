package utils

import (
	"os"
	"sync/atomic"
	"testing"

	"github.com/undefinedlabs/go-mpatch"
)

// PatchedOSExit for patching the os exit
type PatchedOSExit struct {
	called     atomic.Bool
	calledWith atomic.Int32
	patchFunc  *mpatch.Patch
}

// Called reports whether os.Exit was invoked.
func (p *PatchedOSExit) Called() bool {
	return p.called.Load()
}

// CalledWith returns the exit code passed to os.Exit.
func (p *PatchedOSExit) CalledWith() int {
	return int(p.calledWith.Load())
}

// PatchOSExit patch the os exit
func PatchOSExit(t *testing.T, mockOSExitImpl func(int)) *PatchedOSExit {
	patchedExit := &PatchedOSExit{}

	patchFunc, err := mpatch.PatchMethod(os.Exit, func(code int) {
		patchedExit.called.Store(true)
		patchedExit.calledWith.Store(int32(code))

		mockOSExitImpl(code)
	})

	if err != nil {
		t.Errorf("Failed to patch os.Exit due to an error: %v", err)

		return nil
	}

	patchedExit.patchFunc = patchFunc

	return patchedExit
}

// Unpatch the os registration
func (p *PatchedOSExit) Unpatch() {
	_ = p.patchFunc.Unpatch()
}
