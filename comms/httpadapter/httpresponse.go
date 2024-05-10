package httpadapter

import "net/http"

type responseWriterAdapter struct {
	http.ResponseWriter
	statusCode int
}

func (r *responseWriterAdapter) Write(data []byte) (int, error) {
	return r.ResponseWriter.Write(data)
}

func (r *responseWriterAdapter) WriteHeader(statusCode int) {
	r.statusCode = statusCode
	r.ResponseWriter.WriteHeader(statusCode)
}

func (r *responseWriterAdapter) Status() int {
	return r.statusCode
}
