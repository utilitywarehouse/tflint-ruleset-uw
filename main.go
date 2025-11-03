package main

import (
	"github.com/terraform-linters/tflint-plugin-sdk/plugin"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"

	"github.com/utilitywarehouse/tflint-ruleset-uw/rules"
)

// set by goreleaser at build time: https://goreleaser.com/cookbooks/using-main.version/
var version = "dev"

func main() {
	plugin.Serve(&plugin.ServeOpts{
		RuleSet: &tflint.BuiltinRuleSet{
			Name:    "template",
			Version: version,
			Rules: []tflint.Rule{
				&rules.AwsSecurityGroupMissingOwnerTagRule{},
			},
		},
	})
}
