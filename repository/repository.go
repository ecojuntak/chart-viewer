package repository

import (
	"encoding/json"
	"chart-viewer/model"
	"os"

	"github.com/go-redis/redis"
)

type Repository struct {
	redisClient *redis.Client
}

func NewRepository() *Repository {
	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	return &Repository{
		redisClient: redis.NewClient(&redis.Options{
			Addr: host + ":" + port,
		}),
	}
}

func (r *Repository) GetRepos() []model.Repo {
	reposString, _ := r.redisClient.Get("repos").Result()

	var repos []model.Repo
	_ = json.Unmarshal([]byte(reposString), &repos)

	return repos
}

func (r *Repository) AddRepo(newRepo model.Repo) {
	repos := r.GetRepos()
	repos = append(repos, newRepo)

	jsonRepo, _ := json.Marshal(repos)

	_ = r.redisClient.Set("repos", jsonRepo, 0)
}

func (r *Repository) GetRepoDetail(repoName string) []model.Chart {
	repoDetail, _ := r.redisClient.Get(repoName).Result()

	var charts []model.Chart
	_ = json.Unmarshal([]byte(repoDetail), &charts)

	return charts
}

func (r *Repository) GetManifest(cacheKey string) string {
	manifests, _ := r.redisClient.Get(cacheKey).Result()

	return manifests
}

func (r *Repository) StoreRepoDetail(repoName string, charts []model.Chart) {
	jsonRepoDetail, _ := json.Marshal(charts)
	_ = r.redisClient.Set(repoName, jsonRepoDetail, 0)
}

func (r *Repository) GetChartValues(key string) map[string]interface{} {
	chartValues, _ := r.redisClient.Get(key).Result()

	var values map[string]interface{}
	_ = json.Unmarshal([]byte(chartValues), &values)

	return values
}

func (r *Repository) StoreChartValues(key string, values map[string]interface{}) {
	chartValues, _ := json.Marshal(values)
	_ = r.redisClient.Set(key, chartValues, 0)
}

func (r *Repository) GetChartTemplate(key string) []model.Template {
	chartTemplates, _ := r.redisClient.Get(key).Result()

	var templates []model.Template
	_ = json.Unmarshal([]byte(chartTemplates), &templates)

	return templates
}

func (r *Repository) StoreChartTemplate(key string, templates []model.Template) {
	chartTemplates, _ := json.Marshal(templates)
	_ = r.redisClient.Set(key, chartTemplates, 0)
}

func (r *Repository) StoreChartManifest(key, chartManifests string) {
	_ = r.redisClient.Set(key, chartManifests, 0)
}
