package apiclient

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
)

const apiBaseUrl = "https://api.spotify.com/v1/"

type (
	playlistResponse struct {
		Uri string `json:"uri"`
		Tracks tracks `json:"tracks"`
	}
	tracks struct {
		Items []trackItem `json:"items"`
	}
	trackItem struct {
		Track track `json:"track"`
	}
	track struct {
		Uri string `json:"uri"`
		Name string `json:"name"`
	}
)

func GetTrackUris(c http.Client, token string, playlistId string) []string {
	req, err := http.NewRequest(http.MethodGet, apiBaseUrl + "playlists/" + playlistId + "/tracks", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer " + token)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var pr playlistResponse
	err = json.Unmarshal(respBody, &pr)
	if err != nil {
		log.Fatal(err)
	}

	uris := make([]string, len(pr.Tracks.Items))
	for i, item := range pr.Tracks.Items {
		uris[i] = item.Track.Uri
	}
	return uris
}