FROM registry.ci.openshift.org/ocp/builder:rhel-8-golang-1.15-openshift-4.6 AS build
COPY . /go/src/github.com/openshift/eventrouter
RUN cd /go/src/github.com/openshift/eventrouter && go build .
FROM registry.ci.openshift.org/ocp/builder:rhel-8-base-openshift-4.6
COPY --from=build /go/src/github.com/openshift/eventrouter/eventrouter /bin/eventrouter
CMD ["/bin/eventrouter", "-v", "3", "-logtostderr"]
