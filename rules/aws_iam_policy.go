package rules

import (
	"fmt"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AwsIamPolicy struct {
	tflint.DefaultRule
}

func (r *AwsIamPolicy) Name() string {
	return "AwsIamPolicy"
}

func (r *AwsIamPolicy) Enabled() bool {
	return true
}

func (r *AwsIamPolicy) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AwsIamPolicy) Link() string {
	return ReferenceLink(r.Name())
}

func (r *AwsIamPolicy) Check(runner tflint.Runner) error {
	resources := []struct {
		resourceType    string
		label           string
		requireNoPath   bool
		requireBoundary bool
	}{
		{resourceType: "aws_iam_user", label: "user", requireBoundary: true, requireNoPath: true},
		{resourceType: "aws_iam_role", label: "role", requireBoundary: true},
		{resourceType: "aws_iam_group", label: "group"},
		{resourceType: "aws_iam_policy", label: "policy"},
		{resourceType: "aws_iam_instance_profile", label: "instance profile"},
	}

	var teams []string
	var allowedOwners []string

	for _, resource := range resources {
		content, err := runner.GetResourceContent(
			resource.resourceType,
			&hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: "name"}, {Name: "path"}, {Name: "permissions_boundary"}}},
			nil,
		)
		if err != nil {
			return err
		}

		if len(content.Blocks) == 0 {
			continue
		}

		if allowedOwners == nil {
			env, err := GetEnv(runner)
			if err != nil {
				return err
			}
			if env != "dev" && env != "prod" {
				return nil
			}

			teams, err = GetTeams(runner)
			if err != nil {
				return err
			}
			teamPrefixes, err := GetTeamPrefixes(runner)
			if err != nil {
				return err
			}
			allowedOwners = append(teams, teamPrefixes...)
		}

		for _, block := range content.Blocks {
			r.checkResourceName(runner, block, resource.label, allowedOwners)
			if resource.requireNoPath {
				r.checkPathNotSet(runner, block, resource.label)
			}
			if resource.requireBoundary {
				r.checkPermissionBoundary(runner, block, resource.label, teams)
			}
		}
	}

	return nil
}

func (r *AwsIamPolicy) checkPathNotSet(runner tflint.Runner, block *hclext.Block, label string) {
	if attr, ok := block.Body.Attributes["path"]; ok {
		runner.EmitIssue(
			r,
			fmt.Sprintf("%s must not set \"path\".", label),
			attr.Expr.Range(),
		)
	}
}

func (r *AwsIamPolicy) checkPermissionBoundary(runner tflint.Runner, block *hclext.Block, label string, teams []string) {
	if _, ok := block.Body.Attributes["permissions_boundary"]; !ok {
		runner.EmitIssue(
			r,
			fmt.Sprintf("%s is missing required \"permissions_boundary\"; it must be set to \"arn:aws:iam::${var.account_id}:policy/sys-$team-boundary\", with team being one of: %v.", label, teams),
			block.DefRange,
		)
		return
	}

	boundary, err := safelyEvaluateStringAttr(runner, block, "permissions_boundary")

	if err != nil {
		runner.EmitIssue(
			r,
			fmt.Sprintf("%s permissions_boundary must be statically resolvable.", label),
			block.Body.Attributes["permissions_boundary"].Expr.Range(),
		)
		return
	}

	for _, team := range teams {
		expectedName := fmt.Sprintf("sys-%s-boundary", team)
		if boundary == expectedName || strings.HasSuffix(boundary, "/"+expectedName) {
			return
		}
	}

	runner.EmitIssue(
		r,
		fmt.Sprintf("%s has invalid \"permissions_boundary\" %q; it must be set to \"arn:aws:iam::${var.account_id}:policy/sys-$team-boundary\", with team being one of: %v.", label, boundary, teams),
		block.Body.Attributes["permissions_boundary"].Expr.Range(),
	)
}

func (r *AwsIamPolicy) checkResourceName(runner tflint.Runner, block *hclext.Block, label string, allowedOwners []string) {
	nameRequirement := fmt.Sprintf("it must start with \"$team-\", with team being one of: %v", allowedOwners)
	if _, ok := block.Body.Attributes["name"]; !ok {
		runner.EmitIssue(
			r,
			fmt.Sprintf("%s is missing required \"name\"; %s", label, nameRequirement),
			block.DefRange,
		)
		return
	}

	name, err := safelyEvaluateStringAttr(runner, block, "name")

	if err != nil {
		runner.EmitIssue(
			r,
			fmt.Sprintf("%s name permissions_boundary must be statically resolvable.", label),
			block.Body.Attributes["name"].Expr.Range(),
		)
		return
	}

	for _, owner := range allowedOwners {
		if strings.HasPrefix(name, owner+"-") {
			return
		}
	}

	runner.EmitIssue(
		r,
		fmt.Sprintf("%s %q is invalid; %s", label, name, nameRequirement),
		block.Body.Attributes["name"].Expr.Range(),
	)
}

func safelyEvaluateStringAttr(runner tflint.Runner, block *hclext.Block, name string) (value string, err error) {
	attr, ok := block.Body.Attributes[name]
	if !ok {
		return "", nil
	}

	defer func() {
		// EvaluateExpr can panic on some unevaluable Terraform expressions.
		// Recover here so the rule can report a lint issue instead of aborting.
		if recovered := recover(); recovered != nil {
			value = ""
			err = fmt.Errorf("failed to evaluate expression")
		}
	}()

	// Use named returns so the deferred recovery can replace the result on panic.
	err = runner.EvaluateExpr(attr.Expr, &value, nil)
	return value, err
}
