#!/bin/sh
set -x
if ! [ -x /usr/local/bin/ctr ]; then
  tar  -xvzf ../containerd/cri-containerd-cni-linux.tar.gz -C /
  [ -f /usr/lib64/libseccomp.so.2 ] || cp -rf ../lib64/lib* /usr/lib64/
  systemctl enable  containerd.service
  systemctl restart containerd.service
fi
# 已经安装了containerd并且运行了, 就不去重启.
ctr version || systemctl restart containerd.service
[ -f ../images/images.tar ] && ctr -n=k8s.io image import ../images/images.tar
if grep "SystemdCgroup = true"  /etc/containerd/config.toml &> /dev/nul; then
  driver=systemd
else
  driver=cgroupfs
fi
echo "driver is ${driver}"
export criDriver=${driver}
