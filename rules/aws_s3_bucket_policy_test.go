package rules

import (
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketPolicy(t *testing.T) {
	rule := &AwsS3BucketPolicy{}

	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name:     "no issues when no buckets",
			Content:  "",
			Expected: helper.Issues{},
		},
		{
			Name: "no issues outside of dev/prod accounts",
			Content: `
			variable "environment" { default = "other"}
			resource "aws_s3_bucket" "bucket" {}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no issues when complying to policy",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			resource "aws_s3_bucket" "bucket" {
				bucket = "uw-dev-teamA-xx"
				tags = {
					Name="any"
				}
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "bucket attribute (aka name) missing",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			resource "aws_s3_bucket" "bucket" {}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: "Bucket is missing the required \"bucket\" attribute.",
				},
			},
		},
		{
			Name: "bucket attribute (aka name) invalid",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			resource "aws_s3_bucket" "bucket" {
				bucket = "invalid"
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `bucket name "invalid" is invalid; must start with 'uw-$env-$team-', with team being one of [teamA teamB tA tB]`,
				},
			},
		},
		{
			Name: "bucket attribute (aka name) invalid but in exceptions",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			variable "bucket_name_exceptions_dev" { default = ["invalid"] }
			resource "aws_s3_bucket" "bucket" {
				bucket = "invalid"
				tags = {
					Name="any"
				}
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "tags missing",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			resource "aws_s3_bucket" "bucket" {
				bucket = "uw-dev-teamA-name"
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Bucket is missing the required "Name" tag.`,
				},
			},
		},
		{
			Name: "tags not statically evaluable",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			variable "bucket_tags" {}
			resource "aws_s3_bucket" "bucket" {
				bucket = "uw-dev-teamA-name"
				tags   = var.bucket_tags
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Bucket tags could not be fully evaluated; ensure the required "Name" tag is statically resolvable.`,
				},
			},
		},
		{
			Name: "Name tag missing",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			resource "aws_s3_bucket" "bucket" {
				bucket = "uw-dev-teamA-name"
				tags = {
					other="tag"
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Bucket is missing the required "Name" tag.`,
				},
			},
		},
		{
			Name: "Name tag empty",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			resource "aws_s3_bucket" "bucket" {
				bucket = "uw-dev-teamA-name"
				tags = {
					Name=""
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Bucket has an empty required "Name" tag.`,
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
