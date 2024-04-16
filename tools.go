package claude

import "encoding/json"

func MarshalMediaContent(in []*MessageContent) string {
	if len(in) < 1 {
		return ""
	}
	b, _ := json.Marshal(in)
	return string(b)
}
