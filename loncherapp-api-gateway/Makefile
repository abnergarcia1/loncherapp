



VERSION := $(shell git describe --tags | sed -e 's/^v//g' | awk -F "-" '{print $$1}')


build:
		docker build -t loncherapp-api-gateway:${VERSION} -t loncherapp-api-gateway:${VERSION} .

build-and-publish: build
		docker_push

docker_push:
	docker push .us-west-1.amazonaws.com/loncherapp-api-gateway:${VERSION}

aws_auth_repo:
	aws ecr get-login-password --region us-west-1 |  docker login --username AWS --password-stdin  ecr.us-west-1.amazonaws.com


run-local:
	./scripts/runlocal.sh

run-deps:
	./scripts/rundockerdeps.sh

run-all:
	run-deps
	run-local

run: