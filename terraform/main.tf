terraform {
  required_version = ">= 1.5"
  required_providers {
    libvirt = {
      source  = "dmacvicar/libvirt"
      version = "~> 0.8.3"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.8"
    }
  }
}

provider "libvirt" {
  uri = "qemu:///system"
}

variable "ssh_key" {
  type        = string
  description = "Caminho da chave SSH pública autorizada na VM"
  default     = "~/.ssh/id_ed25519.pub"
}

variable "ssh_user" {
  type        = string
  description = "Usuário usado na VM"
  default     = "debian"
}

variable "disk_size_gb" {
  type        = number
  description = "Tamanho do disco em GiB"
  default     = 20
}

locals {
  ip             = libvirt_domain.korp.network_interface[0].addresses[0]
  gb             = pow(1024, 3)
  authorized_key = trimspace(file(pathexpand(var.ssh_key)))
}


resource "libvirt_pool" "korp" {
  name = "korp"
  type = "dir"
  target {
    path = "/var/lib/libvirt/images/korp"
  }
}

resource "libvirt_network" "korp" {
  name      = "korp"
  mode      = "nat"
  addresses = ["10.17.3.0/24"]
  autostart = true
  dhcp {
    enabled = true
  }
}

resource "libvirt_volume" "debian" {
  name   = "debian-13.qcow2"
  pool   = libvirt_pool.korp.name
  source = "https://cloud.debian.org/images/cloud/trixie/latest/debian-13-genericcloud-amd64.qcow2"
  format = "qcow2"
}

resource "libvirt_volume" "disk" {
  name           = "korp.qcow2"
  pool           = libvirt_pool.korp.name
  base_volume_id = libvirt_volume.debian.id
  size           = var.disk_size_gb * local.gb
}

resource "libvirt_cloudinit_disk" "init" {
  name      = "korp-init.iso"
  pool      = libvirt_pool.korp.name
  user_data = <<-EOT
    #cloud-config
    hostname: korp
    ssh_authorized_keys:
      - ${local.authorized_key}
  EOT
}

resource "libvirt_domain" "korp" {
  name      = "korp"
  memory    = 4096
  vcpu      = 2
  cloudinit = libvirt_cloudinit_disk.init.id

  disk {
    volume_id = libvirt_volume.disk.id
  }

  network_interface {
    network_id     = libvirt_network.korp.id
    wait_for_lease = true
  }

  console {
    type        = "pty"
    target_type = "serial"
    target_port = "0"
  }
}

resource "local_file" "inventory" {
  filename        = "${path.module}/inventory.ini"
  file_permission = "0644"
  content         = "[korp]\n${local.ip} ansible_user=${var.ssh_user}\n"
}

output "ip" {
  value = local.ip
}
