package analyzer

import (
	"chart-viewer/pkg/model"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type analytic struct{}

type Analytic interface {
	Analyze(templates []model.Template, kubeAPIVersions model.KubernetesAPIVersion) ([]model.AnalyticsResult, error)
}

func New() Analytic {
	return analytic{}
}

func (a analytic) Analyze(templates []model.Template, kubeAPIVersions model.KubernetesAPIVersion) ([]model.AnalyticsResult, error) {
	var results []model.AnalyticsResult

	for _, t := range templates {
		r := model.AnalyticsResult{
			Template: model.Template{
				Name:    t.Name,
				Content: t.Content,
			},
			Compatible: true,
		}

		apiVersion, err := extractAPIVersion(t.Content)
		if err != nil {
			return nil, err
		}

		if apiVersion == "" {
			continue
		}

		r.Compatible = isCompatible(kubeAPIVersions.APIVersions, apiVersion)
		results = append(results, r)
	}

	return results, nil
}

func isCompatible(versions []string, version string) bool {
	exists := false
	for _, v := range versions {
		if v == version {
			exists = true
			break
		}
	}

	if !exists {
		return false
	}

	return true
}

func extractAPIVersion(template string) (string, error) {
	var resource model.KubeResourceCommonSpec
	lines := strings.Split(template, "\n")

	for _, line := range lines {
		if strings.Contains(line, "apiVersion") {
			err := yaml.Unmarshal([]byte(line), &resource)
			if err != nil {
				fmt.Printf("Error unmarshalling template: %s \n\n\n", err)
				return "", err
			}
		}
	}

	return resource.APIVersion, nil
}
