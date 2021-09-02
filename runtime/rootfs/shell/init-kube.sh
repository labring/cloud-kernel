#!/bin/bash
# Open ipvs
modprobe -- ip_vs
modprobe -- ip_vs_rr
modprobe -- ip_vs_wrr
modprobe -- ip_vs_sh
# 1.20需要开启br_netfilter
modprobe -- br_netfilter
modprobe -- bridge

## version_ge 4.19 4.19 true ;
## version_ge 5.4 4.19 true ;
## version_ge 3.10 4.19 false ;

version_ge(){
    test "$(echo "$@" | tr ' ' '\n' | sort -rV | head -n 1)" == "$1"
}

disable_selinux(){
    if [ -s /etc/selinux/config ] && grep 'SELINUX=enforcing' /etc/selinux/config; then
        sed -i 's/SELINUX=enforcing/SELINUX=disabled/g' /etc/selinux/config
        setenforce 0
    fi
}

get_distribution() {
	lsb_dist=""
	# Every system that we officially support has /etc/os-release
	if [ -r /etc/os-release ]; then
		lsb_dist="$(. /etc/os-release && echo "$ID")"
	fi
	# Returning an empty string here should be alright since the
	# case statements don't act unless you provide an actual value
	echo "$lsb_dist"
}

disable_firewalld() {
  lsb_dist=$( get_distribution )
	lsb_dist="$(echo "$lsb_dist" | tr '[:upper:]' '[:lower:]')"
	case "$lsb_dist" in
		ubuntu|deepin|debian)
			command -v ufw &> /dev/null && ufw disable
		;;
		centos|rhel|kylin|neokylin)
			systemctl stop firewalld && systemctl disable firewalld
		;;
		*)
			systemctl stop firewalld && systemctl disable firewalld
		;;
	esac
}

kernel_version=$(uname -r | cut -d- -f1)
if version_ge "${kernel_version}" 4.19; then
  modprobe -- nf_conntrack
else
  modprobe -- nf_conntrack_ipv4
fi

cat <<EOF >  /etc/sysctl.d/k8s.conf
net.bridge.bridge-nf-call-ip6tables = 1
net.bridge.bridge-nf-call-iptables = 1
net.ipv4.conf.all.rp_filter=0
EOF
sysctl --system
sysctl -w net.ipv4.ip_forward=1
disable_firewalld
swapoff -a || true
disable_selinux

cp ../bin/* /usr/bin
# Cgroup driver
cp ../conf/kubelet.service /etc/systemd/system/
[ -d /etc/systemd/system/kubelet.service.d ] || mkdir /etc/systemd/system/kubelet.service.d
cp ../conf/10-kubeadm.conf /etc/systemd/system/kubelet.service.d/

systemctl enable kubelet
