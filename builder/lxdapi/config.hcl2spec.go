// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package lxdapi

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	PackerBuildName     *string           `mapstructure:"packer_build_name" cty:"packer_build_name" hcl:"packer_build_name"`
	PackerBuilderType   *string           `mapstructure:"packer_builder_type" cty:"packer_builder_type" hcl:"packer_builder_type"`
	PackerCoreVersion   *string           `mapstructure:"packer_core_version" cty:"packer_core_version" hcl:"packer_core_version"`
	PackerDebug         *bool             `mapstructure:"packer_debug" cty:"packer_debug" hcl:"packer_debug"`
	PackerForce         *bool             `mapstructure:"packer_force" cty:"packer_force" hcl:"packer_force"`
	PackerOnError       *string           `mapstructure:"packer_on_error" cty:"packer_on_error" hcl:"packer_on_error"`
	PackerUserVars      map[string]string `mapstructure:"packer_user_variables" cty:"packer_user_variables" hcl:"packer_user_variables"`
	PackerSensitiveVars []string          `mapstructure:"packer_sensitive_variables" cty:"packer_sensitive_variables" hcl:"packer_sensitive_variables"`
	Config              map[string]string `mapstructure:"config" required:"false" cty:"config" hcl:"config"`
	OutputImage         *string           `mapstructure:"output_image" required:"true" cty:"output_image" hcl:"output_image"`
	PublishProperties   map[string]string `mapstructure:"publish_properties" required:"false" cty:"publish_properties" hcl:"publish_properties"`
	SourceImage         *string           `mapstructure:"source_image" required:"true" cty:"source_image" hcl:"source_image"`
	VirtualMachine      *bool             `mapstructure:"virtual_machine" required:"true" cty:"virtual_machine" hcl:"virtual_machine"`
	UnixSocketPath      *string           `mapstructure:"unix_socket_path" required:"false" cty:"unix_socket_path" hcl:"unix_socket_path"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"packer_build_name":          &hcldec.AttrSpec{Name: "packer_build_name", Type: cty.String, Required: false},
		"packer_builder_type":        &hcldec.AttrSpec{Name: "packer_builder_type", Type: cty.String, Required: false},
		"packer_core_version":        &hcldec.AttrSpec{Name: "packer_core_version", Type: cty.String, Required: false},
		"packer_debug":               &hcldec.AttrSpec{Name: "packer_debug", Type: cty.Bool, Required: false},
		"packer_force":               &hcldec.AttrSpec{Name: "packer_force", Type: cty.Bool, Required: false},
		"packer_on_error":            &hcldec.AttrSpec{Name: "packer_on_error", Type: cty.String, Required: false},
		"packer_user_variables":      &hcldec.AttrSpec{Name: "packer_user_variables", Type: cty.Map(cty.String), Required: false},
		"packer_sensitive_variables": &hcldec.AttrSpec{Name: "packer_sensitive_variables", Type: cty.List(cty.String), Required: false},
		"config":                     &hcldec.AttrSpec{Name: "config", Type: cty.Map(cty.String), Required: false},
		"output_image":               &hcldec.AttrSpec{Name: "output_image", Type: cty.String, Required: false},
		"publish_properties":         &hcldec.AttrSpec{Name: "publish_properties", Type: cty.Map(cty.String), Required: false},
		"source_image":               &hcldec.AttrSpec{Name: "source_image", Type: cty.String, Required: false},
		"virtual_machine":            &hcldec.AttrSpec{Name: "virtual_machine", Type: cty.Bool, Required: false},
		"unix_socket_path":           &hcldec.AttrSpec{Name: "unix_socket_path", Type: cty.String, Required: false},
	}
	return s
}
