package gin

import (
	"github.com/misakacoder/logger"
	"net/http"
	"path"
	"strings"
)

var (
	allHttpMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace,
	}
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	path   string
	engine *Engine
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.engine.tree.addMiddleware(group.path, middlewares)
}

func (group *RouterGroup) GET(path string, handler HandlerFunc) {
	group.handle(http.MethodGet, path, handler)
}

func (group *RouterGroup) POST(path string, handler HandlerFunc) {
	group.handle(http.MethodPost, path, handler)
}

func (group *RouterGroup) PUT(path string, handler HandlerFunc) {
	group.handle(http.MethodPut, path, handler)
}

func (group *RouterGroup) DELETE(path string, handler HandlerFunc) {
	group.handle(http.MethodDelete, path, handler)
}

func (group *RouterGroup) HEAD(path string, handler HandlerFunc) {
	group.handle(http.MethodHead, path, handler)
}

func (group *RouterGroup) PATCH(path string, handler HandlerFunc) {
	group.handle(http.MethodPatch, path, handler)
}

func (group *RouterGroup) CONNECT(path string, handler HandlerFunc) {
	group.handle(http.MethodConnect, path, handler)
}

func (group *RouterGroup) OPTIONS(path string, handler HandlerFunc) {
	group.handle(http.MethodOptions, path, handler)
}

func (group *RouterGroup) TRACE(path string, handler HandlerFunc) {
	group.handle(http.MethodTrace, path, handler)
}

func (group *RouterGroup) Any(path string, handler HandlerFunc) {
	for _, method := range allHttpMethods {
		group.handle(method, path, handler)
	}
}

func (group *RouterGroup) Group(gpath string) *RouterGroup {
	if !strings.HasPrefix(gpath, "/") {
		logger.Panic("group path '%s' must start with /", gpath)
	}
	return &RouterGroup{
		path:   path.Join(group.path, gpath),
		engine: group.engine,
	}
}

func (group *RouterGroup) handle(method string, p string, handler HandlerFunc) {
	if !strings.HasPrefix(p, "/") {
		logger.Panic("path '%s' must start with /", p)
	}
	group.engine.addRoute(method, path.Join(group.path, p), handler)
}
