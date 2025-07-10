[![Artifact Hub](https://img.shields.io/endpoint?url=https://artifacthub.io/badge/repository/slurm-on-k8s)](https://artifacthub.io/packages/search?repo=slurm-on-k8s)

### Prerequisites
- kubectl version v1.11.3+.
- buildah version v1.33.10+
- Access to a Kubernetes v1.11.3+ cluster.

### Features

Slurm on Kubernetes provides the following features:

- **Resource Management**: Efficiently manages resources in a Kubernetes cluster, ensuring optimal utilization.
- **Job Scheduling**: Advanced job scheduling capabilities to handle various types of workloads.
- **Scalability**: Easily scales to accommodate growing workloads and resources.
- **High Availability**: Supports high availability configurations to ensure continuous operation.
- **Multi-User Support**: Allows multiple users to submit and manage their jobs concurrently.
- **Integration with MPI Libraries**: Supports both Open MPI and Intel MPI libraries for parallel computing.
- **Customizable**: Using `values.yaml` file, you can customizable a slurm cluster, fitting specific needs and configurations.
- **Separated munged daemon**
- **Support GPU nodes deployment**

### Development
- for image developer
    > you might need to build your own images or chart
    1. build images
        ```shell
        MPI_TYPE=open-mpi #intel-mpi
        bash ./container/builder/build.sh
        bash ./container/base/build.sh
        bash ./container/munged/build.sh
        bash ./container/login/build.sh
        bash ./container/slurmctld/build.sh
        bash ./container/slurmd/build.sh
        bash ./container/slurmdbd/build.sh
        ```
    2. push images to dockerhub registry
        ```shell
        export GITHUB_CR_PAT=ghp_xxxxxxxxxxxx
        echo $GITHUB_CR_PAT | podman login ghcr.io -u aaronyang0628 --password-stdin
        podman push ghcr.io/aaronyang0628/slurm-operator:latest
        ```
    3. load images to minikube for further test
        ```shell
        MPI_TYPE=open-mpi #intel-mpi
        bash ./container/load-into-minikube.sh
        ```
- for helm developer
    1. publish helm chart
        ```shell
        helm package --dependency-update --destination /tmp/ ./chart
        ```
    2. test your helm chart
        ```shell
        helm install -f ./chart/values.yaml slurm /tmp/slurm-1.0.X.tgz
        ```
    3. index your chart
        ```shell
        helm repo index ./chart/ && cp -f ./chart/index.yaml ./index.yaml
        ```
- for operator developer
    1. install CRDs
        ```shell
        make install
        ```
    2. test your operator in terminal
        ```shell
        make run
        ```
    3. build and push your operator
        ```shell
        make docker-build docker-push IMG=ghcr.io/aaronyang0628/slurm-operator:25.05
        ```
    4. deploy your operator
        ```shell
        kubectl apply -f https://raw.githubusercontent.com/AaronYang0628/slurm-on-k8s/refs/heads/main/operator/dist/install.yaml
        ```
