terraform {
  backend "s3" {
    region =  "us-east-1"
    bucket = "terraform-oneledger"
  }
}

provider "google" {
  alias = "chronos"
  credentials = "${file("../Chronos.json")}"
  project     = "chronos-225820"
  region      = "us-east1"
  zone = "us-east1-b"
}

provider "google" {
  alias = "devnet"
  credentials = "${file("../DevNet.json")}"
  project     = "atomic-land-223022"
  region      = "us-west1"
}

module "network"{
  source = "../modules/network"
  vpc_ip_range = "${var.vpc_ip_range}"
  name = "${var.name}"
  regions = "${var.regions}"
  providers = {
    google = "google.devnet"
  }
}

module "node"{
  source = "../modules/fullnode"
  name = "${var.name}"
  vmcount = "${var.vmcount}"
  subnets = "${module.network.subnets}"
  regions = "${var.regions}"
  vm_machine_type = "${var.vm_machine_type}"
  providers = {
    google = "google.devnet"
  }
}

output "public_ip" {
  value = "${module.node.public_ip}"
}