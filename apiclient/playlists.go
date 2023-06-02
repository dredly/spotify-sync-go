package apiclient

import (
	"bytes"
	"dredly/spotify-sync/cli"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const apiBaseUrl = "https://api.spotify.com/v1/"

type (
	syncRequestBody struct {
		Uris []string `json:"uris"`
	}

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

func Sync(c http.Client, token string, pip cli.PlaylistIdPair) {
	sourceUris := getTrackUris(c, token, pip.SourceId)
	srb := syncRequestBody{ Uris: sourceUris }
	jsonData, err := json.Marshal(srb)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, apiBaseUrl + "playlists/" + pip.DestId + "/tracks", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer " + token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode == 201 {
		fmt.Println("Sync complete")
	} else {
		fmt.Println("Problem with sync response status was " + resp.Status)
	}
}

func getTrackUris(c http.Client, token string, playlistId string) []string {
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