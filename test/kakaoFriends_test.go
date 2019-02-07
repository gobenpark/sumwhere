package test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestKakaoFriends(t *testing.T) {

	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v1/friends?limit=3", nil)
	assert.Nil(t, err, err.Error())

	req.Header.Add("Authorization", "Bearer kPAlA9aqU_Cvg940QDlGtHeG0bKaOdVTZ1iZCQopdaYAAAFmxQ31yg")

	client := &http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err, err.Error())

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)
	str := string(bytes)
	fmt.Println(str)
}
