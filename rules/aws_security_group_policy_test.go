package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecurityGroupPolicy(t *testing.T) {
	rule := &AwsSecurityGroupPolicy{}

	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name:     "no issues when no security groups",
			Content:  "",
			Expected: helper.Issues{},
		},
		{
			Name: "no issues outside of dev/prod accounts",
			Content: `
			variable "environment" { default = "other"}
			resource "aws_security_group" "sg" {}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues when valid tag",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			resource "aws_security_group" "sg" {
				tags = {
					owner="teamA"
				}
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "tags missing",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			resource "aws_security_group" "sg" {
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Security group is missing the required "owner" tag.`,
				},
			},
		},
		{
			Name: "tags not statically evaluable",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "security_group_tags" {}
			resource "aws_security_group" "sg" {
				tags = var.security_group_tags
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Security group tags could not be fully evaluated; ensure the required "owner" tag is statically resolvable.`,
				},
			},
		},
		{
			Name: "owner tag missing",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			resource "aws_security_group" "sg" {
				tags = {
					other="tag"
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Security group is missing the required "owner" tag.`,
				},
			},
		},
		{
			Name: "owner tag empty",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			resource "aws_security_group" "sg" {
				tags = {
					owner=""
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Security group has an empty required "owner" tag.`,
				},
			},
		},
		{
			Name: "owner tag invalid",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			resource "aws_security_group" "sg" {
				tags = {
					owner="invalid"
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "Security group has an invalid \"owner\" tag; Found \"invalid\", but it has to be one of: [teamA teamB]",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			if err := rule.Check(runner); err != nil {
				t.Fatalf("Unexpected error occurred: %s", err)
			}

			helper.AssertIssuesWithoutRange(t, test.Expected, runner.Issues)
		})
	}
}
