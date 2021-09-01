#!/bin/bash
# Install docker
chmod a+x docker.sh
#./docker.sh  /var/docker/lib  127.0.0.1
bash docker.sh
chmod a+x init-kube.sh

bash init-kube.sh
