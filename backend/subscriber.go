package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	zmq4 "github.com/pebbe/zmq4"
)

// lastLabels and lastSaved keep track of last processed event for deduplication/throttling.
var lastEvents = make(map[string]string)
var lastSaved = make(map[string]time.Time)


// runSubscriber connects to the ZeroMQ PUB socket and writes each detection to SQLite.
func (app *App) runSubscriber() {
	fmt.Println("[ZeroMQSubscriber] Connecting to tcp://localhost:5555...")

	// Create ZeroMQ context and SUB socket
	context, err := zmq4.NewContext()
	if err != nil {
		log.Fatalf("Failed to create ZeroMQ context: %v", err)
	}
	defer context.Term()

	subscriber, err := context.NewSocket(zmq4.SUB)
	if err != nil {
		log.Fatalf("Failed to create SUB socket: %v", err)
	}
	defer subscriber.Close()

	// Connect to publisher
	err = subscriber.Connect("tcp://localhost:5555")
	if err != nil {
		log.Fatalf("Failed to connect to publisher: %v", err)
	}

	// Subscribe to ALL messages
	subscriber.SetSubscribe("")

	fmt.Println("[ZeroMQSubscriber] Connected! Waiting for messages...")

	// Make sure snapshot folder exists
	if _, err := os.Stat("./snapshots"); os.IsNotExist(err) {
		err := os.MkdirAll("./snapshots", os.ModePerm)
		if err != nil {
			log.Fatalf("Failed to create snapshots folder: %v", err)
		}
	}

	// Loop forever: receive -> parse -> filter -> save
	for {
		msg, err := subscriber.RecvMessage(0)
		if err != nil {
			log.Printf("Failed to receive message: %v", err)
			continue
		}

		if len(msg) == 0 {
			continue
		}

		raw := msg[0]
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(raw), &event); err != nil {
			log.Printf("Failed to parse JSON: %v", err)
			continue
		}

		// Extract metadata
		timestamp, _ := event["timestamp"].(float64)
		cameraID, _ := event["camera_id"].(string)
		labelsJSON, _ := json.Marshal(event["labels"])
		labelsStr := string(labelsJSON)
		boxesJSON, _ := json.Marshal(event["boxes"])

		lastEvent, found := lastEvents[cameraID]
		lastTime := lastSaved[cameraID]

		if app.Config.Subscriber.Deduplicate {
			if labelsStr == lastEvent  && found {
				if app.Config.Subscriber.ThrottleN > 0 {
					if time.Since(lastTime) < time.Duration(app.Config.Subscriber.ThrottleN)*time.Second {
						continue // Skip duplicate within throttle window
					}
				} else {
					continue // Skip all duplicates
				}
			}
		}
		// Handle snapshot (optional)
		snapshotPath := ""
		if snapshotB64, ok := event["snapshot"].(string); ok && snapshotB64 != "" {
			jpgBytes, err := base64.StdEncoding.DecodeString(snapshotB64)
			if err != nil {
				log.Printf("Failed to decode snapshot: %v", err)
			} else {
				filename := fmt.Sprintf("./snapshots/%s_%.0f.jpg", cameraID, timestamp)
				err := os.WriteFile(filename, jpgBytes, 0644)
				if err != nil {
					log.Printf("Failed to save snapshot: %v", err)
				} else {
					snapshotPath = filename
				}
			}
		}

		// Insert into SQLite
		err = insertDetection(timestamp, cameraID, labelsStr, string(boxesJSON), snapshotPath)
		if err != nil {
			log.Printf("Failed to insert detection: %v", err)
		} else {
			fmt.Printf("[ZeroMQSubscriber] Logged event: cam=%s labels=%s\n", cameraID, labelsStr)
		}

		// Update dedup state
		lastEvents[cameraID] = labelsStr
		// lastEvents[cameraID] = labelsAndBoxesStr
		lastSaved[cameraID] = time.Now()
	}
}
