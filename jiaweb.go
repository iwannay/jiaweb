package jiaweb

import (
	"sync"

	"github.com/iwannay/jiaweb/config"
	"github.com/iwannay/jiaweb/logger"
)

type (
	JiaWeb struct {
		HttpServer              *HttpServer
		Config                  *config.Config
		Middlewares             []Middleware
		ExceptionHandler        ExceptionHandle
		NotFoundHandler         StandardHandle
		MethodNotAllowedHandler StandardHandle
		Mutex                   sync.RWMutex
	}

	// 自定义异常处理
	ExceptionHandle func(Context, error)

	StandardHandle func(Context)
	HttpHandle     func(httpCtx Context) error
)

const (
	DefaultHTTPPort    = 8080
	RunModeDevelopment = "development"
	RunModeProduction  = "production"
)

func New() *JiaWeb {
	app := &JiaWeb{
		HttpServer:  NewHttpServer(),
		Config:      config.New(),
		Middlewares: make([]Middleware, 0),
	}
	// app.HttpServer
	logger.InitJiaLog()

	return app
}

func Classic() *JiaWeb {
	app := New()

	app.SetEnableLog(true)

	return app

}

func (app *JiaWeb) SetEnableLog(enableLog bool) {
	logger.SetEnableLog(enableLog)

}

func (app *JiaWeb) SetLogPath(path string) {
	logger.SetLogPath(path)
}

// func (app *JiaWeb) RegisterMiddlewareFunc(name string,)
