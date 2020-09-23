package handler_test

import (
	"chart-viewer/handler"
	"chart-viewer/model"
	"chart-viewer/service/mocks"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/kinbiko/jsonassert"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_GetRepos(t *testing.T) {
	repos := []model.Repo {
		{Name: "stable", URL: "https://repo.stable"},
	}
	serviceMock := new(mocks.Service)
	serviceMock.On("GetRepos").Return(repos).Once()
	appHandler := handler.NewHandler(serviceMock)

	req, err := http.NewRequest("GET", "/repos", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	h := http.HandlerFunc(appHandler.GetReposHandler)
	h.ServeHTTP(recorder, req)

	content, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(content), `[
		{"name": "stable","url": "https://repo.stable"}
	]`)
}

func TestHandler_GetChartsHandler(t *testing.T) {
	charts := []model.Chart{
		{Name: "app-deployment", Versions: []string{"v0.0.1", "v0.0.2"}},
		{Name: "job-deployment", Versions: []string{"v0.2.0", "v0.2.1"}},
	}
	serviceMock := new(mocks.Service)
	serviceMock.On("GetCharts", "stable").Return(nil, charts).Once()
	appHandler := handler.NewHandler(serviceMock)

	req, err := http.NewRequest("GET", "/charts/stable", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/charts/{repo-name}", appHandler.GetChartsHandler)
	router.ServeHTTP(recorder, req)

	content, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(content), `[
		{"name": "app-deployment","versions": ["v0.0.1", "v0.0.2"]},
		{"name": "job-deployment","versions": ["v0.2.0", "v0.2.1"]}
	]`)
}

func TestHandler_GetChartHandler(t *testing.T) {
	chart := model.ChartDetail{
		Values: map[string]interface{}{"appPort": float64(8080) },
		Templates: []model.Template{
			{
				Name:    "deployment.yaml",
				Content: "kind: Deployment",
			},
		},
	}
	serviceMock := new(mocks.Service)
	serviceMock.On("GetChart", "repo-name", "chart-name", "chart-version").Return(nil, chart).Once()
	appHandler := handler.NewHandler(serviceMock)

	req, err := http.NewRequest("GET", "/charts/repo-name/chart-name/chart-version", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/charts/{repo-name}/{chart-name}/{chart-version}", appHandler.GetChartHandler)
	router.ServeHTTP(recorder, req)

	content, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(content), `{
	   "values":{
		  "appPort":8080
	   },
	   "templates":[
		  {
			 "name":"deployment.yaml",
			 "content":"kind: Deployment"
		  }
	   ]
	}`)
}

func TestHandler_GetValuesHandler(t *testing.T) {
	values := map[string]interface{}{
		"values": map[string]interface{}{
			"apiVersion": "app/Deployment",
			"cpuRequest": 11,
			"enableService": true,
		},
	}
	serviceMock := new(mocks.Service)
	serviceMock.On("GetValues", "repo-name", "chart-name", "chart-version").Return(nil, values).Once()
	appHandler := handler.NewHandler(serviceMock)

	req, err := http.NewRequest("GET", "/charts/values/repo-name/chart-name/chart-version", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/charts/values/{repo-name}/{chart-name}/{chart-version}", appHandler.GetValuesHandler)
	router.ServeHTTP(recorder, req)

	content, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(content), `
		{
			"values": {
				"apiVersion": "app/Deployment",
				"cpuRequest": 11,
				"enableService": true
			}
		}
	`)
}

func TestHandler_GetTemplateHandler(t *testing.T) {
	templates := []model.Template{
		{Name: "deployment.yaml", Content: "apiVersion: app/Deployment"},
		{Name: "service.yaml", Content: "kind: Service"},
	}
	serviceMock := new(mocks.Service)
	serviceMock.On("GetTemplates", "repo-name", "chart-name", "chart-version").Return(templates).Once()
	appHandler := handler.NewHandler(serviceMock)

	req, err := http.NewRequest("GET", "/charts/templates/repo-name/chart-name/chart-version", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/charts/templates/{repo-name}/{chart-name}/{chart-version}", appHandler.GetTemplatesHandler)
	router.ServeHTTP(recorder, req)

	content, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(content), `
		[
			{"name": "deployment.yaml", "content": "apiVersion: app/Deployment"},
			{"name": "service.yaml", "content": "kind: Service"}
		]
	`)
}

func TestHandler_GetManifestsHandler(t *testing.T) {
	stringfiedManifests := `
---
# Source: nginx/templates/server-block-configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: nginx-server-block
  labels:
    app.kubernetes.io/name: nginx
    helm.sh/chart: nginx-6.2.1
    app.kubernetes.io/instance: nginx
    app.kubernetes.io/managed-by: Helm
data:
  server-blocks-paths.conf: |-
    include  "/opt/bitnami/nginx/conf/server_blocks/ldap/*.conf";
    include  "/opt/bitnami/nginx/conf/server_blocks/common/*.conf";
`
	serviceMock := new(mocks.Service)
	serviceMock.On("GetStringifiedManifests", "repo-name", "chart-name", "chart-version", "hash").Return(stringfiedManifests).Once()
	appHandler := handler.NewHandler(serviceMock)

	req, err := http.NewRequest("GET", "/charts/manifests/repo-name/chart-name/chart-version/hash", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/charts/manifests/{repo-name}/{chart-name}/{chart-version}/{hash}", appHandler.GetManifestsHandler)
	router.ServeHTTP(recorder, req)

	content, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, stringfiedManifests, string(content))
}

func TestHandler_RenderManifestsHandler(t *testing.T) {
	manifests := model.ManifestResponse{
		URL:       "/charts/manifests/repo-name/chart-name/chart-version/hash",
		Manifests: []model.Manifest{
			{Name: "deployment.yaml", Content: "apiVersion: app/Deployment"},
			{Name: "service.yaml", Content: "kind: Service"},
		},
	}
	fileLocation := fmt.Sprintf("/tmp/%s-values.yaml", time.Now().Format("20060102150405"))
	serviceMock := new(mocks.Service)
	serviceMock.On("RenderManifest", "repo-name", "chart-name", "chart-version", []string{fileLocation}).Return(nil, manifests).Once()
	appHandler := handler.NewHandler(serviceMock)

	req, err := http.NewRequest("GET", "/charts/templates/render/repo-name/chart-name/chart-version?values=affinity:{}", nil)
	assert.NoError(t, err)

	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/charts/templates/render/{repo-name}/{chart-name}/{chart-version}", appHandler.RenderManifestsHandler)
	router.ServeHTTP(recorder, req)

	content, err := ioutil.ReadAll(recorder.Body)
	if err != nil {
		t.Error(err)
	}

	ja := jsonassert.New(t)
	ja.Assertf(string(content), `
		{
			"url" : "/charts/manifests/repo-name/chart-name/chart-version/hash",
			"manifests": [
				{"name": "deployment.yaml", "content": "apiVersion: app/Deployment"},
				{"name": "service.yaml", "content": "kind: Service"}
			]
		}
	`)
}