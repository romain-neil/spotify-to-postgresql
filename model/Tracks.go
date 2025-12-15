package model

import "encoding/json"

type Tracks struct {
	Tracks []SongStreaming
}

// UnmarshalJSON supports both top-level array and an object wrapper with a "tracks" property.
func (t *Tracks) UnmarshalJSON(data []byte) error {
	// If the JSON starts with '[', it's a top-level array of SongStreaming
	if len(data) > 0 && data[0] == '[' {
		var items []SongStreaming
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		t.Tracks = items
		return nil
	}

	// Otherwise, try to unmarshal from an object with a `tracks` field
	var aux struct {
		Tracks []SongStreaming `json:"tracks"`
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	t.Tracks = aux.Tracks
	return nil
}
