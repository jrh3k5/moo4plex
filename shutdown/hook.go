package shutdown

import (
	"context"
	"fmt"
	"sync"
)

// Hook is the definition of an action that can be executed on shutdown
type Hook func(ctx context.Context) error

// HookRegister is a registrar of shutdown hooks
type HookRegistrar interface {
	// AddHook adds the given hook to be executed on application shutdown
	AddHook(hook Hook)
}

// SliceHookRegistrar is an implementation of HookRegistrar backed by a slice
type SliceHookRegistrar struct {
	hooksMutex *sync.Mutex
	hooks      []Hook
}

func NewSliceHookRegistrar() *SliceHookRegistrar {
	return &SliceHookRegistrar{
		hooksMutex: &sync.Mutex{},
	}
}

func (s *SliceHookRegistrar) AddHook(hook Hook) {
	s.hooksMutex.Lock()
	defer s.hooksMutex.Unlock()

	s.hooks = append(s.hooks, hook)
}

func (s *SliceHookRegistrar) ExecuteHooks(ctx context.Context) {
	s.hooksMutex.Lock()
	defer s.hooksMutex.Unlock()

	for _, hook := range s.hooks {
		if hookErr := hook(ctx); hookErr != nil {
			fmt.Printf("SHUTDOWN HOOK ERROR: %v", hookErr)
		}
	}
}
