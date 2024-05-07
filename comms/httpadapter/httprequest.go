package httpadapter

import "net/http"

type requestAdapter struct {
	*http.Request
}

func (r *requestAdapter) Header(key string) string {
	// Get the header value from the request
	return r.Request.Header.Get(key)
}

func (r *requestAdapter) Body() []byte {
	// Read the body
	body := make([]byte, r.Request.ContentLength)
	r.Request.Body.Read(body)

	return body
}
