package service

import (
	"bytes"
	"chart-viewer/helm"
	"chart-viewer/model"
	"chart-viewer/repository"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"net/http"
)

type Service interface {
	GetRepos() []model.Repo
	GetCharts(repoName string) (error, []model.Chart)
	GetValues(repoName, chartName, chartVersion string) (error, map[string]interface{})
	GetTemplates(repoName, chartName, chartVersion string) []model.Template
	RenderManifest(repoName, chartName, chartVersion string, values []string) (error, model.ManifestResponse)
	GetStringifiedManifests(repoName, chartName, chartVersion, hash string) string
	GetChart(repoName string, chartName string, chartVersion string) (error, model.ChartDetail)
}

type service struct {
	helmClient helm.Helm
	repository repository.Repository
}

func NewService(helmClient helm.Helm, repository repository.Repository) Service {
	return &service{
		helmClient: helmClient,
		repository: repository,
	}
}

func (s *service) GetRepos() []model.Repo {
	stringifiedRepos := s.repository.Get("repos")
	var repos []model.Repo
	_ = json.Unmarshal([]byte(stringifiedRepos), &repos)

	return repos
}

func (s *service) GetCharts(repoName string) (error, []model.Chart) {
	stringifiedCharts := s.repository.Get(repoName)
	var cachedCharts []model.Chart
	_ = json.Unmarshal([]byte(stringifiedCharts), &cachedCharts)

	if len(cachedCharts) != 0 {
		log.Printf("%s chart detail fetched from cache\n", repoName)
		return nil, cachedCharts
	}

	url := s.getUrl(repoName)

	log.Printf("out going call: %s\n", url)
	response, err := http.Get(url + "/index.yaml")
	if err != nil {
		return err, nil
	}

	content, err := ioutil.ReadAll(response.Body)

	repoDetail := new(model.RepoDetailResponse)
	err = yaml.Unmarshal(content, &repoDetail)
	if err != nil {
		return err, nil
	}

	var chartNames []string

	for name, _ := range repoDetail.Entries {
		chartNames = append(chartNames, name)
	}

	var charts []model.Chart

	for _, name := range chartNames {
		charts = append(charts, model.Chart{
			Name:     name,
			Versions: getVersion(name, repoDetail.Entries),
		})
	}

	chartsByte, _ := json.Marshal(charts)
	s.repository.Set(repoName, string(chartsByte))

	return nil, charts
}

func (s *service) GetValues(repoName, chartName, chartVersion string) (error, map[string]interface{}) {
	cacheKey := fmt.Sprintf("value-%s-%s-%s", repoName, chartName, chartVersion)
	stringifiedValues := s.repository.Get(cacheKey)
	var cachedValues map[string]interface{}
	_ = json.Unmarshal([]byte(stringifiedValues), &cachedValues)
	if len(cachedValues) != 0 {
		log.Printf("value-%s-%s-%s chart values fetched from cache\n", repoName, chartName, chartVersion)
		return nil, cachedValues
	}

	var url string
	repos := s.GetRepos()

	for _, r := range repos {
		if r.Name == repoName {
			url = r.URL
		}
	}

	err, values := s.helmClient.GetValues(url, chartName, chartVersion)
	if err != nil {
		return err, nil
	}
	valuesByte, _ := json.Marshal(values)
	s.repository.Set(cacheKey, string(valuesByte))

	return nil, values
}

func (s *service) GetTemplates(repoName, chartName, chartVersion string) []model.Template {
	cacheKey := fmt.Sprintf("template-%s-%s-%s", repoName, chartName, chartVersion)
	stringifiedTemplates := s.repository.Get(cacheKey)
	var cachedTemplates []model.Template
	_ = json.Unmarshal([]byte(stringifiedTemplates), &cachedTemplates)
	if len(cachedTemplates) != 0 {
		log.Printf("template-%s-%s-%s chart values fetched from cache\n", repoName, chartName, chartVersion)
		return cachedTemplates
	}

	var url string
	repos := s.GetRepos()

	for _, r := range repos {
		if r.Name == repoName {
			url = r.URL
		}
	}

	templates := s.helmClient.GetManifest(url, chartName, chartVersion)
	templatesByte, _ := json.Marshal(templates)
	s.repository.Set(cacheKey, string(templatesByte))

	return templates
}

func (s *service) RenderManifest(repoName, chartName, chartVersion string, values []string) (error, model.ManifestResponse) {
	hash := hashFileContent(values[0])
	cacheKey := fmt.Sprintf("manifests-%s-%s-%s-%s", repoName, chartName, chartVersion, hash)
	stringifiedManifest := s.repository.Get(cacheKey)
	var cachedManifests model.ManifestResponse
	_ = json.Unmarshal([]byte(stringifiedManifest), &cachedManifests)

	generatedUrl := fmt.Sprintf("/charts/manifests/%s/%s/%s/%s", repoName, chartName, chartVersion, hash)
	if stringifiedManifest != "" {
		log.Printf("manifest fetched from cache with key: %s\n", cacheKey)
		return nil, cachedManifests
	}

	var url string
	repos := s.GetRepos()

	for _, r := range repos {
		if r.Name == repoName {
			url = r.URL
		}
	}

	err, manifests := s.helmClient.RenderManifest(url, chartName, chartVersion, values)
	if err != nil {
		return err, model.ManifestResponse{}
	}

	manifestsResponse := model.ManifestResponse{
		URL:       generatedUrl,
		Manifests: manifests,
	}

	manifestsByte, _ := json.Marshal(manifestsResponse)
	s.repository.Set(cacheKey, string(manifestsByte))

	return nil, manifestsResponse
}

func (s *service) GetStringifiedManifests(repoName, chartName, chartVersion, hash string) string {
	cacheKey := fmt.Sprintf("manifests-%s-%s-%s-%s", repoName, chartName, chartVersion, hash)
	var cachedManifests model.ManifestResponse
	stringifiedManifest := s.repository.Get(cacheKey)

	_ = json.Unmarshal([]byte(stringifiedManifest), &cachedManifests)
	return stringfyManifest(cachedManifests.Manifests)
}

func (s *service) GetChart(repoName string, chartName string, chartVersion string) (error, model.ChartDetail) {
	err, values := s.GetValues(repoName, chartName, chartVersion)
	if err != nil {
		return err, model.ChartDetail{}
	}
	templates := s.GetTemplates(repoName, chartName, chartVersion)

	return nil, model.ChartDetail{
		Values:    values,
		Templates: templates,
	}
}

func hashFileContent(fileLocation string) string {
	valuesFileContent, _ := ioutil.ReadFile(fileLocation)
	hash := md5.Sum(valuesFileContent)
	return fmt.Sprintf("%x", hash)
}

func getVersion(name string, entries map[string][]model.ChartResponse) []string {
	cs := entries[name]

	var versions []string

	for _, c := range cs {
		versions = append(versions, c.Version)
	}

	return versions
}

func stringfyManifest(manifests []model.Manifest) string {
	var buffer bytes.Buffer
	var delimiter = "---\n"
	for _, m := range manifests {
		buffer.WriteString(delimiter + m.Content + "\n")
	}

	return buffer.String()
}

func (s *service) getUrl(repoName string) string {
	repos := s.GetRepos()
	for _, r := range repos {
		if r.Name == repoName {
			return r.URL
		}
	}

	return ""
}
