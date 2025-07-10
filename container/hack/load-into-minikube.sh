#!/bin/bash

load_image_to_minikube() {
    local image_name=$1
    echo "Processing image: $image_name"
    local file_name=$(echo "$image_name" | tr -s '/' '_' | tr -s ':' '_')
    podman save -o "/tmp/${file_name}.tar.dim" "$image_name"
    minikube image load "/tmp/${file_name}.tar.dim"
    rm "/tmp/${file_name}.tar.dim"
    echo "Image $image_name has been loaded into Minikube."
}

if [ "$MPI_TYPE" = "open-mpi" ]; then
    images=(
        "ghcr.io/aaronyang0628/slurm-login:25.05-open-mpi"
        "ghcr.io/aaronyang0628/slurm-slurmd:25.05-open-mpi"
        "ghcr.io/aaronyang0628/slurm-slurmctld:25.05"
        "ghcr.io/aaronyang0628/slurm-slurmdbd:25.05"
        "ghcr.io/aaronyang0628/slurm-munged:25.05"
    )
elif [ "$MPI_TYPE" = "intel-mpi" ]; then
    images=(
        "localhost/ay-dev/slurm-login:intel-mpi"
        "localhost/ay-dev/slurm-slurmd:intel-mpi"
        # "localhost/ay-dev/slurm-slurmctld:latest"
        # "localhost/ay-dev/slurm-slurmdbd:latest"
        # "localhost/ay-dev/slurm-munged:latest"
    )
else
    echo "unknow mpi type, please check ENV MPI_TYPE"
fi

set -e

for image in "${images[@]}"; do
    load_image_to_minikube "$image"
done
