REPOSITORY = ghcr.io/jtcressy

mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
project_path := $(dir $(mkfile_path))
root_path := $(realpath $(project_path)/../..)
image_tag = $(REPOSITORY)/agones-space-engineers:0.1

#   _____                    _
#  |_   _|_ _ _ __ __ _  ___| |_ ___
#    | |/ _` | '__/ _` |/ _ \ __/ __|
#    | | (_| | | | (_| |  __/ |_\__ \
#    |_|\__,_|_|  \__, |\___|\__|___/
#                 |___/

# Build a docker image for the server, and tag it
build:
	cd $(project_path) && docker build -f Dockerfile --tag=$(image_tag) .

push:
	docker push $(image_tag)

docker-login:
	echo $(GH_PAT) | docker login ghcr.io -u jtcressy --password-stdin

#output the tag for this image
echo-image-tag:
	@echo $(image_tag)