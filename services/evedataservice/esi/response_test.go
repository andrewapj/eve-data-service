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

func Test_newResponse(t *testing.T) {

	testhelper.SetTestConfig()

	expected := buildTestResponse(t)

	actual, err := newResponse(buildTestHttpResponse())
	require.Nil(t, err, "unexpected error")

	assert.Equal(t, expected, actual)
}

func Test_newResponse_DefaultsSetWithMissingHeaders(t *testing.T) {

	testhelper.SetTestConfig()

	httpResponse := buildTestHttpResponse()
	httpResponse.Header = map[string][]string{}

	actual, err := newResponse(httpResponse)
	require.Nil(t, err)

	assert.Equal(t, 0, actual.contentLength)

	defaultFutureTime := actual.expires.Add(time.Duration(config.EsiDateAdditionalTime()) * time.Second)
	assert.Greater(t, defaultFutureTime, clock.GetTime())

	assert.Equal(t, 1, actual.pages)
}

func Test_newResponse_DefaultsSetWithEmptyHeaders(t *testing.T) {

	testhelper.SetTestConfig()

	httpResponse := buildTestHttpResponse()
	httpResponse.Header[config.EsiHeaderContentLength()] = []string{""}
	httpResponse.Header[config.EsiHeaderExpiresKey()] = []string{""}
	httpResponse.Header[config.EsiHeaderPagesKey()] = []string{""}

	actual, err := newResponse(httpResponse)
	require.Nil(t, err)

	assert.Equal(t, 0, actual.contentLength)

	defaultFutureTime := actual.expires.Add(time.Duration(config.EsiDateAdditionalTime()) * time.Second)
	assert.Greater(t, defaultFutureTime, clock.GetTime())

	assert.Equal(t, 1, actual.pages)
}

func buildTestResponse(t *testing.T) *response {

	date, err := time.ParseInLocation(config.EsiDateLayout(), "Sun, 31 Mar 2024 11:05:00 GMT", time.UTC)
	if err != nil {
		assert.Fail(t, "unable to generate time")
	}

	return &response{
		body:          []byte("{}"),
		contentLength: 2,
		expires:       date,
		pages:         5,
		statusCode:    200,
	}
}

func buildTestHttpResponse() *http.Response {

	return &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/2.0",
		ProtoMajor: 2,
		ProtoMinor: 0,
		Header: map[string][]string{
			config.EsiHeaderExpiresKey():    {"Sun, 31 Mar 2024 11:05:00 GMT"},
			config.EsiHeaderPagesKey():      {"5"},
			config.EsiHeaderContentLength(): {"2"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte("{}"))),
	}
}

func BenchmarkNewResponse(b *testing.B) {

	testhelper.SetTestConfig()

	for i := 0; i < b.N; i++ {
		httpResponse := buildTestHttpResponse()
		_, err := newResponse(httpResponse)
		if err != nil {
			b.Fatal(err.Error())
		}
	}

	b.ReportAllocs()
}
