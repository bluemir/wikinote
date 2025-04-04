package cache

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/bluemir/wikinote/internal/buildinfo"
)

const (
	HeaderCacheControl = "Cache-Control"
	HeaderETag         = "ETag"
	HeaderIfNoneMatch  = "If-None-Match"
)

/*
https://developer.mozilla.org/ko/docs/Web/HTTP/Caching
https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Cache-Control

	cache-control
		no-store         // cache 하지 않음.
		no-cache         // cache 하지만 매번 검증 후 사용(=`must-revalidata, max-age=0`)
		private          // 단일 사용자를 위해 cache
		public           // 공유 cache 에 의해 cache
		must-revalidate  // cache의 유효기간(max-age)이 지난 후에는 검증 후에 사용
		                 //   만일 `504 gateway timeout` 응답시에는 cache 된 contents 사용할수 있음.
		max-age          // cache의 유효기간 지정
*/

/*
자주 사용되는 cache 설정

- Static file
  - `cache-control: public, max-age=86400`
    - 상대적으로 긴 expire time, shared cache 허용

- 캐시 가능한 API응답
  - `cache-control: private, max-age=60`
    - eg) 권한 check API
    - 짧은 유효기간 + 단일 사용자를 위해서만 cache
    - 1분간은 cache 된 응답을 사용

- 캐시를 disable
  - `cache-control: no-store`

- versioning 된 static File
  - `cache-control: public, max-age=86400000`
    - static file 이 항상 다른 이름으로 제공 될떄
    - eg) `/static/20210815/index.js`
    - 매우 긴 유효기간. expire 해야 할경우 다른 URL 로 제공됨

- 개발시 static file
  - `cache-control: public, max-age=30`
    - 개발도중에는 언제든 static file이 변경될수 있으므로 max-age 를 짧게 설정
  - `cache-control: no-cache`
    - 개발 도중에는 언제나 새 resource 인지 확인하고 사용
*/

type OptionFn func(*Option)
type Option struct {
	MaxAge         int64
	Directive      CacheDirective
	MustRevalidate bool
	// stale-while-revalidate int
	// stale-if-error int
	// immutable
	Extra []string
}

type CacheDirective int

const (
	Undefined CacheDirective = iota
	NoStore
	NoCache
	Public
	Private
)

func Middleware(opts ...OptionFn) gin.HandlerFunc {
	opt := Option{}

	for _, fn := range opts {
		fn(&opt)
	}
	return MiddlewareWithOption(opt)
}

func MiddlewareWithOption(opt Option) gin.HandlerFunc {
	directives := []string{}

	switch opt.Directive {
	case NoStore:
		directives = append(directives, "no-store")
	case NoCache:
		directives = append(directives, "no-cache")
	case Public:
		directives = append(directives, "public")
	case Private:
		directives = append(directives, "private")
	}

	switch opt.Directive {
	case Undefined, Public, Private:
		if opt.MaxAge > 0 {
			directives = append(directives, fmt.Sprintf("max-age=%d", opt.MaxAge))
		}

		if opt.MustRevalidate {
			directives = append(directives, "must-revalidate")
		}
	}

	directives = append(directives, opt.Extra...)

	d := strings.Join(directives, ", ")

	return func(c *gin.Context) {
		c.Writer.Header().Set(HeaderCacheControl, d)
	}
}
func MaxAge(d time.Duration) OptionFn {
	s := int64(d.Seconds())
	return func(opt *Option) {
		opt.MaxAge = s
	}
}
func AlwaysCheckBeforeUseCache(opt *Option) {
	opt.Directive = NoCache
}
func Disable(opt *Option) {
	opt.Directive = NoStore
}
func Shared(opt *Option) {
	opt.Directive = Public
}
func ForLocalCache(opt *Option) {
	opt.Directive = Private
}
func ForStaicFile(opt *Option) {
	//MaxAge(24 * time.Hour)(opt)
	opt.MaxAge = 60 * 60 * 24 // 1d
	opt.Directive = Public
	opt.Extra = []string{
		fmt.Sprintf("stale-while-revalidate=%d", 30*24*60*60), // 30d
	}
}
func ForRevvedResource(opt *Option) {
	//MaxAge(365 * 24 * time.Hour)(opt)
	opt.MaxAge = 60 * 60 * 24 * 365 // 1y
	opt.Directive = Public
}

func ETag() gin.HandlerFunc {
	etag := buildinfo.Signature()[:16]

	return func(c *gin.Context) {
		c.Writer.Header().Set(HeaderETag, etag)

		if match := c.GetHeader(HeaderIfNoneMatch); match != "" {
			if strings.Contains(match, etag) {
				c.Status(http.StatusNotModified)
				c.Abort()
				return
			}
		}

		c.Request.Header.Del("If-Modified-Since") // only accept etag
	}
}
