package esi

import (
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/config"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mutex          sync.Mutex
	server         *httptest.Server
	esiExpiresTime = clock.GetTime().Add(5 * time.Minute).Truncate(time.Second)
)

// mockRequestResponses represents a mocked request and the possible responses it can return.
// A request to the same path may return a different response based upon a query parameter.
type mockRequestResponses struct {
	request  request
	response func(r *http.Request) response
}

// start starts the server and configures it.
func startMockServer() {

	mutex.Lock()

	mux := http.NewServeMux()
	server = httptest.NewServer(mux)

	err := os.Setenv(config.EsiDomainKey(), strings.Replace(server.URL, "http://", "", 1))
	if err != nil {
		panic(err)
	}

	mockReqRes := buildMockRequestResponses()
	for _, mock := range mockReqRes {
		mux.HandleFunc(mock.request.pathWithParams(), func(w http.ResponseWriter, r *http.Request) {

			resp := mock.response(r)

			w.Header().Add(config.EsiHeaderExpiresKey(), resp.expires.Format(config.EsiDateLayout()))
			w.Header().Add(config.EsiHeaderPagesKey(), strconv.Itoa(resp.pages))
			w.WriteHeader(resp.statusCode)

			_, err := w.Write(resp.body)
			if err != nil {
				panic(err.Error())
			}
		})
	}
}

// stop stops the server.
func stopMockServer() {

	err := os.Unsetenv(config.EsiDomainKey())
	if err != nil {
		panic(err.Error())
	}
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
					return response{body: []byte(esiStatusResponse), expires: esiExpiresTime, pages: 0, statusCode: 200}
				}
				return response{body: []byte(esiStatusResponse), expires: esiExpiresTime, pages: 0, statusCode: 404}
			},
		},
		{
			request: esiTypeIdsRequest(),
			response: func(r *http.Request) response {
				if !strings.Contains(r.RequestURI, "page") || strings.Contains(r.RequestURI, "page=1") {
					return response{body: []byte(esiTypeIdsPage1Response), expires: esiExpiresTime, pages: 2, statusCode: 200}
				}
				if strings.Contains(r.RequestURI, "page=2") {
					return response{
						body: []byte(esiTypeIdsPage2Response), expires: esiExpiresTime, pages: 2, statusCode: 200}
				}
				return response{
					body: []byte(esiTypeIdsPage2Response), expires: esiExpiresTime, pages: 2, statusCode: 404}
			},
		},
	}
}
