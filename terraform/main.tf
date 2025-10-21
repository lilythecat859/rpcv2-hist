# AGPL-3.0
terraform {
  required_version = ">= 1.5"
  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }
}

variable "project_id" {
  description = "GCP project"
  type        = string
}

variable "region" {
  description = "GCP region"
  type        = string
  default     = "us-central1"
}

variable "zone" {
  description = "GCP zone"
  type        = string
  default     = "us-central1-a"
}

variable "machine_type" {
  description = "AMD 16-core"
  type        = string
  default     = "n2d-standard-16"
}

resource "google_compute_instance" "rpcv2_hist" {
  name         = "rpcv2-hist"
  machine_type = var.machine_type
  zone         = var.zone

  boot_disk {
    initialize_params {
      size  = 200
      type  = "pd-ssd"
      image = "ubuntu-2204-jammy-v20240501"
    }
  }

network_interface {
    network = "default"
    access_config {}
  }

  metadata_startup_script = <<-EOF
    #!/bin/bash
    apt-get update && apt-get install -y docker.io
    usermod -aG docker ubuntu
    curl -L https://github.com/docker/compose/releases/latest/download/docker-compose-linux-x86_64 -o /usr/local/bin/docker-compose
    chmod +x /usr/local/bin/docker-compose
    git clone https://github.com/faithful-rpc/rpcv2-hist /opt/rpcv2-hist
    cd /opt/rpcv2-hist
    docker compose up -d
  EOF

  service_account {
    scopes = ["cloud-platform"]
  }
}