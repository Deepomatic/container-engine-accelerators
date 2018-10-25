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
