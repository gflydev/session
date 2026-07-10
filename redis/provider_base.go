package redis

import (
	"context"
	"io"
	"time"

	"github.com/valyala/bytebufferpool"
)

// scanBatchSize is the COUNT hint passed to Redis SCAN while counting sessions.
const scanBatchSize = 100

func (p *Provider) getRedisSessionKey(sessionID []byte) string {
	key := bytebufferpool.Get()
	// bytebuffer writes never fail (they append to an in-memory slice), so the
	// returned errors are intentionally ignored.
	key.SetString(p.keyPrefix)
	_, _ = key.WriteString(":")
	_, _ = key.Write(sessionID)

	keyStr := key.String()

	bytebufferpool.Put(key)

	return keyStr
}

// Save saves the session data and expiration from the given session id
func (p *Provider) Save(id, data []byte, expiration time.Duration) error {
	key := p.getRedisSessionKey(id)

	return p.db.Set(context.Background(), key, data, expiration).Err()
}

// Regenerate updates the session id and expiration with the new session id
// of the given current session id
func (p *Provider) Regenerate(id, newID []byte, expiration time.Duration) error {
	key := p.getRedisSessionKey(id)
	newKey := p.getRedisSessionKey(newID)

	exists, err := p.db.Exists(context.Background(), key).Result()
	if err != nil {
		return err
	}

	if exists > 0 { // Exist
		if err := p.db.Rename(context.Background(), key, newKey).Err(); err != nil {
			return err
		}

		if err := p.db.Expire(context.Background(), newKey, expiration).Err(); err != nil {
			return err
		}
	}

	return nil
}

// Destroy destroys the session from the given id
func (p *Provider) Destroy(id []byte) error {
	key := p.getRedisSessionKey(id)

	return p.db.Del(context.Background(), key).Err()
}

// Count returns the total of stored sessions.
//
// It uses a cursor-based SCAN instead of KEYS so it does not block the Redis
// server while iterating over the keyspace.
func (p *Provider) Count() int {
	match := p.getRedisSessionKey(all)
	ctx := context.Background()

	var (
		count  int
		cursor uint64
	)

	for {
		keys, next, err := p.db.Scan(ctx, cursor, match, scanBatchSize).Result()
		if err != nil {
			return 0
		}

		count += len(keys)
		cursor = next

		if cursor == 0 {
			break
		}
	}

	return count
}

// NeedGC indicates if the GC needs to be run
func (p *Provider) NeedGC() bool {
	return false
}

// GC destroys the expired sessions
func (p *Provider) GC() error {
	return nil
}

// Close closes the underlying Redis client and its connection pool.
func (p *Provider) Close() error {
	if closer, ok := p.db.(io.Closer); ok {
		return closer.Close()
	}

	return nil
}
