package rules

import (
	"slices"
	"testing"

	"github.com/terraform-linters/tflint-plugin-sdk/helper"
)

func Test_GetEnv(t *testing.T) {
	tests := []struct {
		Name      string
		Content   string
		Filename  string
		Expected  string
		WantError bool
	}{
		{
			Name:      "no issues",
			Content:   `variable "environment" { default = "dev"}`,
			Filename:  "resource.tf",
			Expected:  "dev",
			WantError: false,
		},
		{
			Name:      "detect environment if it's dev/prod",
			Content:   ``,
			Filename:  "aws/prod/resource.tf",
			Expected:  "prod",
			WantError: false,
		},
		{
			Name:      "detect environment only if not explicitely set",
			Content:   `variable "environment" { default = "dev"}`,
			Filename:  "aws/prod/resource.tf",
			Expected:  "dev",
			WantError: false,
		},
		{
			Name:      "environment missing",
			Content:   ``,
			Filename:  "resource.tf",
			Expected:  "",
			WantError: true,
		},
		{
			Name:      "environment without default value",
			Content:   `variable "environment" {}`,
			Filename:  "resource.tf",
			Expected:  "",
			WantError: true,
		},
		{
			Name:      "environment empty",
			Content:   `variable "environment" { default = ""}`,
			Filename:  "resource.tf",
			Expected:  "",
			WantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{test.Filename: test.Content})

			result, err := GetEnv(runner)

			if test.WantError && err == nil {
				t.Fatal("Expected an error but got nil")
			}
			if !test.WantError && err != nil {
				t.Fatalf("Expected no error but got: %s", err)
			}
			if result != test.Expected {
				t.Errorf("Got %q, want %q", result, test.Expected)
			}

		})
	}
}

func Test_GetTeams(t *testing.T) {
	tests := []struct {
		Name      string
		Content   string
		Expected  []string
		WantError bool
	}{
		{
			Name: "no issues",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = ["teamA", "teamB"] }
			`,
			Expected:  []string{"teamA", "teamB"},
			WantError: false,
		},
		{
			Name: "environment missing",
			Content: `
			variable "teams_dev" { default = ["teamA", "teamB"] }
			`,
			Expected:  nil,
			WantError: true,
		},
		{
			Name: "teams missing",
			Content: `
			variable "environment" { default = "dev"}
			`,
			Expected:  nil,
			WantError: true,
		},
		{
			Name: "teams without default value",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" {}
			`,
			Expected:  nil,
			WantError: true,
		},
		{
			Name: "teams empty",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_dev" { default = [] }
			`,
			Expected:  nil,
			WantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			result, err := GetTeams(runner)

			if test.WantError && err == nil {
				t.Fatal("Expected an error but got nil")
			}
			if !test.WantError && err != nil {
				t.Fatalf("Expected no error but got: %s", err)
			}
			if (result == nil) && (test.Expected != nil) {
				t.Errorf("Got nil result, want non-nil")
			}
			if (result != nil) && (test.Expected == nil) {
				t.Errorf("Got non-nil result, want nil")
			}
			if !slices.Equal(result, test.Expected) {
				t.Errorf("Got %q, want %q", result, test.Expected)
			}

		})
	}
}

func Test_GetTeamPrefixes(t *testing.T) {
	tests := []struct {
		Name      string
		Content   string
		Expected  []string
		WantError bool
	}{
		{
			Name: "no issues",
			Content: `
			variable "environment" { default = "dev"}
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			`,
			Expected:  []string{"tA", "tB"},
			WantError: false,
		},
		{
			Name: "environment missing",
			Content: `
			variable "team_prefixes_dev" { default = ["tA", "tB"] }
			`,
			Expected:  nil,
			WantError: true,
		},
		{
			Name: "teams_prefixes missing",
			Content: `
			variable "environment" { default = "dev"}
			`,
			Expected:  nil,
			WantError: true,
		},
		{
			Name: "teams_prefixes without default value",
			Content: `
			variable "environment" { default = "dev"}
			variable "team_prefixes_dev" {}
			`,
			Expected:  nil,
			WantError: true,
		},
		{
			Name: "teams_prefixes empty",
			Content: `
			variable "environment" { default = "dev"}
			variable "teams_prefixes_dev" { default = [] }
			`,
			Expected:  nil,
			WantError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			result, err := GetTeamPrefixes(runner)

			if test.WantError && err == nil {
				t.Fatal("Expected an error but got nil")
			}
			if !test.WantError && err != nil {
				t.Fatalf("Expected no error but got: %s", err)
			}
			if (result == nil) && (test.Expected != nil) {
				t.Errorf("Got nil result, want non-nil")
			}
			if (result != nil) && (test.Expected == nil) {
				t.Errorf("Got non-nil result, want nil")
			}
			if !slices.Equal(result, test.Expected) {
				t.Errorf("Got %q, want %q", result, test.Expected)
			}
		})
	}
}

func Test_GetBucketNameExceptions(t *testing.T) {
	tests := []struct {
		Name     string
		Content  string
		Expected []string
	}{
		{
			Name: "exceptions have values",
			Content: `
			variable "environment" { default = "dev"}
			variable "bucket_name_exceptions_dev" { default = ["legacy", "other"] }
			`,
			Expected: []string{"legacy", "other"},
		},
		{
			Name: "environment missing",
			Content: `
			variable "bucket_name_exceptions_dev" { default = ["legacy", "other"] }
			`,
			Expected: []string{},
		},
		{
			Name: "exceptions missing",
			Content: `
			variable "environment" { default = "dev"}
			`,
			Expected: []string{},
		},
		{
			Name: "exceptions without default",
			Content: `
			variable "environment" { default = "dev"}
			variable "bucket_name_exceptions_dev" {}
			`,
			Expected: []string{},
		},
		{
			Name: "exceptions empty",
			Content: `
			variable "environment" { default = "dev"}
			variable "bucket_name_exceptions_dev" { default = [] }
			`,
			Expected: []string{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			runner := helper.TestRunner(t, map[string]string{"resource.tf": test.Content})

			result := GetBucketNameExceptions(runner)
			if (result == nil) && (test.Expected != nil) {
				t.Errorf("Got nil result, want non-nil")
			}
			if (result != nil) && (test.Expected == nil) {
				t.Errorf("Got non-nil result, want nil")
			}
			if !slices.Equal(result, test.Expected) {
				t.Errorf("Got %q, want %q", result, test.Expected)
			}
		})
	}
}
