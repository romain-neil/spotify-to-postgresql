# Spotify Data

## About the project

This project reads your Spotify listening history (song streaming events) from a Spotify data export and loads it into a PostgreSQL database for analysis. It is intended for:
- Individuals who want to explore and visualize their personal listening habits.
- Data enthusiasts and engineers who prefer a reproducible, local pipeline over third‑party dashboards.
- Anyone who wants a structured, queryable store of their Spotify play history.

### What it does
- Parses Spotify export files (the JSON files you get from Spotify's "Download your data" request — especially the Extended Streaming History).
- Maps plays to strongly-typed Go models.
- Inserts streaming events and track metadata into PostgreSQL tables for querying.

### What gets stored
- Streaming events (each play) with timestamps and duration (e.g., `played_at`, `ms_played`).
- Track metadata (e.g., `track_name`, `artist`, `album`) when available.

Model definitions live in:
- `model/SongStreaming.go` – streaming/listening events
- `model/Tracks.go` – track metadata

## Requirements
- Go (1.20+ recommended)
- PostgreSQL (13+ recommended)

## Setup
1. Prepare your Spotify export
   - Request your data from Spotify.
   - After receiving the ZIP from Spotify, extract it. The Extended Streaming History JSON files are commonly named like `Streaming_History_A.json`, `..._B.json`, etc.

2. Configure database access
   - Create a `.env` file at the project root with the following variables:
     - `DB_USER`
     - `DB_NAME`
     - `DB_PASSWORD`
     - `DB_HOST`
     - `DB_PORT`

## Usage
- Quick start during development:
  - `go run ./...`
- Or build and run:
  - `go build -o spotify-data-loader`
  - `./spotify-data-loader`

1. Ensure PostgreSQL is running and accessible with your `.env` values.
2. Run the program by passing the filename within argument `--file=ListeningHistory.json`

If the file is not found, the program will exit with an error.

## Querying your data
Once loaded, you can run SQL queries such as:
- Total listening time by month
- Top artists/tracks over a period
- Plays per hour of day

Use your favorite SQL client or visualization tool (e.g., Metabase, Grafana, Superset) pointed at your PostgreSQL instance.

## Notes & limitations
- This project targets the official Spotify data export (local files), not the live Spotify Web API.
- Some fields in the export may vary over time; the loader focuses on commonly present fields (`trackName`, `artistName`, `msPlayed`, `endTime`, etc.).
- Duplicate handling and idempotency depend on how you run the loader; consult the code for the current strategy.

## Project structure (high level)
- `main.go` — entry point and orchestration
- `model/` — Go structs for tracks and streaming events

## License
The project is licensed under the GPLv3
