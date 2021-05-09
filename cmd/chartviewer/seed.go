package chartviewer

import (
	"chart-viewer/pkg/helm"
	"chart-viewer/pkg/model"
	"chart-viewer/pkg/repository"
	"chart-viewer/pkg/server/service"
	"fmt"
	"io/ioutil"
	"log"
	"sync"

	"github.com/spf13/cobra"
)

var wg = &sync.WaitGroup{}

func NewSeedCommand() *cobra.Command {
	var (
		redisHost          string
		redisPort          string
		repoSeedPath       string
		apiVersionSeedPath string
	)

	command := cobra.Command{
		Use:     "seed",
		Short:   "Seed the redis with chart info",
		Example: "chart-viewer seed --redis-host 127.0.0.1 --redis-port 6379 --seed-file ./seed.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			redisAddress := fmt.Sprintf("%s:%s", redisHost, redisPort)

			repo, err := repository.NewRepository(redisAddress)
			if err != nil {
				log.Printf("cannot connect to redis: %s\n", err)
				return err
			}

			log.Printf("connected to redis on %s:%s\n", redisHost, redisPort)
			log.Println("starting to populate redis...")

			err = seedKubeVersion(repo, apiVersionSeedPath)
			if err != nil {
				log.Printf("failed to seed api version: %s\n", err)
			}
			log.Println("Kubernetes API version seeded")

			err = seedRepo(repo, repoSeedPath)
			if err != nil {
				log.Printf("failed to seed chart repository: %s\n", err)
				return err
			}
			seedChart(repo)
			wg.Wait()

			return nil
		},
	}

	command.Flags().StringVar(&redisHost, "redis-host", "127.0.0.1", "Redis host address")
	command.Flags().StringVar(&redisPort, "redis-port", "6379", "Redis host port")
	command.Flags().StringVar(&repoSeedPath, "repo-seed", "./seed.json", "Path to JSON file that contain array of repositories.")
	command.Flags().StringVar(&apiVersionSeedPath, "kube-version-seed", "./api_versions.json", "Path to JSON file that contain list of Kubernetes API version for each Kubernetes version")
	return &command
}

func seedKubeVersion(repo repository.Repository, path string) error {
	apiVersions, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	stringifiedApiVersion := string(apiVersions)
	repo.Set("api-versions", stringifiedApiVersion)
	return nil
}

func seedRepo(repo repository.Repository, seedPath string) error {
	repos, err := ioutil.ReadFile(seedPath)
	if err != nil {
		return err
	}

	log.Printf("populating reposistories from %s\n", seedPath)
	stringifiedRepos := string(repos)
	repo.Set("repos", stringifiedRepos)
	return nil
}

func seedChart(repo repository.Repository) {
	h := helm.NewHelmClient(repo)
	svc := service.NewService(h, repo, nil)

	chartRepos := svc.GetRepos()

	for _, repo := range chartRepos {
		wg.Add(1)
		go pullChart(svc, repo)
	}
}

func pullChart(svc service.Service, repo model.Repo) {
	defer wg.Done()
	err, charts := svc.GetCharts(repo.Name)
	if err != nil {
		log.Printf("error populating charts from repo %s: %s", repo.Name, err)
		return
	}

	for _, chart := range charts {
		versions := chart.Versions
		for _, version := range versions {
			log.Printf("populating %s/%s:%s\n", repo.Name, chart.Name, version)
			err, _ := svc.GetChart(repo.Name, chart.Name, version)

			if err != nil {
				log.Printf("error populating charts %s: %s", repo.Name, err)
			}
		}
	}
}
