variable "name" {}
variable "vpc_ip_range" {}
variable "regions" {
  default = []
}

// Create VPC
resource "google_compute_network" "vpc" {
 name                    = "${var.name}-vpc"
 auto_create_subnetworks = "false"
}

// Create Subnet
resource "google_compute_subnetwork" "subnets" {
  count = "${length(var.regions)}"
  region = "${element(var.regions, count.index)}"
  name          = "${var.name}-subnet-${count.index}"
  ip_cidr_range = "${cidrsubnet(var.vpc_ip_range, 4, count.index)}"
  network       = "${var.name}-vpc"
  depends_on    = ["google_compute_network.vpc"]
}

resource "google_compute_route" "internet-route" {
  name        = "${var.name}-internet-route"
  dest_range  = "0.0.0.0/0"
  network     = "${google_compute_network.vpc.self_link}"
  priority    = 100
  tags = ["${var.name}"]
  next_hop_gateway = "global/gateways/default-internet-gateway"
}

// VPC firewall configuration
resource "google_compute_firewall" "firewall" {
  name    = "${var.name}-firewall"
  network = "${google_compute_network.vpc.name}"
  target_tags = ["${var.name}"]
  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "tcp"
    ports    = ["22","3389","80","8080","110","26600-26699"]
  }
  source_ranges = ["0.0.0.0/0"]
}

resource "google_compute_firewall" "internal-firewall" {
  name    = "${var.name}-internal-firewall"
  network = "${google_compute_network.vpc.name}"
  target_tags = ["${var.name}"]
  allow {
    protocol = "icmp"
  }
  allow {
    protocol = "udp"
    ports    = ["0-65535"]
  }
  allow {
    protocol = "tcp"
    ports    = ["0-65535"]
  }
  source_ranges = ["${var.vpc_ip_range}"]
}

output "subnets" {
  value = "${google_compute_subnetwork.subnets.*.self_link}"
}

output "vpc_network" {
  value = "${google_compute_network.vpc.self_link}"
}