package test

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
)

func TestKeyUpdate(t *testing.T) {

	key := "100"

	resultKey, err := strconv.Atoi(key)
	assert.Nil(t, err, err.Error())

	assert.Equal(t, 100, resultKey)

	var increase float32 = 0.1

	totalKey := float32(resultKey) + (float32(resultKey) * increase)

	assert.Equal(t, 110, int(totalKey))

}
