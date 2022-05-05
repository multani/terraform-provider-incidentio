variable "incidentio_api_key" {}

terraform {
  required_providers {
    incidentio = {
      source = "multani/incidentio"
    }
  }
}

provider "incidentio" {
  api_key = var.incidentio_api_key
}
