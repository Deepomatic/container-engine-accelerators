Fork of the [GCP `nvidia-gpu-device-plugin`](https://github.com/GoogleCloudPlatform/container-engine-accelerators).

# About this fork
This fork advertises multiple fake GPU for each real GPU, allowing to share a GPU between multiple pods using the Kubernetes device plugin api.

The goal is to schedule pods on GPUs until the GPU memory is full (GPU memory bin-packing).

See [here](cmd/nvidia_gpu/README.md) for more details (notably [how to use it on GKE](cmd/nvidia_gpu/README.md#how-to-use-in-gke)).

See also [deepomatic/shared-gpu-nvidia-k8s-device-plugin](https://github.com/Deepomatic/shared-gpu-nvidia-k8s-device-plugin) for the same feature based on [NVIDIA own Kubernetes GPU Device Plugin](https://github.com/NVIDIA/k8s-device-plugin).

## Limits
This is a big workaround given the current situation. It has many drawbacks:
The kubernetes scheduler doesn't know how the underlying real GPUs are shared between the `deepomatic.com/shared-gpu` resources it allocates among Pods.

- there is now way to control/guarantee spreading the pods among real GPUs: the current workaround is to limit to one real GPU per node and to indirectly schedule via other resources such as `memory` (assuming there is a correlation between `memory` and GPU (memory) usage.
- in the case of multiple real GPUs per node, asking for multiple shared GPUs for one Pod doesn't make sense as there is no guarantee the pod will be allocated shared GPUs from different real GPUs

## Roadmap
For proper scheduling, this device plugin will advertise `SharedGPUMemory` as [Kubernetes Extended Resources](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#extended-resources). Since the `SharedGPUMemory` resource is at the Node level (instead of at the Device level), we effectively support only one GPU per node.

## Configuration
You can control the number of fake GPU device declared by this device plugin by adding `-gpu-duplication-factor N` to the container `command` (default: `100`) in the [daemonset definition](cmd/nvidia_gpu/daemonset.yaml).


# Hardware Accelerators in GKE

This repository is a collection of installation recipes and integration utilities for consuming Hardware Accelerators in Google Kubernetes Engine.

This is not an official Google product.

More details on the nvidia-gpu-device-plugin are [here](cmd/nvidia_gpu/README.md).
