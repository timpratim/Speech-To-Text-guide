package main

import (
	"Speech-To-Text/models"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/go-audio/wav"
)

func transcribe(modelPath string, audioFilename string) (*[]models.RawTranscription, error) {
	model, err := whisper.New(modelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load model: %w", err)
	}
	defer model.Close()

	log.Println("Successfully loaded the model")

	// Create processing context
	context, err := model.NewContext()
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}

	var data []float32
	// Decode the WAV file - load the full buffer
	data, err = decodePCMBuffer(audioFilename, data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode audio file: %w", err)
	}
	dataLen := len(data)
	//print data len
	log.Println("Audio data length: ", dataLen)
	// if data len is 0 apply ffmpeg

	fmt.Println("Starting the transcription...")

	// Segment callback when -tokens is specified
	results := make([]models.RawTranscription, 0)
	cb := func(segment whisper.Segment) {
		transcription := models.RawTranscription{
			StartTs: segment.Start,
			StopTs:  segment.End,
			Text:    segment.Text,
			Index:   len(results) + 1,
		}
		results = append(results, transcription)
	}
	var pc whisper.ProgressCallback

	if err := context.Process(data, cb, pc); err != nil {
		return nil, err
	}

	log.Println("Got raw transcriptions: %v", results)

	return &results, nil
}

func decodePCMBuffer(audioFilename string, data []float32) ([]float32, error) {
	fh, err := os.Open(audioFilename)

	if err != nil {
		return nil, err
	}
	defer fh.Close()
	dec := wav.NewDecoder(fh)
	if buf, err := dec.FullPCMBuffer(); err != nil {
		return nil, err
	} else if dec.SampleRate != whisper.SampleRate {
		return nil, fmt.Errorf("unsupported sample rate: %d", dec.SampleRate)
	} else if dec.NumChans != 1 {
		return nil, fmt.Errorf("unsupported number of channels: %d", dec.NumChans)
	} else {
		data = buf.AsFloat32Buffer().Data
	}
	return data, nil
}

func OutputSRT(transcriptions *[]models.RawTranscription) {
	for _, transcription := range *transcriptions {
		fmt.Println(srtTimestamp(transcription.StartTs), "-->", srtTimestamp(transcription.StopTs))
		fmt.Println(transcription.Text)
		fmt.Println("n: ", transcription.Index)
	}
}

func srtTimestamp(t time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d", t/time.Hour, (t%time.Hour)/time.Minute, (t%time.Minute)/time.Second, (t%time.Second)/time.Millisecond)
}
