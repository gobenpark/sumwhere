package utils

import (
	"encoding/json"
	"errors"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"reflect"
	"sumwhere/models"
)

// 페이스북 모델 반환
func FacebookUtil(c echo.Context) (*models.FaceBookUser, error) {
	req := c.Request()
	token := req.PostFormValue("access_token")

	url := "https://graph.facebook.com/v3.0/me?fields=id,email,name&access_token="
	res, err := http.Get(url + token)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var facebookuser models.FaceBookUser
	if err := json.Unmarshal(data, &facebookuser); err != nil {
		return nil, err
	}

	if reflect.DeepEqual(facebookuser, models.FaceBookUser{}) {
		return nil, errors.New("empty facebook model")
	}
	facebookuser.Token = token

	return &facebookuser, nil
}

func KakaoUtil(c echo.Context) (*models.KakaoUser, error) {
	req := c.Request()
	token := req.PostFormValue("access_token")
	req, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "bearer "+token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	user, err := ioutil.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var kakaoUser models.KakaoUser
	json.Unmarshal(user, &kakaoUser)

	kakaoUser.Token = token
	if reflect.DeepEqual(kakaoUser, models.KakaoUser{}) {
		return nil, errors.New("empty kakao model")
	}

	return &kakaoUser, nil
}
