variable "name" {
    default = "devnet"
}
variable "gcp-vm-image" {
  default = "debian-cloud/debian-9"
}
variable "vpc_ip_range" {
  default = "10.10.0.0/20"
}
variable "regions" {
  default = ["us-east-1", "us-east-4", "us-central1"]
}
variable "vmcount" {
  default = 5
}
variable "startup-script-url" {
  default = ""
}
variable "vm_machine_type" {
  default = "n1-standard-1"
}

