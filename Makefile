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

build: fmt
	go build -mod=mod -o $(TARGET)
.PHONY: build

fmt:
	@echo gofmt

image:
	podman build $(CONTAINER_BUILD_ARGS) --build-arg BUILD_VERSION=$(BUILD_VERSION)  -f Dockerfile .
	podman tag localhost/eventrouter $(IMAGE_REPOSITORY_NAME)

image-push: image
	podman manifest push --all $(IMAGE_REPOSITORY_NAME) docker://$(IMAGE_REPOSITORY_NAME)

.PHONY: image

test:
	go test -mod=mod $(TEST_PKGS) $(TESTARGS)
.PHONY: test
