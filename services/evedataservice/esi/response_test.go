package esi

import (
	"bytes"
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/config"
	"github.com/andrewapj/arcturus/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

func Test_newResponse_BuildsResponse(t *testing.T) {

	testhelper.SetTestConfig()

	expected := buildTestResponse()

	actual, err := newResponse(buildTestHttpResponse())
	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func Test_newResponse_BuildsResponseWithMissingHeaders(t *testing.T) {

	testhelper.SetTestConfig()

	httpResponse := buildTestHttpResponse()
	httpResponse.Header = map[string][]string{}

	actual, err := newResponse(httpResponse)
	require.NoError(t, err)

	defaultFutureTime := actual.expires.Add(time.Duration(config.EsiDateAdditionalTime()) * time.Second)
	assert.Greater(t, defaultFutureTime, clock.GetTime())

	assert.Equal(t, 1, actual.pages)
}

func Test_newResponse_BuildsResponseWithEmptyHeaders(t *testing.T) {

	testhelper.SetTestConfig()

	httpResponse := buildTestHttpResponse()
	httpResponse.Header[config.EsiHeaderExpiresKey()] = []string{""}
	httpResponse.Header[config.EsiHeaderPagesKey()] = []string{""}

	actual, err := newResponse(httpResponse)
	require.NoError(t, err)

	defaultFutureTime := actual.expires.Add(time.Duration(config.EsiDateAdditionalTime()) * time.Second)
	assert.Greater(t, defaultFutureTime, clock.GetTime())

	assert.Equal(t, 1, actual.pages)
}

func Test_response_isError(t *testing.T) {

	testhelper.SetTestConfig()

	assert.False(t, buildTestResponse().isError())

	resp := buildTestResponse()
	resp.statusCode = http.StatusBadRequest
	assert.True(t, resp.isError())

	resp = buildTestResponse()
	resp.statusCode = http.StatusGatewayTimeout
	assert.True(t, resp.isError())
}

func buildTestResponse() *response {

	date := clock.ParseWithDefault(config.EsiDateLayout(), "Sun, 31 Mar 2024 11:05:00 GMT", clock.GetTime())
	return &response{
		body:       []byte("{}"),
		expires:    date,
		pages:      5,
		statusCode: 200,
	}
}

func buildTestHttpResponse() *http.Response {

	return &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/2.0",
		ProtoMajor: 2,
		ProtoMinor: 0,
		Header: map[string][]string{
			config.EsiHeaderExpiresKey(): {"Sun, 31 Mar 2024 11:05:00 GMT"},
			config.EsiHeaderPagesKey():   {"5"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte("{}"))),
	}
}
