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
- **munged**

### Usage

> if you wanna change slurm configuration ,please check slurm configuration generator [click](https://slurm.schedmd.com/configurator.html)

- for helm user
    > just run for fun!
    1. `helm repo add slurm https://aaronyang0628.github.io/slurm-on-k8s/`
    2. `helm install slurm slurm/chart --version 1.0.2`

- for developer
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
    2. load images to minikube for test
        ```shell
        MPI_TYPE=open-mpi #intel-mpi
        bash ./container/load-into-minikube.sh
        ```
    3. publish helm chart
        ```shell
        helm package --dependency-update --destination /tmp/ ./chart
        ```
    4. test your helm chart
        ```shell
        helm install slurm /tmp/slurm-1.0.X.tgz
        ```
    5. index your chart
        ```shell
        helm repo index ./chart/ && cp -f ./chart/index.yaml ./index.yaml
        ```
