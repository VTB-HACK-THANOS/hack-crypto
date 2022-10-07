package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHashAndSalt(t *testing.T) {
	password := "test"
	hash, err := HashAndSalt([]byte(password))
	require.NoError(t, err)
	t.Logf(hash)
}
