package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"spotifyData/model"

	_ "github.com/lib/pq"
)

const insertSQL = "INSERT INTO song_streaming (" +
	"ts, username, platform, ms_played, conn_country, ip_addr_decrypted," +
	" user_agent_decrypted, master_metadata_track_name, master_metadata_album_artist_name," +
	"master_metadata_album_name, spotify_track_uri, episode_name, episode_show_name," +
	"spotify_episode_uri, reason_start, reason_end, shuffle, skipped, offline," +
	"offline_timestamp, incognito_mode" +
	") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21);"

func parseFileArg() string {
	// Custom help/usage text shown on -h / -help
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "spotifyData imports Spotify streaming history JSON into PostgreSQL.\n\n")
		fmt.Fprintf(os.Stderr, "Usage:\n  %s -file path/to/Streaming_History.json\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Environment variables (required):\n  DB_USER, DB_NAME, DB_PASSWORD, DB_HOST, DB_PORT\n\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		flag.PrintDefaults()

		os.Exit(2)
	}

	fileName := flag.String("file", "", "JSON file to read from (required)")
	flag.Parse()
	if *fileName == "" {
		// Show usage and exit with a clear message when -file is missing
		flag.Usage()
		log.Fatal("missing required -file argument")
	}
	return *fileName
}

func openPostgres() *sql.DB {
	var (
		dbUser     = os.Getenv("DB_USER")
		dbName     = os.Getenv("DB_NAME")
		dbPassword = os.Getenv("DB_PASSWORD")
		dbHost     = os.Getenv("DB_HOST")
		dbPort     = os.Getenv("DB_PORT")
	)

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s port=%s", dbUser, dbName, dbPassword, dbHost, dbPort)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func loadTracksFromFile(filePath string) (model.Tracks, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return model.Tracks{}, err
	}
	defer func(f *os.File) {
		err = f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	b, err := io.ReadAll(f)
	if err != nil {
		return model.Tracks{}, err
	}

	var export model.Tracks
	if err := json.Unmarshal(b, &export); err != nil {
		return model.Tracks{}, err
	}
	return export, nil
}

func prepareInsertStatement(tx *sql.Tx) (*sql.Stmt, error) {
	return tx.Prepare(insertSQL)
}

func insertTracksInBatches(db *sql.DB, export model.Tracks, batchSize int) error {
	if batchSize <= 0 {
		batchSize = 50
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		// In case of panic, attempt rollback to not leave TX open
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	stmt, err := prepareInsertStatement(tx)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {
			_ = tx.Rollback()
		}
	}(stmt)

	count := 0
	for _, track := range export.Track {
		if _, err = stmt.Exec(
			track.Ts,
			track.Username,
			track.Platform,
			track.MsPlayed,
			track.ConnCountry,
			track.IpAddrDecrypted,
			track.UserAgentDecrypted,
			track.MasterMetadataTrackName,
			track.MasterMetadataAlbumArtistName,
			track.MasterMetadataAlbumName,
			track.SpotifyTrackUri,
			track.EpisodeName,
			track.EpisodeShowName,
			track.SpotifyEpisodeUri,
			track.ReasonStart,
			track.ReasonEnd,
			track.Shuffle,
			track.Skipped,
			track.Offline,
			track.OfflineTimestamp,
			track.IncognitoMode,
		); err != nil {
			_ = tx.Rollback()
			return err
		}
		count++

		if count%batchSize == 0 {
			if err = tx.Commit(); err != nil {
				_ = tx.Rollback()
				return err
			}
			// start new batch
			tx, err = db.Begin()
			if err != nil {
				return err
			}
			err = stmt.Close()
			if err != nil {
				return err
			}
			stmt, err = prepareInsertStatement(tx)
			if err != nil {
				_ = tx.Rollback()
				return err
			}
		}
	}

	// commit remaining
	if err := tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

func main() {
	// Parse CLI
	filePath := parseFileArg()

	// Open DB
	db := openPostgres()
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	// Load tracks from file
	export, err := loadTracksFromFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// Insert in batches
	if err := insertTracksInBatches(db, export, 50); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Données importées avec succès !")
}
