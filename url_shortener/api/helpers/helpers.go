package helpers

import (
	"os"
	"strings"
)

func RemoveDomainError(URL string) bool {
	if URL == os.Getenv("DOMAIN") {
		return false
	}
	newURL := strings.Replace(URL, "http://", "", 1)
	newURL = strings.Replace(newURL, "https://", "", 1)
	newURL = strings.Replace(newURL, "www.", "", 1)
	newURL = strings.Split(newURL, "/")[0]

	if newURL == os.Getenv("DOMAIN") {
		return false
	}
	
	return true
}

func EnforceHTTP(URL string) string{
	if URL[:5] != "http" {
		return "http://" + URL
	}
	return URL
}