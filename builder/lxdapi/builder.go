package lxdapi

import (
	"context"

	"github.com/google/uuid"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/multistep"
	"github.com/hashicorp/packer-plugin-sdk/multistep/commonsteps"
	"github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
)

const BuilderId = "lxdapi.builder"



type Builder struct {
	config Config
	runner multistep.Runner
	instanceName string
}

func (b *Builder) ConfigSpec() hcldec.ObjectSpec { return b.config.FlatMapstructure().HCL2Spec() }

func (b *Builder) Prepare(raws ...interface{}) (generatedVars []string, warnings []string, err error) {
	err = config.Decode(&b.config, &config.DecodeOpts{
		PluginType:  "packer.builder.lxdapi",
		Interpolate: true,
	}, raws...)
	if err != nil {
		return nil, nil, err
	}
	b.instanceName = "lxdapi-builder-" + uuid.New().String()
	buildGeneratedData := []string{"InstanceName"}
	return buildGeneratedData, nil, nil
}

func (b *Builder) Run(ctx context.Context, ui packer.Ui, hook packer.Hook) (packer.Artifact, error) {
	steps := []multistep.Step{
		&stepLaunch{},
		&stepProvision{},
		&stepPublish{},
		new(commonsteps.StepProvision),
	}
	// Setup the state bag and initial state for the steps
	state := new(multistep.BasicStateBag)
	state.Put("config", &b.config)
	state.Put("hook", hook)
	state.Put("ui", ui)
	state.Put("instanceName", b.instanceName)

	// Set the value of the generated data that will become available to provisioners.
	// To share the data with post-processors, use the StateData in the artifact.
	state.Put("generated_data", map[string]interface{}{
		"InstanceName": b.instanceName,
	})

	// Run!
	b.runner = commonsteps.NewRunner(steps, b.config.PackerConfig, ui)
	b.runner.Run(ctx, state)

	// If there was an error, return that
	if err, ok := state.GetOk("error"); ok {
		return nil, err.(error)
	}

	artifact := &Artifact{
		// Add the builder generated data to the artifact StateData so that post-processors
		// can access them.
		StateData: map[string]interface{}{"generated_data": state.Get("generated_data")},
	}
	return artifact, nil
}
