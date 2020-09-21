package handler

import (
	"chart-viewer/model"
	"chart-viewer/service"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type handler struct {
	service service.Service
}

func NewHandler(service service.Service) handler {
	return handler{
		service: service,
	}
}

func (h *handler) GetReposHandler(w http.ResponseWriter, r *http.Request) {
	chartRepo := h.service.GetRepos()
	respondWithJSON(w, http.StatusOK, chartRepo)
}

func (h *handler) GetChartsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	charts := h.service.GetCharts(repoName)

	respondWithJSON(w, http.StatusOK, charts)
}

func (h *handler) GetChartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	chartName := vars["chart-name"]
	chartVersion := vars["chart-version"]
	values := h.service.GetValues(repoName, chartName, chartVersion)
	templates := h.service.GetTemplates(repoName, chartName, chartVersion)

	chart := model.ChartDetail{
		Values:    values,
		Templates: templates,
	}

	respondWithJSON(w, http.StatusOK, chart)
}

func (h *handler) GetValuesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	chartName := vars["chart-name"]
	chartVersion := vars["chart-version"]
	values := h.service.GetValues(repoName, chartName, chartVersion)

	respondWithJSON(w, http.StatusOK, values)
}

func (h *handler) GetTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	chartName := vars["chart-name"]
	chartVersion := vars["chart-version"]
	templates := h.service.GetTemplates(repoName, chartName, chartVersion)

	respondWithJSON(w, http.StatusOK, templates)
}

func (h *handler) GetManifestsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	chartName := vars["chart-name"]
	chartVersion := vars["chart-version"]
	hash := vars["hash"]

	manifest := h.service.GetStringifiedManifests(repoName, chartName, chartVersion, hash)

	respondWithText(w, http.StatusOK, manifest)
}

func (h *handler) RenderManifestsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	values := r.FormValue("values")

	repoName := vars["repo-name"]
	chartName := vars["chart-name"]
	chartVersion := vars["chart-version"]

	valueBytes := []byte(values)
	fileLocation := fmt.Sprintf("/tmp/%s-values.yaml", time.Now().Format("20060102150405"))
	err := ioutil.WriteFile(fileLocation, valueBytes, 0644)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Cannot store values to file: "+err.Error())
		return
	}

	valueFile := []string{fileLocation}
	err, manifests := h.service.RenderManifest(repoName, chartName, chartVersion, valueFile)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error rendering manifest: "+err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, manifests)
}

func (h *handler) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Headers:", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
		return
	})
}

func (h *handler) LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		log.Printf("in coming request: %s\n", r.URL.Path)

		next.ServeHTTP(w, r)
		return
	})
}
