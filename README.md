# Speech-To-Text-guide



**Workshop Guide Outline**

**1. Introduction to CLI Applications using `cli` package in Go**

In this part, we'll cover how to set up a Command Line Interface (CLI) application using the "cli" package in Go. We will start with a simple program that accepts the "get" command for fetching transcriptions using a supplied YouTube link. We'll just print back the youtube link for now.

**Task**: Create a new Go file in a new git branch named "cli-setup". Now, write a basic program that defines a CLI application which accepts the "get" command.

**Depedencies**
`go get github.com/urfave/cli/v2`


**Code Snippets**

```go
package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2" //this imports the cli package
)

func main() {
	app := &cli.App{
		Name:  "ytt",
		Usage: "Transcribe YouTube videos",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get transcriptions by ytlink",
				Action: func(c *cli.Context) error {
					//Print YouTube link
					youtubelink := c.Args().Get(0)
					if youtubelink == "" {
						return cli.NewExitError("Please provide a YouTube link", 1)
				    }
					fmt.Println("YouTube link:", youtubelink)
					return nil
				},
			},
		},
	
   	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

```

**2. Downloading YouTube Video and Audio Streams**

Now, we'll use the `youtube/v2` package to fetch video details and download the audio stream for the supplied YouTube link. 

**Task**: In a new git branch named "downloading-video", append code to your program from Step 1 for downloading a video using a YouTube link. Add a "download" command to ur CLI for testing this functionality.

**Depedencies**

`go get github.com/kkdai/youtube/v2`

**Code Snippets**

```go
//youtube.go
package main

import (
	"context"
	"fmt"

	"github.com/kkdai/youtube/v2"
	"github.com/kkdai/youtube/v2/downloader"
)

func YoutubeDL(ytlink string) error {
	client := youtube.Client{
		Debug: true,
	}
	ctx := context.Background()


	video, err := client.GetVideoContext(ctx, ytlink)
	if err != nil {
		return fmt.Errorf("Error getting video: %w", err)
	}
	downloader := downloader.Downloader{Client: client, OutputDir: "./"}

	return downloader.DownloadComposite(ctx, "", video, "hd1080", "mp4")

}
```

```go
//main.go
err := YoutubeDL(youtubelink)
if err != nil {
    return err
}
```
**3. Audio Conversion Using FFmpeg**

Learn the use of the FFmpeg tool to convert the audio file to a suitable format for transcription. This module may also cover how to install FFmpeg, if not already available on participants' machines. 

**Task**: In a new git branch "ffmpeg-conversion", append the current program you have with the file conversion using FFmpeg code. Add a "convert" command to your CLI for testing this step.

**Code Snippets**

```go
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func ConvertFile(inputFile string, outputFile string) error {
	cmd := exec.Command("ffmpeg", "-i", inputFile, "-vn", "-ac", "1", "-ar", "16000", "-codec:a", "pcm_s16le", "-f", "wav", outputFile)
	cmd.Stderr = os.Stderr
	err := cmd.Run()

	if err != nil {
		return fmt.Errorf("ffmpeg conversion failed: %w", err)
	}

	return nil
}
```
```go
//main.go
fmt.Println("Audio file:", audieofilename)
err = ConvertFile(audieofilename+".mp4", audieofilename+".wav")
if err != nil {
    return err
}
```
