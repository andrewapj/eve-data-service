package esi

import (
	"context"
	"github.com/andrewapj/arcturus/testhelper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestClient_FetchIds(t *testing.T) {

	testhelper.SetTestConfig()
	startMockServer()
	defer stopMockServer()

	client := NewClient()

	ids, err := client.FetchTypeIds(context.Background(), NewPageRequest(1, 2))
	require.NoError(t, err)

	assert.Equal(t, Ids{Ids: []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}}, ids)
}
