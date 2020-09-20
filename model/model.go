package model

import "fmt"

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (r Repo) MarshalBinary() ([]byte, error) {
	s := fmt.Sprintf("%s:%s", r.Name, r.URL)
	return []byte(s), nil
}

type Chart struct {
	Name     string   `json:"name"`
	Versions []string `json:"versions"`
}

type Template struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ChartDetail struct {
	Values    map[string]interface{} `json:"values"`
	Templates []Template             `json:"templates"`
}

type RepoDetailResponse struct {
	ApiVersion string                     `yaml:"apiVersion"`
	Entries    map[string][]ChartResponse `yaml:"entries"`
}

type ChartResponse struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	URLs    []string `yaml:"urls"`
}

type Manifest struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ManifestResponse struct {
	URL       string     `json:"url"`
	Manifests []Manifest `json:"manifests"`
}
