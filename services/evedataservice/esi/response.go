package esi

import (
	"bytes"
	"fmt"
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/config"
	"io"
	"net/http"
	"strconv"
	"time"
)

// response represents a response from the ESI.
type response struct {
	body       []byte
	expires    time.Time
	pages      int
	statusCode int
}

// newResponse builds a new response.
func newResponse(resp *http.Response) (*response, error) {

	expires := resp.Header.Get(config.EsiHeaderExpiresKey())

	pages, err := strconv.Atoi(resp.Header.Get(config.EsiHeaderPagesKey()))
	if err != nil {
		pages = 1
	}

	buff := new(bytes.Buffer)
	_, err = io.Copy(buff, resp.Body)
	if err != nil {
		return &response{}, fmt.Errorf("unable to build responses. %w", err)
	}

	return &response{
		body: buff.Bytes(),
		expires: clock.ParseWithDefault(
			config.EsiDateLayout(),
			expires,
			clock.GetTime().Add(time.Duration(config.EsiDateAdditionalTime())*time.Second)),
		pages:      pages,
		statusCode: resp.StatusCode,
	}, nil
}

// isError checks if the http response is an error.
func (r *response) isError() bool {
	if r.statusCode >= http.StatusBadRequest && r.statusCode <= http.StatusNetworkAuthenticationRequired {
		return true
	}
	return false
}
