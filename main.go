package main

import (
	"fmt"
	"os"
	"packer-plugin-lxdapi/builder/lxdapi"
	lxdapiVersion "packer-plugin-lxdapi/version"
	"github.com/hashicorp/packer-plugin-sdk/plugin"
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterBuilder("builder", new(lxdapi.Builder))
	pps.SetVersion(lxdapiVersion.PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
