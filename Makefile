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

.PHONY: all vet fmt lint gocyclo megacheck misspell ginkgo test install build

all: test

clean:
	@rm -rf $(dir $(realpath $(firstword $(MAKEFILE_LIST))))/binaries
	@go clean -i -cache $(shell go list ./...)

vet:
	$(dir $(realpath $(firstword $(MAKEFILE_LIST))))scripts/go-vet.sh

fmt:
	$(dir $(realpath $(firstword $(MAKEFILE_LIST))))scripts/go-fmt.sh

lint:
	$(dir $(realpath $(firstword $(MAKEFILE_LIST))))scripts/go-lint.sh

gocyclo:
	$(dir $(realpath $(firstword $(MAKEFILE_LIST))))scripts/go-cyclo.sh

misspell:
	$(dir $(realpath $(firstword $(MAKEFILE_LIST))))scripts/misspell.sh

ginkgo:
	GO111MODULE=on ginkgo -r --randomizeAllSpecs --randomizeSuites --failOnPending --trace --race --nodes=4 --compilers=2 --cover

test: vet fmt lint gocyclo misspell ginkgo

install: test
	@$(dir $(realpath $(firstword $(MAKEFILE_LIST))))/scripts/compile-version.sh --only-local

build: test
	@$(dir $(realpath $(firstword $(MAKEFILE_LIST))))/scripts/compile-version.sh --no-local
