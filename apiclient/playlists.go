package apiclient

import (
	"dredly/spotify-sync/cli"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
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
		Next  string `json:"next"`
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
	sourceUrl := apiBaseUrl + "playlists/" + pip.SourceId + "/tracks"
	sourceUris, paginationOpts := getTrackUris(c, token, sourceUrl)
	fmt.Println(sourceUris[0:5])
	counter := 0
	for paginationOpts != "" && counter < 3 {
		fmt.Println("paginationOpts", paginationOpts)
		moreUris, po := getTrackUris(c, token, sourceUrl + "?" + paginationOpts)
		fmt.Println(moreUris[0:5])
		fmt.Println("moreUris has length", len(moreUris))
		fmt.Println("po =", po)
		sourceUris = append(sourceUris, moreUris...)
		paginationOpts = po
		counter ++
	}

	fmt.Println("sourceUris has length", len(sourceUris))
	// srb := syncRequestBody{ Uris: sourceUris }
	// fmt.Println(srb)
	// jsonData, err := json.Marshal(srb)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// req, err := http.NewRequest(http.MethodPost, apiBaseUrl + "playlists/" + pip.DestId + "/tracks", bytes.NewBuffer(jsonData))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// req.Header.Add("Authorization", "Bearer " + token)
	// req.Header.Add("Content-Type", "application/json")

	// resp, err := c.Do(req)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if resp.StatusCode == 201 {
	// 	fmt.Println("Sync complete")
	// } else {
	// 	fmt.Println("Problem with sync response status was " + resp.Status)
	// }
}

func getTrackUris(c http.Client, token string, url string) (uris []string, paginationOpts string) {
	fmt.Println("url =", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	uris = make([]string, len(pr.Tracks.Items))
	for i, item := range pr.Tracks.Items {
		uris[i] = item.Track.Uri
	}

	// TODO: refactor this into a method receiver on playlistResponse struct
	if pr.Tracks.Next != "" {
		spl := strings.Split(pr.Tracks.Next, "?")
		paginationOpts = spl[len(spl) - 1]
	}

	return uris, paginationOpts
}