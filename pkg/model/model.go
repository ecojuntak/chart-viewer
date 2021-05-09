package model

import (
	"fmt"
	"strings"
)

type Repo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func (r Repo) MarshalBinary() ([]byte, error) {
	s := fmt.Sprintf("%s:%s", r.Name, r.URL)
	return []byte(s), nil
}

func (r Repo) GetURL() string {
	if strings.HasSuffix(r.URL, "/") {
		return r.URL[:len(r.URL)-1]
	}

	return r.URL
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

type KubernetesAPIVersion struct {
	KubeVersion string   `json:"kube_version"`
	APIVersions []string `json:"api_versions"`
}

type KubeResourceCommonSpec struct {
	APIVersion string `yaml:"apiVersion"`
}

type AnalyticsResult struct {
	Template
	Compatible bool `json:"compatible"`
}

type AnalyticResponse struct {
	Values    map[string]interface{} `json:"values"`
	Templates []AnalyticsResult      `json:"templates"`
}
