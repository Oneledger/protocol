terraform {
  backend "s3" {
    region =  "us-east-1"
    bucket = "terraform-oneledger"
    key =  "devnet"
  }
}

provider "google-beta" {
  alias = "chronos"
  credentials = "${file("Chronos.json")}"
  project     = "chronos-225820"
  region      = "us-east1"
  zone = "us-east1-b"
}

provider "google" {
  alias = "devnet"
  credentials = "${file("DevNet.json")}"
  project     = "atomic-land-223022"
  region      = "us-west1"
}

module "network"{
  source = "./network"
  vpc_ip_range = "${var.vpc_ip_range}"
  name = "${var.name}"
  regions = "${var.regions}"
  providers = {
    google = "google.devnet"
  }
}

module "node"{
  source = "./fullnode"
  name = "${var.name}"
  vmcount = "${var.vmcount}"
  subnets = "${module.network.subnets}"
  regions = "${var.regions}"
  providers = {
    google = "google.devnet"
  }
}

output "public_ip" {
  value = "${module.node.public_ip}"
}