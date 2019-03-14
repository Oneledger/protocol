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
  credentials = "${file("/home/steven/git/infrastructure/gcp/DevNet.json")}"
  project     = "atomic-land-223022"
  region      = "us-west1"
  zone = "us-west1-b"
}

module "network"{
  source = "./network"
  subnet_cidr = "${var.subnet_cidr}"
  name = "${var.name}"
  providers = {
    google = "google.devnet"
  }
}

module "explorer"{
  source = "./fullnode"
  name = "${var.name}"
  subnet = "${module.network.subnet}"
  providers = {
    google = "google.devnet"
  }
}

output "public_ip" {
  value = "${module.explorer.public_ip}"
}