//go:generate packer-sdc mapstructure-to-hcl2 -type Config

package lxdapi

import (
	"context"
	"fmt"
	utils "packer-plugin-lxdapi/utils"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/hashicorp/packer-plugin-sdk/template/interpolate"
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
	ui.Say("=================================================")
	ui.Say(" Running ShellProvisioner.Provision()")
	ui.Say("=================================================")
	ui.Say("Fetching InstanceName from generatedData")
	instanceName := generatedData["InstanceName"].(string)
	unixSocketPath := generatedData["UnixSocketPath"].(string)

	ui.Say("Environment variables:")
	for k, v := range p.config.Environment {
		ui.Say(fmt.Sprintf("Environment variable: %s=%s", k, v))
	}

	u, err := utils.NewLXDUtilStruct(unixSocketPath)
	if err != nil {
		return err
	}

	// iterate over the inline commands
	for _, command := range p.config.Inline {
		ui.Say(fmt.Sprintf("shell_exec: %s", command))
		exec_result, err := u.Exec(instanceName, command, p.config.Environment)
		if err != nil {
			return err
		}
		ui.Say(fmt.Sprintf("rcode: %v", exec_result.ReturnCode))
		ui.Say(fmt.Sprintf("stdout: %v", exec_result.Stdout))
		ui.Say(fmt.Sprintf("stderr: %v", exec_result.Stderr))
		if exec_result.ReturnCode != 0 {
			return fmt.Errorf("shell_exec returned non-zero exit code: %v", exec_result.ReturnCode)
		}
	}

	return nil
}
	