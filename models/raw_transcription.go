package models

import "time"

type RawTranscription struct {
	StartTs time.Duration
	StopTs  time.Duration
	Text    string
	Index   int
}
