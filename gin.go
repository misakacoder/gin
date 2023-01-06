package gin

import (
	"fmt"
	"github.com/misakacoder/logger"
	"net"
	"net/http"
	"net/url"
	"path"
	"strings"
	"sync"
	"time"
)

var default404Body = "404 page not found"

type Engine struct {
	RouterGroup
	pool     sync.Pool
	tree     *node
	notFound HandlerFunc
}

func New() *Engine {
	engine := &Engine{
		RouterGroup: RouterGroup{
			path: "/",
		},
		tree:     &node{},
		notFound: defaultNotFound,
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() any {
		return &Context{}
	}
	return engine
}

func Default() *Engine {
	engine := New()
	engine.Use(Network, Recovery)
	return engine
}

func (engine *Engine) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	context := engine.pool.Get().(*Context)
	context.reset()
	context.Request = request
	context.Writer = responseWriter{ResponseWriter: writer}
	handlers := make([]HandlerFunc, 0)
	middlewares, handler, params := engine.tree.getRoute(request.Method, path.Clean(request.URL.Path))
	handlers = append(handlers, *middlewares...)
	if handler != nil {
		handlers = append(handlers, handler)
	} else {
		handlers = append(handlers, engine.notFound)
	}
	context.Params = params
	context.handlers = handlers
	context.Next()
	engine.pool.Put(context)
}

func (engine *Engine) NotFound(handler HandlerFunc) {
	if handler != nil {
		engine.notFound = handler
	}
}

func (engine *Engine) Run(port int) {
	logger.Info("Running on http://localhost:%d", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), engine)
	if err != nil {
		logger.Panic(err.Error())
	}
}

func (engine *Engine) addRoute(method string, path string, handler HandlerFunc) {
	if handler == nil {
		logger.Panic("handler for path '%s' cannot be null", path)
	}
	err := engine.tree.addRoute(method, path, handler)
	if err != nil {
		logger.Panic(err.Error())
	} else {
		logger.Debug("add route: %s %s", method, path)
	}
}

func Recovery(context *Context) {
	defer func() {
		if err := recover(); err != nil {
			context.String(http.StatusInternalServerError, fmt.Sprintf("%v", err))
			context.Abort()
		}
	}()
	context.Next()
}

func Network(context *Context) {
	request := context.Request
	ip, _, _ := net.SplitHostPort(strings.TrimSpace(request.RemoteAddr))
	if ip == "::1" {
		ip = "localhost"
	}
	start := time.Now()
	context.Next()
	milliseconds := time.Now().Sub(start).Milliseconds()
	uri, _ := url.PathUnescape(request.RequestURI)
	logger.Info("%s %s %s %d %dms", ip, request.Method, uri, context.Writer.Status(), milliseconds)
}

func defaultNotFound(context *Context) {
	context.String(http.StatusNotFound, default404Body)
}
