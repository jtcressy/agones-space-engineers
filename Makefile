REPOSITORY = ghcr.io/jtcressy/agones-space-engineers

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
project_path := $(dir $(mkfile_path))
root_path := $(realpath $(project_path)/../..)
image_tag = $(REPOSITORY):0.9

#   _____                    _
#  |_   _|_ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                 |___/

# Build a docker image for the server, and tag it
build:
	cd $(project_path) && \
	docker build -f Dockerfile \
	--pull --cache-from=$(REPOSITORY):latest \
	--tag=$(REPOSITORY):latest \
	--tag=$(image_tag) \
	.

push:
	docker push $(image_tag) && docker push $(REPOSITORY):latest

deploy:
	kubectl apply -f 01-storage.yaml && \
	kubectl create secret -n gameservers generic \
		space-engineers-config \
		--from-file=SpaceEngineers-Dedicated.cfg \
		--save-config=true --dry-run=client -o yaml | kubectl apply -f - && \
	cat 20-fleet.yaml | sed 's=__IMAGE_TAG__=$(image_tag)=g' | kubectl apply -f - && \
	kubectl get all -n gameservers

destroy:
	kubectl delete -f 20-fleet.yaml && \
	kubectl delete secret -n gameservers space-engineers-config && \
	kubectl get all -n gameservers

destroy-storage:
	kubectl delete -f 01-storage.yaml

destroy-all: destroy destroy-storage

docker-login:
	echo $(GH_PAT) | docker login ghcr.io -u jtcressy --password-stdin

all: build docker-login push deploy

#output the tag for this image
echo-image-tag:
	@echo $(image_tag)