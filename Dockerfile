FROM quay.io/centos/centos:stream9

USER root

RUN dnf install -y centos-release-nfv-openvswitch
RUN INSTALL_PKGS=" \
    openvswitch2.17 ovn23.09-host ovn23.09-central \
    iptables iproute iputils tcpdump socat procps \
    make go git \
        " && \
    dnf install --best --refresh -y --setopt=tsflags=nodocs $INSTALL_PKGS && \
    dnf clean all && rm -rf /var/cache/dnf/*

ENV GOPATH=$HOME/go

ADD . /root/ovsdb-mon
RUN cd /root/ovsdb-mon && go install github.com/ovn-kubernetes/libovsdb/cmd/modelgen && go mod vendor

RUN ln -s /root/ovsdb-mon/dist/entrypoint.sh /root/entrypoint.sh
WORKDIR /root
ENTRYPOINT /root/entrypoint.sh
