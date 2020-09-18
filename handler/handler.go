package handler

import (
	"encoding/json"
	"fmt"
	"chart-viewer/model"
	"chart-viewer/service"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) handler {
	return handler{
		service: service,
	}
}

func (h *handler) RepoListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		chartRepo := h.service.GetRepoList()

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chartRepo)
	} else if r.Method == "POST" {
		repo := model.Repo{
			Name: r.FormValue("name"),
			URL:  r.FormValue("url"),
		}

		h.service.AddRepo(repo)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	}
}

func (h *handler) RepoDetailHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	repoName := vars["repo-name"]
	charts := h.service.GetRepoDetail(repoName)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(charts)
}

func (h *handler) ChartHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	repoName := vars["repo-name"]
	chartName := vars["chart-name"]
	chartVersion := vars["chart-version"]

	values := h.service.GetChartValues(repoName, chartName, chartVersion)
	templates := h.service.GetChartTemplates(repoName, chartName, chartVersion)

	chartDetail := model.ChartDetail{
		Values:    values,
		Templates: templates,
	}

	response, _ := json.Marshal(chartDetail)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/yaml")
	w.Write(response)
}

func (h *handler) RenderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	values := r.FormValue("values")

	repoName := vars["repo-name"]
	chartName := vars["chart-name"]
	chartVersion := vars["chart-version"]

	valueBytes := []byte(values)
	fileLocation := fmt.Sprintf("/tmp/%s-values.yaml", time.Now().Format("20060102150405"))
	err := ioutil.WriteFile(fileLocation, valueBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}

	valueFile := []string{fileLocation}
	url, manifests := h.service.GenerateManifest(repoName, chartName, chartVersion, valueFile)

	res := &model.ManifestResponse{
		URL:       url,
		Manifests: manifests,
	}

	response, err := json.Marshal(res)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *handler) DownloadManifestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	repoName := vars["repo-name"]
	chartName := vars["chart-name"]
	chartVersion := vars["chart-version"]
	hash := vars["hash"]

	manifests := h.service.DownloadManifest(repoName, chartName, chartVersion, hash)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text/yaml")
	w.Write([]byte(manifests))
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
