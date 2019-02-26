resource "google_compute_instance" "default" {
  name         = "explorer-fullnode"
  machine_type = "n1-standard-1"
  tags = ["p2p"]
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
      size = 100
    }
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }
}

output "public_ip" {
  value = "${google_compute_instance.default.network_interface.0.access_config.0.nat_ip}"
}