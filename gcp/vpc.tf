provider "google" {
  credentials = "${file("DevNet.json")}"
  project     = "atomic-land-223022"
  region      = "northamerica-northeast1"
}

resource "google_compute_instance" "default" {
  name         = "test"
  machine_type = "n1-standard-1"
  zone         = "northamerica-northeast1-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
    }
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral IP
    }
  }

}