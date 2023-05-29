package apiclient

import (
	"fmt"
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

	// Just to test
	fmt.Println("Token is " + token)
}