terraform {
  required_providers {
    mezmo = {
      source = "registry.terraform.io/mezmo/mezmo"
    }
  }
}

output "my_pipeline" {
  value = mezmo_pipeline.my_pipeline
}

resource "mezmo_pipeline" "my_pipeline" {
  title = "A pipeline to publish when there are changes"
}
