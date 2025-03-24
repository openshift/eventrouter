FROM registry.ci.openshift.org/ocp/builder:rhel-9-golang-1.23-openshift-4.19 AS builder
WORKDIR  /go/src/github.com/openshift/eventrouter
USER 0

COPY ./go.mod ./go.sum ./
RUN go mod download
COPY Makefile *.go ./
COPY sinks ./sinks

RUN make build

FROM registry.access.redhat.com/ubi9/ubi-minimal

ARG BUILD_VERSION=0.5.0
USER 1000
COPY --from=builder /go/src/github.com/openshift/eventrouter/eventrouter /bin/eventrouter
CMD ["/bin/eventrouter", "-v", "3", "-logtostderr"]
LABEL version="v${BUILD_VERSION}"
