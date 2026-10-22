package structs

import "strings"

// muxHost prefixes every stored Mux url: https://stream.mux.com/{playback_id}.m3u8
const muxHost = "https://stream.mux.com/"

// SetProviderIDs fills mux_playback_id from a Mux video url. Mirrors
// hillviewVideoAPI's streamurl package.
func (v *Video) SetProviderIDs() {
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
