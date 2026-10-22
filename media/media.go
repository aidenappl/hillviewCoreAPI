// Package media sniffs uploaded files. It has no env dependency, so its tests
// run under a plain `go test ./...`.
package media

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrUnsupportedImage means the upload isn't one of imageExtensions' types.
var ErrUnsupportedImage = errors.New("unsupported image type")

// imageExtensions maps the sniffed content type of an upload to the extension
// it is stored under. Everything used to be stored as .jpeg / image/jpeg.
var imageExtensions = map[string]string{
	"image/jpeg": ".jpeg",
	"image/png":  ".png",
	"image/gif":  ".gif",
	"image/webp": ".webp",
}

// ImageType sniffs an upload and returns its content type and file extension.
func ImageType(b []byte) (string, string, error) {
	contentType := http.DetectContentType(b)
	ext, ok := imageExtensions[contentType]
	if !ok {
		return "", "", fmt.Errorf("%w: %s", ErrUnsupportedImage, contentType)
	}
	return contentType, ext, nil
}
