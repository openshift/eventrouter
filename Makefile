# Copyright 2017 Heptio Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

TARGET = eventrouter
GOTARGET = github.com/openshift/$(TARGET)
BUILD_VERSION?=0.5.0
CONTAINER_BUILD_ARGS ?=
IMAGE_REPOSITORY_NAME ?= quay.io/openshift/logging-eventrouter:v${BUILD_VERSION}

ifneq ($(VERBOSE),)
VERBOSE_FLAG = -v
endif
TESTARGS ?= $(VERBOSE_FLAG) -timeout 60s
TEST_PKGS ?= $(GOTARGET)/sinks/...
TEST = go test $(TEST_PKGS) $(TESTARGS)

# check-deps brackets the build for the same reason as in test: the root
# package is not covered by TEST_PKGS, so go build -mod=mod is where its
# new dependencies materialize in go.sum.
build: fmt check-deps
	go build -mod=mod -o $(TARGET)
	@$(MAKE) --no-print-directory check-deps
.PHONY: build

fmt:
	@echo gofmt

image:
	podman build -t $(IMAGE_REPOSITORY_NAME) $(CONTAINER_BUILD_ARGS) --build-arg BUILD_VERSION=$(BUILD_VERSION)  -f Dockerfile .

image-push: image
	podman manifest push --all $(IMAGE_REPOSITORY_NAME) docker://$(IMAGE_REPOSITORY_NAME)

.PHONY: image

# Guard: none of the removed sink dependency families may return to
# go.mod or go.sum. (klauspost/compress is excluded: it remains a
# test-only indirect dep of prometheus/client_golang, which is kept for
# /metrics; it is not linked into the eventrouter binary.)
check-deps:
	@if grep -Eq 'IBM/sarama|crewjam/rfc5424|eapache|sethgrid/pester|jcmturner|pierrec/lz4|golang/snappy|hashicorp/go-uuid|rcrowley/go-metrics' go.mod go.sum; then \
		echo "ERROR: banned dependency family present in go.mod or go.sum"; exit 1; \
	fi
.PHONY: check-deps

# check-deps runs both before (fail fast on a dirty go.sum) and after
# go test: -mod=mod lets go test add checksums to go.sum, so only the
# post-test check can catch a dependency introduced by this run.
test: check-deps
	go test -mod=mod $(TEST_PKGS) $(TESTARGS)
	@$(MAKE) --no-print-directory check-deps
.PHONY: test
