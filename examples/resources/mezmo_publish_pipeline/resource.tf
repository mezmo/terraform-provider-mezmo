terraform {
  required_providers {
    mezmo = {
      source = "registry.terraform.io/mezmo/mezmo"
    }
  }
}

provider "mezmo" {
  auth_key = "my key"
}

module "pipeline" {
  source = "./modules"
}

resource "mezmo_publish_pipeline" "publisher" {
  pipeline_id = module.pipeline.my_pipeline.id
  updated_at  = module.pipeline.my_pipeline.updated_at
  depends_on  = [module.pipeline]
}
