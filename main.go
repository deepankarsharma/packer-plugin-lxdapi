package main

import (
	"fmt"
	"os"
	"packer-plugin-lxdapi/builder/lxdapi"
	fileProvisioner "packer-plugin-lxdapi/provisioner/file"
	shellProvisioner "packer-plugin-lxdapi/provisioner/shell"
	lxdapiVersion "packer-plugin-lxdapi/version"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder(plugin.DEFAULT_NAME, new(lxdapi.Builder))
	pps.RegisterProvisioner("file", new(fileProvisioner.FileProvisioner))
	pps.RegisterProvisioner("shell", new(shellProvisioner.ShellProvisioner))
	pps.SetVersion(lxdapiVersion.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
