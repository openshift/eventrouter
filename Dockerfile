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

LABEL \
        io.k8s.display-name="Openshift EventRouter" \
        io.k8s.description="Captures Kubernetes/OpenShift events and sends them to the logging pipeline so you can store, search, and analyze events like logs." \
        io.openshift.tags="openshift,logging,EventRouter,observability" \
        description="Captures Kubernetes/OpenShift events and sends them to the logging pipeline so you can store, search, and analyze events like logs." \
        maintainer="AOS Logging <team-logging@redhat.com>" \
        summary="Watching the Kubernetes Event API, Writing events as structured log entries" \
        version="v${BUILD_VERSION}"
