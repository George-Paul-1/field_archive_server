package entities

import "time"

type Recording struct {
	ID              int
	Title           string
	AudioLocation   string
	ArtworkLocation *string
	DateUploaded    *time.Time
	RecordingDate   time.Time
	LocationID      int
	Duration        int
	Format          string
	Description     string
	Equipment       string
	Size            float64
	Channels        string
	License         string
}
