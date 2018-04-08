#!/usr/bin/env bash

set -euo pipefail

BASEDIR="$(cd "$(dirname "$0")/.." && pwd)"
VERSION="$(cd "$BASEDIR" && git describe --tags)"
VERFILE="$BASEDIR/cmd/version.go"

on_exit() {
  if [[ -f "${VERFILE}.bak" ]]; then
    mv "${VERFILE}.bak" "${VERFILE}"
  fi
}

# Run on exit function to clean-up
trap on_exit EXIT

# Backup current version of the version subcommand and set current tag as version
cp -p "${VERFILE}" "${VERFILE}.bak"
perl -pi -e "s/const version = \"\(development\)\"/const version = \"${VERSION}\"/g" "${VERFILE}"

TARGET_PATH="${BASEDIR}/binaries"
mkdir -p "$TARGET_PATH"
while read -r OS ARCH; do
  echo "Compiling dyff version ${VERSION} for OS ${OS} and architecture ${ARCH}"
  TARGET_FILE="${TARGET_PATH}/dyff-${OS}-${ARCH}"
  if [[ "${OS}" == "windows" ]]; then
    TARGET_FILE="${TARGET_FILE}.exe"
  fi

  ( cd $BASEDIR && GOOS=$OS GOARCH=$ARCH go build -o "$TARGET_FILE" )

done << EOL
darwin	386
darwin	amd64
linux	386
linux	amd64
linux	s390x
windows	386
windows	amd64
EOL
