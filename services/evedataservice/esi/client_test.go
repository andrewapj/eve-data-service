package esi

import (
	"context"
	"fmt"
	"github.com/andrewapj/arcturus/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_makeRequest(t *testing.T) {
	testhelper.SetTestConfig()
	startMockServer(t)
	defer stopMockServer(t)

	requestWithoutUserAgent := newRequest(esiStatusPath, map[string]string{}, 0)
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
				r:   newRequest(esiStatusPath, map[string]string{}, 0),
			},
			want: &response{
				body:       []byte(esiStatusResponse),
				expires:    esiRequestTime,
				pages:      0,
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
				r:   newRequest(esiStatusPath, map[string]string{}, 2),
			},
			want: &response{},
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
			want: &response{},
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
	startMockServer(t)
	defer stopMockServer(t)

	esi := NewClient()
	ctx, cancel := context.WithCancel(context.Background())

	cancel()

	_, err := esi.makeRequest(ctx, newRequest(esiStatusPath, map[string]string{}, 0))
	require.NotNilf(t, err, "expected an err with a cancelled context")
}
