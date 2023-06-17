package apiclient

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func GetRefreshTokenFromFileIfPresent() string {
	homeDirName, err := os.UserHomeDir()
    if err != nil {
        log.Fatal( err )
    }
	tokenFilePath := filepath.Join(homeDirName, ".spotify-sync", "refresh.txt")
	content, err := os.ReadFile(tokenFilePath)
	if err != nil {
		fmt.Printf("Error reading file, %s", err)
		return ""
	}
	return string(content)
}

// Save the refresh token to file ~/.spotify-sync/refresh.txt
func SaveToken(token string) error {
	homeDirName, err := os.UserHomeDir()
    if err != nil {
        log.Fatal( err )
    }
	syncDirPath := filepath.Join(homeDirName, ".spotify-sync")
	mkDirErr := os.Mkdir(syncDirPath, 0777)
	if mkDirErr != nil {
		fmt.Printf("Error creating dir: %s\n", mkDirErr)
		return mkDirErr
	}
	tokenFilePath := filepath.Join(syncDirPath, "refresh.txt")
	d := []byte(token)
	writeErr := os.WriteFile(tokenFilePath, d, 0777)
	if writeErr != nil {
		fmt.Printf("Error writing to file: %s\n", writeErr)
		return writeErr
	}
	return nil
}