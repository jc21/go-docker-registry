#!/bin/bash
set -eu

export GREEN='\E[1;32m'
export RED='\E[1;31m'
export YELLOW='\E[1;33m'
export RESET='\033[0m'

cd -- "$(dirname -- "$0")/.." || exit 1

export COMPOSE_FILE=test/docker-compose.yml

docker compose pull
docker compose up -d --remove-orphans registry

trap cleanup EXIT
cleanup() {
	cde="$?"
	echo -e "${YELLOW}Cleaning up...${RESET}"
	docker compose down --remove-orphans --volumes
	if [ "$cde" -ne 0 ]; then
		echo -e "${RED}TESTING FAILED${RESET}"
	fi
}

# copy some images from docker over to our local registry
# so we can test the registry itself
declare -a images=(
	'jc21/rpmbuild-centos6:latest'    # classic docker image
	'jc21/dnsrouter:latest'           # multiarch docker image
	'debian:stable-slim'              # rebuilt every day
)

# copy
for i in "${images[@]}"
do
	echo -e "${YELLOW}Copying ${GREEN}${i}${YELLOW} to local registry...${RESET}"
	docker compose run --rm --no-deps skopeo copy --all \
		"docker://docker.io/${i}" \
		"docker://registry.local:5000/${i}" \
		--dest-tls-verify=false
done

# TODO: test
echo -e "${YELLOW}Testing...${RESET}"
docker compose run --rm --no-deps project ./scripts/integration-test.sh

echo -e "${GREEN}Complete${RESET}"
