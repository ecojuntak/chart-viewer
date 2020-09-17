package service

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"ht-ui/helm"
	"ht-ui/model"
	"ht-ui/repository"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type Service struct {
	helmClient *helm.HelmClient
	repository *repository.Repository
}

func NewService(helmClient *helm.HelmClient, repository *repository.Repository) *Service {
	return &Service{
		helmClient: helmClient,
		repository: repository,
	}
}

func (s *Service) GetRepoList() []model.Repo {
	repos := s.repository.GetRepos()
	return repos
}

func (s *Service) GetRepoDetail(repoName string) []model.Chart {

	cs := s.fetchFromCache(repoName)
	if len(cs) != 0 {
		log.Printf("%s chart detail fetched from cache\n", repoName)
		return cs
	}

	url := s.GetUrl(repoName)

	log.Printf("out going call: %s\n", url)
	response, err := http.Get(url + "/index.yaml")
	if err != nil {
		panic(err)
	}

	content, err := ioutil.ReadAll(response.Body)

	repoDetail := new(model.RepoDetailResponse)
	err = yaml.Unmarshal(content, &repoDetail)
	if err != nil {
		panic(err)
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

	s.repository.StoreRepoDetail(repoName, charts)

	return charts
}

func (s *Service) fetchFromCache(repoName string) []model.Chart {
	charts := s.repository.GetRepoDetail(repoName)

	return charts
}

func getVersion(name string, entries map[string][]model.ChartResponse) []string {
	cs := entries[name]

	var versions []string

	for _, c := range cs {
		versions = append(versions, c.Version)
	}

	return versions
}

func (s *Service) GetUrl(repoName string) string {
	repos := s.GetRepoList()
	for _, r := range repos {
		if r.Name == repoName {
			return r.URL
		}
	}

	return ""
}

func (s *Service) GetChartValues(repoName, chartName, chartVersion string) map[string]interface{} {

	cacheKey := fmt.Sprintf("value-%s-%s-%s", repoName, chartName, chartVersion)
	chartValues := s.repository.GetChartValues(cacheKey)
	if len(chartValues) != 0 {
		log.Printf("value-%s-%s-%s chart values fetched from cache\n", repoName, chartName, chartVersion)
		return chartValues
	}

	var url string
	repos := s.GetRepoList()

	for _, r := range repos {
		if r.Name == repoName {
			url = r.URL
		}
	}

	values := s.helmClient.GetValues(url, chartName, chartVersion)

	s.repository.StoreChartValues(cacheKey, values)

	return values
}

func (s *Service) GetChartTemplates(repoName, chartName, chartVersion string) []model.Template {

	cacheKey := fmt.Sprintf("template-%s-%s-%s", repoName, chartName, chartVersion)
	chartTemplates := s.repository.GetChartTemplate(cacheKey)
	if len(chartTemplates) != 0 {
		log.Printf("template-%s-%s-%s chart values fetched from cache\n", repoName, chartName, chartVersion)
		return chartTemplates
	}

	var url string
	repos := s.GetRepoList()

	for _, r := range repos {
		if r.Name == repoName {
			url = r.URL
		}
	}

	templates := s.helmClient.GetManifest(url, chartName, chartVersion)

	s.repository.StoreChartTemplate(cacheKey, templates)

	return templates
}

func (s *Service) GenerateManifest(repoName, chartName, chartVersion string, values []string) (string, []model.ManifestObject) {
	var url string
	repos := s.GetRepoList()

	for _, r := range repos {
		if r.Name == repoName {
			url = r.URL
		}
	}

	rawManifest, manifests := s.helmClient.RenderManifest(url, chartName, chartVersion, values)

	manifest, err := json.Marshal(manifests)
	if err != nil {
		panic(err)
	}

	hash := md5.Sum(manifest)
	path := fmt.Sprintf("%s/%s/%s/%x", repoName, chartName, chartVersion, hash)

	s.repository.StoreChartManifest(path, rawManifest.String())

	return path, manifests
}

func (s *Service) DownloadManifest(repoName, chartName, chartVersion, hash string) string {
	cacheKey := fmt.Sprintf("%s/%s/%s/%s", repoName, chartName, chartVersion, hash)
	manifests := s.repository.GetManifest(cacheKey)

	return manifests
}

func (s *Service) AddRepo(newRepo model.Repo) {
	s.repository.AddRepo(newRepo)
}
