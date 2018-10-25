Kubernetes Device Plugin for NVIDIA GPUs
----------------------------------------

This directory contains the code for a Kubernetes [device plugin](https://kubernetes.io/docs/concepts/cluster-administration/device-plugins/) for NVIDIA GPUs.

This device plugin provides `deepomatic.com/shared-gpu` resources, which are GPU devices mapping to a common, shared NVIDIA GPU.

The [daemonset manifest](daemonset.yaml) can be used to deploy this device plugin to a cluster (1.9 onwards):
```shell
kubectl apply -f https://raw.githubusercontent.com/Deepomatic/container-engine-accelerators/master/cmd/nvidia_gpu/daemonset.yaml
```

This device plugin requires that NVIDIA drivers and libraries are installed in a particular way.

Examples of how driver installation needs to be done can be found at:
- For [COS](https://cloud.google.com/container-optimized-os/):
  - Installer code: https://github.com/GoogleCloudPlatform/cos-gpu-installer
  - Installer daemonset: https://github.com/GoogleCloudPlatform/Deepomatic/blob/master/daemonset.yaml

- For [Ubuntu](https://cloud.google.com/kubernetes-engine/docs/concepts/node-images#ubuntu) (experimental):
  - Installer code: https://github.com/GoogleCloudPlatform/container-engine-accelerators/blob/master/nvidia-driver-installer/ubuntu/entrypoint.sh
  - Installer daemonset: https://github.com/GoogleCloudPlatform/container-engine-accelerators/blob/master/nvidia-driver-installer/ubuntu/daemonset.yaml

In short, this device plugins expects that all the nvidia libraries needed by the containers are present under a single directory on the host. You can specify the directory on the host containing nvidia libraries using `-host-path`. You can specify the location to mount that directory in all the containers using `-container-path`. For example, let's say on the host all nvidia libraries are present under `/var/lib/nvidia/lib64` and you want to make these libraries available to containers under `/usr/local/nvidia/lib64`, then you would use `-host-path=/var/lib/nvidia/lib64` and `-container-path=/usr/local/nvidia/lib64`.


## How to use in GKE

### Beta
We rely on beta GKE features, and not everything works perfectly as of 2018-10-25.

In [GKE](https://cloud.google.com/kubernetes-engine/) the [upstream GCP nvidia device plugin](https://github.com/GoogleCloudPlatform/container-engine-accelerators) is automatically installed via a kubernetes addon as well as a node taint.
This creates [some issues](https://issuetracker.google.com/issues/74386472#comment9) with this device plugin as the [node taints and pods tolerations mechanism with device plugins](https://notes.mindprince.in/2017/12/17/dedicated-node-pools-and-ExtendedResourceToleration-admission-controller.html) essentially assumes only one extended resource type is available per node.
In GKE with this device plugin we end up with two extended resources: `nvidia.com/gpu` (automatically) and the new `deepomatic.com/shared-gpu`.
In short: Pods requesting `deepomatic.com/shared-gpu` must thus explicitely tolerate `nvidia.com/gpu` taints for them to be scheduled.

Furthermore, currently [GKE refuses `/` in node taints](https://issuetracker.google.com/issues/118393036), so we cannot rely on the [ExtendedResourceToleration admission controller](https://notes.mindprince.in/2017/12/17/dedicated-node-pools-and-ExtendedResourceToleration-admission-controller.html) for `deepomatic.com/shared-gpu` either:
Pods requesting `deepomatic.com/shared-gpu` must thus also explicitely tolerate `deepomatic.com--shared-gpu` taints for them to be scheduled.

In fact [setting node taints at node-pool creation is also broken](https://issuetracker.google.com/issues/116842165), so the recommended way for now is to manually taint nodes (each time a new node is create, which effectively breaks node-pool auto-scaling).

### Create the cluster
- Create the Shared GPUs node-pool with one GPU and the following node label: `deepomatic.com/shared-gpu=true`
- Add the taint to the nodes:
  ```shell
  kubectl taint nodes -l deepomatic.com/shared-gpu=true deepomatic.com/shared-gpu=present:NoSchedule
  ```
- Install the `deepomatic-shared-gpu-gcp-k8s-device-plugin` DaemonSet:
  ```shell
  kubectl apply -f https://raw.githubusercontent.com/Deepomatic/container-engine-accelerators/master/cmd/nvidia_gpu/daemonset.yaml
  ```
- Install the nvidia driver (the docker image is preloaded on GKE, we use that with a special daemonset):
  ```shell
  kubectl apply -f https://github.com/Deepomatic/container-engine-accelerators/blob/master/nvidia-driver-installer/cos/daemonset-preloaded.yaml
  ```
- Verify the `deepomatic.com/shared-gpu` resources appear:
  ```
  kubectl describe nodes -l deepomatic.com/shared-gpu=true | grep deepomatic.com/shared-gpu
  ```

### Use the `deepomatic.com/shared-gpu` Extended Resource
See the [example](../../example/README.md) or the [demo](../../demo/).
