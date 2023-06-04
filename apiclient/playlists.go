package apiclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/dredly/spotify-sync-go/cli"
	"github.com/dredly/spotify-sync-go/utils"

	"golang.org/x/exp/slices"
)

const (
	apiBaseUrl = "https://api.spotify.com/v1/"
	chunkSize  = 100
)

type (
	syncRequestBody struct {
		Uris []string `json:"uris"`
	}

	tracks struct {
		Items []trackItem `json:"items"`
		Next  string      `json:"next"`
	}
	trackItem struct {
		Track track `json:"track"`
	}
	track struct {
		Uri  string `json:"uri"`
		Name string `json:"name"`
	}
)

func Sync(c http.Client, token string, pip cli.PlaylistIdPair) {
	sourceUris := getAllTrackUris(c, token, pip.SourceId)
	destUris := getAllTrackUris(c, token, pip.DestId)
	urisToAdd := getUrisToAdd(sourceUris, destUris)

	if len(urisToAdd) == 0 {
		fmt.Printf("No new tracks since last sync from playlist %s to playlist %s\n", pip.SourceId, pip.DestId)
		return
	}

	uriChunks := utils.Chunkinator(urisToAdd, chunkSize)
	for _, uriChunk := range uriChunks {
		err := addUrisToTrack(c, token, pip.DestId, uriChunk)
		if err != nil {
			log.Fatal(err)
		}
	}

	fmt.Printf("Sync successful. Added %d new tracks from playlist %s to playlist %s\n", len(urisToAdd), pip.SourceId, pip.DestId)
}

func addUrisToTrack(c http.Client, token string, playlistId string, uris []string) (err error) {
	srb := syncRequestBody{Uris: uris}

	jsonData, err := json.Marshal(srb)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, apiBaseUrl+"playlists/"+playlistId+"/tracks", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode != 201 {
		err = errors.New("Got unexpected response status " + resp.Status)
	}
	return
}

func getAllTrackUris(c http.Client, token string, playlistId string) []string {
	uris, nextLink := getTrackUrisPage(c, token, apiBaseUrl+"playlists/"+playlistId+"/tracks")
	for nextLink != "" {
		moreUris, nl := getTrackUrisPage(c, token, nextLink)
		uris = append(uris, moreUris...)
		nextLink = nl
	}
	return uris
}

func getTrackUrisPage(c http.Client, token string, url string) (uris []string, nextLink string) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var tr tracks
	err = json.Unmarshal(respBody, &tr)
	if err != nil {
		log.Fatal(err)
	}

	uris = make([]string, len(tr.Items))
	for i, item := range tr.Items {
		uris[i] = item.Track.Uri
	}
	nextLink = tr.Next

	return uris, nextLink
}

func getUrisToAdd(sourceUris []string, destUris []string) []string {
	urisToAdd := []string{}
	for _, uri := range sourceUris {
		if !slices.Contains(destUris, uri) {
			urisToAdd = append(urisToAdd, uri)
		}
	}

	return urisToAdd
}
