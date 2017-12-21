package jiaweb

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/iwannay/jiaweb/logger"
	"github.com/iwannay/jiaweb/proto"
	"github.com/iwannay/jiaweb/utils"
)

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
	// Middleware func(httpCtx *HttpContext)
	Router interface {
		ServeHTTP(ctx *HttpContext)
		ServerFile(path, fileRoot string) RouteNode
	}

	RouteNode interface {
		Use(m ...Middleware) *Node
		Middlewares() []Middleware
		Node() *Node
	}

	route struct {
		handleMap             map[string]HttpHandle
		NodeMap               map[string]*Node
		rwMutex               sync.RWMutex
		RedirectTrailingSlash bool
		server                *HttpServer
		RedirectFixedPath     bool
		HandleOPTIONS         bool
	}

	RouteHandle func(ctx *HttpContext)
)

var (
	SupportHTTPMethod map[string]bool
)

func NewRoute(server *HttpServer) *route {
	return &route{
		handleMap:             make(map[string]HttpHandle),
		NodeMap:               make(map[string]*Node),
		RedirectTrailingSlash: true,
		RedirectFixedPath:     true,
		HandleOPTIONS:         true,
		server:                server,
	}
}

func (r *route) RegisterHandler(name string, handler HttpHandle) {
	r.rwMutex.Lock()
	r.handleMap[name] = handler
	r.rwMutex.Unlock()
}

func (r *route) GetHandler(name string) (HttpHandle, bool) {
	r.rwMutex.RLock()
	h, ok := r.handleMap[name]
	r.rwMutex.RUnlock()
	return h, ok
}

func (r *route) ServeHTTP(ctx *HttpContext) {
	req := ctx.Request().Request
	rw := ctx.Response().ResponseWriter()
	path := req.URL.Path
	if root := r.NodeMap[req.Method]; root != nil {
		if handler, params := root.GetValue(path); handler != nil {
			handler(ctx)
			ctx.params = params
			// ctx.RouteNode =
		} else if req.Method != "CONNECT" && path != "/" {
			code := 301
			if req.Method != "GET" {
				code = 307
			}

			if r.RedirectTrailingSlash {
				if len(path) > 1 && path[len(path)-1] == '/' {
					req.URL.Path = path[:len(path)-1]
				} else {
					req.URL.Path = path + "/"
				}
				http.Redirect(rw, req, req.URL.String(), code)
				return
			}

			if r.RedirectFixedPath {
				// TODO 自动补全斜线
			}

		}

	}

	if req.Method == "OPTIONS" {
		if r.HandleOPTIONS {
			if allow := r.allowed(path, req.Method); len(allow) > 0 {
				rw.Header().Set("Allow", allow)
				return
			}
		}
	} else {
		// 405
		if allow := r.allowed(path, req.Method); len(allow) > 0 {

			ctx.Response().SetHeader("Allow", allow)
			ctx.Response().SetStatusCode(http.StatusMethodNotAllowed)

			// TODO 设置禁止访问handle

		}
	}

	// Handle 404
	ctx.Response().WriteHeader(http.StatusNotFound)

	// TODO 404 handle

}

func (r *route) allowed(path, reqMethod string) (allow string) {
	if path == "*" {
		for method := range r.NodeMap {
			if method == "OPTIONS" {
				continue
			}

			if len(allow) == 0 {
				allow = method
			} else {
				allow += ", " + method
			}
		}
	} else {
		for method := range r.NodeMap {
			if method == reqMethod || method == "OPTIONS" {
				continue
			}
			h, _ := r.NodeMap[method].GetValue(path)
			if h != nil {
				if len(allow) == 0 {
					allow = method
				} else {
					allow += ", " + method
				}
			}
		}
	}

	if len(allow) > 0 {
		allow += ", OPTIONS"
	}
	return
}

func (r *route) wrapRouteHandle(handler HttpHandle, isHijack bool) RouteHandle {
	return func(ctx *HttpContext) {
		ctx.handler = handler

		// TODO do feature

		if isHijack {
			// TODO Hijack
			_, err := ctx.Hijack()
			if err != nil {
				ctx.Response().WriteHeader(http.StatusInternalServerError)
				ctx.Response().Header().Set(HeaderContentType, CharsetUTF8)

			}
		}

		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
			}

		}()

		// do user handle

	}
}

func (r *route) ServerFile(path string, fileroot string) RouteNode {
	node := &Node{}
	var root http.FileSystem
	root = http.Dir(fileroot)
	if !r.server.ServerConfig().EnableListDir {
		root = &proto.HideDirFS{root}
	}
	fileServer := http.FileServer(root)
	node = r.add(HTTPMethod_GET, path, r.wrapFileHandle(fileServer))
	return node
}

func (r *route) add(method, path string, handle RouteHandle, m ...Middleware) *Node {
	if path[0] != '/' {
		panic("path must begin with '/' in path '" + path + "'")
	}

	if r.NodeMap == nil {
		r.NodeMap = make(map[string]*Node)
	}

	root := r.NodeMap[method]
	if root == nil {
		root = &Node{}
		r.NodeMap[method] = root
	}

	root.insertChild(path, handle)

	return root
}

func (r *route) wrapFileHandle(fHandler http.Handler) RouteHandle {
	return func(httpCtx *HttpContext) {
		startTime := time.Now()
		// TODO not supprot read dir by params

		fHandler.ServeHTTP(httpCtx.Response().rw, httpCtx.Request().Request)
		timeTaken := int64(time.Now().Sub(startTime) / time.Millisecond)

		logger.Logger().Debug(httpCtx.Request().Url()+" "+logRequest(httpCtx, timeTaken), LogTarget_HttpRequest)
	}
}

func logRequest(ctx Context, timeTaken int64) string {
	var reqByteLen, resByteLen, method, proto, status, userip string
	reqByteLen = utils.Int642String(ctx.Request().ContentLength)
	resByteLen = ""
	method = ctx.Request().Method
	proto = ctx.Request().Proto
	status = "200"
	userip = ctx.Request().RemoteIP()

	return fmt.Sprintf("%s %s %s %s %s %s %s",
		method,
		userip,
		proto,
		status,
		reqByteLen,
		resByteLen,
		utils.Int642String(timeTaken))
}

func init() {
	SupportHTTPMethod[HTTPMethod_Any] = true
	SupportHTTPMethod[HTTPMethod_GET] = true
	SupportHTTPMethod[HTTPMethod_POST] = true
	SupportHTTPMethod[HTTPMethod_PUT] = true
	SupportHTTPMethod[HTTPMethod_DELETE] = true
	SupportHTTPMethod[HTTPMethod_PATCH] = true
	SupportHTTPMethod[HTTPMethod_HiJack] = true
	SupportHTTPMethod[HTTPMethod_WebSocket] = true
	SupportHTTPMethod[HTTPMethod_HEAD] = true
	SupportHTTPMethod[HTTPMethod_OPTIONS] = true
}
