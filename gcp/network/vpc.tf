variable "name" {
    default = "devnet"
}
variable "subnet_cidr" {}

// Create VPC
resource "google_compute_network" "vpc" {
 name                    = "${var.name}-vpc"
 auto_create_subnetworks = "false"
}

// Create Subnet
resource "google_compute_subnetwork" "subnet" {
 name          = "${var.name}-subnet"
 ip_cidr_range = "${var.subnet_cidr}"
 network       = "${var.name}-vpc"
 depends_on    = ["google_compute_network.vpc"]
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
    ports    = ["22","3389","80","8080","110"]
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
  source_ranges = ["${var.subnet_cidr}"]
}

output "subnet" {
  value = "${google_compute_subnetwork.subnet.self_link}"
}

output "vpc_network" {
  value = "${google_compute_network.vpc.self_link}"
}