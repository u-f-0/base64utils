// Package base64-utils contains utility functions for converting images to base64 format.
package base64utils

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Default Image in case the requested image
// does not exist.
var df string = ""

// Get Default Image
func DefaultImage() string {
	return df
}

// Set Default Image
func SetDefaultImage(img string) {
	df = img
}

// encode is our main function for
// base64 encoding a passed []byte
func encode(bin []byte) []byte {
	e64 := base64.StdEncoding

	maxEncLen := e64.EncodedLen(len(bin))
	encBuf := make([]byte, maxEncLen)

	e64.Encode(encBuf, bin)
	return encBuf
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

// DEPRECATED
// Begin a NewImage to fetch
// TODO: Deprecate NewImage
func NewImage(url string) string {
	return FromRemote(url)
}

// FromRemote is a better named function that
// presently calls NewImage which will be deprecated.
// Function accepts an RFC compliant URL and returns
// a base64 encoded result.
func FromRemote(url string) string {
	image, _ := get(cleanUrl(url))
	enc := encode(image)

	return enc
}

// cleanUrl converts whitespace in urls to %20
func cleanUrl(s string) string {
	return strings.Replace(s, " ", "%20", -1)
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
