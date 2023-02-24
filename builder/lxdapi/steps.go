// https://linuxcontainers.org/lxd/docs/master/api/#/instances/instance_put
// https://pkg.go.dev/github.com/lxc/lxd@v0.0.0-20230223142449-e78b7f2b47d0/shared/api
// lxc image copy images:ubuntu/jammy/amd64 local: --copy-aliases --auto-update --alias jammy-amd64
package lxdapi
import (
	"context"
	"fmt"
	utils "packer-plugin-lxdapi/utils"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	sdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/lxc/lxd/shared/api"
)

type stepLaunch struct {
	u* utils.LXDUtils
	instanceName string
}

const lxdConfigKey = "user._dt"

func (s *stepLaunch) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(sdk.Ui)
	ui.Say("Launching container...");
	config := state.Get("config").(*Config)
	
	ui.Say("Instance name: " + state.Get("instanceName").(string))
	s.instanceName = state.Get("instanceName").(string)

	ui.Say("Config: " + config.SourceImage)

	// Connect to LXD over the Unix socket
	new_u, err := utils.NewLXDUtilStruct(config.UnixSocketPath)
	if err != nil {
		ui.Error("Error connecting to LXD: " + err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}
	s.u = new_u

	instanceType := api.InstanceTypeContainer
	if config.VirtualMachine {
		instanceType = api.InstanceTypeVM
	}

	ui.Say("Connected to LXD")
	req := api.InstancesPost{
		Type: instanceType,
		Name: s.instanceName,
		Source: api.InstanceSource{
			Type:        "image",
			Alias: config.SourceImage,
		},
		InstancePut: api.InstancePut{
			Ephemeral: false,
			Profiles:  []string{"default"},
			Config: map[string]string{
				lxdConfigKey: fmt.Sprintf("%v", config.Config),
			},
		},
	}

	ui.Say("Creating instance...")
	err = s.u.CreateInstanceHL(req)
	if err != nil {
		ui.Error("Error creating instance: " + err.Error())
		state.Put("error", err)
		return multistep.ActionHalt
	}

	return multistep.ActionContinue
}

func (s *stepLaunch) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(sdk.Ui)
	ui.Say("Unregistering and deleting container...")
	
	err := s.u.DeleteInstance(s.instanceName)
	if err != nil {
		ui.Error("Error deleting instance: " + err.Error())
		state.Put("error", err)
	}

}

type stepProvision struct {}

func (s *stepProvision) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(sdk.Ui)
	ui.Say("Provisioning container...");
	return multistep.ActionContinue
}

func (s *stepProvision) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(sdk.Ui)
	ui.Say("Running stepProvision.Cleanup ...")
}



type stepPublish struct {}

func (s *stepPublish) Run(ctx context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(sdk.Ui)
	ui.Say("Publishing container...");
	return multistep.ActionContinue
}

func (s *stepPublish) Cleanup(state multistep.StateBag) {
	ui := state.Get("ui").(sdk.Ui)
	ui.Say("Running stepPublish.Cleanup ...")
}