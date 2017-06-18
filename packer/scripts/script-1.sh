#!/bin/bash -eux

USER_NAME="devops"

sudo yum install ntp -y && sudo systemctl start ntpd && sudo systemctl enable ntpd

cd /home/$USER_NAME
tar zxfv authorized_keys.tar.gz && rm -f authorized_keys.tar.gz
