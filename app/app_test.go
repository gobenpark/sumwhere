package app

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"path/filepath"
	"runtime"
	"testing"
)

func TestNewApp(t *testing.T) {

	app := NewApp()
	require.NotNil(t, app)
}

func TestPath(t *testing.T) {
	var (
		_, b, _, _ = runtime.Caller(0)
		basepath   = filepath.Dir(b)
	)

	fmt.Println(b)
	fmt.Println(basepath)
}
