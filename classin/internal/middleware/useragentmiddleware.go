package middleware

import (
	"context"
	"net/http"
)

type UserAgentMiddleware struct {
}

func NewUserAgentMiddleware() *UserAgentMiddleware {
	return &UserAgentMiddleware{}
}

func (m *UserAgentMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		val := r.Header.Get("User-Agent")
		if len(val) == 0 {
			val = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36"
		}
		reqCtx := r.Context()
		ctx := context.WithValue(reqCtx, "User-Agent", val)
		newReq := r.WithContext(ctx)
		next(w, newReq)
	}
}
