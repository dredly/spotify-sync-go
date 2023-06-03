package cli

import (
	"fmt"
	"log"
	"os"
)

type PlaylistIdPair struct {
	SourceId string
	DestId   string
}

func GetPlaylistIdPairs() []PlaylistIdPair {
	playlistIds := os.Args[1:]
	if len(playlistIds) == 0 {
		// TODO: Add usage info here
		log.Fatal("No playlist ids Provided")
	}
	if len(playlistIds)%2 != 0 {
		log.Fatal("Each source playlist must have a destionation")
	}

	fmt.Printf("Running spotify-sync with playlist ids %v", playlistIds)
	return getPlaylistPairsFromIds(playlistIds)
}

func getPlaylistPairsFromIds(ids []string) []PlaylistIdPair {
	pairs := make([]PlaylistIdPair, len(ids)/2)
	for i := 0; i < len(ids); i += 2 {
		pairs[i/2] = PlaylistIdPair{SourceId: ids[i], DestId: ids[i+1]}
	}
	return pairs
}
