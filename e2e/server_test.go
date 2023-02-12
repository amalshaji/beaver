package main

import (
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const PROXY_URL = "http://test.localhost:8080"

func getUrlPath(path string) string {
	return PROXY_URL + path
}

func stringifyResBody(body io.ReadCloser) string {
	var p []byte

	_, err := body.Read(p)
	defer body.Close()

	if err != nil {
		return ""
	}

	return string(p)
}

func TestGetRequest(t *testing.T) {
	res, err := http.Get(getUrlPath("/"))

	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 200)
	assert.Equal(t, stringifyResBody(res.Body), `{"message": "ok"}`)
}

func TestGetRequestOnUnregisteredSubdomainShouldFail(t *testing.T) {
	res, err := http.Get("http://xxyyzz.localhost:8080")

	assert.NoError(t, err)
	assert.Equal(t, res.StatusCode, 526)
}
