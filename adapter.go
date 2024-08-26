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
		log.Fatal(err)
	}

	defer func() {
		if err := sessionManager.Save(c.Root(), store); err != nil {
			log.Fatal(err)
		}
	}()

	store.Set(key, value)
}

func (v *Adapter) Get(c *core.Ctx, key string) interface{} {
	store, err := sessionManager.Get(c.Root())
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := sessionManager.Save(c.Root(), store); err != nil {
			log.Fatal(err)
		}
	}()

	return store.Get(key)
}
