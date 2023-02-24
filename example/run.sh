#!/bin/bash

pushd ~/packer-plugin-lxdapi
# Check if calling make dev is successful, if so print success
# else print error and exit
make dev || (echo "Error: make dev failed" && exit 1)
popd

/usr/bin/packer build -force ./build.pkr.hcl
