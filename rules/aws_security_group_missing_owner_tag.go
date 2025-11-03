package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsSecurityGroupMissingOwnerTagRule checks whether a security group has
// the owner tag, needed by teams to manage the security group
type AwsSecurityGroupMissingOwnerTagRule struct {
	tflint.DefaultRule
}

// Name returns the rule name.
func (r *AwsSecurityGroupMissingOwnerTagRule) Name() string {
	return "AwsSecurityGroupMissingOwnerTag"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsSecurityGroupMissingOwnerTagRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsSecurityGroupMissingOwnerTagRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsSecurityGroupMissingOwnerTagRule) Link() string {
	return ReferenceLink(r.Name())
}

func (r *AwsSecurityGroupMissingOwnerTagRule) Check(runner tflint.Runner) error {
	securityGroups, err := runner.GetResourceContent(
		"aws_security_group",
		&hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: "tags"}}},
		nil)

	if err != nil {
		return err
	}

	for _, securityGroup := range securityGroups.Blocks {
		tagsAttr, exists := securityGroup.Body.Attributes["tags"]
		if !exists {
			runner.EmitIssue(
				r,
				"Security group is missing the required \"owner\" tag.",
				securityGroup.DefRange,
			)
			continue
		}

		var tags map[string]string
		err := runner.EvaluateExpr(tagsAttr.Expr, &tags, nil)
		if err != nil {
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
		}
	}

	return nil
}
