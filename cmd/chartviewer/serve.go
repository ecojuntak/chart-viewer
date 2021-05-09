package chartviewer

import (
	"chart-viewer/pkg/analyzer"
	"chart-viewer/pkg/helm"
	"chart-viewer/pkg/repository"
	"chart-viewer/pkg/server/handler"
	"chart-viewer/pkg/server/service"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
)

func NewServeCommand() *cobra.Command {
	var (
		defaultHost      string
		defaultPort      string
		defaultRedisHost string
		defaultRedisPort string
	)

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

			repo, err := repository.NewRepository(redisAddress)
			if err != nil {
				fmt.Printf("cannot connect to redis: %s\n", err)
				return
			}

			helmClient := helm.NewHelmClient(repo)
			analyser := analyzer.New()
			svc := service.NewService(helmClient, repo, analyser)
			r := createRouter(svc)

			log.Printf("server run on http://%s\n", address)
			log.Fatal(http.ListenAndServe(address, r))
		},
	}

	command.Flags().StringVar(&defaultHost, "host", "0.0.0.0", "[Optional] App host address")
	command.Flags().StringVar(&defaultPort, "port", "9999", "[Optional] App host port")
	command.Flags().StringVar(&defaultRedisHost, "redis-host", "127.0.0.1", "[Optional] Redis host address")
	command.Flags().StringVar(&defaultRedisPort, "redis-port", "6379", "[Optional] Redis host port")

	return &command
}

func createRouter(svc service.Service) *mux.Router {
	r := mux.NewRouter()

	appHandler := handler.NewHandler(svc)

	r.Use(appHandler.CORS)
	apiV1 := r.PathPrefix("/api/v1/").Subrouter()
	apiV1.Use(appHandler.LoggerMiddleware)
	apiV1.HandleFunc("/repos", appHandler.GetReposHandler).Methods("GET")
	apiV1.HandleFunc("/charts/{repo-name}", appHandler.GetChartsHandler).Methods("GET")
	apiV1.HandleFunc("/charts/{repo-name}/{chart-name}/{chart-version}", appHandler.GetChartHandler).Methods("GET")
	apiV1.HandleFunc("/charts/values/{repo-name}/{chart-name}/{chart-version}", appHandler.GetValuesHandler).Methods("GET")
	apiV1.HandleFunc("/charts/templates/{repo-name}/{chart-name}/{chart-version}", appHandler.GetTemplatesHandler).Methods("GET")
	apiV1.HandleFunc("/charts/manifests/render/{repo-name}/{chart-name}/{chart-version}", appHandler.RenderManifestsHandler).Methods("POST", "OPTIONS")
	apiV1.HandleFunc("/charts/manifests/{repo-name}/{chart-name}/{chart-version}/{hash}", appHandler.GetManifestsHandler).Methods("GET")

	fileServer := http.FileServer(http.Dir("ui/dist"))
	r.PathPrefix("/js").Handler(http.StripPrefix("/", fileServer))
	r.PathPrefix("/css").Handler(http.StripPrefix("/", fileServer))
	r.PathPrefix("/img").Handler(http.StripPrefix("/", fileServer))
	r.PathPrefix("/favicon.ico").Handler(http.StripPrefix("/", fileServer))
	r.PathPrefix("/fonts").Handler(http.StripPrefix("/", fileServer))
	r.PathPrefix("/").HandlerFunc(indexHandler("ui/dist/index.html"))

	return r
}

func indexHandler(entrypoint string) func(w http.ResponseWriter, r *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, entrypoint)
	}
	return fn
}
