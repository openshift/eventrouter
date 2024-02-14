FROM registry.redhat.io/ubi9/go-toolset:latest AS builder
WORKDIR  /go/src/github.com/openshift/eventrouter
USER 0
COPY Makefile *.go go.mod go.sum ./
COPY sinks ./sinks

RUN make build

FROM registry.access.redhat.com/ubi9/ubi-minimal
USER 1000
COPY --from=builder /go/src/github.com/openshift/eventrouter/eventrouter /bin/eventrouter
CMD ["/bin/eventrouter", "-v", "3", "-logtostderr"]
LABEL version=release-5.9
