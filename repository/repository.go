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



