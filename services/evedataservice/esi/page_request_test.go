package esi

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewPageRequest(t *testing.T) {

	assert.Equal(t, PageRequest{page: 0, size: 1}, NewPageRequest(-1, -1))
	assert.Equal(t, PageRequest{page: 0, size: 1}, NewPageRequest(-1, 0))
	assert.Equal(t, PageRequest{page: 0, size: 1}, NewPageRequest(0, 0))
	assert.Equal(t, PageRequest{page: 1, size: 1}, NewPageRequest(1, -1))
	assert.Equal(t, PageRequest{page: 1, size: 2}, NewPageRequest(1, 2))
}

func TestPageRequest_pages1(t *testing.T) {

	assert.Equal(t, []int{0}, NewPageRequest(0, 5).pages())
	assert.Equal(t, []int{1, 2, 3, 4, 5}, NewPageRequest(1, 5).pages())
	assert.Equal(t, []int{6, 7, 8, 9, 10}, NewPageRequest(6, 5).pages())
}
