package testhelper

import (
	"github.com/andrewapj/arcturus/config"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_SetTestConfig(t *testing.T) {

	SetTestConfig()

	v := config.LogLevel()
	require.NotEmpty(t, v, "expected key was empty")
}

func Test_GetRootDir(t *testing.T) {

	dir, err := GetRootDir()
	require.Nil(t, err, "unexpected error")
	require.NotEmpty(t, dir, "unexpected empty value")
}
