package session

import (
	"github.com/gflydev/core/utils"
	"github.com/valyala/fasthttp"
)

// NewDefaultConfig returns a new default configuration
func NewDefaultConfig() Config {
	config := Config{
		CookieName:              defaultSessionName,
		Domain:                  defaultDomain,
		Expiration:              defaultExpiration,
		GCLifetime:              defaultGCLifetime,
		Secure:                  defaultSecure,
		SessionIDInURLQuery:     defaultSessionIDInURLQuery,
		SessionNameInURLQuery:   defaultSessionName,
		SessionIDInHTTPHeader:   defaultSessionIDInHTTPHeader,
		SessionNameInHTTPHeader: defaultSessionName,
		cookieLen:               defaultCookieLen,
	}

	// default sessionIdGeneratorFunc
	config.SessionIDGeneratorFunc = config.defaultSessionIDGenerator

	// default isSecureFunc
	config.IsSecureFunc = config.defaultIsSecureFunc

	return config
}

func (c *Config) defaultSessionIDGenerator() []byte {
	return utils.RandByte(make([]byte, c.cookieLen))
}

func (c *Config) defaultIsSecureFunc(ctx *fasthttp.RequestCtx) bool {
	return ctx.IsTLS()
}
