package rest

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type TrackData struct {
	Item struct {
		Album struct {
			Images []struct {
				Url    string `json:"url"`
				Height int    `json:"height"`
				Width  int    `json:"width"`
			} `json:"images"`
		} `json:"album"`
	} `json:"item"`
}

func GetTrackImage() string {
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/me/player", nil)
	req.Header.Set("Authorization", "Bearer "+AccessToken)

	res, _ := http.DefaultClient.Do(req)

	body := new(bytes.Buffer)
	body.ReadFrom(res.Body)

	var data TrackData
	json.Unmarshal(body.Bytes(), &data)

	if len(data.Item.Album.Images) < 1 {
		return ""
	}

	return data.Item.Album.Images[0].Url
}
