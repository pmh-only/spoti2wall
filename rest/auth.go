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

	"github.com/pmh-only/spoti2wall/state"
)

const ClientId = "d963921426b6417188c8a66e17c1bc97"
const ClientSecret = "f40b3c7258ce485690c279a38d2db9d7"

const authorizeUri = "http://localhost:49823/authorize"
const redirectUri = "http://localhost:49823/callback"

var AccessToken = ""
var RefreshToken = ""

type authResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func getClientId() string {
	clientID := state.GlobalConfig.Section("").Key("client_id").String()
	if clientID == "" {
		return ClientId
	} else {
		return clientID
	}
}

func getClientSecret() string {
	secret := state.GlobalConfig.Section("").Key("client_secret").String()
	if secret == "" {
		return ClientSecret
	} else {
		return secret
	}
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

	var credential = base64.StdEncoding.EncodeToString(
		[]byte(getClientId() + ":" + getClientSecret()))
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

	var credential = base64.StdEncoding.EncodeToString(
		[]byte(getClientId() + ":" + getClientSecret()))
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
