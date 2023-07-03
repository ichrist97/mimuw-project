# Setup of deployment

This documents explains all the steps for an automated deployment using Ansible.

## Steps

Create password file for ansible:

```bash
cd deployment/ansible
echo "mypassword" > password
```

Initialize ansible control center node:

```bash
bash init.sh
```

Check that ansible works:

```bash
ansible docker -m ping  --extra-vars "ansible_user=<username> ansible_password=<password>"
```

Install docker on all machines:

```bash
cd docker
ansible-playbook --extra-vars "ansible_user=<username> ansible_password=<password>" docker-playbook.yml
```

Inspired by this [guide](https://www.digitalocean.com/community/tutorials/how-to-use-sharding-in-mongodb), the following playbook tries to setup a sharded MongoDB cluster with five shards, one config server and one mongos router server. The playbooks fails at the part of initializing the shard and config servers in the mongo shell, because of the connection to the mongo server cannot be made. Although, the connection can be directly replicated on the machine itself. Because of this issue, the automated playbook does not work fully in the cluster has to deployed partially in a manual way. In this sense, the MongoDB cluster is the only component that is not fully automated.

Start sharded mongodb cluster:

```bash
cd ../mongodb
ansible-playbook --extra-vars "ansible_user=<username> ansible_password=<password>" mongodb-playbook.yml
```

Start loadbalancer:

```bash
cd ../loadbalancer
ansible-playbook --extra-vars "ansible_user=<username> ansible_password=<password>" loadbalancer-playbook.yml
```

Start two instances of the tag service:

```bash
cd ../services
ansible-playbook --extra-vars "ansible_user=<username> ansible_password=<password>" services-playbook.yml
```
