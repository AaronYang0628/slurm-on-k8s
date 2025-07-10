ARG OS_BASE_IMAGE=m.daocloud.io/docker.io/library/debian:bookworm
FROM ${OS_BASE_IMAGE}

ARG TZ=Asia/Shanghai
ARG SLURM_VERSION=25.05.0
ARG GOSU_VERSION=1.17
ARG MUNGE_UID=1108
ARG SLURM_UID=1109
ARG WORKER_UID=1110
ARG BUILDER_IMAGE=ghcr.io/aaronyang0628/slurm-builder:25.05
ENV TZ=${TZ} \
    SLURM_VERSION=${SLURM_VERSION} \
    MUNGE_UID=1108 \
    SLURM_UID=1109 \
    WORKER_UID=1110 \
    DEBIAN_FRONTEND=noninteractive
RUN groupadd -g $MUNGE_UID munge \
  && useradd  -m -c "MUNGE Uid 'N' Gid Emporium" -d /var/lib/munge -u $MUNGE_UID -g munge  -s /sbin/nologin munge \
  && groupadd -g $SLURM_UID slurm \
  && useradd  -m -c "Slurm workload manager" -d /var/lib/slurm -u $SLURM_UID -g slurm  -s /bin/bash slurm \
  && groupadd -g $WORKER_UID worker \
  && useradd  -m -c "Workflow user" -d /home/worker -u $WORKER_UID -g worker  -s /bin/bash worker

RUN apt-get update \
    && apt-get install -y openssh-server iputils-ping telnet iproute2 wget \
    && mkdir -p /run/sshd && chmod 755 /run/sshd \
    && echo "sshd:ALL" >> /etc/hosts.allow \
    && echo "PasswordAuthentication no" >> /etc/ssh/sshd_config \
    && echo "PermitRootLogin yes" >> /etc/ssh/sshd_config
RUN apt-get clean && apt-get update \
    && apt-get install -y gosu libmunge-dev libmunge2 munge\
    && gosu nobody true

RUN mkdir -p /run/munge \
    && chmod 755 /run/munge \
    && chown munge:munge /run/munge
RUN chmod 775 /var/log \
    && mkdir /var/spool/slurmd \
    && chown slurm: /var/spool/slurmd \
    && mkdir /var/log/slurm \
    && chown slurm: /var/log/slurm \
    && mkdir /var/spool/slurmctld \
    && chown slurm: /var/spool/slurmctld \
    && mkdir -p /etc/slurm \
    && rm -rf /slurm-${SLURM_VERSION}
COPY --from=${BUILDER_IMAGE} /opt/slurm/debs/slurm-smd_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd_${SLURM_VERSION}-1_amd64.deb
RUN apt install -y /opt/slurm/debs/slurm-smd_${SLURM_VERSION}-1_amd64.deb

