package utils

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test(t *testing.T) {
	adapter, err := NewFireBaseApp()
	assert.NoError(t, err)
	assert.NotNil(t, adapter)

	result, err := adapter.SetSubscribe(context.Background(), true, []string{":APA91bFzm-wwJQucq_PifhaeafNF7fzM8sTSd2baZmduHJ5lPH1gjipvQxQ1Y6vtxaui3dH74NmaFXfUrd52XeFd8YYlAcB1r3WxM96PXUY1IOVL79PKYvGTGzLNbYIpUUyBG4IQW1SK"}, "ChatAlert")
	assert.Equal(t, 1, result)
	assert.NoError(t, err)
}
