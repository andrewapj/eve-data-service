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

// Client represents an ESI client.
type Client struct {
	client     http.Client
	guard      chan struct{}
	maxRetries int
}

// NewClient builds a new Client.
func NewClient() Client {

	return Client{
		client: http.Client{
			Timeout: time.Duration(config.EsiTimeout()) * time.Second,
		},
		guard:      make(chan struct{}, config.EsiConcurrency()),
		maxRetries: config.EsiMaxRetries(),
	}
}

func (c *Client) getPages(ctx context.Context, r request, pr PageRequest) ([]*response, error) {

	var responses = make([]*response, 0, len(pr.pages()))
	for _, page := range pr.pages() {
		resp, err := c.retryRequest(ctx, r.withPage(page))
		if err != nil {
			return nil, err
		}
		responses = append(responses, resp)
		if resp.pages == page {
			break
		}
	}
	return responses, nil
}

// retryRequest will retry a specific request up to maxRetries.
func (c *Client) retryRequest(ctx context.Context, r request) (*response, error) {

	attempts := 1
	for {
		resp, err := c.makeRequest(ctx, r)
		if err == nil {
			return resp, nil
		}
		if err != nil && attempts == c.maxRetries {
			return nil, err
		}

		attempts++
		time.Sleep(time.Duration(attempts) * time.Second)
	}
}

// makeRequest will make a new http request. It checks for a valid user agent, context completion and a valid
// status code. The request can be cancelled by cancelling the context.
func (c *Client) makeRequest(ctx context.Context, r request) (*response, error) {

	c.guard <- struct{}{}
	defer func() { <-c.guard }()

	req, err := r.toHttpRequestWithCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to build request %s. %w", r.url(), err)
	}

	if req.Header.Get(config.EsiUserAgent()) == "" {
		return nil, fmt.Errorf("unable to make request to %s. Missing user agent header", r.url())
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to make request to %s. %w", r.url(), err)
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
		return nil, fmt.Errorf("unable to build response from %s. %w", r.url(), err)
	}
	if res.isError() {
		return nil, fmt.Errorf("received an http error code from esi, url: %s, code: %d", r.url(), res.statusCode)
	}
	return res, nil
}
