
variable "name" {
    default = "devnet"
}
variable "subnet" {}
variable "vpc_network" {}
resource "google_compute_global_address" "tcp-proxy-lb-ip" {
  name = "${var.name}-ipv4-address"
}

resource "google_compute_instance" "default" {
  name         = "${var.name}-vm"
  machine_type = "n1-standard-1"
  tags = ["${var.name}"]
  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-9"
      size = 100
    }
  }

  network_interface {
    subnetwork = "${var.subnet}"
    access_config {
    }
  }
}

resource "google_compute_health_check" "internal-health-check" {
 name = "${var.name}-health-check"

 timeout_sec        = 1
 check_interval_sec = 1

 tcp_health_check {
   port = "80"
 }
}

resource "google_compute_instance_group" "instance_group" {
  name        = "${var.name}-instance-group"
  network     = "${var.vpc_network}"
  instances = [
    "${google_compute_instance.default.self_link}",
  ]
  named_port {
    name = "tcp"
    port = "110"
  }
}

resource "google_compute_backend_service" "app" {
  name        = "${var.name}-backend"
  port_name   = "tcp"
  protocol    = "TCP"
  timeout_sec = 10
  enable_cdn  = false

  backend {
    group = "${google_compute_instance_group.instance_group.self_link}"
  }

  health_checks = ["${google_compute_health_check.internal-health-check.self_link}"]
}

resource "google_compute_target_tcp_proxy" "tcp_proxy_lb" {
  name            = "${var.name}-proxy"
  backend_service = "${google_compute_backend_service.app.self_link}"
}

resource "google_compute_firewall" "firewall" {
  name    = "${var.name}-allow-tcp-lb-and-health"
  network = "${var.vpc_network}"
  target_tags = ["${var.name}"]
  allow {
    protocol = "icmp"
  }

  allow {
    protocol = "tcp"
    ports    = ["80","8080"]
  }

  source_ranges = ["0.0.0.0/0"]
  priority = 999
}

resource "google_compute_global_forwarding_rule" "tcp-lb-forwarding-rule" {
  name       = "${var.name}-forwarding-rule"
  ip_protocol = "TCP"
  ip_address = "${google_compute_global_address.tcp-proxy-lb-ip.address}"
  target     = "${google_compute_target_tcp_proxy.tcp_proxy_lb.self_link}"
  port_range = "110"
}

output "public_ip"{
    value = "${google_compute_global_forwarding_rule.tcp-lb-forwarding-rule.ip_address}"
}