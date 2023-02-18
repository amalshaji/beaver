package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// HTTPResponse is a serializable version of http.Response ( with only useful fields )
type HTTPResponse struct {
	StatusCode    int
	Header        http.Header
	ContentLength int64
}

// SerializeHTTPResponse create a new HTTPResponse from a http.Response
func SerializeHTTPResponse(resp *http.Response) *HTTPResponse {
	r := new(HTTPResponse)
	r.StatusCode = resp.StatusCode
	r.Header = resp.Header
	r.ContentLength = resp.ContentLength
	return r
}

// NewHTTPResponse creates a new HTTPResponse
func NewHTTPResponse() (r *HTTPResponse) {
	r = new(HTTPResponse)
	r.Header = make(http.Header)
	return
}

// ProxyError log error and return a HTTP 526 error with the message
func ProxyError(c echo.Context, err error) error {
	log.Println(err)
	return c.JSON(526, map[string]string{"error": err.Error()})
}

// ProxyErrorf log error and return a HTTP 526 error with the message
func ProxyErrorf(c echo.Context, format string, args ...interface{}) error {
	return ProxyError(c, fmt.Errorf(format, args...))
}

func HttpBadRequest(c echo.Context, format string, args ...interface{}) error {
	return c.JSON(
		http.StatusBadRequest,
		map[string]string{"error": fmt.Errorf(format, args...).Error()},
	)
}

func HttpUnauthorized(c echo.Context, format string, args ...interface{}) error {
	return c.JSON(
		http.StatusUnauthorized,
		map[string]string{"error": fmt.Errorf(format, args...).Error()},
	)
}
