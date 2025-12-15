package model

type SongStreaming struct {
	Ts                            string `json:"ts"`
	Username                      string `json:"username"`
	Platform                      string `json:"platform"`
	MsPlayed                      int    `json:"ms_played"`
	ConnCountry                   string `json:"conn_country"`
	IpAddrDecrypted               string `json:"ip_addr_decrypted"`
	UserAgentDecrypted            string `json:"user_agent_decrypted"`
	MasterMetadataTrackName       string `json:"master_metadata_track_name"`
	MasterMetadataAlbumArtistName string `json:"master_metadata_album_artist_name"`
	// Note: Spotify exports use "master_metadata_album_album_name" (double "album")
	MasterMetadataAlbumName string  `json:"master_metadata_album_album_name"`
	SpotifyTrackUri         string  `json:"spotify_track_uri"`
	EpisodeName             *string `json:"episode_name"`
	EpisodeShowName         *string `json:"episode_show_name"`
	SpotifyEpisodeUri       *string `json:"spotify_episode_uri"`
	ReasonStart             string  `json:"reason_start"`
	ReasonEnd               string  `json:"reason_end"`
	Shuffle                 bool    `json:"shuffle"`
	Skipped                 *bool   `json:"skipped"`
	Offline                 bool    `json:"offline"`
	OfflineTimestamp        int64   `json:"offline_timestamp"`
	IncognitoMode           bool    `json:"incognito_mode"`
}
