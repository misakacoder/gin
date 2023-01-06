package gin

import (
	"encoding/json"
	"net/http"
)

type Body map[string]any

type Context struct {
	Request  *http.Request
	Writer   responseWriter
	Params   map[string]string
	handlers []HandlerFunc
	index    int
	aborted  bool
}

func (context *Context) JSON(status int, body Body) {
	writer := &context.Writer
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	writer.WriteHeader(status)
	data, _ := json.Marshal(body)
	writer.Write(data)
}

func (context *Context) String(status int, text string) {
	writer := &context.Writer
	writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
	writer.WriteHeader(status)
	writer.WriteString(text)
}

func (context *Context) Next() {
	context.index++
	for !context.aborted && context.index < len(context.handlers) {
		context.handlers[context.index](context)
		context.index++
	}
}

func (context *Context) Abort() {
	context.aborted = true
}

func (context *Context) reset() {
	context.Params = nil
	context.handlers = nil
	context.index = -1
	context.aborted = false
}
