# gFly Session

    Copyright © 2023, gFly
    https://www.gfly.dev
    All rights reserved.

# gFly Session - Memory

### Usage

Install
```bash
go get -u github.com/gflydev/session/memory@v1.0.0
```


Quick usage `main.go`
```go
import (
    "github.com/gflydev/session"
    _ "github.com/gflydev/session/memory"	
)

// Setup session
session.Setup()
```

### Controller (Page/API)
```go
// Note: `c` is `*core.Ctx`

// Set session
c.SetSession("foo", utils.UnsafeStr(utils.RandByte(make([]byte, 128))))

// Get session parameter `foo`
foo := c.GetSession("foo")
```