terraform {
  required_providers {
    mezmo = {
      source = "registry.terraform.io/mezmo/mezmo"
    }
  }
  required_version = ">= 1.1.0"
}

provider "mezmo" {
  auth_key = "my secret"
}

resource "mezmo_pipeline" "pipeline1" {
  title = "My pipeline"
}

resource "mezmo_demo_source" "source1" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "My source"
  description = "This is the point of entry for our data"
  format      = "nginx"
}

resource "mezmo_http_destination" "webhook" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Webhook"
  description = "This URL is an API that acts as a webhook"
  uri         = "https://example.org"
  inputs      = [mezmo_demo_source.source1.id]
}

resource "mezmo_http_destination" "storage" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Storage"
  description = "This is an API that stores a copy of the data on disk"
  uri         = "https://example.org"
  inputs      = [mezmo_demo_source.source1.id]
  auth = {
    strategy = "bearer"
    token    = "<shhh secret token>"
  }
}

resource "mezmo_http_destination" "some_api" {
  pipeline_id = mezmo_pipeline.pipeline1.id
  title       = "Auth endpoint"
  description = "This API is authenticated with username/password"
  uri         = "https://example.org"
  inputs      = [mezmo_demo_source.source1.id]
  auth = {
    strategy = "basic"
    user     = "guest"
    password = "abc123"
  }
  headers = {
    x-some-header = "some header value"
  }
  max_bytes      = 5000
  timeout_secs   = 600
  method         = "post"
  payload_prefix = "{\"extra_prop\": true"
  payload_suffix = "\"extra_prop2\": true }"
  tls_protocols  = ["TLSPPP"]
  proxy = {
    enabled            = true
    endpoint           = "http://myproxy.com"
    hosts_bypass_proxy = ["0.0.0.0", "1.1.1.1"]
  }
  rate_limiting = {
    request_limit = 600
    duration_secs = 900
  }
}
