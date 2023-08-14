//Install go get github.com/urfave/cli/v2

package main

import (
	"Speech-To-Text/models"
	"Speech-To-Text/repository"
	"fmt"
	"os"

	"github.com/urfave/cli/v2" //this imports the cli package
)

const port = 8090

const dbUrl = "ws://192.168.1.33:8000/rpc"
const namespace = "surrealdb-conference-content"
const database = "yttranscriber"

func main() {
	app := &cli.App{
		Name:  "ytt",
		Usage: "Transcribe YouTube videos",
		Commands: []*cli.Command{
			{
				Name:  "get",
				Usage: "Get transcriptions by ytlink",
				Action: func(c *cli.Context) error {
					repository, err := repository.NewTranscriptionsRepository(dbUrl, "root", "root", namespace, database)
					if err != nil {
						return err
					}
					//Print YouTube link
					youtubelink := c.Args().Get(0)
					if youtubelink == "" {
						return cli.NewExitError("Please provide a YouTube link", 1)
					}

					transcriptions, err := repository.GetTranscriptionsByYtlink(youtubelink)
					if err != nil {
						return err
					}
					//Check if transcriptions is empty
					if len(transcriptions.([]interface{})[0].(map[string]interface{})["result"].([]interface{})) == 0 {
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
						rawTranscription, err := transcribe(modelfile, audieofilename+".wav")
						if err != nil {
							return err
						}
						_, err = repository.SaveTranscriptions(youtubelink, models.ToModel(models.RawTranscriptions(rawTranscription)))
						if err != nil {
							return err
						}
					} else {
						fmt.Println("YouTube link already exists in database")
						fmt.Println(transcriptions)
					}

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
