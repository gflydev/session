package session

import (
	"testing"
	"time"
)

func TestStoreSetGetDelete(t *testing.T) {
	s := NewStore()

	s.Set("foo", "bar")
	if got := s.Get("foo"); got != "bar" {
		t.Errorf(`Get("foo") = %v, want "bar"`, got)
	}

	s.SetBytes([]byte("baz"), 1)
	if got := s.GetBytes([]byte("baz")); got != 1 {
		t.Errorf(`GetBytes("baz") = %v, want 1`, got)
	}

	s.Delete("foo")
	if got := s.Get("foo"); got != nil {
		t.Errorf(`Get("foo") after Delete = %v, want nil`, got)
	}

	s.DeleteBytes([]byte("baz"))
	if got := s.Get("baz"); got != nil {
		t.Errorf(`Get("baz") after DeleteBytes = %v, want nil`, got)
	}
}

func TestStoreFlush(t *testing.T) {
	s := NewStore()
	s.Set("a", 1)
	s.Set("b", 2)

	s.Flush()

	if len(s.GetAll().KV) != 0 {
		t.Errorf("after Flush, store has %d entries, want 0", len(s.GetAll().KV))
	}
}

func TestStoreSessionID(t *testing.T) {
	s := NewStore()
	id := []byte("session-123")
	s.SetSessionID(id)

	if got := string(s.GetSessionID()); got != "session-123" {
		t.Errorf("GetSessionID() = %q, want %q", got, "session-123")
	}
}

func TestStoreExpirationDefault(t *testing.T) {
	s := NewStore()
	s.defaultExpiration = 5 * time.Minute

	if s.HasExpirationChanged() {
		t.Error("HasExpirationChanged() = true for a fresh store, want false")
	}
	if got := s.GetExpiration(); got != 5*time.Minute {
		t.Errorf("GetExpiration() = %v, want default 5m", got)
	}
}

func TestStoreExpirationOverride(t *testing.T) {
	s := NewStore()
	s.defaultExpiration = 5 * time.Minute

	if err := s.SetExpiration(30 * time.Second); err != nil {
		t.Fatalf("SetExpiration: %v", err)
	}

	if !s.HasExpirationChanged() {
		t.Error("HasExpirationChanged() = false after SetExpiration, want true")
	}
	if got := s.GetExpiration(); got != 30*time.Second {
		t.Errorf("GetExpiration() = %v, want 30s", got)
	}
}

func TestStoreReset(t *testing.T) {
	s := NewStore()
	s.Set("k", "v")
	s.SetSessionID([]byte("id"))
	s.defaultExpiration = time.Minute

	s.Reset()

	if len(s.GetAll().KV) != 0 {
		t.Errorf("after Reset, %d entries remain, want 0", len(s.GetAll().KV))
	}
	if len(s.GetSessionID()) != 0 {
		t.Errorf("after Reset, sessionID = %q, want empty", s.GetSessionID())
	}
	if s.defaultExpiration != 0 {
		t.Errorf("after Reset, defaultExpiration = %v, want 0", s.defaultExpiration)
	}
}
