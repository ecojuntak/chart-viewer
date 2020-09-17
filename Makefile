include .env

APP_NAME=chart-viewer-server
APP_VERSION = $(shell cat VERSION)
CHART_REPOS = $(shell cat ./seed.json)

build:
	go build -o bin/app .

run:build
	CHART_REPOS='${CHART_REPOS}' ./bin/app

package:
	docker build . -t ecojuntak/${APP_NAME}:${APP_VERSION} -t ecojuntak/${APP_NAME}:latest

publish-image:
	docker push ecojuntak/${APP_NAME}:${APP_VERSION} 
	docker push ecojuntak/${APP_NAME}:latest
