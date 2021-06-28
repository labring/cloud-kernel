#!/bin/bash
# Install docker
chmod a+x docker.sh
#./docker.sh  /var/docker/lib  127.0.0.1
sh docker.sh
chmod a+x init-kube.sh

sh init-kube.sh
