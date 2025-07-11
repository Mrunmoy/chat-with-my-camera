package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// SubscriberConfig holds PUB/SUB related settings.
type SubscriberConfig struct {
	ThrottleN   int  `yaml:"throttle_n"`   // seconds to throttle same events
	Deduplicate bool `yaml:"deduplicate"`  // deduplication flag
}

// Config holds all global settings for the backend.
type Config struct {
	Subscriber    SubscriberConfig `yaml:"subscriber"`
	RetentionDays int              `yaml:"retention_days"`
}

var config Config

// loadConfig reads and parses the YAML config file.
func loadConfig() {
	// Adjust the path if your config lives under /config
	data, err := os.ReadFile("../config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config.yaml: %v", err)
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Failed to parse config.yaml: %v", err)
	}

	fmt.Printf("[Config] Loaded: %+v\n", config)
}
