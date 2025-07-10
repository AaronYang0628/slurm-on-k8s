ARG OS_BASE_IMAGE=ghcr.io/aaronyang0628/slurm-base:25.05
FROM ${OS_BASE_IMAGE}

ARG SLURM_VERSION=25.05.0
ARG BUILDER_IMAGE=ghcr.io/aaronyang0628/slurm-builder:25.05
COPY --from=${BUILDER_IMAGE} /opt/slurm/debs/slurm-smd-slurmctld-dbgsym_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd-slurmctld-dbgsym_${SLURM_VERSION}-1_amd64.deb
COPY --from=${BUILDER_IMAGE} /opt/slurm/debs/slurm-smd-slurmctld_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd-slurmctld_${SLURM_VERSION}-1_amd64.deb

RUN apt install -y /opt/slurm/debs/slurm-smd-slurmctld-dbgsym_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd-slurmctld_${SLURM_VERSION}-1_amd64.deb

COPY entrypoint.sh /opt/slurm/entrypoint.sh
RUN chmod 755 /opt/slurm/entrypoint.sh

CMD ["/opt/slurm/entrypoint.sh"]