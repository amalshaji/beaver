package main

import (
	"bytes"
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

	p, err := io.ReadAll(body)
	defer body.Close()

	if err != nil {
		return ""
	}

	return string(p)
}

func httpRequest(method, url string, data io.Reader) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func putRequest(url string, data io.Reader) (*http.Response, error) {
	return httpRequest(http.MethodPut, url, data)
}

func patchRequest(url string, data io.Reader) (*http.Response, error) {
	return httpRequest(http.MethodPatch, url, data)
}

func deleteRequest(url string, data io.Reader) (*http.Response, error) {
	return httpRequest(http.MethodDelete, url, data)
}

func optionsRequest(url string, data io.Reader) (*http.Response, error) {
	return httpRequest(http.MethodOptions, url, data)
}

func headRequest(url string, data io.Reader) (*http.Response, error) {
	return httpRequest(http.MethodHead, url, data)
}

func connectRequest(url string, data io.Reader) (*http.Response, error) {
	return httpRequest(http.MethodConnect, url, data)
}

func traceRequest(url string, data io.Reader) (*http.Response, error) {
	return httpRequest(http.MethodTrace, url, data)
}

func TestGetRequest(t *testing.T) {
	res, err := http.Get(getUrlPath("/"))

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}

func TestGetRequestOnUnregisteredSubdomainShouldFail(t *testing.T) {
	res, err := http.Get("http://xxyyzz.localhost:8080")

	assert.NoError(t, err)
	assert.Equal(t, 526, res.StatusCode)
	assert.Equal(t, "{\"error\":\"unregistered tunnel subdomain\"}\n", stringifyResBody(res.Body))
}

func TestPostRequestAsJson(t *testing.T) {
	var jsonData = []byte(`{
		"username": "morpheus",
		"password": "dingdong"
	}`)

	res, err := http.Post(getUrlPath("/"), "application/json", bytes.NewBuffer(jsonData))

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}

func TestPostRequestAsForm(t *testing.T) {
	var jsonData = []byte(`{
		"username": "morpheus",
		"password": "dingdong"
	}`)

	res, err := http.Post(getUrlPath("/"), "application/x-www-form-urlencoded", bytes.NewBuffer(jsonData))

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}

func TestPostRequestOnUnregisteredSubdomainShouldFail(t *testing.T) {
	res, err := http.Post("http://xxyyzz.localhost:8080", "application/json", nil)

	assert.NoError(t, err)
	assert.Equal(t, 526, res.StatusCode)
	assert.Equal(t, "{\"error\":\"unregistered tunnel subdomain\"}\n", stringifyResBody(res.Body))
}

func TestPutRequest(t *testing.T) {
	res, err := putRequest(getUrlPath("/"), nil)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}

func TestPutRequestOnUnregisteredSubdomainShouldFail(t *testing.T) {
	res, err := putRequest("http://xxyyzz.localhost:8080", nil)

	assert.NoError(t, err)
	assert.Equal(t, 526, res.StatusCode)
	assert.Equal(t, "{\"error\":\"unregistered tunnel subdomain\"}\n", stringifyResBody(res.Body))
}

func TestPatchRequest(t *testing.T) {
	res, err := patchRequest(getUrlPath("/"), nil)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}

func TestPatchRequestOnUnregisteredSubdomainShouldFail(t *testing.T) {
	res, err := patchRequest("http://xxyyzz.localhost:8080", nil)

	assert.NoError(t, err)
	assert.Equal(t, 526, res.StatusCode)
	assert.Equal(t, "{\"error\":\"unregistered tunnel subdomain\"}\n", stringifyResBody(res.Body))
}

func TestDeleteRequest(t *testing.T) {
	res, err := deleteRequest(getUrlPath("/"), nil)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}

func TestDeleteRequestOnUnregisteredSubdomainShouldFail(t *testing.T) {
	res, err := deleteRequest("http://xxyyzz.localhost:8080", nil)

	assert.NoError(t, err)
	assert.Equal(t, 526, res.StatusCode)
	assert.Equal(t, "{\"error\":\"unregistered tunnel subdomain\"}\n", stringifyResBody(res.Body))
}

func TestOptionsRequest(t *testing.T) {
	res, err := optionsRequest(getUrlPath("/"), nil)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}

func TestOptionsRequestOnUnregisteredSubdomainShouldFail(t *testing.T) {
	res, err := optionsRequest("http://xxyyzz.localhost:8080", nil)

	assert.NoError(t, err)
	assert.Equal(t, 526, res.StatusCode)
	assert.Equal(t, "{\"error\":\"unregistered tunnel subdomain\"}\n", stringifyResBody(res.Body))
}

func TestHeadRequest(t *testing.T) {
	res, err := headRequest(getUrlPath("/"), nil)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "beaver-server", res.Header.Get("custom-server"))
}

func TestConnectRequest(t *testing.T) {
	res, err := connectRequest(getUrlPath("/"), nil)

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}

func TestRedirect302Request(t *testing.T) {
	res, err := http.Get(getUrlPath("/redirect-302"))

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}

func TestRedirect307Request(t *testing.T) {
	res, err := http.Get(getUrlPath("/redirect-307"))

	assert.NoError(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "{\"message\":\"ok\"}\n", stringifyResBody(res.Body))
}
