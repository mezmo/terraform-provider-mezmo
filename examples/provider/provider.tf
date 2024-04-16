terraform {
  required_providers {
    mezmo = {
      source  = "mezmo/mezmo"
      version = "~> 1.0"
    }
  }
}

provider "mezmo" {
  auth_key = "my secret"
}
