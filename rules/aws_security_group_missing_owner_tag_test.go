package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsSecurityGroupMissingOwnerTagRule(t *testing.T) {
	rule := &AwsSecurityGroupMissingOwnerTagRule{}

	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "no issues",
			Content: `
			resource "aws_security_group" "sg" {
				tags = {
					owner="team"
				}
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no tags",
			Content: `
			resource "aws_security_group" "sg" {
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Security group is missing the required "owner" tag.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 38},
					},
				},
			},
		},
		{
			Name: "owner tag missing",
			Content: `
			resource "aws_security_group" "sg" {
				tags = {
					other="tag"
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Security group is missing the required "owner" tag.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 5, Column: 6},
					},
				},
			},
		},
		{
			Name: "owner tag empty",
			Content: `
			resource "aws_security_group" "sg" {
				tags = {
					owner=""
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Security group has an empty required "owner" tag.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 5, Column: 6},
					},
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

			helper.AssertIssues(t, test.Expected, runner.Issues)
		})
	}
}
