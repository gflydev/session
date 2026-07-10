package session

import (
	"github.com/gflydev/core"
	"github.com/gflydev/core/log"
)

// ======================================================================================
// 								Session and Providers
// ======================================================================================

var provider Provider

// Register assign session provider type `memory`, `redis`....
func Register(p Provider) {
	provider = p
}

// ======================================================================================
// 								Session and Core
// ======================================================================================

var sessionManager *Session

// New create session manager and session adapter
func New() *Adapter {
	// Create session manager
	cfg := NewDefaultConfig()
	cfg.EncodeFunc = MSGPEncode
	cfg.DecodeFunc = MSGPDecode
	sessionManager = NewSession(cfg)

	if err := sessionManager.SetProvider(provider); err != nil {
		log.Fatal(err)
	}

	return &Adapter{}
}

// Adapter instance for Session with Core
type Adapter struct {
}

func (v *Adapter) Set(c *core.Ctx, key string, value interface{}) {
	store, err := sessionManager.Get(c.Root())
	if err != nil {
		log.Errorf("session set failed to load store: %v", err)
		return
	}

	defer func() {
		if err := sessionManager.Save(c.Root(), store); err != nil {
			log.Errorf("session set failed to save store: %v", err)
		}
	}()

	store.Set(key, value)
}

func (v *Adapter) Get(c *core.Ctx, key string) interface{} {
	store, err := sessionManager.Get(c.Root())
	if err != nil {
		log.Errorf("session get failed to load store: %v", err)
		return nil
	}

	defer func() {
		if err := sessionManager.Save(c.Root(), store); err != nil {
			log.Errorf("session get failed to save store: %v", err)
		}
	}()

	return store.Get(key)
}

// Close releases resources held by the session manager, including the GC
// goroutine and the underlying provider connection.
func (v *Adapter) Close() error {
	if sessionManager == nil {
		return nil
	}

	return sessionManager.Close()
}
