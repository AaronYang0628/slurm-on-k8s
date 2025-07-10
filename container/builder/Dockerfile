ARG OS_BASE_IMAGE=m.daocloud.io/docker.io/library/debian:bookworm

FROM ${OS_BASE_IMAGE}
ARG SLURM_BINARY_BASE_URL=https://download.schedmd.com/slurm
ARG SLURM_VERSION=25.05.0
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get clean && apt-get update \
    && apt-get install -y build-essential fakeroot devscripts equivs curl
WORKDIR /opt/slurm/build
RUN curl -o slurm-${SLURM_VERSION}.tar.bz2 ${SLURM_BINARY_BASE_URL}/slurm-${SLURM_VERSION}.tar.bz2 \
    && tar -jxf slurm-${SLURM_VERSION}.tar.bz2 \
    && rm slurm-${SLURM_VERSION}.tar.bz2 \
    && cd slurm-${SLURM_VERSION} \
    && mk-build-deps -i --tool 'apt-get -y --no-install-recommends' debian/control \
    && debuild -b -uc -us \
    && cd - \
    && mkdir -p /opt/slurm/debs \
    && mv slurm*.deb /opt/slurm/debs \
    && rm -rf slurm-${SLURM_VERSION}
