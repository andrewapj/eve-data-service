package esi

import (
	"context"
	"fmt"
	"github.com/andrewapj/arcturus/config"
	"net/http"
	"strings"
)

// request represents a request made to the ESI.
type request struct {
	headers    map[string][]string
	method     string
	protocol   string
	domain     string
	path       string
	pathParams map[string]string
	datasource string
	language   string
	page       int
}

// newRequest builds a new request struct.
func newRequest(path string, pathParams map[string]string, page int) request {

	if page < 0 {
		page = 0
	}

	return request{
		headers: map[string][]string{
			"Accept":          {"application/json"},
			"Accept-Language": {"en"},
			"Cache-Control":   {"no-cache"},
			"User-Agent":      {config.EsiUserAgent()},
		},
		method:     "GET",
		protocol:   config.EsiProtocol(),
		domain:     config.EsiDomain(),
		path:       path,
		pathParams: pathParams,
		datasource: config.EsiDatasource(),
		language:   config.EsiLanguage(),
		page:       page,
	}
}

// url generates a full URL from a request including query string parameters.
func (r request) url() string {

	url := fmt.Sprintf("%s://%s%s?datasource=%s&language=%s",
		r.protocol, r.domain, r.pathWithParams(), r.datasource, r.language)
	if r.page > 0 {
		url += fmt.Sprintf("&page=%d", r.page)
	}

	return url
}

// pathWithParams generates a path with the path parameters included.
func (r request) pathWithParams() string {

	path := r.path

	for param, val := range r.pathParams {
		param = fmt.Sprintf("{%s}", param)
		path = strings.ReplaceAll(path, param, val)
	}

	return path
}

// toHttpRequestWithCtx maps a request to a http.Request. It includes a context in the http.Request.
func (r request) toHttpRequestWithCtx(ctx context.Context) (*http.Request, error) {

	req, err := http.NewRequestWithContext(ctx, r.method, r.url(), nil)
	if err != nil {
		return &http.Request{}, err
	}

	for key, values := range r.headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	return req, nil
}
