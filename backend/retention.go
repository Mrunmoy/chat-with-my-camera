package main

import (
	// "database/sql"
	"fmt"
	"log"
	"os"
	"time"
)

// runRetentionJob deletes old rows + files beyond retention window.
func runRetentionJob() {
	retentionDays := 5 // Keep only last 5 days
	ticker := time.NewTicker(1 * time.Hour) // How often to run
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println("[Retention] Running cleanup...")
			now := float64(time.Now().Unix())
			cutoff := now - (float64(retentionDays) * 86400) // 5 days in seconds

			rows, err := db.Query("SELECT id, snapshot_file FROM detections WHERE timestamp < ?", cutoff)
			if err != nil {
				log.Printf("Retention query failed: %v", err)
				continue
			}

			var deleteIDs []int64
			for rows.Next() {
				var id int64
				var snapshot string
				if err := rows.Scan(&id, &snapshot); err != nil {
					log.Printf("Retention row scan failed: %v", err)
					continue
				}

				// Delete snapshot file if it exists
				if snapshot != "" {
					err := os.Remove(snapshot)
					if err != nil && !os.IsNotExist(err) {
						log.Printf("Failed to delete snapshot %s: %v", snapshot, err)
					} else {
						log.Printf("[Retention] Deleted snapshot: %s", snapshot)
					}
				}

				deleteIDs = append(deleteIDs, id)
			}
			rows.Close()

			// Delete rows by ID
			for _, id := range deleteIDs {
				_, err := db.Exec("DELETE FROM detections WHERE id = ?", id)
				if err != nil {
					log.Printf("Failed to delete row id=%d: %v", id, err)
				} else {
					log.Printf("[Retention] Deleted row id=%d", id)
				}
			}

			fmt.Printf("[Retention] Done. Cutoff: %.0f, Deleted: %d rows.\n", cutoff, len(deleteIDs))
		}
	}
}
