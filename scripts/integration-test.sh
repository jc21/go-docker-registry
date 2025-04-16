#!/bin/bash
set -eu

# This script should be run within the jc21/gotools image
# as part of the scripts/test.sh suite

export RED='\E[1;31m'
export YELLOW='\E[1;33m'
export RESET='\033[0m'

cd -- "$(dirname -- "$0")/.." || exit 1

trap cleanup EXIT
cleanup() {
	if [ "$?" -ne 0 ]; then
		echo -e "${RED}INTEGRATION TESTING FAILED${RESET}"
	fi
	rm -f cover.out
}

go test -tags=integration -json -cover -coverprofile="./cover.out" ./... | tparse
