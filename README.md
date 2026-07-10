# gFly Session

    Copyright © 2023, gFly
    https://www.gfly.dev
    All rights reserved.

# gFly Session - Memory

### Usage

Install
```bash
go get -u github.com/gflydev/session@v1.0.1
go get -u github.com/gflydev/session/memory@v1.0.1
```


Quick usage `main.go`
```go
import (
    "github.com/gflydev/session"
    sessionMemory "github.com/gflydev/session/memory"
)

// Setup session
session.Register(sessionMemory.New())
core.RegisterSession(session.New())
```

Redis provider
```go
import (
    "github.com/gflydev/session"
    sessionRedis "github.com/gflydev/session/redis"
)

session.Register(sessionRedis.New())
core.RegisterSession(session.New())
```

### Controller (Page/API)
```go
// Set session
c.SetSession("foo", utils.UnsafeStr(utils.RandByte(make([]byte, 128))))

// Get session parameter `foo`
foo := c.GetSession("foo")
```

### Shutdown

Call `Close()` on graceful shutdown to stop the background GC goroutine and
release the provider's resources (e.g. the Redis connection pool):
```go
adapter := session.New()
defer adapter.Close()
```