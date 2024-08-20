package memory

import "github.com/gflydev/session"

// ========================================================================================
// 										Initial
// ========================================================================================

// Auto initial redis session and register to session manager
func init() {
	provider, _ := New()

	session.Register(provider)
}
