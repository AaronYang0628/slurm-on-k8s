ARG OS_BASE_IMAGE=ghcr.io/aaronyang0628/slurm-base:25.05
FROM ${OS_BASE_IMAGE}

COPY entrypoint.sh /opt/slurm/entrypoint.sh
RUN chmod 755 /opt/slurm/entrypoint.sh

CMD ["/opt/slurm/entrypoint.sh"]