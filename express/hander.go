package express

import (
	"net/http"
)

type (
	HandlerFunc    func(ctx *HttpCtx)
	HanderRegister struct {
	}
)

func NewHanderRegister() *HanderRegister {
	return &HanderRegister{}
}

func (h *HanderRegister) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	// host := r.URL.Host
	// path := r.URL.Path

}
