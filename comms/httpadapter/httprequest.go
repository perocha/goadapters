package httpadapter

import "net/http"

type requestAdapter struct {
	*http.Request
}

func (r *requestAdapter) Header(key string) string {
	return r.Request.Header.Get(key)
}

func (r *requestAdapter) Body() []byte {
	// Implement your logic to read the request body if needed
	return nil
}
