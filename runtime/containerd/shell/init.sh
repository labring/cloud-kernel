#!/bin/bash
# Install containerd
chmod a+x containerd.sh
bash containerd.sh


# 修改kubelet
mkdir -p /etc/systemd/system/kubelet.service.d
cat > /etc/systemd/system/kubelet.service.d/containerd.conf << eof
[Service]
Environment="KUBELET_EXTRA_ARGS=--container-runtime=remote --cgroup-driver=${criDriver} --runtime-request-timeout=15m --container-runtime-endpoint=unix:///run/containerd/containerd.sock --image-service-endpoint=unix:///run/containerd/containerd.sock"
eof

cat /etc/systemd/system/kubelet.service.d/containerd.conf

chmod a+x init-kube.sh
bash init-kube.sh

cat /var/lib/kubelet/config.yaml
