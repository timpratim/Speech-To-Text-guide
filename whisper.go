package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"Speech-To-Text/models"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/go-audio/wav"
)

func transcribe(modelPath string, audioFilename string) error {
	model, err := whisper.New(modelPath)
	if err != nil {
		return fmt.Errorf("failed to load model: %w", err)
	}
	defer model.Close()

	log.Println("Successfully loaded the model")

	// Create processing context
	context, err := model.NewContext()
	if err != nil {
		return fmt.Errorf("failed to create context: %w", err)
	}

	var data []float32
	// Decode the WAV file - load the full buffer
	data, err = decodePCMBuffer(audioFilename, data)
	if err != nil {
		return fmt.Errorf("failed to decode audio file: %w", err)
	}
	dataLen := len(data)
	//print data len
	log.Println("Audio data length: ", dataLen)
	// if data len is 0 apply ffmpeg

	fmt.Println("Starting the transcription...")
	// Segment callback when -tokens is specified
	var cb whisper.SegmentCallback
	var pc whisper.ProgressCallback

	if err := context.Process(data, cb, pc); err != nil {
		return err
	}

	// Print out the results
	transcriptions, err := OutputSRT(context)
	if err != nil {
		return fmt.Errorf("failed to output SRT: %w", err)
	}

	// print got transcriptions
	log.Println("Got raw transcriptions: %v", transcriptions)
	return nil

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

func OutputSRT(context whisper.Context) (*[]models.RawTranscription, error) {
	n := 1
	results := make([]models.RawTranscription, 0)

	for {
		segment, err := context.NextSegment()
		if err != nil {
			break
		}
		transcription := models.RawTranscription{
			StartTs: segment.Start,
			StopTs:  segment.End,
			Text:    segment.Text,
			Index:   n,
		}
		fmt.Println(srtTimestamp(segment.Start), "-->", srtTimestamp(segment.End))
		fmt.Println(segment.Text)
		fmt.Println("n: ", n)
		results = append(results, transcription)
		n++
	}
	return &results, nil
}

func srtTimestamp(t time.Duration) string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d", t/time.Hour, (t%time.Hour)/time.Minute, (t%time.Minute)/time.Second, (t%time.Second)/time.Millisecond)
}
