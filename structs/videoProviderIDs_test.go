package structs

import (
	"testing"
)

func TestVideoSetProviderIDs(t *testing.T) {
	tests := []struct {
		name           string
		url            string
		wantCloudflare string
		wantMux        string
	}{
		{"stream customer subdomain", "https://customer-nakrsdfbtn3mdz5z.cloudflarestream.com/30a81e2078e321a9784b97d6a7074398/manifest/video.m3u8", "30a81e2078e321a9784b97d6a7074398", ""},
		{"videodelivery.net", "https://videodelivery.net/e7e125f35c62ae1a7b9a3eb354e56cb9/manifest/video.m3u8", "e7e125f35c62ae1a7b9a3eb354e56cb9", ""},
		{"mux", "https://stream.mux.com/abc123.m3u8", "", "abc123"},
		{"mux with query", "https://stream.mux.com/abc123.m3u8?max_resolution=720p", "", "abc123"},
		{"s3", "https://content.hillview.tv/videos/uploads/a.mp4", "", ""},
	}
	deref := func(s *string) string {
		if s == nil {
			return ""
		}
		return *s
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Video{URL: tt.url}
			v.SetProviderIDs()
			if got := deref(v.CloudflareID); got != tt.wantCloudflare {
				t.Errorf("CloudflareID = %q, want %q", got, tt.wantCloudflare)
			}
			if got := deref(v.MuxPlaybackID); got != tt.wantMux {
				t.Errorf("MuxPlaybackID = %q, want %q", got, tt.wantMux)
			}
		})
	}
}
