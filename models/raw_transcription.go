package models

import "time"

type RawTranscriptions *[]RawTranscription

func ToModel(transcriptions RawTranscriptions) Transcriptions {
	var modelTranscriptions Transcriptions
	if transcriptions != nil {
		for _, transcription := range *transcriptions {
			modelTranscriptions = append(modelTranscriptions, transcription.ToModel())
		}
		return modelTranscriptions
	}
	return nil

}

type RawTranscription struct {
	StartTs time.Duration
	StopTs  time.Duration
	Text    string
	Index   int
}

// ToModel
func (t *RawTranscription) ToModel() Transcription {
	return Transcription{
		//to string
		StartTime: t.StartTs.String(),
		EndTime:   t.StopTs.String(),
		Text:      t.Text,
		Index:     t.Index,
	}
}