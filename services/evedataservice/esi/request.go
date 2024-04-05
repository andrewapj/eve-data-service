package esi

import (
	"context"
	"fmt"
	"github.com/andrewapj/arcturus/config"
	"net/http"
	"strings"
)

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

// makeUrl generates a full URL from a request.
func (r request) makeUrl() string {

	url := fmt.Sprintf("%s://%s%s?datasource=%s&language=%s",
		r.protocol, r.domain, r.path, r.datasource, r.language)
	if r.page > 0 {
		url += fmt.Sprintf("&page=%d", r.page)
	}

	for param, val := range r.pathParams {
		param = fmt.Sprintf("{%s}", param)
		url = strings.ReplaceAll(url, param, val)
	}

	return url
}

func (r request) toHttpRequestWithCtx(ctx context.Context) (*http.Request, error) {

	req, err := http.NewRequestWithContext(ctx, r.method, r.makeUrl(), nil)
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
