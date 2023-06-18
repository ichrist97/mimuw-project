#!/bin/bash

# Variables
host_vm1="st103vm102.rtb-lab.pl"
host_vm2="st103vm103.rtb-lab.pl"
remote_user="st103"
container_name="mimuw"
mongo_host="10.112.103.106"

# check if sshpass is installed
package_name="sshpass"
if dpkg -s "$package_name" >/dev/null 2>&1; then
    echo "$package_name is installed."
else
    echo "$package_name is not installed."
    sudo apt install -y $package_name
fi

# first vm
sshpass -f pass_file ssh "$remote_user@$host_vm1" "sudo docker run -d -e MONGO_HOST=$mongo_host -e DEBUG=1 -p 3000:3000 $container_name"
echo "Started tag service on $host_vm1"

# second vm
sshpass -f pass_file ssh "$remote_user@$host_vm2" "sudo docker run -d -e MONGO_HOST=$mongo_host -e DEBUG=1 -p 3000:3000 $container_name"
echo "Started tag service on $host_vm2"

# start load balancer
sudo docker run -d --net="host" --privileged load-balancer
echo "Started load balancer"