package main

import (
	"context"
	"fmt"

	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
)

func YoutubeDL(ytID string) error {
	client := youtube.Client{}
	ctx := context.Background()

	// youtube-dl test video
	video, err := client.GetVideoContext(ctx, ytID)
	if err != nil {
		return fmt.Errorf("Error getting video: %w", err)
	}
	downloader := downloader.Downloader{Client: client, OutputDir: "./"}

	downloader.Download(ctx, video, &video.Formats[0], ytID+".mp4")
	return nil

}
