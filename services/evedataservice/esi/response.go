package esi

import (
	"fmt"
	"github.com/andrewapj/arcturus/clock"
	"github.com/andrewapj/arcturus/config"
	"io"
	"net/http"
	"strconv"
	"time"
)

type response struct {
	body          []byte
	contentLength int
	expires       time.Time
	pages         int
	statusCode    int
}

// newResponse builds a new response.
func newResponse(resp *http.Response) (*response, error) {

	contentLength, err := strconv.Atoi(resp.Header.Get(config.EsiHeaderContentLength()))
	if err != nil {
		contentLength = 0
	}

	expires := resp.Header.Get(config.EsiHeaderExpiresKey())

	pages, err := strconv.Atoi(resp.Header.Get(config.EsiHeaderPagesKey()))
	if err != nil {
		pages = 1
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &response{}, fmt.Errorf("unable to build response. %w", err)
	}

	return &response{
		body:          bytes,
		contentLength: contentLength,
		expires: clock.ParseWithDefault(
			config.EsiDateLayout(),
			expires,
			clock.GetTime().Add(time.Duration(config.EsiDateAdditionalTime())*time.Second)),
		pages:      pages,
		statusCode: resp.StatusCode,
	}, nil
}
