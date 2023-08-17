# Speech-To-Text-guide



**Workshop Guide Outline**

Link to presentation slides: https://docs.google.com/presentation/d/1e5PEJu6yn3tTYsO_zEAJ8hCgg7VW11cE/preview
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
**4. Building Docker Images and running Docker Containers**

This part of the workshop dives into Docker, a platform that uses OS-level virtualization to deliver software in isolated packages known as containers. Docker can simplify the setup process for application development and distribution.

**Understanding the Dockerfile**

```Dockerfile
# Use the official Golang base image
FROM golang:1.19

# Set the working directory inside the container
WORKDIR /app

# Install whisper
RUN git clone https://github.com/ggerganov/whisper.cpp.git &&\
    cd whisper.cpp && make &&\
    make libwhisper.so libwhisper.a &&\
    cp whisper.h /usr/local/include &&\
    cp ggml.h /usr/local/include &&\
    cp libwhisper.a /usr/local/lib &&\
    cp libwhisper.so /usr/local/lib &&\
    cd ..

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download all necessary dependencies
RUN go mod download

# Copy the source code to the working directory
COPY . .

# Build the Go app
RUN go build -o ytt

# Install ffmpeg
RUN apt-get update && apt-get install -y ffmpeg

# execute cli help to check if everything is ok
RUN ./ytt -h


# Run the compiled binary with a default command
ENTRYPOINT ["/app/ytt"]
```
    
We will go through the Dockerfile and explain each directive:

- `FROM` specifies the base image. Here, we are using the official Golang image.
- `WORKDIR` sets the current working directory inside the Docker container.
- `RUN` is used to execute commands. In this Dockerfile, we use it to download and compile the whisper library, download and install FFmpeg, and build our Go application.
- `COPY` copies new files or directories from "<source>" and adds them to the filesystem of the image at the path "<destination>".
- `ENTRYPOINT` sets the command and parameters that will be executed first when a container is run.

**Task: Building the Docker Image**

> **Note:** If you don't want to build everything from scratch and just want to see the application running, you can pull the Docker image directly from the GitHub registry. Use the following command:
> 
>  ```bash
>  docker run --platform linux/amd64 -v "$(pwd)":/data -it --rm ghcr.io/timpratim/speech-to-text-guide:4-dockerfile get "JzPfMbG1vrE" 
>  ```
    
We will use Docker build command to create a Docker image from the Dockerfile. Here's the command we'll use:

```bash
docker buildx build --platform linux/amd64 -t ytt-amd64 --load -f Dockerfile .
```

This command tells Docker to build an image using the Dockerfile in the current directory (the "." at the end). The flag `--platform linux/amd64` specifies the platform the image is being built for. The `-t ytt-amd64` flag tags the image with the name "ytt-amd64". The `--load` flag tells Docker to load the built image into Docker's locally accessible image store.
    
Run this command so that Docker can build your image.

**Task: Running the Docker Container**

After successfully building the Docker image, we can now create and run a Docker container from it.

```bash
docker run --platform linux/amd64 -v "$(pwd)":/data -it --rm ytt-amd64 get "JzPfMbG1vrE"
```

The `docker run` command creates and starts a Docker container. The `--platform linux/amd64` flag specifies the platform of the container. The `-v "$(pwd)":/data` flag mounts the current directory from the host into the container at "/data". The `-it` flag ensures that we can interact with the container via the terminal. `--rm` tells Docker to automatically clean up the container and remove the file system when the container exits. The `get "JzPfMbG1vrE"` at the end of the command is the command arguments that will be passed to the ENTRYPOINT command inside the Kubernetes container.

Run this command to interact with your image. Check to see that it works as expected, given its output, and familiarize yourself with the process of building and running Docker containers. This should conclude our workshop.

**5. Transcribing a Video**

This section of the workshop focuses on transcribing the converted audio file. We will use the whisper package to do this.

**Task**: In a new git branch "video-transcription", use the whisper package in your current program to perform the transcription. Add a "transcribe" command to your CLI for testing this step.

**Code Snippets**
```go
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
```

```diff
//main.go
modelfile := prefix + "ggml-tiny.en.bin"
err = transcribe(modelfile, audieofilename+".wav")
if err != nil {
    return err
}
```
**6. Managing the Transcription Data in a Database**

The final section will explain how to store and retrieve transcriptions in a database using the repository pattern.

**Task**: In the git branch named "db-layer", extend your current program with the database operations. Add the necessary commands like "save" and "get" to the CLI for testing.

**Code Snippets**
    
```go
//repository/repository.go
package repository

import (
	"fmt"

	"Speech-To-Text/models"

	surreal "github.com/surrealdb/surrealdb.go"
)

type TranscriptionsRepository struct {
	db *surreal.DB
}

func NewTranscriptionsRepository(address, user, password, namespace, database string) (*TranscriptionsRepository, error) {
	//surreal.New is
	db, err := surreal.New(address)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %s", err)
	}
	_, err = db.Signin(map[string]interface{}{
		"user": user,
		"pass": password,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to sign in: %w", err)
	}

	_, err = db.Use(namespace, database)
	if err != nil {
		return nil, err
	}

	return &TranscriptionsRepository{db}, nil
}

func (r TranscriptionsRepository) Close() {
	r.db.Close()
}

func (r TranscriptionsRepository) SaveTranscriptions(ytlink string, transcriptions models.Transcriptions) (interface{}, error) {
	fmt.Println("saving transcriptions %v", transcriptions)
	return r.db.Create("transcriptions", map[string]interface{}{
		"ytlink":         ytlink,
		"transcriptions": transcriptions.ToMap(),
	})
}

// get transcriptions by ytlink
func (r TranscriptionsRepository) GerTranscriptionsByYtlink(ytlink string) (interface{}, error) {

	return r.db.Query("SELECT * FROM transcriptions WHERE ytlink = $ytlink limit 1", map[string]interface{}{
		"ytlink": ytlink,
	})
}

```

**Task: Deploy a Database using Docker**
    
We'll now deploy our database using Docker, which will give each participant their own isolated database instance. We'll be using the SurrealDB container.

Make sure the IP address in main.go is same as your IP address of your laptop
    
The Docker command to set this up is:
    
```bash
    docker run --name surrealdb --network surrealdb-network --restart always -v "$(pwd)/data:/data" -p 8000:8000 surrealdb/surrealdb:latest start --user root --pass root
```

This command deploys a Docker container based on the surrealdb image. 
    
- `docker run` creates and starts a Docker container.
    
- `--name surrealdb` names the container "surrealdb".

- `--network surrealdb-network` specifies that this container is part of the network "surrealdb-network".

- `--restart always` will automatically restart the container if it stops.

- `-v "$(pwd)/data:/data"` mounts the current directory's "data" folder from the host into the container at "/data".
    
- `-p 8000:8000` maps the Docker host's port 8000 to the container's port 8000, which SurrealDB listens on.

- `start --user root --pass root` is the command we want the container to run upon start. We're starting the SurrealDB with the user "root" and password "root".
    
Please run this command to deploy your database. After successful deployment, your application should be able to connect and interact with the database.
    
## Using curl to read the database

```bash
curl -X POST \
         -u "root:root" \
         -H "NS: surrealdb-conference-content" \
         -H "DB: yttranscriber" \
         -H "Accept: application/json" \
         -d "SELECT * FROM transcriptions" \
         http://192.168.1.22:8000/sql
```
