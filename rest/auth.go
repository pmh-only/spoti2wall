package rest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const clientId = "cec309f90be94ee6bc2f90f95c9882e0"
const clientSecret = "be14629cf6e6477d85fd61443538e4d8"
const authorizeUri = "http://localhost:49823/authorize"
const redirectUri = "http://localhost:49823/callback"

var AccessToken = ""
var RefreshToken = ""

var credential = base64.StdEncoding.EncodeToString(
	[]byte(clientId + ":" + clientSecret))

type authResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func KeepRefreshToken() {
	for {
		time.Sleep(3000 * time.Second)
		AccessToken = RefreshAccessToken(RefreshToken)
	}
}

func getAuthToken(code string) (access, refresh string) {
	authUrl := url.URL{
		Scheme: "https",
		Host:   "accounts.spotify.com",
		Path:   "/api/token",
	}

	formData := url.Values{
		"code":         {code},
		"redirect_uri": {redirectUri},
		"grant_type":   {"authorization_code"},
	}

	req, _ := http.NewRequest(
		"POST",
		authUrl.String(),
		strings.NewReader(formData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", credential))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	body := new(bytes.Buffer)
	body.ReadFrom(res.Body)

	var data authResponse
	_ = json.Unmarshal(body.Bytes(), &data)

	return data.AccessToken, data.RefreshToken
}

func RefreshAccessToken(refresh string) string {
	authUrl := url.URL{
		Scheme: "https",
		Host:   "accounts.spotify.com",
		Path:   "/api/token",
	}

	formData := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {refresh},
	}

	req, _ := http.NewRequest(
		"POST",
		authUrl.String(),
		strings.NewReader(formData.Encode()))

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", credential))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	body := new(bytes.Buffer)
	body.ReadFrom(res.Body)

	var data authResponse
	_ = json.Unmarshal(body.Bytes(), &data)

	return data.AccessToken
}
