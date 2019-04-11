variable "name" {}
variable "subnets" {
  default = []
}
variable "vmcount" {}
variable "regions" {
  default = []
}
resource "google_compute_address" "static-ips"{
  count = "${var.vmcount}"
  name = "${var.name}-static-ip-${count.index}"
  region = "${element(var.regions,count.index % length(var.regions))}"
}

resource "google_compute_instance" "default" {
  count = "${var.vmcount}"
  name = "${var.name}-vm-${count.index}"
  machine_type = "n1-standard-2"
  tags = ["${var.name}"]
  zone = "${element(var.regions,count.index % length(var.regions))}-b"
  boot_disk {
    initialize_params {
      image = "ubuntu-os-cloud/ubuntu-1604-lts"
      size = 100
    }
  }
  allow_stopping_for_update = true
  network_interface {
    subnetwork = "${element(var.subnets,count.index % length(var.subnets))}"
    access_config {
      nat_ip = "${element(google_compute_address.static-ips.*.address,count.index)}"
    }
  }
}

output "public_ip"{
  value = "${google_compute_address.static-ips.*.address}"
}