package main

import (
	"fmt"
	"ht-ui/handler"
	"ht-ui/helm"
	"ht-ui/repository"
	"ht-ui/service"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

func seedDatabase() {
	repos := os.Getenv("CHART_REPOS")

	host := os.Getenv("REDIS_HOST")
	port := os.Getenv("REDIS_PORT")

	redisClient := redis.NewClient(&redis.Options{
		Addr: host + ":" + port,
	})

	_ = redisClient.Set("repos", repos, 0)
}

func main() {
	r := mux.NewRouter()

	repository := repository.NewRepository()
	helmClient := helm.NewHelmClient(repository)
	service := service.NewService(helmClient, repository)
	handler := handler.NewHandler(service)

	r.Use(handler.CORS)
	r.Use(handler.LoggerMiddleware)
	r.HandleFunc("/repos", handler.RepoListHandler).Methods("get", "post")
	r.HandleFunc("/repos/{repo-name}", handler.RepoDetailHandler)
	r.HandleFunc("/repos/{repo-name}/{chart-name}/{chart-version}", handler.ChartHandler)
	r.HandleFunc("/repos/{repo-name}/{chart-name}/{chart-version}/render", handler.RenderHandler)
	r.HandleFunc("/manifest/{repo-name}/{chart-name}/{chart-version}/{hash}", handler.DownloadManifestHandler)

	seedDatabase()

	host := os.Getenv("HT_HOST")
	port := os.Getenv("HT_PORT")

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("server run on http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, r))
}
