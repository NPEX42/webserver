package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Project struct {
	Name        string   `json:"name"`
	Description string   `json:"desc"`
	Languages   []string `json:"languages"`
	Repo        string   `json:"repo"`
}

const PROJECT_CARD = `
	<div>
		<p class="langTag">$LANGS</p>
		<h3>$NAME</h3>
		<p>$DESC</p>
	<a href="$REPO" target="_blank"><img src="https://github.com/favicon.ico"></img></a>
	</div>
`

func LoadProjects(filepath string) ([]Project, error) {
	var projects []Project
	data, err := os.ReadFile(filepath)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = json.Unmarshal(data, &projects)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return projects, nil
}

func (proj *Project) Render() string {
	output := PROJECT_CARD
	output = strings.ReplaceAll(output, "$NAME", proj.Name)
	output = strings.ReplaceAll(output, "$DESC", proj.Description)
	output = strings.ReplaceAll(output, "$REPO", proj.Repo)
	langs := RenderLangs(proj.Languages)
	output = strings.ReplaceAll(output, "$LANGS", langs)
	return output
}

func RenderLangs(langs []string) string {
	var s bytes.Buffer

	for _, lang := range langs {
		s.WriteString(fmt.Sprintf("%s ", lang))
	}

	return s.String()
}
