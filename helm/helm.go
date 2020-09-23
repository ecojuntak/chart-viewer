package helm

import (
	"bytes"
	"chart-viewer/model"
	"chart-viewer/repository"
	"fmt"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/releaseutil"
	"log"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

type Helm interface {
	GetValues(chartUrl, chartName, chartVersion string) (error, map[string]interface{})
	GetManifest(chartUrl, chartName, chartVersion string) []model.Template
	RenderManifest(chartUrl, chartName, chartVersion string, files []string) (error, []model.Manifest)
}

type helm struct {
	client     *action.Install
	repository repository.Repository
}

var settings = cli.New()

func debug(format string, v ...interface{}) {}

func NewHelmClient(repository repository.Repository) Helm {
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "", debug)
	if err != nil {
		panic(err)
	}

	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.ClientOnly = true
	client.UseReleaseName = true

	return &helm{
		client:     client,
		repository: repository,
	}
}

func (h helm) GetValues(chartUrl, chartName, chartVersion string) (error, map[string]interface{}) {
	log.Printf("getting %s:%s from remote\n", chartName, chartVersion)

	h.client.ChartPathOptions.Version = chartVersion
	h.client.ReleaseName = chartName
	h.client.RepoURL = chartUrl
	cp, err := h.client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		return err, nil
	}

	chartRequested, _ := loader.Load(cp)

	return nil, chartRequested.Values
}

func (h helm) GetManifest(chartUrl, chartName, chartVersion string) []model.Template {
	h.client.ChartPathOptions.Version = chartVersion
	h.client.ReleaseName = chartName
	h.client.RepoURL = chartUrl
	cp, err := h.client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		panic(err)
	}

	chartRequested, err := loader.Load(cp)

	templates := chartRequested.Templates

	var templateStrings []model.Template

	for _, t := range templates {
		templateStrings = append(templateStrings, model.Template{
			Name:    t.Name,
			Content: string(t.Data),
		})
	}

	return templateStrings
}

func (h helm) RenderManifest(chartUrl, chartName, chartVersion string, files []string) (error, []model.Manifest) {
	h.client.ChartPathOptions.Version = chartVersion
	h.client.ReleaseName = chartName
	h.client.RepoURL = chartUrl
	cp, err := h.client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		return err, nil
	}

	chartRequested, err := loader.Load(cp)

	valueOption := &values.Options{
		ValueFiles: files,
	}

	p := getter.All(settings)
	vals, err := valueOption.MergeValues(p)
	if err != nil {
		return err, nil
	}

	rel, err := h.client.Run(chartRequested, vals)
	if err != nil {
		return err, nil
	}

	var manifests bytes.Buffer
	fmt.Fprintln(&manifests, strings.TrimSpace(rel.Manifest))

	if !h.client.DisableHooks {
		for _, h := range rel.Hooks {
			fmt.Fprintln(&manifests, fmt.Sprintf("---\n # Source: %s\n%s", h.Path, h.Manifest))
		}
	}

	splitManifests := releaseutil.SplitManifests(manifests.String())
	manifestsKeys := make([]string, 0, len(splitManifests))
	for k := range splitManifests {
		manifestsKeys = append(manifestsKeys, k)
	}

	sort.Sort(releaseutil.BySplitManifestsOrder(manifestsKeys))

	var finalManifests []model.Manifest

	for _, manifestKey := range manifestsKeys {
		manifest := splitManifests[manifestKey]

		manifestNameRegex := regexp.MustCompile("# Source: [^/]+/(.+)")
		submatch := manifestNameRegex.FindStringSubmatch(manifest)
		if len(submatch) == 0 {
			continue
		}
		manifestName := submatch[1]
		manifestPath := filepath.Join(strings.Split(manifestName, "/")[1:]...)
		manifestPathSplit := strings.Split(manifestPath, "/")

		if manifestPathSplit[0] == "tests" {
			continue
		}

		finalManifests = append(finalManifests, model.Manifest{
			Name:    manifestPath,
			Content: manifest,
		})
	}

	return nil, finalManifests
}
