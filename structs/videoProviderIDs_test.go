package structs

import "testing"

func TestVideoSetProviderIDs(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantMux string
	}{
		{"mux", "https://stream.mux.com/abc123.m3u8", "abc123"},
		{"mux with query", "https://stream.mux.com/abc123.m3u8?max_resolution=720p", "abc123"},
		{"old stream url", "https://videodelivery.net/e7e125f35c62ae1a7b9a3eb354e56cb9/manifest/video.m3u8", ""},
		{"s3", "https://content.hillview.tv/videos/uploads/a.mp4", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := Video{URL: tt.url}
			v.SetProviderIDs()
			got := ""
			if v.MuxPlaybackID != nil {
				got = *v.MuxPlaybackID
			}
			if got != tt.wantMux {
				t.Errorf("MuxPlaybackID = %q, want %q", got, tt.wantMux)
			}
		})
	}
}
