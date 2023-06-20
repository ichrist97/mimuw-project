#!/bin/bash

# Install Ansible and sshpass
sudo apt -y install ansible sshpass

# Add all of your machines to ssh known hosts
for i in `seq -w 01 10`; do sshpass -f password ssh st103@st103vm1$i.rtb-lab.pl -o StrictHostKeyChecking=no -C "/bin/true"; done

# define ansible hosts
sudo cp hosts /etc/ansible/hosts

echo "Initialzed ansible."