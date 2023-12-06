package utils

import (
	"os"
	"testing"

	"github.com/undefinedlabs/go-mpatch"
)

// PatchedOSExit for patching the os exit
type PatchedOSExit struct {
	Called     bool
	CalledWith int
	patchFunc  *mpatch.Patch
}

// PatchOSExit patch the os exit
func PatchOSExit(t *testing.T, mockOSExitImpl func(int)) *PatchedOSExit {
	patchedExit := &PatchedOSExit{Called: false}

	patchFunc, err := mpatch.PatchMethod(os.Exit, func(code int) {
		patchedExit.Called = true
		patchedExit.CalledWith = code

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
