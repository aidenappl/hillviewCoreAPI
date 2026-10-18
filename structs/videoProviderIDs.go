package structs

import "strings"

// muxHost prefixes every stored Mux url: https://stream.mux.com/{playback_id}.m3u8
const muxHost = "https://stream.mux.com/"

// SetProviderIDs fills the streaming provider id derived from the video url:
// cloudflare_id for Cloudflare Stream videos, mux_playback_id for Mux videos.
// Mirrors hillviewVideoAPI's streamurl package.
func (v *Video) SetProviderIDs() {
	if strings.Contains(v.URL, "cloudflare") || strings.Contains(v.URL, "videodelivery.net") {
		if parts := strings.Split(v.URL, "/"); len(parts) >= 3 {
			v.CloudflareID = &parts[len(parts)-3]
		}
	}
	if strings.HasPrefix(v.URL, muxHost) {
		id := strings.TrimPrefix(v.URL, muxHost)
		if i := strings.IndexAny(id, "?#"); i >= 0 {
			id = id[:i]
		}
		id = strings.TrimSuffix(id, ".m3u8")
		if id != "" && !strings.Contains(id, "/") {
			v.MuxPlaybackID = &id
		}
	}
}
