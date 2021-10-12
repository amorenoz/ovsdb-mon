FROM quay.io/centos/centos:stream8

USER root

RUN dnf install -y centos-release-nfv-openvswitch
RUN INSTALL_PKGS=" \
    openvswitch2.15 ovn-2021-host ovn-2021-central \
    iptables iproute iputils tcpdump socat procps \
    make go git \
        " && \
    dnf install --best --refresh -y --setopt=tsflags=nodocs $INSTALL_PKGS && \
    dnf clean all && rm -rf /var/cache/dnf/*

ENV GOPATH=$HOME/go

ADD dist/entrypoint.sh /root/entrypoint.sh

WORKDIR /root
ENTRYPOINT /root/entrypoint.sh
