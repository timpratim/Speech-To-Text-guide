# Speech-To-Text-transcriptions using Whisper.cpp-guide

**Workshop Guide Outline**

Link to presentation slides: https://docs.google.com/presentation/d/1e5PEJu6yn3tTYsO_zEAJ8hCgg7VW11cE/preview

This workshop advances as we switch to other branches with the 6th branch having the complete code and documentation. 
You can directly switch to the 6th branch if you wish to test the application.

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
