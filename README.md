# Chart Viewer

You can try the demo [here](https://chart-viewer-84975.web.app)

A simple web app to help you inspect helm chart. So far, you can use this tool for:
- Inspect helm chart; as simple as showing all the chart templates
- Compare template between two versions; showing changes on git-like view.
- Compare rendered manifest between two versions; showing changes on git-like view.
- Render kubernetes manifest; modify the `values.yaml` and render the manifest. You will get a link that you can use to directly create the kubernetes resources on your cluster, as simple as `kubectl apply -f http://the.given.link` 

## Prerequisite
- Docker

## Run Instruction

### Run on docker
```shell script
$ git clone https://github.com/ecojuntak/chart-viewer.git
$ cd chart-viewer/
$ CHART_REPOS=`cat ./seed.json` docker-compose up
```
It will run three containers on your local, `redis`, `server`, and `ui`.
Then access http://localhost:8080 on your browser.

### Configuration
You can add more chart repo on the `seed.json` file.
```json
[
  {
    "name": "stable",
    "url": "https://kubernetes-charts.storage.googleapis.com"
  },
  {
    "name": "incubator",
    "url": "https://kubernetes-charts-incubator.storage.googleapis.com"
  },
  {
    "name": "bitnami",
    "url": "https://charts.bitnami.com/bitnami"
  }
]
```

## Roadmap
No roadmap yet. Still looking others feature that can be implemeted here.

## Contribute
Pull requests are welcome!
