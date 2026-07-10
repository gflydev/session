package session

import (
	"sync"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

// fakeProvider is a minimal in-test Provider backed by a map.
type fakeProvider struct {
	mu     sync.Mutex
	data   map[string][]byte
	needGC bool
	closed bool
	gcRuns int
}

func newFakeProvider(needGC bool) *fakeProvider {
	return &fakeProvider{data: make(map[string][]byte), needGC: needGC}
}

func (p *fakeProvider) Get(id []byte) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.data[string(id)], nil
}

func (p *fakeProvider) Save(id, data []byte, _ time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	cp := make([]byte, len(data))
	copy(cp, data)
	p.data[string(id)] = cp
	return nil
}

func (p *fakeProvider) Destroy(id []byte) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.data, string(id))
	return nil
}

func (p *fakeProvider) Regenerate(id, newID []byte, _ time.Duration) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	if v, ok := p.data[string(id)]; ok {
		p.data[string(newID)] = v
		delete(p.data, string(id))
	}
	return nil
}

func (p *fakeProvider) Count() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return len(p.data)
}

func (p *fakeProvider) NeedGC() bool { return p.needGC }

func (p *fakeProvider) GC() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.gcRuns++
	return nil
}

func (p *fakeProvider) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.closed = true
	return nil
}

func newTestSession(t *testing.T, p Provider) *Session {
	t.Helper()
	s := NewSession(NewDefaultConfig())
	if err := s.SetProvider(p); err != nil {
		t.Fatalf("SetProvider: %v", err)
	}
	return s
}

func TestGetWithoutProviderErrors(t *testing.T) {
	s := NewSession(NewDefaultConfig())
	if _, err := s.Get(&fasthttp.RequestCtx{}); err != ErrNotSetProvider {
		t.Errorf("Get without provider = %v, want ErrNotSetProvider", err)
	}
}

func TestGetGeneratesNewSessionID(t *testing.T) {
	s := newTestSession(t, newFakeProvider(false))

	ctx := &fasthttp.RequestCtx{}
	store, err := s.Get(ctx)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if len(store.GetSessionID()) == 0 {
		t.Fatal("Get on a new visitor produced an empty session id")
	}
}

func TestSaveThenGetRoundTrip(t *testing.T) {
	p := newFakeProvider(false)
	s := newTestSession(t, p)

	// First request: new visitor, store a value.
	ctx := &fasthttp.RequestCtx{}
	store, err := s.Get(ctx)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	id := append([]byte(nil), store.GetSessionID()...)
	store.Set("user", "alice")
	if err := s.Save(ctx, store); err != nil {
		t.Fatalf("Save: %v", err)
	}

	if p.Count() != 1 {
		t.Fatalf("provider Count() = %d, want 1", p.Count())
	}

	// Second request: carry the session id back in as a cookie.
	ctx2 := &fasthttp.RequestCtx{}
	ctx2.Request.Header.SetCookie(NewDefaultConfig().CookieName, string(id))
	store2, err := s.Get(ctx2)
	if err != nil {
		t.Fatalf("Get (2nd): %v", err)
	}
	if got := store2.Get("user"); got != "alice" {
		t.Errorf(`round-trip Get("user") = %v, want "alice"`, got)
	}
}

func TestRegenerate(t *testing.T) {
	p := newFakeProvider(false)
	s := newTestSession(t, p)

	ctx := &fasthttp.RequestCtx{}
	store, _ := s.Get(ctx)
	id := append([]byte(nil), store.GetSessionID()...)
	store.Set("k", "v")
	if err := s.Save(ctx, store); err != nil {
		t.Fatalf("Save: %v", err)
	}

	ctx2 := &fasthttp.RequestCtx{}
	ctx2.Request.Header.SetCookie(NewDefaultConfig().CookieName, string(id))
	if err := s.Regenerate(ctx2); err != nil {
		t.Fatalf("Regenerate: %v", err)
	}

	if _, ok := p.data[string(id)]; ok {
		t.Error("old session id still present after Regenerate")
	}
	if p.Count() != 1 {
		t.Errorf("provider Count() = %d after Regenerate, want 1", p.Count())
	}
}

func TestDestroy(t *testing.T) {
	p := newFakeProvider(false)
	s := newTestSession(t, p)

	ctx := &fasthttp.RequestCtx{}
	store, _ := s.Get(ctx)
	id := append([]byte(nil), store.GetSessionID()...)
	store.Set("k", "v")
	if err := s.Save(ctx, store); err != nil {
		t.Fatalf("Save: %v", err)
	}

	ctx2 := &fasthttp.RequestCtx{}
	ctx2.Request.Header.SetCookie(NewDefaultConfig().CookieName, string(id))
	if err := s.Destroy(ctx2); err != nil {
		t.Fatalf("Destroy: %v", err)
	}

	if p.Count() != 0 {
		t.Errorf("provider Count() = %d after Destroy, want 0", p.Count())
	}
}

func TestCloseStopsGCAndClosesProvider(t *testing.T) {
	p := newFakeProvider(true)
	s := newTestSession(t, p)

	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	if !p.closed {
		t.Error("provider.Close() was not called by Session.Close()")
	}

	// Second Close must not panic (double close of the GC channel).
	if err := s.Close(); err != nil {
		t.Fatalf("second Close: %v", err)
	}
}

func TestGetRecoversFromCorruptData(t *testing.T) {
	p := newFakeProvider(false)
	s := newTestSession(t, p)

	// Plant an undecodable payload under a known id.
	id := []byte("corrupt-id")
	if err := p.Save(id, []byte("!!!not-valid!!!"), time.Minute); err != nil {
		t.Fatalf("seed Save: %v", err)
	}

	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetCookie(NewDefaultConfig().CookieName, string(id))

	store, err := s.Get(ctx)
	if err != nil {
		t.Fatalf("Get on corrupt data returned error %v, want graceful recovery", err)
	}
	if len(store.GetAll().KV) != 0 {
		t.Errorf("corrupt session yielded %d entries, want a clean store", len(store.GetAll().KV))
	}
}
