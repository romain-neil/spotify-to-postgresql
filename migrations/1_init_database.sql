-- Initial database schema for SongStreaming
-- Generated on 2025-11-09

-- Using PostgreSQL-compatible SQL
CREATE TABLE IF NOT EXISTS song_streaming (
    ts BIGINT NOT NULL,
    username TEXT,
    platform TEXT,
    ms_played INTEGER,
    conn_country TEXT,
    ip_addr_decrypted TEXT,
    user_agent_decrypted TEXT,
    master_metadata_track_name TEXT,
    master_metadata_album_artist_name TEXT,
    master_metadata_album_name TEXT,
    spotify_track_uri TEXT,
    episode_name TEXT,
    episode_show_name TEXT,
    spotify_episode_uri TEXT,
    reason_start TEXT,
    reason_end TEXT,
    shuffle BOOLEAN,
    skipped BOOLEAN,
    offline BOOLEAN,
    offline_timestamp TEXT,
    incognito_mode BOOLEAN
);

-- Helpful indexes (optional but common for analytics/queries)
CREATE INDEX IF NOT EXISTS idx_song_streaming_ts ON song_streaming (ts);
CREATE INDEX IF NOT EXISTS idx_song_streaming_username ON song_streaming (username);
CREATE INDEX IF NOT EXISTS idx_song_streaming_spotify_track_uri ON song_streaming (spotify_track_uri);
