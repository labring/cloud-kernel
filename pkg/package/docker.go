package _package

//k8s docker docker k8s
var dockerShell = `yum install -y git conntrack && \
git clone https://github.com/sealyun/cloud-kernel && \
cd cloud-kernel && mkdir -p kube && cp -rf runtime/docker/* kube/ && \
wget https://dl.k8s.io/v%s/kubernetes-server-linux-amd64.tar.gz && \
wget https://download.docker.com/linux/static/stable/x86_64/docker-%s.tgz && \
cp  docker-%s.tgz kube/docker/docker.tgz && \
tar zxvf kubernetes-server-linux-amd64.tar.gz && \
cp  kubernetes/server/bin/kubectl kube/bin/ && \
cp  kubernetes/server/bin/kubelet kube/bin/ && \
cp  kubernetes/server/bin/kubeadm kube/bin/ && \
sed s/k8s_version/%s/g -i kube/conf/kubeadm.yaml && \
cd kube/shell && chmod a+x docker.sh && sh docker.sh && \
rm -rf /etc/docker/daemon.json && systemctl restart docker && \
sh init.sh && sh master.sh && \
docker pull fanux/lvscare &&  \
cp /usr/sbin/conntrack ../bin/`
