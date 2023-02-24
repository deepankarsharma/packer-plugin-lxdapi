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
  source_image        = "jammy-amd64"
  output_image = "jammy-output"
  publish_properties = {
    description = "Jammy container image"
  }
  config = {
    "security.secureboot": "false"
  }
  virtual_machine = false
}

build {
  sources = [
    "source.lxdapi.instance",
  ]
  
  provisioner "lxdapi-file" {
    source = "foo.txt"
    destination = "/tmp/app.tar.gz"
  }

  provisioner "lxdapi-shell" {
    environment = {
      "FOO" = "bar"
    }

    inline = [
      "sudo dnf install -y git",
      "git clone"
    ]
  }
}
