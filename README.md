# Speech-To-Text-transcriptions using Whisper.cpp-guide

**Workshop Guide Outline**

The [presentation slides](https://sdb.li/STTpresentation) are open to public.

This workshop advances as we switch to other branches with the 6th branch having the complete code and documentation. 
You can directly switch to the 6th branch if you wish to test the application.

> **Note for Windows users**
> Refer to the following blogs to download the necessary packages on Windows. 
> - [Docker](https://medium.com/devops-with-valentine/how-to-install-docker-on-windows-10-11-step-by-step-83074a80e6f9)
> - [git](https://zaycodes.medium.com/how-to-install-git-on-windows-f6031afef08c)
> - [ffmpeg](https://www.geeksforgeeks.org/how-to-install-ffmpeg-on-windows/)



**1. Building the basic CLI Application using `cli` package in Go**

In this part, we'll cover how to set up a Command Line Interface (CLI) application using the "cli" package in Go. We will start with a simple program that accepts the "get" command for fetching transcriptions using a supplied YouTube link. We'll just print back the youtube link for now.

**Task**: Create a new Go file in a new git branch named "cli-setup". Now, write a basic program that defines a CLI application which accepts the "get" command.

**Depedencies**
`go get github.com/urfave/cli/v2`


**Code Snippets**

```go
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
}```
**Code breakdown** 

cli.App creates the command line application with details about the Application like its
Name: ytt
Usage: Transcribe YouTube videos
Commands: An array of cli.Command structs that represent the commands the application accepts. We first accept only "get". 

**Executing the code**

1. Go ahead and run the main.go file: `go run main.go`
   You get all info about your cli app. 
2. Run `go run main.go get "https://www.youtube.com/watch?v=ltmInTalwXQ"`

