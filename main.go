//Install go get github.com/urfave/cli/v2

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
					audieofilename, err := YoutubeDL(youtubelink)

					if err != nil {
						return err
					}
					fmt.Println("Audio file:", audieofilename)
					prefix := "/data/"
					audieofilename = prefix + audieofilename
					err = ConvertFile(audieofilename+".mp4", audieofilename+".wav")
					if err != nil {
						return err
					}
					modelfile := prefix + "ggml-tiny.en.bin"
					context, err := transcribe(modelfile, audieofilename+".wav")
					if err != nil {
						return err
					}
					OutputSRT(context)
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
