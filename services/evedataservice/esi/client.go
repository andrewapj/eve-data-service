package esi

import (
	"context"
	"fmt"
	"github.com/andrewapj/arcturus/config"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const esiStatusPath = "/latest/status/"

// Client represents an ESI client.
type Client struct {
	client    http.Client
	guardChan chan struct{}
}

// NewClient builds a new Client.
func NewClient() Client {

	return Client{
		client: http.Client{
			Timeout: time.Duration(config.EsiTimeout()) * time.Second,
		},
		guardChan: make(chan struct{}, config.EsiConcurrency()),
	}
}

// makeRequest will make a new http request. It checks for a valid user agent, context completion and a valid
// status code. The request can be cancelled by cancelling the context.
func (e *Client) makeRequest(ctx context.Context, r request) (*response, error) {

	e.guardChan <- struct{}{}
	defer func() { <-e.guardChan }()

	req, err := r.toHttpRequestWithCtx(ctx)
	if err != nil {
		return &response{}, fmt.Errorf("unable to build request %s. %w", r.url(), err)
	}

	if req.Header.Get(config.EsiUserAgent()) == "" {
		return &response{}, fmt.Errorf("unable to make request to %s. Missing user agent header", r.url())
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return &response{}, fmt.Errorf("unable to make request to %s. %w", r.url(), err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			slog.Error("unable to close request body. %w", err)
		}
	}(resp.Body)
	slog.Debug("received a response from the esi", slog.String("url", r.url()), slog.Int("code", resp.StatusCode))

	res, err := newResponse(resp)
	if err != nil {
		return &response{}, fmt.Errorf("unable to build response from %s. %w", r.url(), err)
	}
	if res.isError() {
		return &response{}, fmt.Errorf("received an http error code from esi, url: %s, code: %d", r.url(), res.statusCode)
	}
	return res, nil
}
