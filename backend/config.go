package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

// SubscriberConfig holds PUB/SUB related settings.
type SubscriberConfig struct {
	ThrottleN   int  `yaml:"throttle_n"`
	Deduplicate bool `yaml:"deduplicate"`
}

// CameraConfig represents each camera entry from your YAML config.
type CameraConfig struct {
	ID        string `yaml:"id"`
	Type      string `yaml:"type"`
	Index     int    `yaml:"index,omitempty"`
	URL       string `yaml:"url,omitempty"`
	Thumbnail string `yaml:"thumbnail"`
}

// Config holds all global settings for the backend.
type Config struct {
	Subscriber    SubscriberConfig `yaml:"subscriber"`
	RetentionDays int              `yaml:"retention_days"`
	Cameras       []CameraConfig   `yaml:"cameras"`
}

// loadConfig reads and parses the YAML config file and RETURNS it.
func loadConfig() Config {
	var cfg Config

	data, err := os.ReadFile("../config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config.yaml: %v", err)
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Failed to parse config.yaml: %v", err)
	}

	fmt.Printf("[Config] Loaded: %+v\n", cfg)
	return cfg
}
