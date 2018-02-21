#!/usr/bin/env bash
echo
echo -e "\n#### Docke Installation"
echo -e "\n#### Remove old docker versions ####"
yum remove docker \
    docker-common \
    container-selinux \
    docker-selinux \
    docker-engine

# yum -y remove docker-selinux

echo -e "\n#### Adding Docker repositories ####"
yum install -y yum-utils

yum-config-manager \
    --add-repo \
    https://download.docker.com/linux/centos/docker-ce.repo

yum makecache fast

echo -e "\n#### Installing Docker ####"
yum -y install docker-ce

echo -e "\n#### Starting Docker ####"
systemctl start docker

echo -e "\n#### Enable Docker after reboot####"
systemctl enable docker

echo -e "\n#### Add vagrant user to docker group####"
usermod -a -G docker vagrant

# https://github.com/eclipse/che/issues/795
#sudo chmod 777 /var/run/docker.sock


#echo -e "#### Test docker installation ####"
# docker run hello-world

#echo -e "\n#### Installing Docker-Compose ####"
#curl -L https://github.com/docker/compose/releases/download/1.11.2/docker-compose-`uname -s`-`uname -m` > /usr/local/bin/docker-compose

#chmod +x /usr/local/bin/docker-compose
#docker-compose --version

#echo -e "\n#### Installing Command Completion for Docker-Compose ####"
#curl -L https://raw.githubusercontent.com/docker/compose/$(docker-compose version --short)/contrib/completion/bash/docker-compose -o /etc/bash_completion.d/docker-compose
#


