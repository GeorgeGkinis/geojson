#!/usr/bin/env bash
set -e
echo -e "\n#### Enabling password authentication ####"
sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/g' /etc/ssh/sshd_config