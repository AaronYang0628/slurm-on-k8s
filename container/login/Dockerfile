ARG OS_BASE_IMAGE=ghcr.io/aaronyang0628/slurm-base:25.05
FROM ${OS_BASE_IMAGE}

ARG SLURM_VERSION=25.05.0
ARG MPI_TYPE=open-mpi
ARG BUILDER_IMAGE=ghcr.io/aaronyang0628/slurm-builder:25.05
COPY --from=${BUILDER_IMAGE} /opt/slurm/debs/slurm-smd-client-dbgsym_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd-client-dbgsym_${SLURM_VERSION}-1_amd64.deb
COPY --from=${BUILDER_IMAGE} /opt/slurm/debs/slurm-smd-client_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd-client_${SLURM_VERSION}-1_amd64.deb

RUN apt install -y curl vim build-essential cmake gcc-11 g++-11 gfortran 

ENV MPI_TYPE=${MPI_TYPE}

RUN if [ "$MPI_TYPE" = "open-mpi" ]; then \
        apt install -y vim libopenmpi-dev; \
    elif [ "$MPI_TYPE" = "intel-mpi" ]; then \
        curl -k -o intel-mpi-2021.14.1.7_offline.sh -L https://data-and-computing.oss-cn-hangzhou-zjy-d01-a.res.cloud.zhejianglab.com/projects%2Fslurm-on-k8s%2Fintel-mpi-libs%2Fintel-mpi-2021.14.1.7_offline.sh \
        && bash intel-mpi-2021.14.1.7_offline.sh -s -a -s --eula accept \
        && rm -rf intel-mpi-2021.14.1.7_offline.sh; \
    else \
        echo "you need to spcify ENV 'MPI_TYPE', either 'open-mpi' or 'intel-mpi'"; \
    fi

RUN curl -k -o intel-dpcpp-cpp-compiler-2025.0.4.20_offline.sh https://data-and-computing.oss-cn-hangzhou-zjy-d01-a.res.cloud.zhejianglab.com/projects%2Fslurm-on-k8s%2Fintel-mpi-libs%2Fintel-dpcpp-cpp-compiler-2025.0.4.20_offline.sh
RUN bash intel-dpcpp-cpp-compiler-2025.0.4.20_offline.sh -s -a -s --eula accept

ENV PATH="/opt/intel/oneapi/compiler/2025.0/bin:${PATH}"
RUN rm -rf intel-dpcpp-cpp-compiler-2025.0.4.20_offline.sh

RUN curl -k -o intel-fortran-compiler-2025.0.4.21_offline.sh https://data-and-computing.oss-cn-hangzhou-zjy-d01-a.res.cloud.zhejianglab.com/projects%2Fslurm-on-k8s%2Fintel-mpi-libs%2Fintel-fortran-compiler-2025.0.4.21_offline.sh
RUN bash intel-fortran-compiler-2025.0.4.21_offline.sh -s -a -s --eula accept

RUN rm -rf intel-fortran-compiler-2025.0.4.21_offline.sh


RUN apt install -y /opt/slurm/debs/slurm-smd-client-dbgsym_${SLURM_VERSION}-1_amd64.deb /opt/slurm/debs/slurm-smd-client_${SLURM_VERSION}-1_amd64.deb
COPY entrypoint.sh /opt/slurm/entrypoint.sh
RUN chmod 755 /opt/slurm/entrypoint.sh

CMD ["/opt/slurm/entrypoint.sh"]