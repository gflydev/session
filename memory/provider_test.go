package memory

import (
	"testing"
	"time"
)

func TestSaveGet(t *testing.T) {
	p := New()

	id := []byte("id-1")
	if err := p.Save(id, []byte("payload"), time.Minute); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := p.Get(id)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if string(got) != "payload" {
		t.Errorf("Get = %q, want %q", got, "payload")
	}
}

func TestGetMissing(t *testing.T) {
	p := New()
	got, err := p.Get([]byte("missing"))
	if err != nil {
		t.Fatalf("Get missing: %v", err)
	}
	if got != nil {
		t.Errorf("Get missing = %q, want nil", got)
	}
}

func TestDestroy(t *testing.T) {
	p := New()
	id := []byte("id-2")
	_ = p.Save(id, []byte("x"), time.Minute)

	if err := p.Destroy(id); err != nil {
		t.Fatalf("Destroy: %v", err)
	}
	if p.Count() != 0 {
		t.Errorf("Count after Destroy = %d, want 0", p.Count())
	}
}

func TestRegenerate(t *testing.T) {
	p := New()
	old := []byte("old")
	fresh := []byte("new")
	_ = p.Save(old, []byte("data"), time.Minute)

	if err := p.Regenerate(old, fresh, time.Minute); err != nil {
		t.Fatalf("Regenerate: %v", err)
	}

	if got, _ := p.Get(old); got != nil {
		t.Errorf("old id still resolves after Regenerate: %q", got)
	}
	if got, _ := p.Get(fresh); string(got) != "data" {
		t.Errorf("new id = %q, want %q", got, "data")
	}
}

func TestCount(t *testing.T) {
	p := New()
	_ = p.Save([]byte("a"), []byte("1"), time.Minute)
	_ = p.Save([]byte("b"), []byte("2"), time.Minute)

	if p.Count() != 2 {
		t.Errorf("Count = %d, want 2", p.Count())
	}
}

func TestGCRemovesExpired(t *testing.T) {
	p := New()

	// Expired: a tiny TTL with a lastActiveTime in the distant past.
	_ = p.Save([]byte("expired"), []byte("x"), time.Nanosecond)
	// Non-expiring: expiration 0 means keep forever.
	_ = p.Save([]byte("keep"), []byte("y"), 0)

	time.Sleep(2 * time.Millisecond)

	if err := p.GC(); err != nil {
		t.Fatalf("GC: %v", err)
	}

	if got, _ := p.Get([]byte("expired")); got != nil {
		t.Errorf("expired session survived GC: %q", got)
	}
	if got, _ := p.Get([]byte("keep")); string(got) != "y" {
		t.Errorf("non-expiring session was collected, got %q", got)
	}
}

func TestNeedGC(t *testing.T) {
	if !New().NeedGC() {
		t.Error("memory provider NeedGC() = false, want true")
	}
}

func TestClose(t *testing.T) {
	p := New()
	_ = p.Save([]byte("a"), []byte("1"), time.Minute)

	if err := p.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if p.Count() != 0 {
		t.Errorf("Count after Close = %d, want 0", p.Count())
	}
}
