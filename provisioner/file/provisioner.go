//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package lxdapi

import (
	"context"

	utils "packer-plugin-lxdapi/utils"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
)

type Config struct {
	ctx        interpolate.Context
	Source string `mapstructure:"source" required:"true"`
	Destination string `mapstructure:"destination" required:"true"`
}

type FileProvisioner struct {
	config Config
}

func (p *FileProvisioner) ConfigSpec() hcldec.ObjectSpec {
	return p.config.FlatMapstructure().HCL2Spec()
}

func (p *FileProvisioner) Prepare(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		PluginType:         "packer.provisioner.lxdapi.file",
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

func (p *FileProvisioner) Provision(_ context.Context, ui packer.Ui, _ packer.Communicator, generatedData map[string]interface{}) error {
	ui.Say("=================================================")
	ui.Say(" Running FileProvisioner.Provision()")
	ui.Say("=================================================")
	instanceName := generatedData["InstanceName"].(string)
	unixSocketPath := generatedData["UnixSocketPath"].(string)

	u, err := utils.NewLXDUtilStruct(unixSocketPath)
	if err != nil {
		return err
	}

	err = u.UploadFile(instanceName, p.config.Source, p.config.Destination)
	if err != nil {
		return err
	}

	return nil
}
