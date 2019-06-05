#!/usr/bin/env bash

# Copyright © 2019 The Homeport Team
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

if ! hash curl 2>/dev/null; then
  echo "Required tool curl is not installed."
  exit 1
fi

if [[ "$(uname -m)" != "x86_64" ]]; then
  echo -e "Unsupported machine type \\033[1m$(uname -m)\\033[0m: Please check \\033[4;94mhttps://api.github.com/repos/homeport/dyff/releases\\033[0m manually"
  exit 1
fi

if [[ $# -eq 0 ]]; then
  if ! hash jq 2>/dev/null; then
    echo -e 'Required tool \033[1mjq\033[0m is not installed.'
    exit 1
  fi

  # Find the latest version using the GitHub API
  SELECTED_TAG="$(curl --silent --location https://api.github.com/repos/homeport/dyff/releases | jq --raw-output '.[0].tag_name')"
else
  # Use provided argument as tag to download
  SELECTED_TAG="$1"
fi

SYSTEM_UNAME="$(uname | tr '[:upper:]' '[:lower:]')"

# Find a suitable install location
if [[ -w /usr/local/bin ]]; then
  TARGET_DIR=/usr/local/bin

elif [[ -w "$HOME/bin" ]] && grep -q -e "$HOME/bin" -e '\~/bin' <<<"$PATH"; then
  TARGET_DIR=$HOME/bin

else
  echo -e "Unable to determine a writable install location. Make sure that you have write access to either \\033[1m/usr/local/bin\\033[0m or \\033[1m$HOME/bin\\033[0m and that is in your PATH."
  exit 1
fi

# Download and install
case "${SYSTEM_UNAME}" in
  darwin | linux)
    DOWNLOAD_URI="https://github.com/homeport/dyff/releases/download/${SELECTED_TAG}/dyff-${SYSTEM_UNAME}-amd64"

    echo -e "Downloading \\033[4;94m${DOWNLOAD_URI}\\033[0m to place it into \\033[1m${TARGET_DIR}\\033[0m"
    if curl --progress-bar --location "${DOWNLOAD_URI}" --output "${TARGET_DIR}/dyff"; then
      chmod a+rx "${TARGET_DIR}/dyff"
      echo -e "\\nSuccessfully installed \\033[1mdyff\\033[0m into \\033[1m${TARGET_DIR}\\033[0m\\n"
    fi
    ;;

  *)
    echo "Unsupported operating system: ${SYSTEM_UNAME}"
    exit 1
    ;;
esac
