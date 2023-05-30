package apiclient

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const apiBaseUrl = "https://api.spotify.com/v1/"

func GetDestinationTrackUris(c http.Client, token string, playlistId string) {
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

	fmt.Println("respBody")
	fmt.Println(string(respBody))
}