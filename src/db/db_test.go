package db

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var localDBAddr = "localhost:28015"

func TestNewClient(t *testing.T) {
	client, err := NewClient(localDBAddr)
	require.NoError(t, err)
	require.NotNil(t, client)
	require.NotNil(t, client.Session)
}

func TestClose(t *testing.T) {
	client, err := NewClient(localDBAddr)
	require.NoError(t, err)
	require.NoError(t, client.Close())
}
