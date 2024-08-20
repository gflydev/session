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

func Setup() {
	// Create session manager
	cfg := NewDefaultConfig()
	cfg.EncodeFunc = MSGPEncode
	cfg.DecodeFunc = MSGPDecode
	sessionManager = New(cfg)

	if err := sessionManager.SetProvider(provider); err != nil {
		log.Fatal(err)
	}

	// Register session into core
	core.RegisterSession(&adapter{})
}

// Adapter instance for Session with Core
type adapter struct {
}

func (v *adapter) Set(c *core.Ctx, key string, value interface{}) {
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

func (v *adapter) Get(c *core.Ctx, key string) interface{} {
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
