package utils

import "os"

func GetEnvWithFallback(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func Chunkinator(slice []string, chunkSize int) [][]string {
	length := len(slice)
	chunkCount := (length + chunkSize - 1) / chunkSize
	chunks := make([][]string, chunkCount)

	for i := 0; i < chunkCount; i++ {
		start := i * chunkSize
		end := start + chunkSize

		if end > length {
			end = length
		}

		chunks[i] = slice[start:end]
	}

	return chunks
}
