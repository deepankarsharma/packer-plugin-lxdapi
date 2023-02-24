// packer {
//   required_plugins {
//     lxdapi = {
//       source  = "github.com/deepankarsharma/lxdapi"
//       version = "0.0.4"
//     }
//   }

// }

source "lxdapi" "instance" {
  unix_socket_path = "/var/snap/lxd/common/lxd/unix.socket"
  source_image     = "jammy-amd64"
  output_image = "jammy-output"
  output_image_description = "Jammy container image"
  publish_properties = {
    description = "Jammy container image"
  }
  config = {
    "security.secureboot": "false"
  }
  virtual_machine = false
  compression_algorithm = "zstd"
}

build {
  sources = [
    "source.lxdapi.instance",
  ]

  provisioner "lxdapi-shell" {
    environment = {
      "HELLO": "WORLD"
    }

    inline = [
      "env",
      "echo $HELLO",
    ]
  }
}
