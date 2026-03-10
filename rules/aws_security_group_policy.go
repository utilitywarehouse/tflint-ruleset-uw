package rules

import (
	"fmt"
	"slices"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AwsSecurityGroupPolicy struct {
	tflint.DefaultRule
}

func (r *AwsSecurityGroupPolicy) Name() string {
	return "AwsSecurityGroupPolicy"
}

func (r *AwsSecurityGroupPolicy) Enabled() bool {
	return true
}

func (r *AwsSecurityGroupPolicy) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AwsSecurityGroupPolicy) Link() string {
	return ReferenceLink(r.Name())
}

func (r *AwsSecurityGroupPolicy) Check(runner tflint.Runner) error {
	securityGroups, err := runner.GetResourceContent(
		"aws_security_group",
		&hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: "tags"}}},
		nil)

	if err != nil {
		return err
	}

	for _, securityGroup := range securityGroups.Blocks {
		env, err := GetEnv(runner)
		if err != nil {
			return err
		}
		if env != "dev" && env != "prod" {
			return nil
		}
		teams, err := GetTeams(runner)
		if err != nil {
			return err
		}
		teamPrefixes, err := GetTeamPrefixes(runner)
		if err != nil {
			return err
		}
		allowedOwners := append(teams, teamPrefixes...)

		tagsAttr, ok := securityGroup.Body.Attributes["tags"]
		if !ok {
			runner.EmitIssue(
				r,
				"Security group is missing the required \"owner\" tag.",
				securityGroup.DefRange,
			)
			continue
		}

		var tags map[string]string
		err = runner.EvaluateExpr(tagsAttr.Expr, &tags, nil)
		if err != nil {
			runner.EmitIssue(
				r,
				"Security group tags could not be fully evaluated; ensure the required \"owner\" tag is statically resolvable.",
				tagsAttr.Expr.Range(),
			)
			continue
		}

		value, ok := tags["owner"]
		if !ok {
			runner.EmitIssue(
				r,
				"Security group is missing the required \"owner\" tag.",
				tagsAttr.Expr.Range(),
			)
			continue
		}
		if value == "" {
			runner.EmitIssue(
				r,
				"Security group has an empty required \"owner\" tag.",
				tagsAttr.Expr.Range(),
			)
			continue
		}

		if !slices.Contains(allowedOwners, value) {
			runner.EmitIssue(
				r,
				fmt.Sprintf("Security group has an invalid \"owner\" tag; Found \"%s\", but it has to be one of: %v", value, teams),
				tagsAttr.Expr.Range(),
			)
		}
	}

	return nil
}
