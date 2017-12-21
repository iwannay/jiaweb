package jiaweb

import (
	"io/ioutil"
	"net/http"
)

type (
	Request struct {
		*http.Request
		httpCtx     *HttpContext
		body        []byte
		isReadyBody bool
	}
)

func (req *Request) reset(r *http.Request, ctx *HttpContext) {
	req.Request = r
	req.httpCtx = ctx
}

func (req *Request) Get(key string) string {
	return req.URL.Query().Get(key)
}

func (req *Request) Post(key string) string {
	return req.PostFormValue(key)
}

func (req *Request) Body() []byte {
	if !req.isReadyBody {
		bts, err := ioutil.ReadAll(req.Request.Body)
		if err != nil {
			return []byte{}
		}
		req.isReadyBody = true
		req.body = bts
	}
	return req.body
}

func (req *Request) RemoteIP() string {
	return req.Request.RemoteAddr
}

func (req *Request) Url() string {
	return req.URL.String()
}
