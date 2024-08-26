# gFly Session - Redis

`Use for load balancing`

### Usage

Install
```bash
go get -u github.com/gflydev/session@v1.0.1
go get -u github.com/gflydev/session/redis@v1.0.1
```

Quick usage `main.go`
```go
import (
    "github.com/gflydev/session"
    sessionRedis "github.com/gflydev/session/redis"	
)

// Setup session
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