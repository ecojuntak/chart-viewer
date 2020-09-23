package service_test

import (
	helmMock "chart-viewer/helm/mocks"
	"chart-viewer/model"
	repoMock "chart-viewer/repository/mocks"
	"chart-viewer/service"
	"crypto/md5"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestService_GetRepos(t *testing.T) {
	repository := new(repoMock.Repository)
	helm := new(helmMock.Helm)
	stringifiedRepos := "[{\"name\":\"stable\",\"url\":\"https://chart.stable.com\"}]"
	repository.On("Get", "repos").Return(stringifiedRepos).Once()
	svc := service.NewService(helm, repository)
	charts := svc.GetRepos()

	expectedCharts := []model.Repo{
		{
			Name: "stable",
			URL:  "https://chart.stable.com",
		},
	}

	assert.Equal(t, expectedCharts, charts)
}

func TestService_GetChartsFromCache(t *testing.T) {
	stringifiedChart := "[{\"name\":\"discourse\",\"versions\":[\"0.3.5\",\"0.3.4\",\"0.3.3\",\"0.3.2\"]}]"

	repository := new(repoMock.Repository)
	helm := new(helmMock.Helm)
	repository.On("Get", "stable").Return(stringifiedChart)
	svc := service.NewService(helm, repository)
	err, charts := svc.GetCharts("stable")
	assert.NoError(t, err)

	expectedCharts := []model.Chart{
		{
			Name: "discourse",
			Versions: []string{
				"0.3.5", "0.3.4", "0.3.3", "0.3.2",
			},
		},
	}

	assert.Equal(t, expectedCharts, charts)
}

func TestService_GetValuesFromCache(t *testing.T) {
	stringifiedValues := "{\"affinity\":{},\"cloneHtdocsFromGit\":{\"enabled\":false,\"interval\":60}}"

	repository := new(repoMock.Repository)
	helm := new(helmMock.Helm)
	repository.On("Get", "value-stable-app-deploy-v0.0.1").Return(stringifiedValues)
	svc := service.NewService(helm, repository)
	err, values := svc.GetValues("stable", "app-deploy", "v0.0.1")
	assert.NoError(t, err)

	expectedValues := map[string]interface{}{
		"affinity": map[string]interface{}{},
		"cloneHtdocsFromGit": map[string]interface{}{
			"enabled":  false,
			"interval": float64(60),
		},
	}

	assert.Equal(t, expectedValues, values)
}

func TestService_GetTemplatesFromCache(t *testing.T) {
	stringifiedTemplates := "[{\"name\":\"deployment.yaml\",\"content\":\"kind: Deployment\"}]"

	repository := new(repoMock.Repository)
	helm := new(helmMock.Helm)
	repository.On("Get", "template-stable-app-deploy-v0.0.1").Return(stringifiedTemplates)
	svc := service.NewService(helm, repository)
	templates := svc.GetTemplates("stable", "app-deploy", "v0.0.1")

	expectedTemplates := []model.Template{
		{
			Name:    "deployment.yaml",
			Content: "kind: Deployment",
		},
	}

	assert.Equal(t, expectedTemplates, templates)
}

func TestService_GetStringifiedManifestsFromCache(t *testing.T) {
	stringifiedManifest := "{\"url\":\"http://chart-viewer.com\",\"manifests\":[{\"name\":\"deployment.yaml\",\"content\":\"kind: Deployment\"}]}"

	repository := new(repoMock.Repository)
	helm := new(helmMock.Helm)
	repository.On("Get", "manifests-stable-app-deploy-v0.0.1-hash").Return(stringifiedManifest)
	svc := service.NewService(helm, repository)
	manifest := svc.GetStringifiedManifests("stable", "app-deploy", "v0.0.1", "hash")

	expectedManifests := "---\nkind: Deployment\n"

	assert.Equal(t, expectedManifests, manifest)
}

func TestService_RenderManifest(t *testing.T) {
	createValuesTestFile()
	hash := getValuesHash()

	stringifiedManifest := "{\"url\":\"/charts/manifests/stable/app-deploy/v0.0.1/"+hash+"\",\"manifests\":[{\"name\":\"deployment.yaml\",\"content\":\"kind: Deployment\"}]}"
	rawManifest := "---\nkind: Deployment\n"
	manifest := model.ManifestResponse{
		URL: "/charts/manifests/stable/app-deploy/v0.0.1/" + hash,
		Manifests: []model.Manifest{
			{
				Name:    "deployment.yaml",
				Content: "kind: Deployment",
			},
		},
	}

	repository := new(repoMock.Repository)
	helm := new(helmMock.Helm)
	repository.On("Get", "manifests-stable-app-deploy-v0.0.1-" + hash).Return(stringifiedManifest)
	repository.On("Set", "manifests-stable-app-deploy-v0.0.1-"+hash, rawManifest)
	helm.On("RenderManifest", "stable", "app-deploy", "v0.0.1", []string{"/tmp/values.yaml"}).Return(manifest)
	svc := service.NewService(helm, repository)
	err, manifest := svc.RenderManifest("stable", "app-deploy", "v0.0.1", []string{"/tmp/values.yaml"})
	assert.NoError(t, err)

	expectedManifests := model.ManifestResponse{
		URL: "/charts/manifests/stable/app-deploy/v0.0.1/" + hash,
		Manifests: []model.Manifest{
			{
				Name:    "deployment.yaml",
				Content: "kind: Deployment",
			},
		},
	}

	assert.Equal(t, expectedManifests, manifest)
}

func getValuesHash() string {
	valuesFileContent, _ := ioutil.ReadFile("/tmp/values.yaml")
	hash := md5.Sum(valuesFileContent)
	return fmt.Sprintf("%x", hash)
}

func createValuesTestFile() {
	valueBytes := []byte("affinity: {}")
	fileLocation := "/tmp/values.yaml"
	_ = ioutil.WriteFile(fileLocation, valueBytes, 0644)
}
