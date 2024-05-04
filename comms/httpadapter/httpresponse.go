package httpadapter

import "net/http"

type responseWriterAdapter struct {
	http.ResponseWriter
}

func (r *responseWriterAdapter) Write(data []byte) (int, error) {
	return r.ResponseWriter.Write(data)
}

func (r *responseWriterAdapter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
}
