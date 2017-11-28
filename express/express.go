package express

import (
	"net/http"
)

type Express struct {
	route     Route
	headerMap map[string]string
}

func (e *Express) ServeHTTP(rw http.ResponseWriter, r *http.Request) {

}

func (e *Express) Use(path string, fn func(ctx *HttpCtx)) {

}

func New() {

}
