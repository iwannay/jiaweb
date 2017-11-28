package express

const (
	HTTPMethod_Any       = "ANY"
	HTTPMethod_GET       = "GET"
	HTTPMethod_POST      = "POST"
	HTTPMethod_PUT       = "PUT"
	HTTPMethod_DELETE    = "DELETE"
	HTTPMethod_PATCH     = "PATCH"
	HTTPMethod_HiJack    = "HIJACK"
	HTTPMethod_WebSocket = "WEBSOCKET"
	HTTPMethod_HEAD      = "HEAD"
	HTTPMethod_OPTIONS   = "OPTIONS"
)

type (
	RouteIf interface {
		Middleware(func(ctx *HttpCtx)) *Route
		Get(path string, fn func(ctx *HttpCtx)) *Route
		Put(path string, handle)
		Parse(path string) func(ctx *HttpCtx)
		Register(path string, fn func(c *HttpCtx))
	}

	RouterNodeIf interface {
		Use(m ...Middleware) *Node
		Middlewares() []MiddlewareIf
		Node() *Node
	}

	Route struct {
		hashMap map[string]func(*HttpCtx)
	}
)


func (r *Route) Get(path string, fn func(ctx *HttpCtx)) *Route {
	r.hashMap[path] = fn
	return r
}

func (r *Route) Register(path string, fn func(c *HttpCtx)) {
	r.hashMap[path] = fn
}
