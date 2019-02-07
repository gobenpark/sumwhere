package test

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"sumwhere/models"
	"testing"
)

func TestReportAPI(t *testing.T) {
	t.Run("post test", func(t *testing.T) {
		report := models.Report{
			TargetUserID: 99,
			ReportType:   1,
			Comment:      "안녕",
		}
		pbytes, _ := json.Marshal(report)
		buff := bytes.NewBuffer(pbytes)
		req, err := http.NewRequest("POST", "http://localhost:8080/restrict/report", buff)
		assert.Nil(t, err)

		req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6ODcsImVtYWlsIjoiIiwiYWRtaW4iOmZhbHNlLCJleHAiOjE1Njk5OTk4Mjd9.R-p_kH5BGoB2dzelaa7q7pmC268L1jR7ZmFjxTBcpV4")
		client := &http.Client{}

		resp, err := client.Do(req)
		defer resp.Body.Close()
		bytes, _ := ioutil.ReadAll(resp.Body)
		var data map[string]interface{}
		err = json.Unmarshal(bytes, &data)
		assert.Nil(t, err)

		b, err := json.MarshalIndent(data, "", "  ")
		assert.Nil(t, err)

		t.Log(string(b))

	})
}
