
### push to dockerhub registry
```shell
export DOCKER_CR_PAT=YOUR_TOKEN
echo $DOCKER_CR_PAT | docker login docker.io -u aaron666 --password-stdin
podman push docker.io/aaronyang0628/slurm-operator:latest
```


### push to github registry
```shell
export GITHUB_CR_PAT=YOUR_TOKEN
echo $GITHUB_CR_PAT | docker login ghcr.io -u aaronyang0628 --password-stdin
podman push ghcr.io/aaronyang0628/slurm-operator:latest
```
