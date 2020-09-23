package cmd

import (
	"chart-viewer/handler"
	"chart-viewer/helm"
	"chart-viewer/repository"
	"chart-viewer/service"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	defaultHost := os.Getenv("APP_HOST")
	defaultPort := os.Getenv("APP_PORT")
	defaultRedisHost := os.Getenv("REDIS_HOST")
	defaultRedisPort := os.Getenv("REDIS_PORT")

	command := cobra.Command{
		Use:     "serve",
		Short:   "Start the http server",
		Example: "chart-viewer serve --host 127.0.0.1 --port 9999 --redis-host 127.0.0.1 --redis-port 6379",
		Run: func(cmd *cobra.Command, args []string) {
			host := defaultHost
			port := defaultPort
			redisHost := defaultRedisHost
			redisPort := defaultRedisPort

			redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)
			address := fmt.Sprintf("%s:%s", host, port)

			err, repo := repository.NewRepository(redisAddress)
			if err != nil {
				fmt.Printf("cannot connect to redis: %s\n", err)
				return
			}

			helmClient := helm.NewHelmClient(repo)
			svc := service.NewService(helmClient, repo)
			r := createRouter(svc)

			log.Printf("server run on http://%s\n", address)
			log.Fatal(http.ListenAndServe(address, r))
		},
	}

	command.Flags().StringVar(&defaultHost, "host", "127.0.0.1", "[Optional] App host address")
	command.Flags().StringVar(&defaultPort, "port", "9999", "[Optional] App host port")
	command.Flags().StringVar(&defaultRedisHost, "redis-host", "127.0.0.1", "[Optional] Redis host address")
	command.Flags().StringVar(&defaultRedisPort, "redis-port", "6379", "[Optional] Redis host port")

	return &command
}

func createRouter(svc service.Service) *mux.Router {
	r := mux.NewRouter()

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

	return r
}
