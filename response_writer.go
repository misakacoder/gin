package gin

import (
	"net/http"
)

type responseWriter struct {
	http.ResponseWriter
	status int
}

func (writer *responseWriter) Status() int {
	return writer.status
}

func (writer *responseWriter) WriteHeader(code int) {
	writer.status = code
	writer.ResponseWriter.WriteHeader(code)
}

func (writer *responseWriter) Write(data []byte) (n int, err error) {
	n, err = writer.ResponseWriter.Write(data)
	return
}

func (writer *responseWriter) WriteString(data string) (n int, err error) {
	n, err = writer.Write([]byte(data))
	return
}
