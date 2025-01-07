FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.15-openshift-4.6 AS builder
WORKDIR  /go/src/github.com/openshift/eventrouter
COPY . .
RUN go build .

FROM registry.ci.openshift.org/ocp/builder:rhel-8-base-openshift-4.6
COPY --from=builder /go/src/github.com/openshift/eventrouter/eventrouter /bin/eventrouter
CMD ["/bin/eventrouter", "-v", "3", "-logtostderr"]
