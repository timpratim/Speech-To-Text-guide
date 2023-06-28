package main

import "io"

type ProgressWriter struct {
	writer        io.Writer
	totalBytes    int64
	contentLength int64
	onProgress    func(written int64, progressPercentage float64)
}

func (pw *ProgressWriter) Write(p []byte) (int, error) {
	n, err := pw.writer.Write(p)
	pw.totalBytes += int64(n)

	if pw.onProgress != nil {
		progressPercentage := float64(pw.totalBytes) / float64(pw.contentLength) * 100
		pw.onProgress(pw.totalBytes, progressPercentage)
	}

	return n, err
}

func NewProgressWriter(w io.Writer, contentLength int64, onProgress func(written int64, progressPercentage float64)) *ProgressWriter {
	return &ProgressWriter{
		writer:        w,
		totalBytes:    0,
		contentLength: contentLength,
		onProgress:    onProgress,
	}
}