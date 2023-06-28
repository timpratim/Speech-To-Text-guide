package models

type Transcription struct {
	Index     int    `json:"index"`
	Text      string `json:"text"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

type Transcriptions []Transcription

func (t Transcriptions) ToMap() []map[string]interface{} {
	results := make([]map[string]interface{}, len(t))
	for i, transcription := range t {
		results[i] = transcription.ToMap()
	}
	return results
}

// to map[string]interface{}
func (t Transcription) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"index":     t.Index,
		"text":      t.Text,
		"startTime": t.StartTime,
		"endTime":   t.EndTime,
	}
}