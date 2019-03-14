variable "name" {
    default = "devnet"
}
variable "subnet" {}
variable "vmcount" {
  default = 5
}

resource "google_compute_address" "static-ips"{
  count = "${var.vmcount}"
  name = "static-ip-${count.index}"
}

resource "google_compute_instance" "default" {
  count = "${var.vmcount}"
  name         = "${var.name}-vm-${count.index}"
  machine_type = "n1-highcpu-8"
  tags = ["${var.name}"]
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
      size = 25
    }
  }
  allow_stopping_for_update = true
  network_interface {
    subnetwork = "${var.subnet}"
    access_config {
      nat_ip = "${element(google_compute_address.static-ips.*.address,count.index)}"
    }
  }
}

output "public_ip"{
  value = "${google_compute_address.static-ips.*.address}"
}