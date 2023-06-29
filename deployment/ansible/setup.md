# Setup of deployment

This documents explains all the steps for an automated deployment using Ansible.

## Steps

Initialize ansible control center node:

```bash
cd deployment/ansible
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

Start services:

```bash
cd ../services
ansible-playbook --extra-vars "ansible_user=<username> ansible_password=<password>" services-playbook.yml
```
