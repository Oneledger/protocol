provider "google-beta" {
  alias = "test-region"
  credentials = "${file("Chronos.json")}"
  project     = "chronos-225820"
  region      = "us-east1"
  zone = "us-east1-b"
}

provider "google" {
  alias = "test-region2"
  credentials = "${file("DevNet.json")}"
  project     = "atomic-land-223022"
  region      = "us-west1"
  zone = "us-west1-b"
}

module "node" {
  source  = "./nfs"
  providers = {
    google = "google-beta.test-region"
  }
}

output "public_ip" {
  value = "${module.node.public_ip}"
}