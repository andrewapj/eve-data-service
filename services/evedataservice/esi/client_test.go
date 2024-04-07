package esi

import (
	"context"
	"fmt"
	"github.com/andrewapj/arcturus/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_getPages_WithMissingPageInRequest(t *testing.T) {

	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	client := NewClient()

	pageReq := NewPageRequest(0, 5)
	resp, err := client.getPages(context.Background(), esiTypeIdsRequest(), pageReq)
	require.NoError(t, err)

	require.Equal(t, 1, len(resp))
	assert.Equal(t, &response{
		body:       []byte(esiTypeIdsPage1Response),
		expires:    esiExpiresTime,
		pages:      2,
		statusCode: 200}, resp[0])
}

func TestClient_getPages_WithPage1Only(t *testing.T) {

	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	client := NewClient()

	pageReq := NewPageRequest(1, 1)
	resp, err := client.getPages(context.Background(), esiTypeIdsRequest(), pageReq)
	require.NoError(t, err)

	require.Equal(t, 1, len(resp))
	assert.Equal(t, &response{
		body:       []byte(esiTypeIdsPage1Response),
		expires:    esiExpiresTime,
		pages:      2,
		statusCode: 200}, resp[0])
}

func TestClient_getPages_WithMultiplePages_DoesNotExceedMaxPages(t *testing.T) {

	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	client := NewClient()

	// Generate a page request for 10 pages. The client should not request pages that do not exist.
	pageReq := NewPageRequest(1, 10)
	resp, err := client.getPages(context.Background(), esiTypeIdsRequest(), pageReq)
	require.NoError(t, err)

	require.Equal(t, 2, len(resp))
	assert.Equal(t, &response{
		body:       []byte(esiTypeIdsPage1Response),
		expires:    esiExpiresTime,
		pages:      2,
		statusCode: 200}, resp[0])
	assert.Equal(t, &response{
		body:       []byte(esiTypeIdsPage2Response),
		expires:    esiExpiresTime,
		pages:      2,
		statusCode: 200}, resp[1])
}

func TestClient_retryRequest(t *testing.T) {
	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	client := NewClient()

	resp, err := client.retryRequest(context.Background(), esiStatusRequest())

	require.NoError(t, err)
	assert.Equal(t, 200, resp.statusCode)
}

func TestClient_retryRequest_WithErr(t *testing.T) {
	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	client := NewClient()
	req := esiStatusRequest()
	req.page = 2

	resp, err := client.retryRequest(context.Background(), req)
	require.NotNil(t, err)
	require.Nil(t, resp)
}

func TestClient_makeRequest(t *testing.T) {
	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	requestWithoutUserAgent := esiStatusRequest()
	requestWithoutUserAgent.headers["User-Agent"] = []string{}

	type args struct {
		ctx context.Context
		r   request
	}
	tests := []struct {
		name    string
		args    args
		want    *response
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "should make valid request",
			args: args{
				ctx: context.Background(),
				r:   esiStatusRequest(),
			},
			want: &response{
				body:       []byte(esiStatusResponse),
				expires:    esiExpiresTime,
				pages:      1,
				statusCode: 200,
			},
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.Nil(t, err)
			},
		},
		{
			name: "should get a 404 for a missing page",
			args: args{
				ctx: context.Background(),
				r:   esiStatusRequest().withPage(2),
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNilf(t, err, "expected an error with a 404 response")
			},
		},
		{
			name: "should get an error with a missing user agent header",
			args: args{
				ctx: context.Background(),
				r:   requestWithoutUserAgent,
			},
			want: nil,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.NotNilf(t, err, "expected an error with a missing user agent header")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			esi := NewClient()
			got, err := esi.makeRequest(tt.args.ctx, tt.args.r)
			if !tt.wantErr(t, err, fmt.Sprintf("makeRequest(%v, %v)", tt.args.ctx, tt.args.r)) {
				return
			}
			assert.Equalf(t, tt.want, got, "makeRequest(%v, %v)", tt.args.ctx, tt.args.r)
		})
	}
}

func TestClient_makeRequest_WithCtxCancel(t *testing.T) {

	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	esi := NewClient()
	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	_, err := esi.makeRequest(ctx, esiStatusRequest())
	require.NotNilf(t, err, "expected an err with a cancelled context")
}
