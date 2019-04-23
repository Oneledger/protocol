terraform {
  backend "s3" {
    region =  "us-east-1"
    bucket = "terraform-oneledger"
  }
}

provider "google" {
  alias = "chronos"
  credentials = "${file("../gcp/Chronos.json")}"
  project     = "chronos-225820"
  region      = "us-east1"
  zone = "us-east1-b"
}

provider "google" {
  alias = "devnet"
  credentials = "${file("../gcp/DevNet.json")}"
  project     = "atomic-land-223022"
  region      = "us-west1"
}

module "network"{
  source = "../gcp/network"
  vpc_ip_range = "${var.vpc_ip_range}"
  name = "${var.name}"
  regions = "${var.regions}"
  providers = {
    google = "google.devnet"
  }
}

# module "node"{
#   source = "../gcp/fullnode"
#   name = "${var.name}"
#   vmcount = "${var.vmcount}"
#   gcp-vm-image= "${var.gcp-vm-image}"
#   subnets = "${module.network.subnets}"
#   regions = "${var.regions}"
#   providers = {
#     google = "google.devnet"
#   }
# }
resource "google_compute_instance" "default" {
  provider = "google.devnet"
  count = "${var.vmcount}"
  name = "${var.name}-vm-${count.index}"
  machine_type = "n1-standard-1"
  tags = ["${var.name}"]
  zone = "${element(var.regions,count.index % length(var.regions))}-b"
  boot_disk {
    initialize_params {
      image = "${var.gcp-vm-image}"
      size = 100
    }
  }
  allow_stopping_for_update = true
  network_interface {
    subnetwork = "${element(module.network.subnets,count.index % length(module.network.subnets))}"
    access_config {}
  }
}
