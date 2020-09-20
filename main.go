package main

import (
	"fmt"
	"chart-viewer/handler"
	"chart-viewer/helm"
	"chart-viewer/repository"
	"chart-viewer/service"
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

	repo := repository.NewRepository()
	helmClient := helm.NewHelmClient(repo)
	svc := service.NewService(helmClient, repo)
	appHandler := handler.NewHandler(svc)

	r.Use(appHandler.CORS)
	r.Use(appHandler.LoggerMiddleware)
	r.HandleFunc("/repos", appHandler.GetReposHandler)
	r.HandleFunc("/charts/{repo-name}", appHandler.GetChartsHandler)
	r.HandleFunc("/charts/{repo-name}/{chart-name}/{chart-version}", appHandler.GetChartHandler)
	r.HandleFunc("/charts/values/{repo-name}/{chart-name}/{chart-version}", appHandler.GetValuesHandler)
	r.HandleFunc("/charts/templates/{repo-name}/{chart-name}/{chart-version}", appHandler.GetTemplatesHandler)
	r.HandleFunc("/charts/manifests/render/{repo-name}/{chart-name}/{chart-version}", appHandler.RenderManifestsHandler)
	r.HandleFunc("/charts/manifests/{repo-name}/{chart-name}/{chart-version}/{hash}", appHandler.GetManifestsHandler)


	seedDatabase()

	host := os.Getenv("HT_HOST")
	port := os.Getenv("HT_PORT")

	address := fmt.Sprintf("%s:%s", host, port)

	log.Printf("server run on http://%s\n", address)
	log.Fatal(http.ListenAndServe(address, r))
}
