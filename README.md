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
- **Customizable**: Using [values.yaml](https://raw.githubusercontent.com/AaronYang0628/helm-chart-mirror/refs/heads/main/templates/slurm/slurm.values.yaml) file, you can customizable a slurm cluster, fitting specific needs and configurations.
- **Separated munged daemon**
- **Support GPU nodes deployment**

### Usage

> if you wanna change slurm configuration ,please check slurm configuration generator, check [link](https://slurm.schedmd.com/configurator.html)

- for github helm user
    1. get helm repo and update
        ```
        helm repo add ay-helm-mirror https://aaronyang0628.github.io/helm-chart-mirror/charts
        ```
    2. install slurm chart
        ```
        helm install slurm ay-helm-mirror/chart -f charts/values.yaml --version 1.0.7
        ```
- for artifact helm user
    1.  get helm repo and update
        ```
        helm repo add ay-helm-mirror https://aaronyang0628.github.io/helm-chart-mirror/charts
        ```
    2. install slurm chart
        ```shell
        helm install slurm ay-helm-mirror/chart -f charts/values.yaml --version 1.0.7
        ```
    Or you can get template values.yaml from [link](https://raw.githubusercontent.com/AaronYang0628/helm-chart-mirror/refs/heads/main/templates/slurm/slurm.values.yaml)
- for opertaor user
    1. test pull an image and apply
        ```
        podman pull ghcr.io/aaronyang0628/slurm-operator:25.05
        ```
    2. deploy slurm operator
        ```
        kubectl apply -f https://raw.githubusercontent.com/AaronYang0628/helm-chart-mirror/refs/heads/main/templates/slurm/operator_install.yaml
        ```
    3. apply CRD slurmdeployment 
        ```
        kubectl apply -f https://raw.githubusercontent.com/AaronYang0628/helm-chart-mirror/refs/heads/main/templates/slurm/slurmdeployment.values.yaml
        ```

### Manage Your Slurm Cluster
- check cluster status
    ```shell
    kubectl get slurmdep slurmdeployment-sample
    kubectl -n slurm get pods -w
    ```

When everything is ready, you can login your cluster and submit jobs.
- Add PubKeys to login node
    ```markdown
    you can edit `auth.ssh.configmap.perfabPubKeys` the file chart/values.yaml to add your public keys
    Or you can edit `spec.values.auth.ssh.configmap.perfabPubKeys` in your slurmdeployment CRD
    ```

- reapply your chart or CRD
- login your cluster
    ```shell
    kubectl -n slurm exec -it deploy/slurm-login -c login -- bin/bash
    ```
    Or
    ```shell
    ssh root@slurm-login.svc.cluster.local
    ```
