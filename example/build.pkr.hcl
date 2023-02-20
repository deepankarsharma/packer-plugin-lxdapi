packer {
  required_plugins {
    lxdapi = {
      source  = "github.com/deepankarsharma/lxdapi"
      version = ">= 0.0.1"
    }
  }
}

source "lxdapi" "vm" {
  mock = local.foo
  virtual_machine = true
  image        = "images:rockylinux/8/cloud/amd64"
  output_image = "rocky8-lxdapi-phase0"
  publish_properties = {
    description = "Rocky Linux 8 LXD API Phase 0"
  }
}

build {
  sources = [
    "source.lxdapi.vm",
  ]

  source "source.scaffolding-my-builder.bar-example" {
    name = "bar"
  }
}
