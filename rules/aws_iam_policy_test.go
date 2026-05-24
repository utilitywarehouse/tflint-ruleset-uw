package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsIamPolicy(t *testing.T) {
	rule := &AwsIamPolicy{}

	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name:     "no issues when no IAM resources",
			Content:  "",
			Expected: helper.Issues{},
		},
		{
			Name: "no issues outside of dev/prod accounts",
			Content: `
			variable "environment" { default = "other" }
			resource "aws_iam_user" "user" { name = "invalid" }`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues when all IAM resources use valid name prefixes",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }

			resource "aws_iam_user" "user" {
				name                 = "teamA-admin"
				permissions_boundary = "sys-teamA-boundary"
			}
			resource "aws_iam_role" "role" {
				name                 = "tb-ci-role"
				permissions_boundary = "sys-teamB-boundary"
			}
			resource "aws_iam_group" "group" { name = "teamB-readonly" }
			resource "aws_iam_policy" "policy" { name   = "ta-read" }
			resource "aws_iam_instance_profile" "profile" { name = "teamA-profile" } `,
			Expected: helper.Issues{},
		},
		{
			Name: "user name invalid",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }

			resource "aws_iam_user" "user" {
				name                 = "invalid"
				permissions_boundary = "sys-teamA-boundary"
			}
			`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `user "invalid" is invalid; it must start with "$team-", with team being one of: [teamA teamB ta tb]`,
				},
			},
		},
		{
			Name: "role name invalid",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }

			resource "aws_iam_role" "role" {
				name                 = "invalid"
				permissions_boundary = "sys-teamA-boundary"
			}
			`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `role "invalid" is invalid; it must start with "$team-", with team being one of: [teamA teamB ta tb]`,
				},
			},
		},
		{
			Name: "name missing",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }

			resource "aws_iam_role" "role" {
				permissions_boundary = "sys-teamA-boundary"
			}
			`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `role is missing required "name"; it must start with "$team-", with team being one of: [teamA teamB ta tb]`,
				},
			},
		},
		{
			Name: "user with path is invalid",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }

			resource "aws_iam_user" "user" {
				name                 = "teamA-admin"
				path                 = "/engineering/"
				permissions_boundary = "sys-teamA-boundary"
			}
			`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `user must not set "path".`,
				},
			},
		},
		{
			Name: "user with permissions boundary missing",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }

			resource "aws_iam_user" "user" {
				name = "teamA-admin"
			}
			`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `user is missing required "permissions_boundary"; it must be set to "arn:aws:iam::${var.account_id}:policy/sys-$team-boundary", with team being one of: [teamA teamB].`,
				},
			},
		},
		{
			Name: "role with permissions boundary missing",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }

			resource "aws_iam_role" "role" {
				name = "teamA-role"
			}
			`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `role is missing required "permissions_boundary"; it must be set to "arn:aws:iam::${var.account_id}:policy/sys-$team-boundary", with team being one of: [teamA teamB].`,
				},
			},
		},
		{
			Name: "invalid permissions boundary",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }

			resource "aws_iam_role" "role" {
				name                 = "teamA-role"
				permissions_boundary = "invalid"
			}
			`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `role has invalid "permissions_boundary" "invalid"; it must be set to "arn:aws:iam::${var.account_id}:policy/sys-$team-boundary", with team being one of: [teamA teamB].`,
				},
			},
		},
		{
			Name: "permissions boundary allows any account",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }

			resource "aws_iam_role" "role" {
				name                 = "teamA-role"
				permissions_boundary = "arn:aws:iam::XXXXXXXXXXXX:policy/sys-teamB-boundary"
			}
			`,
			Expected: helper.Issues{},
		},
		{
			Name: "permissions boundary not statically evaluable doesn't panic",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }
			variable "boundary" {}

			resource "aws_iam_user" "user" {
				name                 = "teamA-admin"
				permissions_boundary = var.boundary
			}
			`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `user permissions_boundary must be statically resolvable.`,
				},
			},
		},
		{
			Name: "name not statically evaluable doesn't panic",
			Content: `
			variable "environment" { default = "dev" }
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["ta", "tb"] }
			variable "iam_name" {}

			resource "aws_iam_group" "group" {
				name = var.iam_name
			}
			`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `group name permissions_boundary must be statically resolvable.`,
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
