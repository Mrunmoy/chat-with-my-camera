package main

import (
	"log"
	"os"
	"time"
)

// runRetention runs in a loop: every hour it deletes events older than retention_days.
func (app *App) runRetention() {
	retentionDays := app.Config.RetentionDays
	for {
		log.Printf("[Retention] Running... (keeping last %d days)", retentionDays)

		// Calculate cutoff timestamp
		cutoff := time.Now().AddDate(0, 0, -retentionDays).Unix()

		// Query old rows to get snapshot filenames
		rows, err := app.DB.Query("SELECT snapshot_file FROM detections WHERE timestamp < ?", cutoff)
		if err != nil {
			log.Printf("Retention query failed: %v", err)
			time.Sleep(time.Hour)
			continue
		}

		var snapshotsToDelete []string
		for rows.Next() {
			var snap string
			if err := rows.Scan(&snap); err == nil && snap != "" {
				snapshotsToDelete = append(snapshotsToDelete, snap)
			}
		}
		rows.Close()

		// Delete snapshots from disk
		for _, snap := range snapshotsToDelete {
			if err := os.Remove(snap); err != nil {
				log.Printf("Failed to remove snapshot: %s (%v)", snap, err)
			} else {
				log.Printf("Deleted snapshot: %s", snap)
			}
		}

		// Delete old rows from DB
		res, err := app.DB.Exec("DELETE FROM detections WHERE timestamp < ?", cutoff)
		if err != nil {
			log.Printf("Retention delete failed: %v", err)
		} else {
			affected, _ := res.RowsAffected()
			log.Printf("[Retention] Deleted %d rows older than cutoff", affected)
		}

		// Sleep until next run
		time.Sleep(time.Hour)
	}
}
