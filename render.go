package jiaweb

import (
	"fmt"
	"html/template"
	"path/filepath"
	"time"
)

type (
	Viewer interface {
		AppendTpl(tpl ...string)
		AppendFunc(funcMap template.FuncMap)
		RenderHtml(rw *Response, viewPath []string, locals map[string]interface{}) error
		Tpls() []string
	}

	view struct {
		funcMap  template.FuncMap
		innerTpl []string
		server   *HttpServer
		locals   map[string]interface{}
		// mutex    sync.RWMutex
	}
)

func NewView(s *HttpServer) *view {
	return &view{
		server:  s,
		funcMap: make(template.FuncMap),
	}
}

func (v *view) Tpls() []string {
	return v.innerTpl
}

func (v *view) AppendTpl(tpl ...string) {
	v.innerTpl = append(v.innerTpl, tpl...)
}

func (v *view) AppendFunc(funcMap template.FuncMap) {
	v.funcMap = funcMap
}

func (v *view) RenderHtml(rw *Response, viewPath []string, locals map[string]interface{}) error {
	var tplPaths []string
	var tplName string
	var tplPath string
	var startTime = time.Now()

	for _, item := range viewPath {
		tplPath = filepath.Join(".", v.server.TemplateConfig().TplDir, item+v.server.TemplateConfig().TplExt)
		tplPaths = append(tplPaths, tplPath)
		if tplName == "" {
			tplName = filepath.Base(tplPath)
		}
	}

	tplPaths = append(v.innerTpl, tplPaths...)

	if locals == nil {
		locals = make(map[string]interface{})

	}
	for index, item := range v.locals {
		locals[index] = item
	}

	t := template.Must(template.New(tplName).Funcs(v.funcMap).ParseFiles(tplPaths...))

	subTime := time.Now().Sub(startTime).Nanoseconds()

	locals["costTime"] = fmt.Sprintf("%.5fms", subTime)
	err := t.ExecuteTemplate(rw.ResponseWriter(), tplName, locals)
	return err
}
