package helm

import (
	"bytes"
	"fmt"
	"chart-viewer/model"
	"chart-viewer/repository"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/releaseutil"
	repoPkg "helm.sh/helm/v3/pkg/repo"
)

type HelmClient struct {
	client     *action.Install
	repository *repository.Repository
}

var (
	settings = cli.New()
)

func debug(format string, v ...interface{}) {
	if settings.Debug {
		log.Println(v)
	}
}

func NewHelmClient(repository *repository.Repository) *HelmClient {
	actionConfig := new(action.Configuration)
	err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), "", debug)
	if err != nil {
		panic(err)
	}

	client := action.NewInstall(actionConfig)
	client.DryRun = true
	client.ClientOnly = true
	client.UseReleaseName = true

	return &HelmClient{
		client:     client,
		repository: repository,
	}
}

func (h *HelmClient) InitRepos() {
	chartRepos := h.repository.GetRepos()

	for _, r := range chartRepos {
		err := addRepo(r.Name, r.URL)
		if err != nil {
			panic(err)
		}
	}
}

func addRepo(name string, url string) error {

	log.Printf("adding chart repo %s\n", name)

	b, err := ioutil.ReadFile(settings.RepositoryConfig)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	var f repoPkg.File
	if err := yaml.Unmarshal(b, &f); err != nil {
		return err
	}

	c := repoPkg.Entry{
		Name:     name,
		URL:      url,
		Username: "",
		Password: "",
		CertFile: "",
		KeyFile:  "",
		CAFile:   "",
	}

	rep, err := repoPkg.NewChartRepository(&c, getter.All(settings))

	if err != nil {
		return err
	}

	if _, err := rep.DownloadIndexFile(); err != nil {
		return errors.Wrapf(err, "looks like %q is not a valid chart repository or cannot be reached", url)
	}

	f.Update(&c)

	if err := f.WriteFile(settings.RepositoryConfig, 0644); err != nil {
		return err
	}

	log.Printf("%s chart repo added\n", name)

	return nil
}

func (h *HelmClient) GetValues(chartUrl, chartName, chartVersion string) map[string]interface{} {

	log.Printf("getting %s:%s from remote\n", chartName, chartVersion)

	h.client.ChartPathOptions.Version = chartVersion
	h.client.ReleaseName = chartName
	h.client.RepoURL = chartUrl
	cp, err := h.client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		panic(err)
	}

	chartRequested, _ := loader.Load(cp)

	return chartRequested.Values

}

func (h *HelmClient) GetManifest(chartUrl, chartName, chartVersion string) []model.Template {

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

func (h *HelmClient) RenderManifest(chartUrl, chartName, chartVersion string, files []string) (bytes.Buffer, []model.ManifestObject) {
	h.client.ChartPathOptions.Version = chartVersion
	h.client.ReleaseName = chartName
	h.client.RepoURL = chartUrl
	cp, err := h.client.ChartPathOptions.LocateChart(chartName, settings)
	if err != nil {
		panic(err)
	}

	chartRequested, err := loader.Load(cp)

	valueOption := &values.Options{
		ValueFiles: files,
	}

	p := getter.All(settings)
	vals, err := valueOption.MergeValues(p)

	rel, err := h.client.Run(chartRequested, vals)
	if err != nil {
		panic(err)
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

	var finalManifests []model.ManifestObject

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

		finalManifests = append(finalManifests, model.ManifestObject{
			Name:    manifestPath,
			Content: manifest,
		})
	}

	return manifests, finalManifests
}
