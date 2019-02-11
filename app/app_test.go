package app

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewApp(t *testing.T) {

	app := NewApp()
	require.NotNil(t, app)
}
