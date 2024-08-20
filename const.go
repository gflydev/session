package session

import (
	"github.com/gflydev/core/utils"
	"time"
)

// TTL Minutes
var ttl = utils.Getenv("SESSION_TTL", 30)

var defaultSessionName = utils.Getenv("SESSION_ID", "gflyid")

const defaultDomain = ""

var defaultExpiration = time.Duration(ttl) * time.Minute

const defaultGCLifetime = 1 * time.Minute
const defaultSecure = true
const defaultSessionIDInURLQuery = false
const defaultSessionIDInHTTPHeader = false
const defaultCookieLen uint32 = 32

// If set the cookie expiration when the browser is closed (-1), set the expiration as a keep alive (2 days)
// so as not to keep dead sessions for a long time
var keepAliveExpiration = time.Duration(ttl) * time.Minute

const expirationAttrKey = "__store:expiration__"
