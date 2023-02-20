//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package lxdapi
import (
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/common"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/mitchellh/mapstructure"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`
	LaunchConfig map[string]string `mapstructure:"launch_config" required:"false"`
	OutputImage   string `mapstructure:"output_image" required:"true"`
	PublishProperties map[string]string `mapstructure:"publish_properties" required:"false"`
	SourceImage   string `mapstructure:"source_image" required:"true"`
	VirtualMachine bool `mapstructure:"virtual_machine"`
}

func (c *Config) Prepare(raws ...interface{}) error {

	var md mapstructure.Metadata
	err := config.Decode(c, &config.DecodeOpts{
		Metadata:    &md,
		Interpolate: true,
	}, raws...)
	if err != nil {
		return err
	}

	// Accumulate any errors
	var errs *packersdk.MultiError

	if c.SourceImage == "" {
		errs = packersdk.MultiErrorAppend(errs, fmt.Errorf("`image` is a required parameter for LXD. Please specify an image by alias or fingerprint. e.g. `ubuntu-daily:x`"))
	}

	if errs != nil && len(errs.Errors) > 0 {
		return errs
	}

	return nil
}