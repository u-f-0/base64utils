//Package base64utils contains utility functions for converting images to base64 format.
package base64utils

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Default Image in case the requested image
// does not exist.
var df string = ""

// DefaultImage is function to get the default image
func DefaultImage() string {
	return df
}

//SetDefaultImage is a function to set the default image if image
//does not exist.
func SetDefaultImage(img string) {
	df = img
}

// Lightweight HTTP Client to fetch the image
// Note: This will also pull webpages. @todo
// It is up to the user to use valid image urls.
func get(url string) ([]byte, string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error getting url.")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ct := resp.Header.Get("Content-Type")

	if resp.StatusCode == 200 && len(body) > 512 {
		return body, ct
	}

	if DefaultImage() == "" {
		return []byte(""), ct
	}

	if url == DefaultImage() {
		panic("Catching an infinite loop! Default Image doesn't exist or is broken. Please rectify!")
	}

	return get(DefaultImage())
}

//NewImage to fetch
// TODO: Deprecate NewImage
func NewImage(url string) string {
	return FromRemote(url)
}

// FromRemote is a better named function that
// presently calls NewImage which will be deprecated.
// Function accepts an RFC compliant URL and returns
// a base64 encoded result.
func FromRemote(url string) string {
	image, _ := get(cleanURL(url))

	encoded := base64.StdEncoding.EncodeToString(image)
	out := string(encoded)
	return out
}

//cleanURL converts whitespace in urls to %20
func cleanURL(s string) string {
	return strings.Replace(s, " ", "%20", -1)
}
