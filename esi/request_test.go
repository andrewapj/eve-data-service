package esi

import (
	"context"
	"github.com/andrewapj/arcturus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/url"
	"testing"
)

func Test_newRequest(t *testing.T) {

	config.SetTestConfig()

	actual := newRequest("/path/{param1}/path2/{param2}/",
		map[string]string{"param1": "value1", "param2": "value2"})

	assert.Equal(t, buildTestRequest(), actual)
}

func Test_url_WithZeroPage(t *testing.T) {

	config.SetTestConfig()

	r := newRequest("/path/{param1}/path2/{param2}/",
		map[string]string{"param1": "value1", "param2": "value2"})

	actual := r.url()
	expected := config.EsiProtocol() + "://" + config.EsiDomain() + "/path/value1/path2/value2/?datasource=tranquility&language=en"

	assert.Equal(t, expected, actual)
}

func Test_url_WithPage(t *testing.T) {

	config.SetTestConfig()

	r := newRequest("/path/{param1}/path2/{param2}/",
		map[string]string{"param1": "value1", "param2": "value2"}).withPage(1)

	actual := r.url()
	expected := config.EsiProtocol() + "://" + config.EsiDomain() + "/path/value1/path2/value2/?datasource=tranquility&language=en&page=1"

	assert.Equal(t, expected, actual)
}

func Test_path_WithParams(t *testing.T) {

	config.SetTestConfig()

	r := newRequest("/path/{param1}/path2/{param2}/",
		map[string]string{"param1": "value1", "param2": "value2"})

	actual := r.pathWithParams()
	expected := "/path/value1/path2/value2/"

	assert.Equal(t, expected, actual)
}

func Test_toHttpRequestWithCtx(t *testing.T) {

	config.SetTestConfig()

	ctx := context.WithValue(context.Background(), "key", "value")
	expectedUrl, err := url.Parse(
		config.EsiProtocol() +
			"://" +
			config.EsiDomain() +
			"/path/value1/path2/value2/" +
			"?datasource=" + config.EsiDatasource() +
			"&language=" + config.EsiLanguage() +
			"&page=5")
	require.NoError(t, err, "unexpected err: %v", err)
	expected := &http.Request{
		Method:     "GET",
		URL:        expectedUrl,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: map[string][]string{
			"Accept":          {"application/json"},
			"Accept-Language": {"en"},
			"Cache-Control":   {"no-cache"},
			"User-Agent":      {config.EsiUserAgent()},
		},
		Host: config.EsiDomain(),
	}
	expected = expected.WithContext(ctx)

	r := newRequest("/path/{param1}/path2/{param2}/",
		map[string]string{"param1": "value1", "param2": "value2"}).withPage(5)

	actual, err := r.toHttpRequestWithCtx(ctx)
	require.NoError(t, err, "unexpected err: %v", err)
	assert.Equal(t, expected, actual)
}

func buildTestRequest() request {

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
		path:       "/path/{param1}/path2/{param2}/",
		pathParams: map[string]string{"param1": "value1", "param2": "value2"},
		datasource: config.EsiDatasource(),
		language:   config.EsiLanguage(),
		page:       0,
	}
}
