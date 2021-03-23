



VERSION := $(shell git describe --tags | sed -e 's/^v//g' | awk -F "-" '{print $$1}')


build:
		docker build -t loncherapp-profiles-service:${VERSION} -tus-west-1.amazonaws.com/loncherapp-profiles-service:${VERSION} .

build-and-publish: build
		docker_push

docker_push:
	docker push .dkr.ecr.us-west-1.amazonaws.com/loncherapp-profiles-service:${VERSION}

aws_auth_repo:
	aws ecr get-login-password --region us-west-1 |  docker login --username AWS --password-stdin  .dkr.ecr.us-west-1.amazonaws.com


run-local:
	./scripts/runlocal.sh

