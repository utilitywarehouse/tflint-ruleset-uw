package rules

import (
	"fmt"
	"slices"
	"strings"

	"github.com/terraform-linters/tflint-plugin-sdk/hclext"
	"github.com/terraform-linters/tflint-plugin-sdk/tflint"
)

type AwsS3BucketPolicy struct {
	tflint.DefaultRule
}

func (r *AwsS3BucketPolicy) Name() string {
	return "AwsS3BucketPolicy"
}

func (r *AwsS3BucketPolicy) Enabled() bool {
	return true
}

func (r *AwsS3BucketPolicy) Severity() tflint.Severity {
	return tflint.ERROR
}

func (r *AwsS3BucketPolicy) Link() string {
	return ReferenceLink(r.Name())
}

func (r *AwsS3BucketPolicy) checkName(name, env string, teams []string, exceptions []string) error {
	if slices.Contains(exceptions, name) {
		return nil
	}

	for _, team := range teams {
		prefixWithEnv := fmt.Sprintf("uw-%s-%s-", env, team)
		prefixWithoutEnv := fmt.Sprintf("uw-%s-", team)

		if strings.HasPrefix(name, prefixWithEnv) || strings.HasPrefix(name, prefixWithoutEnv) {
			return nil // Found a valid prefix
		}
	}
	return fmt.Errorf("bucket name %q is invalid; must start with 'uw-$env-$team-', with team being one of %v",
		name, teams)
}

func (r *AwsS3BucketPolicy) Check(runner tflint.Runner) error {
	buckets, err := runner.GetResourceContent(
		"aws_s3_bucket",
		&hclext.BodySchema{Attributes: []hclext.AttributeSchema{
			{Name: "bucket"}, // name of the bucket
			{Name: "tags"},
		}},
		nil)

	if err != nil {
		return err
	}

	for _, bucket := range buckets.Blocks {
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
		team_prefixes, err := GetTeamPrefixes(runner)
		if err != nil {
			return err
		}
		exceptions := GetBucketNameExceptions(runner)

		bucketAttr, ok := bucket.Body.Attributes["bucket"]
		if !ok {
			runner.EmitIssue(
				r,
				"Bucket is missing the required \"bucket\" attribute.",
				bucket.DefRange,
			)
			continue
		}
		var name string
		err = runner.EvaluateExpr(bucketAttr.Expr, &name, nil)
		if err != nil {
			return err
		}
		err = r.checkName(name, env, append(teams, team_prefixes...), exceptions)
		if err != nil {
			runner.EmitIssue(
				r,
				err.Error(),
				bucketAttr.Expr.Range(),
			)
			continue
		}

		tagsAttr, ok := bucket.Body.Attributes["tags"]
		if !ok {
			runner.EmitIssue(
				r,
				"Bucket is missing the required \"Name\" tag.",
				bucket.DefRange,
			)
			continue
		}

		var tags map[string]string
		err = runner.EvaluateExpr(tagsAttr.Expr, &tags, nil)
		if err != nil {
			runner.EmitIssue(
				r,
				"Bucket tags could not be fully evaluated; ensure the required \"Name\" tag is statically resolvable.",
				tagsAttr.Expr.Range(),
			)
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
