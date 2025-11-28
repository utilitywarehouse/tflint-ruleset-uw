package rules

import (
	"testing"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_AwsS3BucketMissingNameTagRule(t *testing.T) {
	rule := &AwsS3BucketMissingNameTagRule{}

	tests := []struct {
		Name     string
		Content  string
		Expected helper.Issues
	}{
		{
			Name: "no issues",
			Content: `
			resource "aws_s3_bucket" "bucket" {
				tags = {
					Name="bucket"
				}
			}`,
			Expected: helper.Issues{},
		},
		{
			Name: "no tags",
			Content: `
			resource "aws_s3_bucket" "bucket" {
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Bucket is missing the required "Name" tag.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 2, Column: 4},
						End:      hcl.Pos{Line: 2, Column: 37},
					},
				},
			},
		},
		{
			Name: "Name tag missing",
			Content: `
			resource "aws_s3_bucket" "bucket" {
				tags = {
					other="tag"
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Bucket is missing the required "Name" tag.`,
					Range: hcl.Range{
						Filename: "resource.tf",
						Start:    hcl.Pos{Line: 3, Column: 12},
						End:      hcl.Pos{Line: 5, Column: 6},
					},
				},
			},
		},
		{
			Name: "Name tag empty",
			Content: `
			resource "aws_s3_bucket" "bucket" {
				tags = {
					Name=""
				}
			}`,
			Expected: helper.Issues{
				{
					Rule:    rule,
					Message: `Bucket has an empty required "Name" tag.`,
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
