#!/usr/bin/env bash
#yum update -y

yum install -y wget
wget https://dl.google.com/go/go1.9.4.linux-amd64.tar.gz

tar -xzf go1.9.4.linux-amd64.tar.gz

mv go /usr/local

echo "export GOROOT=/usr/local/go" >> /home/vagrant/.bash_profile
echo "export GOPATH=/vagrant" >> /home/vagrant/.bash_profile
echo "export GOBIN=/vagrant/bin" >> /home/vagrant/.bash_profile
echo "export PATH=/vagrant/bin:/usr/local/go/bin:$PATH" >> /home/vagrant/.bash_profile
