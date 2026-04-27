terraform {
  required_providers {
    mezmo = {
      source  = "mezmo/mezmo"
      version = "~> 6.0.0"
    }
  }
}

provider "mezmo" {
  auth_key = "my secret"
}
