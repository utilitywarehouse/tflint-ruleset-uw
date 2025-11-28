package rules

import (
	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

// AwsS3BucketMissingNameTagRule checks whether a security group has
// the Name tag, needed by teams to manage the security group
type AwsS3BucketMissingNameTagRule struct {
	tflint.DefaultRule
}

// Name returns the rule name.
func (r *AwsS3BucketMissingNameTagRule) Name() string {
	return "AwsS3BucketMissingNameTag"
}

// Enabled returns whether the rule is enabled by default.
func (r *AwsS3BucketMissingNameTagRule) Enabled() bool {
	return true
}

// Severity returns the rule severity.
func (r *AwsS3BucketMissingNameTagRule) Severity() tflint.Severity {
	return tflint.ERROR
}

// Link returns the rule reference link.
func (r *AwsS3BucketMissingNameTagRule) Link() string {
	return ReferenceLink(r.Name())
}

func (r *AwsS3BucketMissingNameTagRule) Check(runner tflint.Runner) error {
	buckets, err := runner.GetResourceContent(
		"aws_s3_bucket",
		&hclext.BodySchema{Attributes: []hclext.AttributeSchema{{Name: "tags"}}},
		nil)

	if err != nil {
		return err
	}

	for _, bucket := range buckets.Blocks {
		tagsAttr, exists := bucket.Body.Attributes["tags"]
		if !exists {
			runner.EmitIssue(
				r,
				"Bucket is missing the required \"Name\" tag.",
				bucket.DefRange,
			)
			continue
		}

		var tags map[string]string
		err := runner.EvaluateExpr(tagsAttr.Expr, &tags, nil)
		if err != nil {
			continue
		}

		value, ok := tags["Name"]
		if !ok {
			runner.EmitIssue(
				r,
				"Bucket is missing the required \"Name\" tag.",
				tagsAttr.Expr.Range(),
			)
			continue
		}
		if value == "" {
			runner.EmitIssue(
				r,
				"Bucket has an empty required \"Name\" tag.",
				tagsAttr.Expr.Range(),
			)
		}
	}

	return nil
}
