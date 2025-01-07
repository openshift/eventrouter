FROM brew.registry.redhat.io/rh-osbs/openshift-golang-builder:v1.15.14-202307071644.el8.g4a85211 AS builder

# Start Konflux-specific steps
RUN mkdir -p /tmp/yum_temp; mv /etc/yum.repos.d/*.repo /tmp/yum_temp/ || true
COPY .oit/unsigned.repo /etc/yum.repos.d/
ADD https://certs.corp.redhat.com/certs/Current-IT-Root-CAs.pem /tmp
# End Konflux-specific steps
ENV __doozer=update BUILD_RELEASE=202501072046.p0.g3134f6a.assembly.stream.el8 BUILD_VERSION=v4.6.0 OS_GIT_MAJOR=4 OS_GIT_MINOR=6 OS_GIT_PATCH=0 OS_GIT_TREE_STATE=clean OS_GIT_VERSION=4.6.0-202501072046.p0.g3134f6a.assembly.stream.el8 SOURCE_GIT_TREE_STATE=clean __doozer_group=openshift-4.6 __doozer_key=logging-eventrouter __doozer_uuid_tag=ose-logging-eventrouter-v4.6.0-20250107.204624 __doozer_version=v4.6.0 
ENV __doozer=merge OS_GIT_COMMIT=3134f6a OS_GIT_VERSION=4.6.0-202501072046.p0.g3134f6a.assembly.stream.el8-3134f6a SOURCE_DATE_EPOCH=1685976826 SOURCE_GIT_COMMIT=3134f6af9e41a57228a22d8071a860e2d92b4278 SOURCE_GIT_TAG=v0.2-10-g3134f6a SOURCE_GIT_URL=https://github.com/openshift/eventrouter 
WORKDIR  /go/src/github.com/openshift/eventrouter
COPY . .
RUN go build .

FROM quay.io/redhat-user-workloads/ocp-art-tenant/art-images:base-rhel8-v4.6.0-20250107.204624

# Start Konflux-specific steps
RUN mkdir -p /tmp/yum_temp; mv /etc/yum.repos.d/*.repo /tmp/yum_temp/ || true
COPY .oit/unsigned.repo /etc/yum.repos.d/
ADD https://certs.corp.redhat.com/certs/Current-IT-Root-CAs.pem /tmp
# End Konflux-specific steps
ENV __doozer=update BUILD_RELEASE=202501072046.p0.g3134f6a.assembly.stream.el8 BUILD_VERSION=v4.6.0 OS_GIT_MAJOR=4 OS_GIT_MINOR=6 OS_GIT_PATCH=0 OS_GIT_TREE_STATE=clean OS_GIT_VERSION=4.6.0-202501072046.p0.g3134f6a.assembly.stream.el8 SOURCE_GIT_TREE_STATE=clean __doozer_group=openshift-4.6 __doozer_key=logging-eventrouter __doozer_uuid_tag=ose-logging-eventrouter-v4.6.0-20250107.204624 __doozer_version=v4.6.0 
ENV __doozer=merge OS_GIT_COMMIT=3134f6a OS_GIT_VERSION=4.6.0-202501072046.p0.g3134f6a.assembly.stream.el8-3134f6a SOURCE_DATE_EPOCH=1685976826 SOURCE_GIT_COMMIT=3134f6af9e41a57228a22d8071a860e2d92b4278 SOURCE_GIT_TAG=v0.2-10-g3134f6a SOURCE_GIT_URL=https://github.com/openshift/eventrouter 
COPY --from=builder /go/src/github.com/openshift/eventrouter/eventrouter /bin/eventrouter
CMD ["/bin/eventrouter", "-v", "3", "-logtostderr"]

# Start Konflux-specific steps
RUN cp /tmp/yum_temp/* /etc/yum.repos.d/ || true
# End Konflux-specific steps

LABEL \
        License="GPLv2+" \
        io.k8s.description="OpenShift logging event router" \
        io.k8s.display-name="Logging Event Router" \
        io.openshift.tags="logging,eventrouter" \
        vendor="Red Hat" \
        name="openshift/ose-logging-eventrouter" \
        com.redhat.component="logging-eventrouter-container" \
        io.openshift.maintainer.project="OCPBUGS" \
        io.openshift.maintainer.component="Logging" \
        version="v4.6.0" \
        release="202501072046.p0.g3134f6a.assembly.stream.el8" \
        io.openshift.build.commit.id="3134f6af9e41a57228a22d8071a860e2d92b4278" \
        io.openshift.build.source-location="https://github.com/openshift/eventrouter" \
        io.openshift.build.commit.url="https://github.com/openshift/eventrouter/commit/3134f6af9e41a57228a22d8071a860e2d92b4278" \
        description="" \
        summary=""

