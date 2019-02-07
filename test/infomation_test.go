package test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func getRequest(t *testing.T, url string) *http.Response {
	req, err := http.NewRequest("GET", "http://localhost:8080/restrict"+url, nil)
	assert.Nil(t, err)

	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ODcsImVtYWlsIjoiIiwiYWRtaW4iOmZhbHNlLCJleHAiOjE1Njk5OTk4Mjd9.R-p_kH5BGoB2dzelaa7q7pmC268L1jR7ZmFjxTBcpV4")
	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)
	return resp
}

func TestInfomationeAPI(t *testing.T) {

	t.Run("notice", func(t *testing.T) {
		resp := getRequest(t, "/notice")

		defer resp.Body.Close()

		bytes, _ := ioutil.ReadAll(resp.Body)
		var data map[string]interface{}
		err := json.Unmarshal(bytes, &data)
		assert.Nil(t, err)

		b, err := json.MarshalIndent(data, "", "  ")
		assert.Nil(t, err)

		t.Log(string(b))
	})

	t.Run("event", func(t *testing.T) {
		resp := getRequest(t, "/event")

		defer resp.Body.Close()

		bytes, _ := ioutil.ReadAll(resp.Body)
		var data map[string]interface{}
		err := json.Unmarshal(bytes, &data)
		assert.Nil(t, err)

		b, err := json.MarshalIndent(data, "", "  ")
		assert.Nil(t, err)

		t.Log(string(b))
	})

}
