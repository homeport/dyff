#!/usr/bin/env bash

# Copyright © 2018 Matthias Diester
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

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
