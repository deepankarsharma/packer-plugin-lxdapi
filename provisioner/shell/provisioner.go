//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package lxdapi

import (
	"context"
	"fmt"
	//utils "packer-plugin-lxdapi/utils"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
	//lxd "github.com/lxc/lxd/client"
)

type Config struct {
	ctx        interpolate.Context
	Environment map[string]string `mapstructure:"environment" required:"false"`
	Inline      []string          `mapstructure:"inline" required:"true"`
}

type ShellProvisioner struct {
	config Config
}

func (p *ShellProvisioner) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}


func (p *ShellProvisioner) Prepare(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		PluginType:         "packer.provisioner.lxdapi.shell",
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
		}, raws...)
		if err != nil {
			return err
		}
		return nil
	}
	
func (p *ShellProvisioner) Provision(_ context.Context, ui packer.Ui, _ packer.Communicator, generatedData map[string]interface{}) error {
	ui.Say("Hello from the LXD API Shell Provisioner")
	// print number of keys in generatedData
	ui.Say("Number of keys in generatedData: " + fmt.Sprintf("%v", len(generatedData)))
	// print keys of generatedData
	for k := range generatedData {
		ui.Say("Key: " + k)
	}
	//instanceName := generatedData["InstanceName"].(string)
	//ui.Say("Instance server: " + fmt.Sprintf("%v", instanceName))
	return nil
}
	