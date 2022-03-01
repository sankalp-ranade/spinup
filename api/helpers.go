package api

import (
	"fmt"
	"github.com/spinup-host/config"
	"math/rand"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spinup-host/templates"
)

//TODO: vicky find how to keep the templates/* outside of the api. ie need to figure how to do relative path.
// check: https://stackoverflow.com/questions/66285635/how-do-you-use-go-1-16-embed-features-in-subfolders-packages

type malformedRequest struct {
	status int
	msg    string
}

func (mr *malformedRequest) Error() string {
	return mr.msg
}

// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// TODO: To remove the duplication here. We don't need separate function for each file
func CreateDockerComposeFile(absolutepath string, s service) error {
	outputPath := filepath.Join(absolutepath, "docker-compose.yml")
	// Create the file:
	f, err := os.Create(outputPath)
	if err != nil {
		panic(err)
	}

	defer f.Close() // don't forget to close the file when finished.
	templ, err := template.ParseFS(templates.DockerTempl, "templates/docker-compose-template.yml")
	if err != nil {
		return fmt.Errorf("ERROR: parsing template file %v", err)
	}
	// TODO: not sure is there a better way to pass data to template
	// A lot of this data is redundant. Already available in Service struct
	data := struct {
		UserID           string
		Architecture     string
		Type             string
		Port             int
		MajVersion       uint
		MinVersion       uint
		PostgresUsername string
		PostgresPassword string
		NetworkName      string
	}{
		s.UserID,
		s.Architecture,
		s.Db.Type,
		s.Db.Port,
		s.Version.Maj,
		s.Version.Min,
		s.Db.Username,
		s.Db.Password,
		config.DefaultNetwork,
	}
	err = templ.Execute(f, data)
	if err != nil {
		return fmt.Errorf("ERROR: executing template file %v", err)
	}
	return nil
}
