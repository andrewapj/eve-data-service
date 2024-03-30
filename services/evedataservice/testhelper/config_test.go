package testhelper

import (
	"github.com/andrewapj/arcturus/configuration"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestSetTestConfig(t *testing.T) {

	SetTestConfig()

	v := os.Getenv(configuration.LogLevel)
	require.NotEmpty(t, v, "expected key was empty")
}

func TestGetRootDir(t *testing.T) {

	dir, err := GetRootDir()
	require.Nil(t, err, "unexpected error")
	require.NotEmpty(t, dir, "unexpected empty value")
}
