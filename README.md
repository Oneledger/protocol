# infrastructure
Deployment Scripts Using Ansible, Terraform and Packer.  
## Requirements

- ansible >= 2.7.9
- terraform >= 0.11.13
- packer >= 1.3.5 (only for baking docker & gcp image)
- docker >= 18.09.3 (only for baking docker & gcp image)
- GCP service account key to access DevNet project

## Instructions

Clone this repo:

```
$ git clone git@github.com:Oneledger/infrastructure.git
```
From the repo directory, place the gcp credential file(service account key) under `gcp` directory locally and rename it to `DevNet.json`. 

### Deploy blockchain cluster with validator and seed node 
From the repo directory, nevigate to gcp directory and initialize terraform project
```
$ cd gcp
$ terraform init
```
This will prompt user to enter the terraform state file path located in aws S3 bucket under **terraform-oneledger**. If you creating a new application deployment, make sure the value entered are different than existing deployment names such as **chronos** and **devnet**. 

From the repo directory, nevigate to ansible directory and run ansible playbook to deploy to DevNet project in GCP:
```
$ cd ansible
$ ansible-playbook main.yml --extra-vars "reset_network=true stage=devnet vpc_ip_range=10.10.0.0/20 vmcount=5 remote_user=steven"
```
The playbook above will deploy 4 validator node and 1 seed node. The deployment can be configured through `--extra-vars` specfied below:

* `reset_network`: takes in `boolean` value, will erase the data, configuration and log file if configure to `true`. Default value: `false`
* `vpc_ip_range`: the ip range assigned to the vpc in cidr notation. Default value: `10.10.0.0/20`. Make sure this value is unique per GCP project.   
* `stage`: the application deployment name. the value of this variable will be used as prefix for all the source generated inside vpc. (ex. devnet-vpc, devnet-firewall and devnet-vm-1) 
* `vmcount`: number of virtual machine provisioned in GCP compute. The first four VM will be assigned as validator and the rest will be assign and seed. Requirement: `>= 5`, default value: `5`. Example: if vmcount=6, 2 seed node will be deployed, if vmcount=7, 3 seed node will be deployed and etc... 
* `remote_user`: the remote user ansible will authenicated with to access GCP VMs. This user must match your local ssh private key that already has remote access to VMs allocated in GCP "DevNet" project 

### Create GCP VM and Docker images for OneLedger Fullnode
From the repo directory, run packer to create a docker image that contains fullnode application dependencies and configurations. 
```
packer build -only=docker-image -var "tag=develop" packer.json
```
The command above will create a docker image locally with REPOSITORY name **oneledger/chronos** and TAG value based on the cli input value of `tag` variable which is develop.

To create a GCP Compute VM image
```
packer build -only=gcp-vm-image -var "tag=develop" packer.json 
```
The command above will create a GCP VM image with name develop in DevNet project. Note: since this image is compiled in a GCP VM remotely, this process can take between 5 to 10 minutes depends on network latency. 

The following are the configurations the packer script can take in as input through `-var` flag: 
* `version`: the branch name or tag version of [protocol repo](https://github.com/Oneledger/protocol), packer will use specfied branch to compile binary executable such as `olfullnode`. Default: `develop`
* `tag`: the gcp image name or docker tag name. Note: don't confuse this with `version`. Default: `develop`
* `app_user`: the linux user that will execute the application. Default: `node`
* `enable_logging`: expecting boolean value. If true, this will install **logstash**. Default: `false`
