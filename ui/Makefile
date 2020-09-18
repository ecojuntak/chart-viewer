include .env.development

APP_NAME=chart-viewer-ui
APP_VERSION=$(shell cat VERSION)

build:
	npm run build

deploy:
	npm run build
	firebase deploy
	
run: 
	npm run serve

package:
	docker build --build-arg VUE_APP_API_SERVER_HOST=${VUE_APP_API_SERVER_HOST} . -t ecojuntak/${APP_NAME}:${APP_VERSION} -t ecojuntak/${APP_NAME}:latest

run-docker: package
	docker run -p 8080:80 ecojuntak/${APP_NAME}:latest

publish-image:
	docker push ecojuntak/${APP_NAME}:${APP_VERSION}
	docker push ecojuntak/${APP_NAME}:latest