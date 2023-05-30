package main

import (
	"dredly/spotify-sync/apiclient"
	"dredly/spotify-sync/browserautomation"
	"dredly/spotify-sync/cli"
	"dredly/spotify-sync/echoserver"
	"fmt"
)

func main() {
	playlistIdPairs := cli.GetPlaylistIdPairs()
	firstPlaylistIdPair := playlistIdPairs[0]

	authCodeChan := make(chan string)
	e := echoserver.SpinUpTempServer(authCodeChan)

	go browserautomation.AutoLogin()
	c := *apiclient.NewHttpClient()

	authCode := <-authCodeChan
	echoserver.GracefulShutdown(e)
	t := apiclient.GetAccessToken(c, authCode)
	sourceUris := apiclient.GetTrackUris(c, t, firstPlaylistIdPair.SourceId)
	destUris := apiclient.GetTrackUris(c, t, firstPlaylistIdPair.DestId)
	fmt.Println("sourceUris", sourceUris)
	fmt.Println("destUris", destUris)
}
