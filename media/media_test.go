package media

import (
	"errors"
	"testing"
)

func TestImageType(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantType string
		wantExt  string
		wantErr  bool
	}{
		{"jpeg", []byte("\xff\xd8\xff\xe0\x00\x10JFIF\x00"), "image/jpeg", ".jpeg", false},
		{"png", []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR"), "image/png", ".png", false},
		{"gif", []byte("GIF89a\x01\x00\x01\x00"), "image/gif", ".gif", false},
		{"webp", []byte("RIFF\x24\x00\x00\x00WEBPVP8 "), "image/webp", ".webp", false},
		{"text", []byte("hello world"), "", "", true},
		{"pdf", []byte("%PDF-1.7\n"), "", "", true},
		{"empty", nil, "", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotType, gotExt, err := ImageType(tt.data)
			if tt.wantErr {
				if !errors.Is(err, ErrUnsupportedImage) {
					t.Fatalf("err = %v, want ErrUnsupportedImage", err)
				}
				return
			}
			if err != nil || gotType != tt.wantType || gotExt != tt.wantExt {
				t.Fatalf("ImageType = (%q, %q, %v), want (%q, %q, nil)", gotType, gotExt, err, tt.wantType, tt.wantExt)
			}
		})
	}
}
