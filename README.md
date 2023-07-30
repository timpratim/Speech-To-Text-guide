# Speech-To-Text-transcriptions using Whisper.cpp-guide

**Workshop Guide Outline**

Link to presentation slides: https://docs.google.com/presentation/d/1e5PEJu6yn3tTYsO_zEAJ8hCgg7VW11cE/preview

This workshop advances as we switch to other branches with the 6th branch having the complete code and documentation. 
You can directly switch to the 6th branch if you wish to test the application.

**1. Building the basic CLI Application using `cli` package in Go**

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
**Code breakdown** 

cli.App creates the command line application with details about the Application like its
Name: ytt
Usage: Transcribe YouTube videos
Commands: An array of cli.Command structs that represent the commands the application accepts. We first accept only "get". 

**Executing the code**

1. Go ahead and run the main.go file: `go run main.go`
   You get all info about your cli app. 
2. Run `go run main.go get "https://www.youtube.com/watch?v=ltmInTalwXQ"`


**2. Downloading YouTube Video and Audio Streams**

Now, we'll use the `youtube/v2` package to fetch video details and download the audio stream for the supplied YouTube link. 

**Task**: In a new git branch named "downloading-video", append code to your program from Step 1 for downloading a video using a YouTube link. Add a "download" command to ur CLI for testing this functionality.

We will provide only the ID of the yt link and not the complete link.

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

The `context` package might seem like it doesn't serve any purpose, but it's actually a placeholder for future enhancements. Suppose you decide to add a feature that allows users to cancel downloads or you want to set a maximum time limit for downloads. You can modify the code to pass a cancelable context or a context with a timeout instead of `context.Background()`

```go
//main.go
err := YoutubeDL(youtubelink)
if err != nil {
    return err
}
```
