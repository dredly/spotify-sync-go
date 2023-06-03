package apiclient

import (
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
	sourceUris, nextLink := getTrackUris(c, token, sourceUrl)
	counter := 0
	for nextLink != "" && counter < 3 {
		fmt.Println("nextLink", nextLink)
		moreUris, nl := getTrackUris(c, token, nextLink)
		sourceUris = append(sourceUris, moreUris...)
		nextLink = nl
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

func getTrackUris(c http.Client, token string, url string) (uris []string, nextLink string) {
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