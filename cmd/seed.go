package cmd

import (
	"chart-viewer/helm"
	"chart-viewer/model"
	"chart-viewer/repository"
	"chart-viewer/service"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var wg = &sync.WaitGroup{}

func NewSeedCommand() *cobra.Command {
	defaultHost := os.Getenv("REDIS_HOST")
	defaultPort := os.Getenv("REDIS_PORT")
	defaultSeedPath := ""

	command := cobra.Command{
		Use:     "seed",
		Short:   "Seed the redis with chart info",
		Example: "chart-viewer seed --redis-host 127.0.0.1 --redis-port 6379 --seed-file ./seed.json",
		RunE: func(cmd *cobra.Command, args []string) error {
			host := defaultHost
			port := defaultPort
			redisAddress := fmt.Sprintf("%s:%s", defaultHost, defaultPort)

			err, repo := repository.NewRepository(redisAddress)
			if err != nil {
				fmt.Printf("cannot connect to redis: %s\n", err)
				return err
			}

			log.Printf("connected to redis on %s:%s\n", host, port)
			log.Println("starting to populate redis...")

			err = seedRepo(repo, defaultSeedPath)
			if err != nil {
				fmt.Printf("cannot connect to redis: %s\n", err)
				return err
			}
			seedChart(repo)
			wg.Wait()

			return nil
		},
	}

	command.Flags().StringVar(&defaultHost, "redis-host", "127.0.0.1", "[Optional] Redis host address")
	command.Flags().StringVar(&defaultPort, "redis-port", "6379", "[Optional] Redis host port")
	command.Flags().StringVar(&defaultSeedPath, "seed-file", "", "[Optional] Path to JSON file that contain array of repositories. Will read config from environment variable CHART_REPOS if not set")

	return &command
}

func seedRepo(repo repository.Repository, seedPath string) error {
	if seedPath != "" {
		repos, err := ioutil.ReadFile(seedPath)
		if err != nil {
			return err
		}

		log.Printf("populating reposistories from %s\n", seedPath)
		stringifiedRepos := string(repos)
		repo.Set("repos", stringifiedRepos)
		return nil
	}

	log.Println("populating reposistories from environment varaible CHART_REPOS")
	stringifiedRepos := os.Getenv("CHART_REPOS")
	repo.Set("repos", stringifiedRepos)
	return nil
}

func seedChart(repo repository.Repository) {
	h := helm.NewHelmClient(repo)
	svc := service.NewService(h, repo)

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
			svc.GetChart(repo.Name, chart.Name, version)
		}
	}
}
