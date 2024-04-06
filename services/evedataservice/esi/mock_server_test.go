package esi

import (
	"github.com/andrewapj/arcturus/config"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"testing"
)

var (
	mutex  sync.Mutex
	server *httptest.Server
)

// mockRequestResponses represents a mocked request and the possible responses it can return.
// A request to the same path may return a different response based upon a query parameter.
type mockRequestResponses struct {
	request  request
	response func(r *http.Request) response
}

// start starts the server and configures it.
func startMockServer(t *testing.T) {

	mutex.Lock()

	mux := http.NewServeMux()
	server = httptest.NewServer(mux)

	err := os.Setenv(config.EsiDomainKey(), strings.Replace(server.URL, "http://", "", 1))
	require.Nilf(t, err, "unable to start server. %v", err)

	mockReqRes := buildMockRequestResponses()
	for _, mock := range mockReqRes {
		mux.HandleFunc(mock.request.pathWithParams(), func(w http.ResponseWriter, r *http.Request) {

			resp := mock.response(r)

			w.Header().Add(config.EsiHeaderExpiresKey(), resp.expires.Format(config.EsiDateLayout()))
			w.Header().Add(config.EsiHeaderPagesKey(), strconv.Itoa(resp.pages))
			w.WriteHeader(resp.statusCode)

			_, err := w.Write(resp.body)
			if err != nil {
				t.Fatalf("unable to write responses body for mock request, %s", err.Error())
			}
		})
	}
}

// stop stops the server.
func stopMockServer(t *testing.T) {

	err := os.Unsetenv(config.EsiDomainKey())
	require.Nilf(t, err, "unable to stop server. %v", err)
	server.CloseClientConnections()
	mutex.Unlock()
}

// buildMockRequestResponses generates the mocked requests and responses.
func buildMockRequestResponses() []mockRequestResponses {
	return []mockRequestResponses{
		{
			request: esiStatusRequest(),
			response: func(r *http.Request) response {
				if !strings.Contains(r.RequestURI, "page=") {
					return response{
						body:       []byte(esiStatusResponse),
						expires:    esiRequestTime,
						pages:      0,
						statusCode: 200,
					}
				}
				return response{
					body:       []byte(esiStatusResponse),
					expires:    esiRequestTime,
					pages:      0,
					statusCode: 404,
				}
			},
		},
	}
}
