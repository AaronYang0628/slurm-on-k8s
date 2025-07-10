## AY Slurm Container
This is a development install of the Slurm cluster running as a controller and multiple types of workers. We split all related services into different containers, such as slurmd, slurmctld, slurmdbd, and munged. So you don't recommend you to run these images directly, you can deploy them using Helm or Operator. Please check [Helm Chart](https://github.com/AaronYang0628/slurm-on-k8s) and [Operator](https://github.com/AaronYang0628/slurm-on-k8s) for more information.

### What is Slurm
Slurm is an open-source cluster resource management and job scheduling system that strives to be simple, scalable, portable, fault-tolerant, and interconnect agnostic. SLURM currently has been tested only under Linux.

As a cluster resource manager, SLURM provides three key functions. First, it allocates exclusive and/or non-exclusive access to resources (compute nodes) to users for some duration of time so they can perform work. Second, it provides a framework for starting, executing, and monitoring work (normally a parallel job) on the set of allocated nodes. Finally, it arbitrates conflicting requests for resources by managing a queue of pending work.

For more information on Slurm, consult the official website⁠, check [link](https://slurm.schedmd.com/configurator.easy.html)

### What's inside
- **Resource Management**: Efficiently manages resources in a Kubernetes cluster, ensuring optimal utilization.
- **Job Scheduling**: Advanced job scheduling capabilities to handle various types of workloads.
- **Scalability**: Easily scales to accommodate growing workloads and resources.
- **High Availability**: Supports high availability configurations to ensure continuous operation.
- **Multi-User Support**: Allows multiple users to submit and manage their jobs concurrently.
- **Integration with MPI Libraries**: Supports both Open MPI and Intel MPI libraries for parallel computing.
- **Customizable**: Using values.yaml file, you can customizable a slurm cluster, fitting specific needs and configurations.
- **Separate Munged Daemon**: munged daemon is a daemon that provides a secure way to communicate between nodes in a cluster.

### How to use this image

As we said, we are not recommend you to run these images directly, you can deploy them using Helm or Operator. Please check [Helm Chart](https://github.com/AaronYang0628/slurm-on-k8s) and [Operator](https://github.com/AaronYang0628/slurm-on-k8s) for more information.

But you still can run it using docker compose.
```
docker compose up
```
This will start several containers with a `login` process which will run a sshd server on exposed port 22.

### To submit jobs
You will need to create an interactive session in order to run jobs in this container. 

There are two ways to do this.

First, you can start a container with the default command and ssh in.

```
docker exec -it slurm-login -- bin/bash
```

Once you have a session in the container, you can submit jobs using the sbatch command. A test script is included in the image at /home/slurm/hello.slurm. You can submit this script to verify the scheduler is working properly.

```shell
sbatch `/home/slurm/hello.slurm`
```

### How to build the image
Build from this directory using the enclosed Dockerfile

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

And then, load images to minikube for test
```shell
MPI_TYPE=open-mpi #intel-mpi
bash ./container/load-into-minikube.sh
```