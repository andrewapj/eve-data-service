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

func Test_newResponse_FutureExpiresWithMissingExpiresHeader(t *testing.T) {

	testhelper.SetTestConfig()

	httpResponse := buildTestHttpResponse()
	httpResponse.Header = map[string][]string{}

	actual, err := newResponse(httpResponse)
	require.Nil(t, err)

	defaultFutureTime := actual.expires.Add(time.Duration(config.EsiDateAdditionalTime()) * time.Second)
	assert.Greater(t, defaultFutureTime, clock.GetTime())

}

func Test_newResponse_FutureExpiresWithEmptyExpiresHeader(t *testing.T) {

	testhelper.SetTestConfig()

	httpResponse := buildTestHttpResponse()
	httpResponse.Header[config.EsiHeaderExpiresKey()] = []string{""}

	actual, err := newResponse(httpResponse)
	require.Nil(t, err)

	defaultFutureTime := actual.expires.Add(time.Duration(config.EsiDateAdditionalTime()) * time.Second).Truncate(time.Second)
	assert.Greater(t, defaultFutureTime, clock.GetTime())
}

func Test_newResponse_DefaultsToPage1WithMissingPageHeader(t *testing.T) {

	testhelper.SetTestConfig()

	httpResponse := buildTestHttpResponse()
	httpResponse.Header = map[string][]string{}

	actual, err := newResponse(httpResponse)
	require.Nil(t, err)

	assert.Equal(t, 1, actual.pages)
}

func Test_newResponse_DefaultsToPage1WithEmptyPageHeader(t *testing.T) {

	testhelper.SetTestConfig()

	httpResponse := buildTestHttpResponse()
	httpResponse.Header[config.EsiHeaderPagesKey()] = []string{""}

	actual, err := newResponse(httpResponse)
	require.Nil(t, err)

	assert.Equal(t, 1, actual.pages)
}

func buildTestResponse(t *testing.T) *response {

	date, err := time.ParseInLocation(config.EsiDateLayout(), "Sun, 31 Mar 2024 11:05:00 GMT", time.UTC)
	if err != nil {
		assert.Fail(t, "unable to generate time")
	}

	return &response{
		body:       []byte("{}"),
		expires:    date,
		pages:      1,
		statusCode: 200,
	}
}

func buildTestHttpResponse() *http.Response {

	return &http.Response{
		StatusCode: 200,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: map[string][]string{
			config.EsiHeaderExpiresKey(): {"Sun, 31 Mar 2024 11:05:00 GMT"},
			config.EsiHeaderPagesKey():   {"1"},
		},
		Body: io.NopCloser(bytes.NewReader([]byte("{}"))),
	}
}
