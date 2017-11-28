package express

import (
	"bytes"
	"context"
	"io"
	"net/http"
)

type (
	HttpCtx struct {
		Context        context.Context
		CancelFunc     context.CancelFunc
		ResponseWriter http.ResponseWriter
		Request        *http.Request
	}
)

func NewhttpCtx() *HttpCtx {
	return &HttpCtx{}
}

func (c *HttpCtx) Reset(rw http.ResponseWriter, r *http.Request) {
	c.ResponseWriter = rw
	c.Request = r
}

func (c *HttpCtx) Copy(buf *bytes.Buffer) (int64, error) {
	return io.Copy(c.ResponseWriter, buf)
}
