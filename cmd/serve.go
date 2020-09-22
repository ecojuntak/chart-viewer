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
	command := cobra.Command{
		Use:     "serve",
		Short:   "Start the http server",
		Example: "chart-viewer serve",
		Run: func(cmd *cobra.Command, args []string) {
			r := createRouter()

			host := os.Getenv("HT_HOST")
			port := os.Getenv("HT_PORT")

			address := fmt.Sprintf("%s:%s", host, port)

			log.Printf("server run on http://%s\n", address)
			log.Fatal(http.ListenAndServe(address, r))
		},
	}

	return &command
}

func createRouter() *mux.Router {
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

	return r
}
