#!/bin/bash

set -e
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
IMAGE_NAME=${IMAGE_NAME:-ghcr.io/aaronyang0628/slurm-munged:25.05}
OS_BASE_IMAGE=${OS_BASE_IMAGE:-ghcr.io/aaronyang0628/slurm-base:25.05}
TLS_VERIFY=${TLS_VERIFY:-false}
podman build \
    --build-arg OS_BASE_IMAGE=${OS_BASE_IMAGE} \
    -f $SCRIPT_DIR/Dockerfile \
    -t $IMAGE_NAME $SCRIPT_DIR
