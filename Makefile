APP_NAME=chart-viewer
APP_VERSION = $(shell cat VERSION)

build-backend:
	CGO_ENABLED=0 go build -o bin/chart-viewer cmd/main.go

build-frontend:
	cd ui; make build

test:
	go test -cover ./...

run:build-backend build-frontend
	./bin/chart-viewer serve --host 0.0.0.0 --redis-host 127.0.0.1

run-backend:build-backend
	./bin/chart-viewer serve --host 0.0.0.0 --redis-host 127.0.0.1

seed:build-backend
	./bin/chart-viewer seed --repo-seed seed.json --kube-version-seed api_versions.json

help:build-backend
	./bin/chart-viewer --help

package:
	docker build . -t ecojuntak/${APP_NAME}:${APP_VERSION} -t ecojuntak/${APP_NAME}:latest

publish-image:
	docker push ecojuntak/${APP_NAME}:${APP_VERSION} 
	docker push ecojuntak/${APP_NAME}:latest
