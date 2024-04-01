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
	body       []byte
	expires    time.Time
	pages      int
	statusCode int
}

// newResponse builds a new response.
func newResponse(resp *http.Response) (*response, error) {

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &response{}, fmt.Errorf("unable to build response. %w", err)
	}

	expires := ""
	if len(resp.Header[config.EsiHeaderExpiresKey()]) == 1 {
		expires = resp.Header[config.EsiHeaderExpiresKey()][0]
	}

	pages := "1"
	if len(resp.Header[config.EsiHeaderPagesKey()]) == 1 && resp.Header[config.EsiHeaderPagesKey()][0] != "" {
		pages = resp.Header[config.EsiHeaderPagesKey()][0]
	}
	pagesInt, err := strconv.Atoi(pages)
	if err != nil {
		return &response{}, fmt.Errorf("unable to build response with invalid pages. %w", err)
	}

	return &response{
		body: bytes,
		expires: clock.ParseWithDefault(
			config.EsiDateLayout(),
			expires,
			clock.GetTime().Add(time.Duration(config.EsiDateAdditionalTime())*time.Second)),
		pages:      pagesInt,
		statusCode: resp.StatusCode,
	}, nil
}
