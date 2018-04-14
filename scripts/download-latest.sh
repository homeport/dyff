#!/usr/bin/env bash

set -euo pipefail

if ! hash curl 2> /dev/null; then
  echo "Required tool curl is not installed."
  exit 1
fi

if ! hash jq 2> /dev/null; then
  echo "Required tool jq is not installed."
  exit 1
fi

LATEST_TAG="$(curl --silent --location https://api.github.com/repos/HeavyWombat/dyff/releases | jq --raw-output '.[0].tag_name')"
SYSTEM_UNAME="$(uname | tr '[:upper:]' '[:lower:]')"
TARGET_FILE=/usr/local/bin/dyff

case "${SYSTEM_UNAME}" in
  darwin|linux)
    DYFF_URI="https://github.com/HeavyWombat/dyff/releases/download/${LATEST_TAG}/dyff-${SYSTEM_UNAME}-amd64"

    echo -e "Downloading \\033[4;94m${DYFF_URI}\\033[0m to \\033[1m${TARGET_FILE}\\033[0m"
    curl --progress-bar --location "${DYFF_URI}" --output "${TARGET_FILE}" && chmod a+rx "${TARGET_FILE}"
    echo -e "\\nSuccessfully installed \\033[1m$(${TARGET_FILE} version)\\033[0m\\n"
    ;;

  *)
    echo "Unsupported operating system: ${SYSTEM_UNAME}"
    exit 1
    ;;
esac
